package reindexer_db

import (
	"auth-service/internal/model"
	"github.com/restream/reindexer/v3"
)

type UserApi interface {
	GetUserByPhoneDB(phone int64) (*model.UserItem, error)
	CheckUserByPhoneDB(phone int64) (bool, error)
	CreateUserDB(req model.CreateUserDBReq) error
	CheckUserByIDDB(id int64) (bool, error)
	ChangePassDB(id int64, password string) error
	CheckUsernameDB(username string) (bool, error)
}

type GeneralApi interface {
	GetRecordByUIDDB(id int64) (*model.GeneralItem, error)
	CreateGeneralDB(req model.CreateGeneralDBReq) error
	ChangePassTimeDB(uid int64) error
}

type GeneralApiDB struct {
	GeneralApi
}

type UserApiDB struct {
	UserApi
}

func NewGeneralApiDB(db *reindexer.Reindexer) *GeneralApiDB {
	return &GeneralApiDB{
		GeneralApi: NewGeneralApi(db),
	}
}

func NewUserApiDB(db *reindexer.Reindexer) *UserApiDB {
	return &UserApiDB{
		UserApi: NewUserApi(db),
	}
}
