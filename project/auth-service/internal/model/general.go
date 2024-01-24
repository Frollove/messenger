package model

type GeneralItem struct {
	ID           int64  `json:"id" reindex:"id,hash,pk"`
	UID          int64  `json:"uid" reindex:"uid,hash"`
	ActiveStatus bool   `json:"active_status" reindex:"active_status,-"`
	EmailConfirm bool   `json:"email_confirm" reindex:"email_confirm,-"`
	LastVisit    int64  `json:"last_visit" reindex:"last_visit,tree"`
	RegDate      int64  `json:"reg_date" reindex:"reg_date,tree"`
	ChangePass   int64  `json:"change_pass" reindex:"change_pass,tree"`
	LastIP       string `json:"last_ip" reindex:"last_ip,hash"`
	DFA          int64  `json:"2fa" reindex:"2fa,hash"`
	IMGLink      string `json:"img_link" reindex:"img_link,hash"`
}

type CreateGeneralDBReq struct {
	UID int64  `json:"uid"`
	IP  string `json:"ip"`
}
