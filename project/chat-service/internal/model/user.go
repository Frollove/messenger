package model

type UserItem struct {
	ID         int64    `json:"id" reindex:"id,hash,pk"`
	Email      string   `json:"email" reindex:"email,hash"`
	Username   string   `json:"username" reindex:"username,hash"`
	Password   string   `json:"password" reindex:"password,hash"`
	Name       string   `json:"name" reindex:"name,hash"`
	Surname    string   `json:"surname" reindex:"surname,hash"`
	Patronymic string   `json:"patronymic" reindex:"patronymic,hash"`
	Birthday   int64    `json:"birthday" reindex:"birthday,tree"`
	Phone      int64    `json:"phone" reindex:"phone,hash"`
	_          struct{} `reindex:"username=search,text,composite"`
}
