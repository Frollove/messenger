package reindexer_db

import (
	"chat-service/internal/model"
	"chat-service/pkg/custom_errors"
	"encoding/json"
	"fmt"
	"github.com/restream/reindexer/v3"
	"time"
)

type MessageApiImpl struct {
	db *reindexer.Reindexer
}

func NewMessageApi(db *reindexer.Reindexer) *MessageApiImpl {
	return &MessageApiImpl{
		db: db,
	}
}

func (a *MessageApiImpl) GetAllChatMessagesDB(chatId int64) ([]*model.ChatMessageItem, error) {
	if err := a.db.OpenNamespace("chat_messages", reindexer.DefaultNamespaceOptions(), model.ChatMessageItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	query := a.db.Query("chat_messages").WhereInt64("chat_id", reindexer.EQ, chatId).Sort("time", true)

	iterator := query.Exec()
	if iterator.Error() != nil {
		return nil, fmt.Errorf(fmt.Errorf("iterator error: %w", iterator.Error()).Error()+": %w", custom_errors.ErrInternal)
	}

	var res []*model.ChatMessageItem

	for iterator.Next() {
		res = append(res, iterator.Object().(*model.ChatMessageItem))
	}

	return res, nil
}

func (a *MessageApiImpl) CreateNewMessageDB(chatID, uid int64, messageStr string, typeMessage string) (*model.ChatMessageItem, error) {
	if err := a.db.OpenNamespace("chat_messages", reindexer.DefaultNamespaceOptions(), model.ChatMessageItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	messageTime := time.Now().Unix()

	if _, err := a.db.Insert("chat_messages", model.ChatMessageItem{
		ChatID:  chatID,
		UID:     uid,
		Message: messageStr,
		Time:    messageTime,
		Type:    typeMessage,
		Read:    false,
	}, "id=serial()"); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("insert message: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	var message *model.ChatMessageItem

	messageJSON, _ := a.db.Query("chat_messages").Where("chat_id", reindexer.EQ, chatID).Where("uid", reindexer.EQ, uid).Where("time", reindexer.EQ, messageTime).GetJson()

	if err := json.Unmarshal(messageJSON, &message); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return message, nil
}

func (a *MessageApiImpl) GetLastMessageDB(chatID int64) (*model.ChatMessageItem, error) {
	if err := a.db.OpenNamespace("chat_messages", reindexer.DefaultNamespaceOptions(), model.ChatMessageItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	messageJson, _ := a.db.Query("chat_messages").Where("chat_id", reindexer.EQ, chatID).Sort("time", true).Limit(1).GetJson()

	var res *model.ChatMessageItem

	if err := json.Unmarshal(messageJson, &res); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("unmarshal: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return res, nil
}

func (a *MessageApiImpl) FullTextMessageSearchDB(search string, chatID int64) ([]*model.ChatMessageItem, error) {
	if err := a.db.OpenNamespace("chat_messages", reindexer.DefaultNamespaceOptions(), model.ChatMessageItem{}); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("open namespace: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	query := a.db.Query("chat_messages")
	if search != "" {
		query.Match("message", "~*"+search+"*~").Limit(5).Where("chat_id", reindexer.EQ, chatID)
	}

	iterator := query.Exec()
	if iterator.Error() != nil {
		return nil, fmt.Errorf(fmt.Errorf("exec: %w", iterator.Error()).Error()+": %w", custom_errors.ErrInternal)
	}

	var res []*model.ChatMessageItem

	for iterator.Next() {
		res = append(res, iterator.Object().(*model.ChatMessageItem))
	}

	return res, nil
}
