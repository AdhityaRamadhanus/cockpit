package cockpit

import (
	"time"

	"github.com/pkg/errors"
)

var (
	//ErrKeyNotFound represent key is not found when searching in KeyValueService
	ErrKeyNotFound = errors.New("Key not found")
)

//KeyValueService provide an interface to get key
type KeyValueService interface {
	Set(key string, value []byte) error
	SetEx(key string, value []byte, expirationInSeconds time.Duration) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}
