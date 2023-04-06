package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Apikey interface {
	Create(ctx context.Context, callerId string, platformId string, withSecret bool) (model.Apikey, error)
	GetAll(ctx context.Context, callerId string, platformId string) ([]model.Apikey, error)
	Get(ctx context.Context, callerId string, platformId string, id string) (model.Apikey, error)
	Deactivate(ctx context.Context, callerId string, platformId string, keyId string) (model.Apikey, error)
	Update(ctx context.Context, request model.RequestApikeyUpdate, callerId string, platformId string, keyId string) (model.Apikey, error)
}

type apikey struct {
	repos repository.Repositories
}

func NewApikey(repos repository.Repositories) Apikey {
	return &apikey{repos}
}

//	type Apikey struct {
//		ID            string     `json:"id,omitempty" db:"id"`
//		CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
//		UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
//		DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
//		Type          string     `json:"type" db:"type"`
//		Data          string     `json:"data" db:"data"`
//		Description   string     `json:"description" db:"description"`
//		CreatedBy     string     `json:"createdBy" db:"created_by"`
//		PlatformId    string     `json:"platformId" db:"platform_id"`
//	}
func (a apikey) Create(ctx context.Context, callerId string, platformId string, withSecret bool) (key model.Apikey, err error) {

	publicKey := "str." + uuidWithoutHyphens()
	key = model.Apikey{Type: "public", Public: publicKey, PlatformId: platformId, CreatedBy: callerId}
	secretKey := "strsk." + uuidWithoutHyphens()

	if withSecret {
		hash := common.ToSha256(secretKey)
		key.Secret = &hash
	}

	key, err = a.repos.Apikey.Create(ctx, key)
	if err != nil {
		return key, common.StringError(err)
	}

	// save a hash but return the secret
	if withSecret {
		key.Secret = &secretKey
	}

	return key, nil
}

func (a apikey) GetAll(ctx context.Context, callerId string, platformId string) (keys []model.Apikey, err error) {
	keys, err = a.repos.Apikey.List(ctx, platformId, 0, 0)
	if err != nil {
		return keys, common.StringError(err)
	}

	return keys, nil
}

func (a apikey) Get(ctx context.Context, callerId string, platformId string, id string) (model.Apikey, error) {
	return model.Apikey{}, nil
}

func (a apikey) Deactivate(ctx context.Context, callerId string, platformId string, keyId string) (key model.Apikey, err error) {
	err = RequireAuthority(a.repos, callerId, "Admin", "Owner")
	if err != nil {
		return key, common.StringError(err)
	}

	key, err = a.repos.Apikey.GetById(ctx, keyId)
	if err != nil {
		return model.Apikey{}, common.StringError(err)
	}
	if key.PlatformId != platformId {
		return model.Apikey{}, common.StringError(fmt.Errorf("key not maintained by accessing platform %s", platformId))
	}

	type DeactivateUpdate struct {
		DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	}

	now := time.Now()
	update := DeactivateUpdate{DeactivatedAt: &now}

	err = a.repos.Apikey.Update(ctx, keyId, update)
	if err != nil {
		return key, common.StringError(err)
	}

	key, err = a.repos.Apikey.GetById(ctx, keyId)
	if err != nil {
		return key, common.StringError(err)
	}

	return key, nil
}

func (a apikey) Update(ctx context.Context, request model.RequestApikeyUpdate, callerId string, platformId string, keyId string) (key model.Apikey, err error) {
	log.Printf("\n\nhit the update service! with keyId: %+v\n", keyId)
	key, err = a.repos.Apikey.GetById(ctx, keyId)
	if err != nil {
		return model.Apikey{}, common.StringError(err)
	}
	if key.PlatformId != platformId {
		return model.Apikey{}, common.StringError(fmt.Errorf("key not maintained by accessing platform %s", platformId))
	}
	type KeyUpdate struct {
		Description string `json:"description" db:"description"`
	}

	update := KeyUpdate{Description: request.Description}

	err = a.repos.Apikey.Update(ctx, keyId, update)
	if err != nil {
		return key, common.StringError(err)
	}

	key, err = a.repos.Apikey.GetById(ctx, keyId)
	if err != nil {
		return key, common.StringError(err)
	}

	return key, nil
}
