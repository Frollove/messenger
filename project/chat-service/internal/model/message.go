package model

type ChatMessageItem struct {
	ID      int64  `json:"id" reindex:"id,hash,pk"`
	ChatID  int64  `json:"chat_id" reindex:"chat_id,hash"`
	UID     int64  `json:"uid" reindex:"uid,hash"`
	Time    int64  `json:"time" reindex:"time,tree"`
	Message string `json:"message" reindex:"message,text"`
	Read    bool   `json:"read" reindex:"read,-"`
	Type    string `json:"type" reindex:"type,hash"`
}
