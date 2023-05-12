package repository

import (
	"time"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
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

func (d denyList) AddMember(memberId string) error {
	expireAt := time.Hour * 24 * 7 // 7 days expiration
	key := KEY_PREFIX + memberId
	return d.redis.Set(key, "true", expireAt)
}

func (d denyList) RemoveMember(memberId string) error {
	key := KEY_PREFIX + memberId
	return d.redis.Delete(key)
}

func (d denyList) IsDenied(memberId string) (bool, error) {
	key := KEY_PREFIX + memberId
	denied, err := d.redis.Get(key)
	if err != nil {
		if serror.Is(err, serror.NOT_FOUND) {
			return false, nil // if we error for this reason, continue
		}
		return false, common.StringError(err)
	}
	if string(denied) == "true" {
		return true, nil
	}
	return false, nil
}
