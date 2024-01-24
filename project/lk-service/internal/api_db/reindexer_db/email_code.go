package reindexer_db

import (
	"encoding/json"
	"fmt"
	"github.com/restream/reindexer/v3"
	"lk-service/internal/model"
	"lk-service/pkg/custom_errors"
)

type EmailCodeApiImpl struct {
	db *reindexer.Reindexer
}

func NewEmailCodeApi(db *reindexer.Reindexer) *EmailCodeApiImpl {
	return &EmailCodeApiImpl{
		db: db,
	}
}

func (a *EmailCodeApiImpl) CreateRecordEmailCodeDB(uid, code int64, email string) error {
	if err := a.db.OpenNamespace("email_code", reindexer.DefaultNamespaceOptions(), model.EmailCodeItem{}); err != nil {
		return fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	if err := a.db.Upsert("email_code", &model.EmailCodeItem{
		Code:  code,
		UID:   uid,
		Email: email,
	}, "id=serial()", "lifetime=now()"); err != nil {
		return fmt.Errorf(fmt.Errorf("upsert: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *EmailCodeApiImpl) GetRecordEmailCodeDB(uid int64) (*model.EmailCodeItem, error) {
	if err := a.db.OpenNamespace("email_code", reindexer.DefaultNamespaceOptions(), model.EmailCodeItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	var res *model.EmailCodeItem

	recordJSON, found := a.db.Query("email_code").Where("uid", reindexer.EQ, uid).GetJson()
	if !found {
		return nil, fmt.Errorf("code lifetime ended: %w", custom_errors.ErrWrongInputData)
	}

	if err := json.Unmarshal(recordJSON, &res); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return res, nil
}

func (a *EmailCodeApiImpl) CreateRecordEmailCodeWithPasswordDB(uid, code int64, email, password string) error {
	if err := a.db.OpenNamespace("email_code", reindexer.DefaultNamespaceOptions(), model.EmailCodeItem{}); err != nil {
		return fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	if err := a.db.Upsert("email_code", &model.EmailCodeItem{
		Code:     code,
		UID:      uid,
		Email:    email,
		Password: password,
	}, "id=serial()", "lifetime=now()"); err != nil {
		return fmt.Errorf(fmt.Errorf("upsert: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}
