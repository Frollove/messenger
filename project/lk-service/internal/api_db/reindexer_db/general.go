package reindexer_db

import (
	"encoding/json"
	"fmt"
	"github.com/restream/reindexer/v3"
	"lk-service/internal/model"
	"lk-service/pkg/custom_errors"
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

func (a *GeneralApiImpl) ChangeProfilePictureLinkDB(path string, id int64) error {
	if err := a.db.OpenNamespace("general", reindexer.DefaultNamespaceOptions(), model.GeneralItem{}); err != nil {
		return fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	a.db.Query("general").Where("id", reindexer.EQ, id).Set("img_link", path).Update()

	return nil
}
