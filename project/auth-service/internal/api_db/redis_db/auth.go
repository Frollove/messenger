package redis_db

import (
	"auth-service/internal/model"
	"auth-service/pkg/custom_errors"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type RedisApiImpl struct {
	db7 *redis.Client
	db8 *redis.Client
}

func NewRedisApi(db7, db8 *redis.Client) *RedisApiImpl {
	return &RedisApiImpl{
		db7: db7,
		db8: db8,
	}
}

func (a *RedisApiImpl) CreateRedisUserSingInDB7(req *model.CreateRedisUserSingInDB7Req) error {
	err := a.db7.HSet(context.Background(), req.Key, map[string]interface{}{
		"password": req.Password,
		"salt":     req.Salt,
		"id":       req.ID,
		"ip":       req.IP,
		"2fa":      req.DFA,
	}).Err()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("post user data to redis: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *RedisApiImpl) GetRedisUserSingInDB7(key string) (map[string]string, error) {
	values := a.db7.HGetAll(context.Background(), key).Val()
	if values == nil {
		return nil, fmt.Errorf("hgetall: %w", custom_errors.ErrNotFound)
	}

	a.db7.Del(context.Background(), key)

	return values, nil
}

func (a *RedisApiImpl) GetJtiArrayDB(key int64) ([]string, error) {
	phone := strconv.Itoa(int(key))

	result, err := a.db8.LRange(context.Background(), phone, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("LRange: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return result, nil
}

func (a *RedisApiImpl) RemoveLastJtiRecordDB(key int64) error {
	phone := strconv.Itoa(int(key))

	err := a.db8.LRem(context.Background(), phone, 1, a.db8.LIndex(context.Background(), phone, 0).Val()).Err()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("LRem: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *RedisApiImpl) CheckExistingJtiDB(key int64, jti int) int64 {
	phone := strconv.Itoa(int(key))

	result := a.db8.LInsert(context.Background(), phone, "Redis:AFTER", jti, jti).Val()

	return result
}

func (a *RedisApiImpl) AddJtiToArrayDB(key int64, jti int) error {
	phone := strconv.Itoa(int(key))

	err := a.db8.RPush(context.Background(), phone, jti).Err()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("RPush: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *RedisApiImpl) RemoveNewJtiDB(key int64, jti int) error {
	phone := strconv.Itoa(int(key))

	err := a.db8.LRem(context.Background(), phone, 1, jti).Err()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("LRem: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *RedisApiImpl) CreateRedisAuthSingInDB8(req *model.CreateRedisAuthSingInDB8Req) error {
	err := a.db8.HSet(context.Background(), req.Key, map[string]interface{}{
		"auth_code": req.AuthCode,
		"id":        req.ID,
		"ip":        req.IP,
		"2fa":       req.DFA,
	}).Err()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("HSet: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *RedisApiImpl) GetRedisAuthSingInDB8(key string) (map[string]string, error) {
	values := a.db8.HGetAll(context.Background(), key).Val()
	if values == nil {
		return nil, fmt.Errorf("hgetall: %w", custom_errors.ErrNotFound)
	}

	a.db8.Del(context.Background(), key)

	return values, nil
}

func (a *RedisApiImpl) CreateRedisAuthSingUpDB7(req model.CreateRedisAuthSingUpDB7Req) error {
	err := a.db7.HSet(context.Background(), req.Key, map[string]interface{}{
		"salt":      req.Salt,
		"ip":        req.IP,
		"auth_code": req.AuthCode,
	}).Err()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("hset: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *RedisApiImpl) CreateRedisUserSingUpDB8(req model.CreateRedisUserSingUpDB8Req) error {
	err := a.db8.HSet(context.Background(), req.Key, map[string]interface{}{
		"name":       req.Name,
		"surname":    req.Surname,
		"patronymic": req.Patronymic,
		"email":      req.Email,
		"birthday":   req.Birthday,
		"phone":      req.Phone,
		"username":   req.Username,
	}).Err()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("hset: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *RedisApiImpl) GetRedisAuthSingUpDB7(key string) (map[string]string, error) {
	values := a.db7.HGetAll(context.Background(), key).Val()
	if values == nil {
		return nil, fmt.Errorf("hgetall: %w", custom_errors.ErrNotFound)
	}

	a.db7.Del(context.Background(), key)

	return values, nil
}

func (a *RedisApiImpl) GetRedisUserSingUpDB8(key string) (map[string]string, error) {
	values := a.db8.HGetAll(context.Background(), key).Val()
	if values == nil {
		return nil, fmt.Errorf("hgetall: %w", custom_errors.ErrNotFound)
	}

	a.db8.Del(context.Background(), key)

	return values, nil
}

func (a *RedisApiImpl) CreateRedisConfirmEmailDB7(req model.CreateRedisConfirmEmailDB7Req) error {
	err := a.db7.HSet(context.Background(), req.Key, map[string]interface{}{
		"salt": req.Salt,
		"ip":   req.IP,
	}).Err()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("hset: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *RedisApiImpl) CreateRedisGeneralRecoveryDB7(req model.CreateRedisGeneralRecoveryDB7Req) error {
	err := a.db7.HSet(context.Background(), req.Key, map[string]interface{}{
		"auth_code":     req.AuthCode,
		"active_status": req.ActiveStatus,
		"ip":            req.IP,
		"id":            req.ID,
		"2fa":           req.DFA,
		"salt":          req.Salt,
	}).Err()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("hset: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *RedisApiImpl) GetRedisGeneralRecoveryDB7(key string) (map[string]string, error) {
	values := a.db7.HGetAll(context.Background(), key).Val()
	if values == nil {
		return nil, fmt.Errorf("hgetall: %w", custom_errors.ErrNotFound)
	}

	a.db7.Del(context.Background(), key)

	return values, nil
}

func (a *RedisApiImpl) CreateRedisConfirmRecoveryDB8(req model.CreateRedisConfirmRecoveryDB8Req) error {
	err := a.db8.HSet(context.Background(), req.Key, map[string]interface{}{
		"auth_code": req.AuthCode,
		"id":        req.ID,
	}).Err()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("hset: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *RedisApiImpl) GetRedisConfirmRecoveryDB8(key string) (map[string]string, error) {
	values := a.db8.HGetAll(context.Background(), key).Val()
	if values == nil {
		return nil, fmt.Errorf("hgetall: %w", custom_errors.ErrNotFound)
	}

	a.db8.Del(context.Background(), key)

	return values, nil
}
