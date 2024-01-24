package model

type EmailCodeItem struct {
	ID       int64  `json:"id" reindex:"id,hash,pk"`
	Lifetime int64  `json:"lifetime" reindex:"lifetime,ttl,,expire_after=120"`
	UID      int64  `json:"uid" reindex:"uid,hash"`
	Code     int64  `json:"code" reindex:"code,hash"`
	Email    string `json:"email" reindex:"email,hash"`
	Password string `json:"password" reindex:"password,hash"`
}
