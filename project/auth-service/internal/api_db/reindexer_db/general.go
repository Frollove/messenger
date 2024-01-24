package reindexer_db

import (
	"auth-service/internal/model"
	"auth-service/pkg/custom_errors"
	"encoding/json"
	"fmt"
	"github.com/restream/reindexer/v3"
	"time"
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

func (a *GeneralApiImpl) CreateGeneralDB(req model.CreateGeneralDBReq) error {
	if err := a.db.OpenNamespace("general", reindexer.DefaultNamespaceOptions(), model.GeneralItem{}); err != nil {
		return fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	err := a.db.Upsert("general", &model.GeneralItem{
		UID:          req.UID,
		ActiveStatus: true,
		EmailConfirm: true,
		LastVisit:    time.Now().Unix(),
		RegDate:      time.Now().Unix(),
		ChangePass:   time.Now().Unix(),
		LastIP:       req.IP,
		DFA:          1,
		IMGLink:      "https://center.web-gen.ru:444/avatars/default.webp",
	}, "id=serial()")
	if err != nil {
		return fmt.Errorf(fmt.Errorf("upsert: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}

func (a *GeneralApiImpl) ChangePassTimeDB(id int64) error {
	if err := a.db.OpenNamespace("general", reindexer.DefaultNamespaceOptions(), model.GeneralItem{}); err != nil {
		return fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	err := a.db.Query("general").Where("uid", reindexer.EQ, id).Set("change_pass", time.Now().Unix()).Update().Error()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("query: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return nil
}
