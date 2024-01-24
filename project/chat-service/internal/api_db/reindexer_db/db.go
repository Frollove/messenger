package reindexer_db

import (
	"chat-service/internal/model"
	"github.com/restream/reindexer/v3"
)

type ChatApi interface {
	GetAllUserChatsDB(id int64) ([]*model.ChatItem, error)
	CreateNewChatDB(hostUID, UID int64) (*model.ChatItem, error)
	CheckExistingChatDB(hostUID, UID int64) (bool, *model.ChatItem, error)
	GetChatByHostUIDAndUIDDB(hostUID, UID int64) (*model.ChatItem, error)
	AddLastMessageToChatDB(message string, id, time int64) error
	UserParticipantOfChatDB(uid int64) (bool, error)
}

type UserApi interface {
	FindUsersWithLoginDB(login string, id int64) ([]*model.UserItem, error)
	GetUserByUsernameDB(username string) (*model.UserItem, error)
	GetUserByIDDB(id int64) (*model.UserItem, error)
}

type MessageApi interface {
	GetAllChatMessagesDB(chatId int64) ([]*model.ChatMessageItem, error)
	CreateNewMessageDB(chatID, uid int64, messageStr string, typeMessage string) (*model.ChatMessageItem, error)
	GetLastMessageDB(chatID int64) (*model.ChatMessageItem, error)
	FullTextMessageSearchDB(search string, chatID int64) ([]*model.ChatMessageItem, error)
}

type ActiveApi interface {
	GetOnlineDB(uid int64) (bool, error)
	SetOnlineDB(uid int64) error
}

type GeneralApi interface {
	GetRecordByUIDDB(id int64) (*model.GeneralItem, error)
}

type ChatApiDB struct {
	ChatApi
}

type UserApiDB struct {
	UserApi
}

type MessageApiDB struct {
	MessageApi
}

type ActiveApiDB struct {
	ActiveApi
}

type GeneralApiDB struct {
	GeneralApi
}

func NewChatApiDB(db *reindexer.Reindexer) *ChatApiDB {
	return &ChatApiDB{
		ChatApi: NewChatApi(db),
	}
}

func NewUserApiDB(db *reindexer.Reindexer) *UserApiDB {
	return &UserApiDB{
		UserApi: NewUserApi(db),
	}
}

func NewMessageApiDB(db *reindexer.Reindexer) *MessageApiDB {
	return &MessageApiDB{
		MessageApi: NewMessageApi(db),
	}
}

func NewActiveApiDB(db *reindexer.Reindexer) *ActiveApiDB {
	return &ActiveApiDB{
		ActiveApi: NewActiveApi(db),
	}
}

func NewGeneralApiDB(db *reindexer.Reindexer) *GeneralApiDB {
	return &GeneralApiDB{
		GeneralApi: NewGeneralApi(db),
	}
}
