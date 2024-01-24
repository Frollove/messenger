package model

type ChatItem struct {
	ID           int64       `json:"id" reindex:"id,hash,pk"`
	HostUID      int64       `json:"host_uid" reindex:"host_uid,hash"`
	UID          int64       `json:"uid" reindex:"uid,hash"`
	Message      MessageType `json:"message"`
	Active       bool        `json:"active" reindex:"active,-"`
	CreateDate   int64       `json:"create_date" reindex:"create_date,tree"`
	LastHostUser bool        `json:"last_host_user" reindex:"last_host_user,-"`
}

type MessageType struct {
	NumberOfUnread int    `json:"number_of_unread" reindex:"number_message,-"`
	LastMessage    string `json:"last_message" reindex:"last_message,-"`
	Time           int64  `json:"time" reindex:"time,tree"`
}
