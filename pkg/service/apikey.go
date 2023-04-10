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
	Create(ctx context.Context, callerId string, platformId string, keyType string) (model.Apikey, error)
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

func (a *apikey) Create(ctx context.Context, callerID, platformId, keyType string) (model.Apikey, error) {
	// Create the base key object
	var keyValue string
	key := model.Apikey{
		Type:       keyType,
		PlatformId: platformId,
		CreatedBy:  callerID,
	}

	// Only admins can create secret keys
	if keyType == "secret" {
		err := RequireAuthority(a.repos, callerID, "Owner", "Admin")
		if err != nil {
			return model.Apikey{}, err
		}
	}

	// Generate the key either as a secret or a public key
	if keyType == "secret" {
		keyValue = "strsk." + uuidWithoutHyphens()
		key.Data = common.ToSha256(keyValue)
	} else { // public
		keyValue = "str." + uuidWithoutHyphens()
		key.Data = keyValue
	}

	// Save the key
	key.Hint = keyValue[0:10]
	createdKey, err := a.repos.Apikey.Create(ctx, key)
	if err != nil {
		return model.Apikey{}, common.StringError(err)
	}

	// Return the key with the correct value (the secret for secret keys)
	if keyType == "secret" {
		createdKey.Data = keyValue
	}

	return createdKey, nil
}

func (a apikey) GetAll(ctx context.Context, callerId string, platformId string) (keys []model.Apikey, err error) {
	keys, err = a.repos.Apikey.List(ctx, platformId, 0, 0)
	if err != nil {
		return keys, common.StringError(err)
	}

	return keys, nil
}

func (a apikey) Get(ctx context.Context, callerId string, platformId string, id string) (model.Apikey, error) {
	return a.repos.Apikey.GetById(ctx, id)
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
