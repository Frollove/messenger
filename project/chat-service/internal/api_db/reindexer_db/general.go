package reindexer_db

import (
	"chat-service/internal/model"
	"chat-service/pkg/custom_errors"
	"encoding/json"
	"fmt"
	"github.com/restream/reindexer/v3"
)

type GeneralApiImpl struct {
	db *reindexer.Reindexer
}

func NewGeneralApi(db *reindexer.Reindexer) *GeneralApiImpl {
	return &GeneralApiImpl{
		db: db,
	}
}

func (a *GeneralApiImpl) GetRecordByUIDDB(id int64) (*model.GeneralItem, error) {
	if err := a.db.OpenNamespace("general", reindexer.DefaultNamespaceOptions(), model.GeneralItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	var general *model.GeneralItem

	result, found := a.db.Query("general").Where("uid", reindexer.EQ, id).GetJson()
	if !found {
		return nil, fmt.Errorf("query: %w", custom_errors.ErrNotFound)
	}

	if err := json.Unmarshal(result, &general); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return general, nil
}
