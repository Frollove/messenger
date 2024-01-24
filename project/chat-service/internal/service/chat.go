package service

import (
	"chat-service/internal/api_db/reindexer_db"
	"chat-service/internal/model"
	"chat-service/pkg/custom_errors"
	"chat-service/pkg/files"
	"fmt"
	"strconv"
	"strings"
)

type ChatServiceImpl struct {
	userApiDB    reindexer_db.UserApi
	chatApiDB    reindexer_db.ChatApi
	messageApiDB reindexer_db.MessageApi
	activeApiDB  reindexer_db.ActiveApi
	generalApiDB reindexer_db.GeneralApi
}

func NewChatService(userApiDB reindexer_db.UserApi, chatApiDB reindexer_db.ChatApi, messageApiDB reindexer_db.MessageApi, activeApiDB reindexer_db.ActiveApi, generalApiDB reindexer_db.GeneralApi) *ChatServiceImpl {
	return &ChatServiceImpl{
		userApiDB:    userApiDB,
		chatApiDB:    chatApiDB,
		messageApiDB: messageApiDB,
		activeApiDB:  activeApiDB,
		generalApiDB: generalApiDB,
	}
}

func (s *ChatServiceImpl) FullTextSearchService(req model.FullTextSearchHandlerReq) (*model.FullTextSearchHandlerRes, error) {
	users, err := s.userApiDB.FindUsersWithLoginDB(req.Search, req.UID)
	if err != nil {
		return nil, fmt.Errorf("find users with login DB: %w", err)
	}

	var res model.FullTextSearchHandlerRes

	for _, user := range users {
		record, err := s.generalApiDB.GetRecordByUIDDB(user.ID)
		if err != nil {
			return nil, fmt.Errorf("get record by UID DB: %w", err)
		}

		online, err := s.activeApiDB.GetOnlineDB(user.ID)
		if err != nil {
			return nil, fmt.Errorf("get online DB : %w", err)
		}

		found, chat, err := s.chatApiDB.CheckExistingChatDB(req.UID, user.ID)

		if !found {
			res.Users = append(res.Users, model.UserSearchItem{
				Username:           user.Username,
				Password:           user.Password,
				ProfilePictureLink: record.IMGLink,
				LastMessage:        nil,
			})

		} else {
			message, err := s.messageApiDB.GetLastMessageDB(chat.ID)
			if err != nil {
				return nil, fmt.Errorf("get last message DB: %w", err)
			}

			res.Users = append(res.Users, model.UserSearchItem{
				Username:           user.Username,
				Password:           user.Password,
				ProfilePictureLink: record.IMGLink,
				LastMessage:        message,
				Online:             online,
			})
		}
	}

	return &res, nil
}

func (s *ChatServiceImpl) GetChatByUsernameService(req model.GetChatByUsernameHandlerReq) (*model.GetChatByUsernameHandlerRes, error) {
	user, err := s.userApiDB.GetUserByUsernameDB(req.Username)
	if err != nil {
		return nil, fmt.Errorf("get user by username DB: %w", err)
	}

	if user.ID == req.UID {
		return nil, fmt.Errorf("check id's: can't create chat with yourself")
	}

	chat, err := s.chatApiDB.GetChatByHostUIDAndUIDDB(req.UID, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get chat by host UID and UID DB: %w", err)
	}

	chat, err = s.chatApiDB.GetChatByHostUIDAndUIDDB(user.ID, req.UID)
	if err != nil {
		return nil, fmt.Errorf("get chat by host UID and UID DB: %w", err)
	}

	if chat == nil {
		chat, err = s.chatApiDB.CreateNewChatDB(user.ID, req.UID)
		if err != nil {
			return nil, fmt.Errorf("create new chat DB: %w", err)
		}
	}

	messages, err := s.messageApiDB.GetAllChatMessagesDB(chat.ID)
	if err != nil {
		return nil, fmt.Errorf("get all chat messages DB: %w", err)
	}

	return &model.GetChatByUsernameHandlerRes{Messages: messages, RequestedUserID: req.UID}, nil
}

func (s *ChatServiceImpl) GetAllUserChatsService(req model.GetAllUserChatsHandlerReq) ([]*model.GetAllUserChatsHandlerRes, error) {
	chats, err := s.chatApiDB.GetAllUserChatsDB(req.ID)
	if err != nil {
		return nil, fmt.Errorf("get all user chats DB: %w", err)
	}

	var res []*model.GetAllUserChatsHandlerRes

	var user *model.UserItem

	for _, elem := range chats {
		var id int64
		if req.ID == elem.UID {
			id = elem.HostUID
		} else {
			id = elem.UID
		}

		user, err = s.userApiDB.GetUserByIDDB(id)
		if err != nil {
			return nil, fmt.Errorf("get user by id DB: %w", err)
		}

		record, err := s.generalApiDB.GetRecordByUIDDB(id)
		if err != nil {
			return nil, fmt.Errorf("get record by uid DB: %w", err)
		}

		online, err := s.activeApiDB.GetOnlineDB(id)
		if err != nil {
			return nil, fmt.Errorf("get online DB : %w", err)
		}

		res = append(res, &model.GetAllUserChatsHandlerRes{Username: user.Username, Chat: elem, Online: online, ProfilePictureLink: record.IMGLink})
	}

	return res, nil
}

func (s *ChatServiceImpl) GetChatService(req model.GetChatHandlerReq) (*model.GetChatHandlerRes, error) {
	found, err := s.chatApiDB.UserParticipantOfChatDB(req.UID)
	if err != nil {
		return nil, fmt.Errorf("user participant of chat DB: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("user isn't participant of this chat: %w", custom_errors.ErrNotParticipant)
	}

	messages, err := s.messageApiDB.GetAllChatMessagesDB(req.ChatID)
	if err != nil {
		return nil, fmt.Errorf("get all chat messages: %w", err)
	}

	return &model.GetChatHandlerRes{Messages: messages, RequestedUserID: req.UID}, nil
}

func (s *ChatServiceImpl) SendMessageService(req model.SendMessageHandlerReq) (*model.SendMessageHandlerRes, error) {
	receiverUser, err := s.userApiDB.GetUserByUsernameDB(req.ReceiverUsername)
	if err != nil {
		return nil, fmt.Errorf("get user by username DB: %w", err)
	}

	if receiverUser.ID == req.UID {
		return nil, fmt.Errorf("check id's: can't send message to yourself")
	}

	found, chat, err := s.chatApiDB.CheckExistingChatDB(req.UID, receiverUser.ID)

	if !found {
		chat, err = s.chatApiDB.CreateNewChatDB(req.UID, receiverUser.ID)
		if err != nil {
			return nil, fmt.Errorf("create new chat DB: %w", err)
		}
	}

	message, err := s.messageApiDB.CreateNewMessageDB(chat.ID, req.UID, req.Message, "string")
	if err != nil {
		return nil, fmt.Errorf("create new message DB: %w", err)
	}

	if err = s.chatApiDB.AddLastMessageToChatDB(message.Message, chat.ID, message.Time); err != nil {
		return nil, fmt.Errorf("add last message to chat DB: %w", err)
	}

	return &model.SendMessageHandlerRes{MessageID: message.ID, RequestedUserID: req.UID}, nil
}

func (s *ChatServiceImpl) SendFileService(req model.SendFilesHandlerReq) (*model.SendFilesHandlerRes, error) {
	receiverUser, err := s.userApiDB.GetUserByUsernameDB(req.ReceiverUsername)
	if err != nil {
		return nil, fmt.Errorf("get user by username DB: %w", err)
	}

	if receiverUser.ID == req.UID {
		return nil, fmt.Errorf("check id's: can't send message to yourself")
	}

	found, chat, err := s.chatApiDB.CheckExistingChatDB(req.UID, receiverUser.ID)

	var messages []int64
	if !found {
		chat, err = s.chatApiDB.CreateNewChatDB(req.UID, receiverUser.ID)
		if err != nil {
			return nil, fmt.Errorf("create new chat DB: %w", err)
		}
	}

	paths, err := files.DownloadFiles(chat.ID, req.Files)
	if err != nil {
		return nil, fmt.Errorf("download files: %w", err)
	}

	for key, elem := range paths {
		message, err := s.messageApiDB.CreateNewMessageDB(chat.ID, req.UID, strings.ReplaceAll(key, "src/chat/", "https://img.web-gen.ru/chat/"), elem)
		if err != nil {
			return nil, fmt.Errorf("create new message DB: %w", err)
		}

		if err = s.chatApiDB.AddLastMessageToChatDB(message.Message, chat.ID, message.Time); err != nil {
			return nil, fmt.Errorf("add last message to chat DB: %w", err)
		}

		messages = append(messages, message.ID)
	}

	return &model.SendFilesHandlerRes{MessagesID: messages, RequestedUserID: req.UID}, nil
}

func (s *ChatServiceImpl) SetOnlineService(req model.SetOnlineHandlerReq) error {
	err := s.activeApiDB.SetOnlineDB(req.UID)
	if err != nil {
		return fmt.Errorf("set online DB: %w", err)
	}

	return nil
}

func (s *ChatServiceImpl) FullTextMessageSearchService(req model.FullTextMessageSearchHandlerReq) (*model.FullTextMessageSearchHandlerRes, error) {
	chatIDint, err := strconv.Atoi(req.ChatID)
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("chat id atoi: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	messages, err := s.messageApiDB.FullTextMessageSearchDB(req.Search, int64(chatIDint))
	if err != nil {
		return nil, fmt.Errorf("full text message search DB: %w", err)
	}

	return &model.FullTextMessageSearchHandlerRes{Messages: messages}, nil
}
