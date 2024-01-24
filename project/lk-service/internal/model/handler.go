package model

import (
	"github.com/golang-jwt/jwt/v5"
	"mime/multipart"
)

type Response struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetMyProfileHandlerReq struct {
	UID int64
}

type GetMyProfileHandlerRes struct {
	ID                 int64  `json:"id"`
	Email              string `json:"email"`
	Username           string `json:"username"`
	Name               string `json:"name"`
	Surname            string `json:"surname"`
	Patronymic         string `json:"patronymic"`
	Birthday           string `json:"birthday"`
	Phone              string `json:"phone"`
	ProfilePictureLink string `json:"profile_picture_link"`
}

type GetUserInfoHandlerReq struct {
	Username string
}

type GetUserInfoHandlerRes struct {
	ID                 int64  `json:"id"`
	Email              string `json:"email"`
	Username           string `json:"username"`
	Name               string `json:"name"`
	Surname            string `json:"surname"`
	Patronymic         string `json:"patronymic"`
	Birthday           string `json:"birthday"`
	Phone              string `json:"phone"`
	ProfilePictureLink string `json:"profile_picture_link"`
}

type ChangeUserInfoHandlerReq struct {
	ID         int64
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	Patronymic string `json:"patronymic,omitempty"`
	Username   string `json:"username,omitempty"`
	Birthday   string `json:"birthday,omitempty"`
	Phone      string `json:"phone,omitempty"`
}

type ChangeUserInfoHandlerRes struct {
	ID                 int64  `json:"id"`
	Email              string `json:"email"`
	Username           string `json:"username"`
	Name               string `json:"name"`
	Surname            string `json:"surname"`
	Patronymic         string `json:"patronymic"`
	Birthday           string `json:"birthday"`
	Phone              string `json:"phone"`
	ProfilePictureLink string `json:"profile_picture_link"`
	Token              string `json:"token"`
}

type ChangeUserEmailConfirmHandlerReq struct {
	UID   int64
	Email string `json:"email"`
}

type ChangeUserEmailCodeHandlerReq struct {
	Code int `json:"code"`
	UID  int64
}

type ChangeUserEmailCodeHandlerRes struct {
	ID                 int64  `json:"id"`
	Email              string `json:"email"`
	Username           string `json:"username"`
	Name               string `json:"name"`
	Surname            string `json:"surname"`
	Patronymic         string `json:"patronymic"`
	Birthday           string `json:"birthday"`
	Phone              string `json:"phone"`
	ProfilePictureLink string `json:"profile_picture_link"`
	Token              string `json:"token"`
}

type ChangeUserPasswordEmailConfirmHandlerReq struct {
	UID         int64
	PasswordNew string `json:"password_new"`
}

type ChangeUserPasswordCodeHandlerReq struct {
	UID  int64
	Code int `json:"code"`
}

type ChangeUserPasswordCodeHandlerRes struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Birthday   string `json:"birthday"`
	Phone      string `json:"phone"`
}

type ChangeUserProfilePictureLinkHandlerReq struct {
	UID  int64
	File *multipart.FileHeader
}

type ChangeUserProfilePictureLinkHandlerRes struct {
	ID                 int64  `json:"id"`
	Email              string `json:"email"`
	Username           string `json:"username"`
	Name               string `json:"name"`
	Surname            string `json:"surname"`
	Patronymic         string `json:"patronymic"`
	Birthday           string `json:"birthday"`
	Phone              string `json:"phone"`
	ProfilePictureLink string `json:"profile_picture_link"`
	Token              string `json:"token"`
}

type JWTCustomClaims struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Birthday   string `json:"birthday"`
	Phone      string `json:"phone"`
	IMGLink    string `json:"img_link"`
	jwt.RegisteredClaims
}
