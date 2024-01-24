package redis_db

import (
	"auth-service/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisApi interface {
	CreateRedisUserSingInDB7(req *model.CreateRedisUserSingInDB7Req) error
	GetRedisUserSingInDB7(key string) (map[string]string, error)
	GetJtiArrayDB(key int64) ([]string, error)
	RemoveLastJtiRecordDB(key int64) error
	CheckExistingJtiDB(key int64, jti int) int64
	AddJtiToArrayDB(key int64, jti int) error
	RemoveNewJtiDB(key int64, jti int) error
	CreateRedisAuthSingInDB8(req *model.CreateRedisAuthSingInDB8Req) error
	GetRedisAuthSingInDB8(key string) (map[string]string, error)
	CreateRedisAuthSingUpDB7(req model.CreateRedisAuthSingUpDB7Req) error
	CreateRedisUserSingUpDB8(req model.CreateRedisUserSingUpDB8Req) error
	GetRedisAuthSingUpDB7(key string) (map[string]string, error)
	GetRedisUserSingUpDB8(key string) (map[string]string, error)
	CreateRedisConfirmEmailDB7(req model.CreateRedisConfirmEmailDB7Req) error
	CreateRedisGeneralRecoveryDB7(req model.CreateRedisGeneralRecoveryDB7Req) error
	GetRedisGeneralRecoveryDB7(key string) (map[string]string, error)
	CreateRedisConfirmRecoveryDB8(req model.CreateRedisConfirmRecoveryDB8Req) error
	GetRedisConfirmRecoveryDB8(key string) (map[string]string, error)
}

type RedisApiDB struct {
	RedisApi
}

func NewRedisApiDB(db7, db8 *redis.Client) *RedisApiDB {
	return &RedisApiDB{
		RedisApi: NewRedisApi(db7, db8),
	}
}
