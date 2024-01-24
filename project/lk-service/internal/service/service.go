package service

import (
	"lk-service/internal/api_db/reindexer_db"
	"lk-service/internal/model"
)

type LKService interface {
	GetMyProfileService(req model.GetMyProfileHandlerReq) (*model.GetMyProfileHandlerRes, error)
	GetUserInfoService(req model.GetUserInfoHandlerReq) (*model.GetUserInfoHandlerRes, error)
	ChangeUserInfoService(req model.ChangeUserInfoHandlerReq) (*model.ChangeUserInfoHandlerRes, error)
	ChangeUserEmailConfirmService(req model.ChangeUserEmailConfirmHandlerReq) error
	ChangeUserEmailCodeService(req model.ChangeUserEmailCodeHandlerReq) (*model.ChangeUserEmailCodeHandlerRes, error)
	ChangeUserPasswordEmailConfirmService(req model.ChangeUserPasswordEmailConfirmHandlerReq) error
	ChangeUserPasswordCodeService(req model.ChangeUserPasswordCodeHandlerReq) (*model.ChangeUserPasswordCodeHandlerRes, error)
	ChangeUserProfilePictureLinkService(req model.ChangeUserProfilePictureLinkHandlerReq) (*model.ChangeUserProfilePictureLinkHandlerRes, error)
}

type Service struct {
	LKService
}

func NewService(userApiDB *reindexer_db.UserApiDB, generalApiDB *reindexer_db.GeneralApiDB, emailCodeApiDB *reindexer_db.EmailCodeApiDB) *Service {
	return &Service{
		LKService: NewLKService(userApiDB.UserApi, generalApiDB.GeneralApi, emailCodeApiDB.EmailCodeApi),
	}
}
