package repository

import (
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	serrors "github.com/String-xyz/go-lib/stringerror"
)

const KEY_PREFIX = "deny_"

type DenyList interface {
	AddMember(memberId string) error
	RemoveMember(memberId string) error
	IsDenied(memberId string) (bool, error)
}

type denyList struct {
	redis database.RedisStore
}

func NewDenyList(r database.RedisStore) DenyList {
	return &denyList{r}
}

func (a denyList) AddMember(memberId string) error {
	expireAt := time.Hour * 24 * 7 // 7 days expiration
	key := KEY_PREFIX + memberId
	return a.redis.Set(key, "true", expireAt)
}

func (a denyList) RemoveMember(memberId string) error {
	key := KEY_PREFIX + memberId
	return a.redis.Delete(key)
}

func (a denyList) IsDenied(memberId string) (bool, error) {
	key := KEY_PREFIX + memberId
	denied, err := a.redis.Get(key)
	if err != nil {
		if serrors.ErrorIs(err, serrors.NOT_FOUND) {
			return false, nil // if we error for this reason, continue
		}
		return false, common.StringError(err)
	}
	if string(denied) == "true" {
		return true, nil
	}
	return false, nil
}
