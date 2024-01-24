package service

import (
	"auth-service/internal/api_db/redis_db"
	"auth-service/internal/api_db/reindexer_db"
	"auth-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	GenerateHashService(req model.GenerateHashHandlerReq) (*model.GenerateHashHandlerRes, error)
	CheckDFAService(req model.CheckDFAHandlerReq) (*model.CheckDFAHandlerRes, error)
	ConfirmEmailSingInService(req model.ConfirmEmailSingInHandlerReq) (*model.ConfirmEmailSingInHandlerRes, error)
	SendEmailSingUpService(req model.SendEmailSingUpHandlerReq) (*model.SendEmailSingUpHandlerRes, error)
	ConfirmEmailSingUpService(req model.ConfirmEmailSingUpHandlerReq) (*model.ConfirmEmailSingUpHandlerRes, error)
	CreateUserService(req model.CreateUserHandlerReq) (*model.CreateUserHandlerRes, error)
	SendEmailRecoverService(req *model.SendEmailRecoverHandlerReq) (*model.SendEmailRecoverHandlerRes, error)
	ConfirmEmailRecoverService(req *model.ConfirmEmailRecoverHandlerReq) (*model.ConfirmEmailRecoverHandlerRes, error)
	ChangePasswordRecoverService(req *model.ChangePasswordRecoverHandlerReq) (*model.ChangePasswordRecoverHandlerRes, error)
	AuthMiddleWareService(claims jwt.MapClaims) (*model.AuthMiddleWareHandlerRes, error)
}

type Service struct {
	UserService
}

func NewService(userApiDB *reindexer_db.UserApiDB, generalApiDB *reindexer_db.GeneralApiDB, userRedisApiDB *redis_db.RedisApiDB) *Service {
	return &Service{
		UserService: NewAuthService(userApiDB.UserApi, generalApiDB.GeneralApi, userRedisApiDB.RedisApi),
	}
}
