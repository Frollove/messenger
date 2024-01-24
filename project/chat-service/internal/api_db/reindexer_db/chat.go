package reindexer_db

import (
	"chat-service/internal/model"
	"chat-service/pkg/custom_errors"
	"encoding/json"
	"fmt"
	"github.com/restream/reindexer/v3"
	"time"
)

type ChatApiImpl struct {
	db *reindexer.Reindexer
}

func NewChatApi(db *reindexer.Reindexer) *ChatApiImpl {
	return &ChatApiImpl{
		db: db,
	}
}

func (a *ChatApiImpl) GetAllUserChatsDB(id int64) ([]*model.ChatItem, error) {
	if err := a.db.OpenNamespace("chat", reindexer.DefaultNamespaceOptions(), model.ChatItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	query := a.db.Query("chat").Where("host_uid", reindexer.EQ, id).Or().Where("uid", reindexer.EQ, id).Sort("message.time", true)

	iterator := query.Exec()
	if iterator.Error() != nil {
		return nil, fmt.Errorf(fmt.Errorf("interator error: %w", iterator.Error()).Error()+": %w", custom_errors.ErrInternal)
	}

	var res []*model.ChatItem

	for iterator.Next() {
		res = append(res, iterator.Object().(*model.ChatItem))
	}

	return res, nil
}

func (a *ChatApiImpl) CreateNewChatDB(hostUID, UID int64) (*model.ChatItem, error) {
	if err := a.db.OpenNamespace("chat", reindexer.DefaultNamespaceOptions(), model.ChatItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	var res *model.ChatItem

	chat, found := a.db.Query("chat").Where("host_uid", reindexer.EQ, hostUID).And().Where("uid", reindexer.EQ, UID).GetJson()

	if found {
		if err := json.Unmarshal(chat, &res); err != nil {
			return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}

		return res, nil
	}

	if _, err := a.db.Insert("chat", model.ChatItem{
		HostUID: hostUID,
		UID:     UID,
		Message: model.MessageType{
			LastMessage:    "",
			NumberOfUnread: 0,
			Time:           0,
		},
		Active:       true,
		CreateDate:   time.Now().Unix(),
		LastHostUser: true,
	}, "id=serial()"); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("upsert: %v", err).Error()+": %w", custom_errors.ErrInternal)
	}

	chat, _ = a.db.Query("chat").Where("host_uid", reindexer.EQ, hostUID).And().Where("uid", reindexer.EQ, UID).GetJson()
	if err := json.Unmarshal(chat, &res); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return res, nil
}

func (a *ChatApiImpl) CheckExistingChatDB(hostUID, UID int64) (bool, *model.ChatItem, error) {
	if err := a.db.OpenNamespace("chat", reindexer.DefaultNamespaceOptions(), model.ChatItem{}); err != nil {
		return false, nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	chatJson, found := a.db.Query("chat").Where("host_uid", reindexer.SET, []int64{hostUID, UID}).Where("uid", reindexer.SET, []int64{hostUID, UID}).GetJson()
	if !found {
		return false, nil, nil
	}

	var chat *model.ChatItem

	if err := json.Unmarshal(chatJson, &chat); err != nil {
		return false, nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return found, chat, nil
}

func (a *ChatApiImpl) GetChatByHostUIDAndUIDDB(hostUID, UID int64) (*model.ChatItem, error) {
	if err := a.db.OpenNamespace("chat", reindexer.DefaultNamespaceOptions(), model.ChatItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	chatJson, found := a.db.Query("chat").Where("host_uid", reindexer.SET, []int64{hostUID, UID}).Where("uid", reindexer.SET, []int64{hostUID, UID}).GetJson()
	if !found {
		return nil, nil
	}

	var chat *model.ChatItem

	if err := json.Unmarshal(chatJson, &chat); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return chat, nil
}

func (a *ChatApiImpl) AddLastMessageToChatDB(message string, id, time int64) error {
	if err := a.db.OpenNamespace("chat", reindexer.DefaultNamespaceOptions(), model.ChatItem{}); err != nil {
		return fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	query := a.db.Query("chat").Where("id", reindexer.EQ, id)

	query.Set("message.time", time).Update()
	query.Set("message.last_message", message).Update()

	return nil
}

func (a *ChatApiImpl) UserParticipantOfChatDB(uid int64) (bool, error) {
	if err := a.db.OpenNamespace("chat", reindexer.DefaultNamespaceOptions(), model.ChatItem{}); err != nil {
		return false, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	_, found := a.db.Query("chat").WhereInt64("host_uid", reindexer.EQ, uid).Limit(0).GetJson()
	if !found {
		_, found = a.db.Query("chat").WhereInt64("uid", reindexer.EQ, uid).Limit(0).GetJson()
	}

	return found, nil
}
