package model

type ActiveItem struct {
	ID     int64 `json:"id" reindex:"id,hash,pk"`
	Online int64 `json:"online" reindex:"online,ttl,,expire_after=30"`
	UID    int64 `json:"uid" reindex:"uid,hash"`
}
