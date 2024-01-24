package reindexer_db

import (
	"auth-service/internal/model"
	"auth-service/pkg/custom_errors"
	"encoding/json"
	"fmt"
	"github.com/restream/reindexer/v3"
	"strconv"
)

type UserApiImpl struct {
	db *reindexer.Reindexer
}

func NewUserApi(db *reindexer.Reindexer) *UserApiImpl {
	return &UserApiImpl{
		db: db,
	}
}

func (a *UserApiImpl) GetUserByPhoneDB(phone int64) (*model.UserItem, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	result, found := a.db.Query("users").Where("phone", reindexer.EQ, phone).GetJson()
	if !found {
		return nil, fmt.Errorf("query: %w", custom_errors.ErrNotFound)
	}

	var user *model.UserItem

	if err := json.Unmarshal(result, &user); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return user, nil
}

func (a *UserApiImpl) CheckUserByPhoneDB(phone int64) (bool, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return false, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	_, found := a.db.Query("users").Where("phone", reindexer.EQ, strconv.Itoa(int(phone))).GetJson()

	return found, nil
}

func (a *UserApiImpl) CreateUserDB(req model.CreateUserDBReq) error {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	err := a.db.Upsert("users", &model.UserItem{
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   req.Password,
		Username:   req.Username,
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
		Birthday:   req.Birthday,
	}, "id=serial()")
	if err != nil {
		return fmt.Errorf(fmt.Errorf("upsert: %v", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *UserApiImpl) CheckUserByIDDB(id int64) (bool, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return false, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	_, found := a.db.Query("users").Where("id", reindexer.EQ, id).GetJson()

	return found, nil
}

func (a *UserApiImpl) ChangePassDB(id int64, password string) error {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	err := a.db.Query("users").Where("id", reindexer.EQ, id).Set("password", password).Update().Error()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("query: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *UserApiImpl) CheckUsernameDB(username string) (bool, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return false, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	_, found := a.db.Query("users").Where("username", reindexer.EQ, username).GetJson()

	return found, nil
}
