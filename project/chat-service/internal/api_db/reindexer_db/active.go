package reindexer_db

import (
	"chat-service/internal/model"
	"chat-service/pkg/custom_errors"
	"fmt"
	"github.com/restream/reindexer/v3"
)

type ActiveApiImpl struct {
	db *reindexer.Reindexer
}

func NewActiveApi(db *reindexer.Reindexer) *ActiveApiImpl {
	return &ActiveApiImpl{
		db: db,
	}
}

func (a *ActiveApiImpl) GetOnlineDB(uid int64) (bool, error) {
	if err := a.db.OpenNamespace("active_user", reindexer.DefaultNamespaceOptions(), model.ActiveItem{}); err != nil {
		return false, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	_, found := a.db.Query("active_user").Where("uid", reindexer.EQ, uid).GetJson()

	return found, nil
}

func (a *ActiveApiImpl) SetOnlineDB(uid int64) error {
	if err := a.db.OpenNamespace("active_user", reindexer.DefaultNamespaceOptions(), model.ActiveItem{}); err != nil {
		return fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	_, err := a.db.Insert("active_user", model.ActiveItem{UID: uid}, "id=serial()", "online=now()")
	if err != nil {
		return fmt.Errorf(fmt.Errorf("insert: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}
