package service

import (
	"chat-service/internal/api_db/reindexer_db"
	"chat-service/internal/model"
)

type ChatService interface {
	FullTextSearchService(req model.FullTextSearchHandlerReq) (*model.FullTextSearchHandlerRes, error)
	GetChatByUsernameService(req model.GetChatByUsernameHandlerReq) (*model.GetChatByUsernameHandlerRes, error)
	GetAllUserChatsService(req model.GetAllUserChatsHandlerReq) ([]*model.GetAllUserChatsHandlerRes, error)
	GetChatService(req model.GetChatHandlerReq) (*model.GetChatHandlerRes, error)
	SendMessageService(req model.SendMessageHandlerReq) (*model.SendMessageHandlerRes, error)
	SendFileService(req model.SendFilesHandlerReq) (*model.SendFilesHandlerRes, error)
	SetOnlineService(req model.SetOnlineHandlerReq) error
	FullTextMessageSearchService(req model.FullTextMessageSearchHandlerReq) (*model.FullTextMessageSearchHandlerRes, error)
}

type Service struct {
	ChatService
}

func NewService(userApiDB *reindexer_db.UserApiDB, chatApiDB *reindexer_db.ChatApiDB, messageApiDB *reindexer_db.MessageApiDB, activeApiDB *reindexer_db.ActiveApiDB, generalApiDB *reindexer_db.GeneralApiDB) *Service {
	return &Service{
		ChatService: NewChatService(userApiDB.UserApi, chatApiDB.ChatApi, messageApiDB.MessageApi, activeApiDB.ActiveApi, generalApiDB.GeneralApi),
	}
}
