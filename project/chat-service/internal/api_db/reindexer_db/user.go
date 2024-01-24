package reindexer_db

import (
	"chat-service/internal/model"
	"chat-service/pkg/custom_errors"
	"encoding/json"
	"fmt"
	"github.com/restream/reindexer/v3"
)

type UserApiImpl struct {
	db *reindexer.Reindexer
}

func NewUserApi(db *reindexer.Reindexer) *UserApiImpl {
	return &UserApiImpl{
		db: db,
	}
}

func (a *UserApiImpl) FindUsersWithLoginDB(login string, id int64) ([]*model.UserItem, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	query := a.db.Query("users")
	if login != "" {
		query.Match("search", "~*"+login+"*~").Limit(5)
	}

	var res []*model.UserItem

	iterator := query.Exec()
	if iterator.Error() != nil {
		return nil, fmt.Errorf(fmt.Errorf("iterator error: %w", iterator.Error()).Error()+": %w", custom_errors.ErrInternal)
	}

	for iterator.Next() {
		user := iterator.Object().(*model.UserItem)
		if user.ID == id {
			continue
		}
		res = append(res, user)
	}

	return res, nil
}

func (a *UserApiImpl) GetUserByUsernameDB(username string) (*model.UserItem, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	result, found := a.db.Query("users").Where("username", reindexer.EQ, username).GetJson()
	if !found {
		return nil, fmt.Errorf("query: %w", custom_errors.ErrNotFound)
	}

	var user *model.UserItem

	if err := json.Unmarshal(result, &user); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return user, nil
}

func (a *UserApiImpl) GetUserByIDDB(id int64) (*model.UserItem, error) {
	if err := a.db.OpenNamespace("users", reindexer.DefaultNamespaceOptions(), model.UserItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	var user *model.UserItem

	userJSON, _ := a.db.Query("users").Where("id", reindexer.EQ, id).GetJson()

	if err := json.Unmarshal(userJSON, &user); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return user, nil
}
