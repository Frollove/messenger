package reindexer_db

import (
	"encoding/json"
	"fmt"
	"github.com/restream/reindexer/v3"
	"lk-service/internal/model"
	"lk-service/pkg/custom_errors"
)

type UserApiImpl struct {
	db *reindexer.Reindexer
}

func NewUserApi(db *reindexer.Reindexer) *UserApiImpl {
	return &UserApiImpl{
		db: db,
	}
}

func (a *UserApiImpl) FindUserByUsernameDB(username string) (*model.UserItem, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	userJSON, found := a.db.Query("users").Where("username", reindexer.EQ, username).GetJson()
	if !found {
		return nil, fmt.Errorf("user doesn't exist: %w", custom_errors.ErrNotFound)
	}

	var res *model.UserItem

	if err := json.Unmarshal(userJSON, &res); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return res, nil
}

func (a *UserApiImpl) FindUserByIDDB(id int64) (*model.UserItem, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	userJSON, _ := a.db.Query("users").Where("id", reindexer.EQ, id).GetJson()

	var res *model.UserItem

	if err := json.Unmarshal(userJSON, &res); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return res, nil
}

func (a *UserApiImpl) UpdateUserInfoDB(req model.ChangeUserInfoHandlerReq, phone, birthday int64) error {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	query := a.db.Query("users").Where("id", reindexer.EQ, req.ID)

	if phone != 0 {
		query.Set("phone", phone).Update()
	}

	if req.Username != "" {
		query.Set("username", req.Username).Update()
	}

	if birthday != 0 {
		query.Set("birthday", birthday).Update()
	}

	if req.Surname != "" {
		query.Set("surname", req.Surname).Update()
	}

	if req.Name != "" {
		query.Set("name", req.Name).Update()
	}

	if req.Patronymic != "" {
		query.Set("patronymic", req.Patronymic).Update()
	}

	return nil
}

func (a *UserApiImpl) CheckUserByUsernameDB(username string) (bool, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return false, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	_, found := a.db.Query("users").Where("username", reindexer.EQ, username).GetJson()

	return found, nil
}

func (a *UserApiImpl) CheckUserByPhoneDB(phone int64) (bool, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return false, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	_, found := a.db.Query("users").Where("phone", reindexer.EQ, phone).GetJson()

	return found, nil
}

func (a *UserApiImpl) CheckUserByEmailDB(email string) (bool, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return false, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	_, found := a.db.Query("users").Where("email", reindexer.EQ, email).GetJson()

	return found, nil
}

func (a *UserApiImpl) UpdateUserEmailDB(id int64, email string) (*model.UserItem, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	a.db.Query("users").Where("id", reindexer.EQ, id).Set("email", email).Update()

	var res *model.UserItem

	userJSON, _ := a.db.Query("users").Where("id", reindexer.EQ, id).GetJson()
	if err := json.Unmarshal(userJSON, &res); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return res, nil
}

func (a *UserApiImpl) UpdateUserPasswordDB(id int64, password string) (*model.UserItem, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	a.db.Query("users").Where("id", reindexer.EQ, id).Set("password", password).Update()

	var res *model.UserItem

	userJSON, _ := a.db.Query("users").Where("id", reindexer.EQ, id).GetJson()
	if err := json.Unmarshal(userJSON, &res); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return res, nil
}
