package repository

import "errors"

var ErrNotFound = errors.New("not found")

type Repositories struct {
	Platform Platform
}
