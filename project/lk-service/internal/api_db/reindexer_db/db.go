package reindexer_db

import (
	"github.com/restream/reindexer/v3"
	"lk-service/internal/model"
)

type UserApi interface {
	FindUserByUsernameDB(username string) (*model.UserItem, error)
	FindUserByIDDB(id int64) (*model.UserItem, error)
	UpdateUserInfoDB(req model.ChangeUserInfoHandlerReq, phone, birthday int64) error
	CheckUserByUsernameDB(username string) (bool, error)
	CheckUserByPhoneDB(phone int64) (bool, error)
	CheckUserByEmailDB(email string) (bool, error)
	UpdateUserEmailDB(id int64, email string) (*model.UserItem, error)
	UpdateUserPasswordDB(id int64, password string) (*model.UserItem, error)
}

type GeneralApi interface {
	GetRecordByUIDDB(id int64) (*model.GeneralItem, error)
	ChangeProfilePictureLinkDB(path string, id int64) error
}

type EmailCodeApi interface {
	CreateRecordEmailCodeDB(uid, code int64, email string) error
	GetRecordEmailCodeDB(uid int64) (*model.EmailCodeItem, error)
	CreateRecordEmailCodeWithPasswordDB(uid, code int64, email, password string) error
}

type UserApiDB struct {
	UserApi
}

type GeneralApiDB struct {
	GeneralApi
}

type EmailCodeApiDB struct {
	EmailCodeApi
}

func NewUserApiDB(db *reindexer.Reindexer) *UserApiDB {
	return &UserApiDB{
		UserApi: NewUserApi(db),
	}
}

func NewGeneralApiDB(db *reindexer.Reindexer) *GeneralApiDB {
	return &GeneralApiDB{
		GeneralApi: NewGeneralApi(db),
	}
}

func NewEmailCodeApiDB(db *reindexer.Reindexer) *EmailCodeApiDB {
	return &EmailCodeApiDB{
		EmailCodeApi: NewEmailCodeApi(db),
	}
}
