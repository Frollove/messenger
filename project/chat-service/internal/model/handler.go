package model

import "mime/multipart"

type Response struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type FullTextSearchHandlerReq struct {
	UID    int64
	Search string
}

type FullTextSearchHandlerRes struct {
	Users []UserSearchItem `json:"users"`
}

type GetChatByUsernameHandlerReq struct {
	Username string `json:"username"`
	UID      int64
}

type GetChatByUsernameHandlerRes struct {
	Messages        []*ChatMessageItem `json:"messages"`
	RequestedUserID int64              `json:"requested_user_id"`
}

type UserSearchItem struct {
	Username           string           `json:"username"`
	Password           string           `json:"password"`
	ProfilePictureLink string           `json:"profile_picture_link"`
	LastMessage        *ChatMessageItem `json:"last_message"`
	Online             bool             `json:"online"`
}

type GetAllUserChatsHandlerReq struct {
	ID int64
}

type GetAllUserChatsHandlerRes struct {
	Username           string    `json:"username"`
	ProfilePictureLink string    `json:"profile_picture_link"`
	Chat               *ChatItem `json:"chats"`
	Online             bool      `json:"online"`
}

type GetChatHandlerReq struct {
	ChatID int64 `json:"chat_id"`
	UID    int64
}

type GetChatHandlerRes struct {
	RequestedUserID int64              `json:"requested_user_id"`
	Messages        []*ChatMessageItem `json:"messages"`
}

type SendMessageHandlerReq struct {
	Message          string `json:"message"`
	ReceiverUsername string `json:"receiver_username"`
	UID              int64
}

type SendMessageHandlerRes struct {
	RequestedUserID int64 `json:"requested_user_id"`
	MessageID       int64 `json:"message_id"`
}

type SendFilesHandlerReq struct {
	Files            []*multipart.FileHeader
	ReceiverUsername string
	UID              int64
}

type SendFilesHandlerRes struct {
	RequestedUserID int64   `json:"requested_user_id"`
	MessagesID      []int64 `json:"message_id"`
}

type SetOnlineHandlerReq struct {
	UID int64
}

type SetOnlineHandlerRes struct {
	Online int64 `json:"online"`
}

type FullTextMessageSearchHandlerReq struct {
	ChatID string
	Search string
}

type FullTextMessageSearchHandlerRes struct {
	Messages []*ChatMessageItem
}
