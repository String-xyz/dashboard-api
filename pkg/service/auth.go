package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	netmail "net/mail"
	"os"
	"strings"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"

	"github.com/String-xyz/platform-admin-api/pkg/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWT struct {
	ExpAt        time.Time            `json:"expAt"`
	IssuedAt     time.Time            `json:"issuedAt"`
	Token        string               `json:"token"`
	RefreshToken RefreshTokenResponse `json:"refreshToken"`
}

type JWTClaims struct {
	MemberId   string
	PlatformId string
	jwt.StandardClaims
}

type JWTStrategy struct {
	ID            string     `json:"id,omitempty" db:"id"`
	Type          string     `json:"authType" db:"type"`
	EntityType    string     `json:"entityType,omitempty"`
	Data          string     `json:"data" data:"data"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	ExpiresAt     time.Time  `json:"expireAt,omitempty" db:"expire_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
}

type RefreshTokenResponse struct {
	Token string    `json:"token"`
	ExpAt time.Time `json:"expAt"`
}

type Auth interface {
	GenerateJWT(memberId string, platformId string) (JWT, error)
	ValidateAPIKey(key string) bool
	RefreshToken(token string, walletAddress string) (MemberCreateResponse, error)
	InvalidateRefreshToken(token string) error
}

type auth struct {
	repos repository.Repositories
	redis database.RedisStore
}

// GenerateJWT generates a jwt token and a refresh token which is saved on redis
func (a auth) GenerateJWT(memberId string, platformId string) (JWT, error) {
	claims := JWTClaims{}
	refreshToken := uuidWithoutHyphens()
	t := &JWT{
		IssuedAt: time.Now(),
		ExpAt:    time.Now().Add(time.Minute * 15),
	}

	claims.MemberId = memberId
	claims.PlatformId = platformId
	claims.ExpiresAt = t.ExpAt.Unix()
	claims.IssuedAt = t.IssuedAt.Unix()
	// replace this signing method with RSA or something similar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return *t, err
	}
	t.Token = signed

	// create and save
	refreshObj, err := a.CreateJWTRefresh(toSha256(refreshToken), memberId)
	if err != nil {
		return *t, err
	}
	t.RefreshToken = RefreshTokenResponse{
		Token: refreshToken,
		ExpAt: refreshObj.ExpiresAt,
	}

	return *t, nil
}

// CreateJWTRefresh creates and persists a refresh jwt token
func (a auth) CreateJWTRefresh(key string, memberId string) (JWTStrategy, error) {
	expireAt := time.Hour * 24 * 7 // 7 days expiration
	m := JWTStrategy{
		ID:         key,
		CreatedAt:  time.Now(),
		Type:       "JWT",
		EntityType: "Member",
		Data:       memberId,
		ExpiresAt:  time.Now().Add(expireAt),
	}

	return m, a.redis.Set(key, m, expireAt)
}

func (a auth) ValidateJWT(token string) (bool, error) {
	var claims = &JWTClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	return t.Valid, err
}

func (a auth) InvalidateRefreshToken(refreshToken string) error {
	return a.Delete(toSha256(refreshToken))
}

func (a auth) RefreshToken(refreshToken string, platformId string) (MemberCreateResponse, error) {
	resp := MemberCreateResponse{}

	// get user id from refresh token
	memberId, err := a.GetUserIdFromRefreshToken(toSha256(refreshToken))
	if err != nil {
		return resp, common.StringError(err)
	}

	// get platforms associated with user
	platforms, err := a.repos.MemberToPlatform.GetByMember(memberId)
	if err != nil {
		return resp, common.StringError(err)
	}
	if len(platforms) == 0 {
		return resp, common.StringError(errors.New("No platforms associated with member"))
	}

	// create new jwt
	jwt, err := a.GenerateJWT(memberId, platforms[0].PlatformID) // TODO: use the correct platform, not the first
	if err != nil {
		return resp, common.StringError(err)
	}
	resp.JWT = jwt

	// delete old refresh token
	err = a.InvalidateRefreshToken(refreshToken)
	if err != nil {
		return resp, common.StringError(err)
	}

	ctx := context.Background() // TODO: allow context into this function
	user, err := a.repos.PlatformMember.GetById(ctx, memberId)
	if err != nil {
		return resp, common.StringError(err)
	}
	resp.User = user

	return resp, nil
}

// return the user id from the refresh token or error if token is invalid or expired
func (a auth) GetUserIdFromRefreshToken(refreshToken string) (string, error) {
	authStrat, err := a.Get(refreshToken)

	if err != nil {
		return "", common.StringError(err)
	}
	// assert token has not expired
	if authStrat.ExpiresAt.Before(time.Now()) {
		return "", common.StringError(fmt.Errorf("refresh token expired"))
	}
	// assert token has not been deactivated
	if authStrat.DeactivatedAt != nil {
		return "", common.StringError(fmt.Errorf("refresh token deactivated at %s", authStrat.DeactivatedAt))
	}
	// if all is well, return the user id
	return authStrat.Data, nil
}

func (a auth) Get(key string) (JWTStrategy, error) {
	m, err := a.redis.Get(key)
	if err != nil {
		return JWTStrategy{}, common.StringError(err)
	}
	authStrat := JWTStrategy{}
	err = json.Unmarshal(m, &authStrat)
	if err != nil {
		return JWTStrategy{}, common.StringError(err)
	}

	return authStrat, nil
}

func (a auth) Delete(key string) error {
	return a.redis.Delete(key)
}

// Use native mail package to check if email a valid email
func validEmail(email string) bool {
	_, err := netmail.ParseAddress(email)
	return err == nil
}

func uuidWithoutHyphens() string {
	s := uuid.New().String()
	return strings.Replace(s, "-", "", -1)
}

func toSha256(v string) string {
	bs := sha256.Sum256([]byte(v))
	return hex.EncodeToString(bs[:])
}
