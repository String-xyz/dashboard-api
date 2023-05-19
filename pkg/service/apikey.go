package service

import (
	"context"
	"fmt"

	"github.com/String-xyz/go-lib/v2/common"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Apikey interface {
	Create(ctx context.Context, keyType string, callerId string, platformId string, organizationId string) (model.Apikey, error)
	GetAll(ctx context.Context, callerId string, platformId string, organizationId string, limit int, offset int) ([]model.Apikey, error)
	Get(ctx context.Context, id string, callerId string, organizationId string) (model.Apikey, error)
	Delete(ctx context.Context, keyId string, callerId string, organizationId string) error
	Update(ctx context.Context, keyId string, request model.RequestApikeyUpdate, callerId string, organizationId string) (model.Apikey, error)
}

type apikey struct {
	repos repository.Repositories
}

func NewApikey(repos repository.Repositories) Apikey {
	return &apikey{repos}
}

func (a *apikey) Create(ctx context.Context, keyType string, callerId string, platformId string, organizationId string) (model.Apikey, error) {
	_, finish := Span(ctx, "service.apikey.Create", SpanTag{"organizationId": organizationId})
	defer finish()

	// Create the base key object
	var keyValue string
	key := model.Apikey{
		Type:           keyType,
		PlatformId:     &platformId,
		OrganizationId: organizationId,
		CreatedBy:      callerId,
	}

	// Only admins can create secret keys
	if keyType == "secret" {
		err := RequireAuthority(a.repos, callerId, "Owner", "Admin")
		if err != nil {
			return model.Apikey{}, err
		}
	}

	// platform must be active
	err := a.repos.Platform.CheckActive(ctx, platformId)
	if err != nil {
		return model.Apikey{}, common.StringError(serror.FORBIDDEN)
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

func (a apikey) GetAll(ctx context.Context, callerId string, platformId string, organizationId string, limit int, offset int) (keys []model.Apikey, err error) {
	_, finish := Span(ctx, "service.apikey.GetAll", SpanTag{"organizationId": organizationId})
	defer finish()

	if platformId != "" {
		keys, err = a.repos.Apikey.ListByPlatform(ctx, platformId, limit, offset)
		if err != nil {
			return keys, common.StringError(err)
		}
	} else {
		keys, err = a.repos.Apikey.ListByOrganization(ctx, organizationId, limit, offset)
		if err != nil {
			return keys, common.StringError(err)
		}
	}

	return keys, nil
}

func (a apikey) Get(ctx context.Context, id string, callerId string, organizationId string) (model.Apikey, error) {
	_, finish := Span(ctx, "service.apikey.Get", SpanTag{"organizationId": organizationId})
	defer finish()

	key, err := a.repos.Apikey.GetById(ctx, id)
	if err != nil {
		return model.Apikey{}, common.StringError(err)
	}
	if key.OrganizationId != organizationId {
		return model.Apikey{}, common.StringError(fmt.Errorf("key not maintained by accessing organization %s", organizationId))
	}
	return key, nil
}

func (a apikey) Delete(ctx context.Context, keyId string, callerId string, organizationId string) error {
	_, finish := Span(ctx, "service.apikey.Deactivate", SpanTag{"organizationId": organizationId})
	defer finish()

	err := RequireAuthority(a.repos, callerId, "Admin", "Owner")
	if err != nil {
		return common.StringError(err)
	}

	key, err := a.repos.Apikey.GetById(ctx, keyId)
	if err != nil {
		return common.StringError(err)
	}
	if key.OrganizationId != organizationId {
		return common.StringError(fmt.Errorf("key not maintained by accessing organization %s", organizationId))
	}

	err = a.repos.Apikey.SoftDelete(ctx, keyId)

	return err
}

func (a apikey) Update(ctx context.Context, keyId string, request model.RequestApikeyUpdate, callerId string, organizationId string) (key model.Apikey, err error) {
	_, finish := Span(ctx, "service.apikey.Update", SpanTag{"organizationId": organizationId})
	defer finish()

	key, err = a.repos.Apikey.GetById(ctx, keyId)
	if err != nil {
		return model.Apikey{}, common.StringError(err)
	}
	if key.OrganizationId != organizationId {
		return model.Apikey{}, common.StringError(fmt.Errorf("key not maintained by accessing organization %s", organizationId))
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
