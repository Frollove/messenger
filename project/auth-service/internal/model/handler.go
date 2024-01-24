package model

import (
	"auth-service/pkg/custom_errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
)

type Response struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GenerateHashHandlerReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

type GenerateHashHandlerRes struct {
	Login string `json:"login"`
	Hash  string `json:"hash"`
	DFA   int64  `json:"DFA"`
}

type CheckDFAHandlerReq struct {
	Login string `json:"login"`
	Hash  string `json:"hash"`
	IP    string `json:"ip"`
}

type CheckDFAHandlerRes struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	DFA        int64  `json:"2fa"`
	Token      string `json:"token"`
}

type ConfirmEmailSingInHandlerReq struct {
	Code  int    `json:"code"`
	Login string `json:"login"`
	Hash  string `json:"hash"`
}

type ConfirmEmailSingInHandlerRes struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	DFA        int64  `json:"2fa"`
	Token      string `json:"token"`
}

type SendEmailSingUpHandlerReq struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Phone      string `json:"phone"`
	Birthday   string `json:"birthday"`
	IP         string `json:"ip"`
}

type SendEmailSingUpHandlerRes struct {
	Login string `json:"login"`
	Hash  string `json:"hash"`
}

type ConfirmEmailSingUpHandlerReq struct {
	Code  int    `json:"code"`
	Login string `json:"login"`
	Hash  string `json:"hash"`
	IP    string `json:"ip"`
}

type ConfirmEmailSingUpHandlerRes struct {
	Login string `json:"login"`
	Hash  string `json:"hash"`
}

type CreateUserHandlerReq struct {
	Password   string `json:"password"`
	Login      string `json:"login"`
	HashFirst  string `json:"hash_first"`
	HashSecond string `json:"hash_second"`
	IP         string `json:"ip"`
}

type CreateUserHandlerRes struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	DFA     int64  `json:"2fa"`
	Token   string `json:"token"`
}

type AuthMiddlewareHandlerReq struct {
	Token string `json:"token"`
}

type AuthMiddleWareHandlerRes struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Birthday   string `json:"birthday"`
}

type SendEmailRecoverHandlerReq struct {
	Login string `json:"login"`
	IP    string `json:"ip"`
}

type SendEmailRecoverHandlerRes struct {
	Login string `json:"login"`
	Hash  string `json:"hash"`
	DFA   int64  `json:"2fa"`
}

type ConfirmEmailRecoverHandlerReq struct {
	Code  int64  `json:"code"`
	Login string `json:"login"`
	Hash  string `json:"hash"`
	IP    string `json:"ip"`
}

type ConfirmEmailRecoverHandlerRes struct {
	Hash string `json:"hash"`
	Code int    `json:"code"`
}

type ChangePasswordRecoverHandlerReq struct {
	Code     int    `json:"code"`
	Login    string `json:"login"`
	Hash     string `json:"hash"`
	Password string `json:"password"`
}

type ChangePasswordRecoverHandlerRes struct {
	Message string `json:"message"`
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

func (req *GenerateHashHandlerReq) Validate() error {
	req.Password = strings.Trim(req.Password, " \n\t")

	if len(req.Login) != 16 || len(req.Password) == 0 {
		return fmt.Errorf("check length: %w", custom_errors.ErrWrongInputData)
	}

	return nil
}

func (req *CheckDFAHandlerReq) Validate() error {
	req.Hash = strings.Trim(req.Hash, " \n\t")

	if len(req.Login) != 16 || len(req.Hash) == 0 {
		return fmt.Errorf("check length: %w", custom_errors.ErrWrongInputData)
	}

	return nil
}

func (req *ConfirmEmailSingInHandlerReq) Validate() error {
	req.Hash = strings.Trim(req.Hash, " \n\t")

	if len(strconv.Itoa(req.Code)) != 4 || len(req.Login) != 16 || len(req.Hash) == 0 {
		return fmt.Errorf("check length: %w", custom_errors.ErrWrongInputData)
	}

	return nil
}

func (req *SendEmailSingUpHandlerReq) Validate() error {
	req.Name, req.Surname, req.Email, req.Patronymic, req.Username = strings.Trim(req.Name, " \n\t"), strings.Trim(req.Surname, " \n\t"), strings.Trim(req.Email, " \n\t"), strings.Trim(req.Patronymic, " \n\t"), strings.Trim(req.Username, " \n\t")

	if len(req.Name) == 0 || len(req.Surname) == 0 || len(req.Email) == 0 || len(req.Birthday) == 0 || len(req.Phone) != 16 || len(req.Username) == 0 {
		return fmt.Errorf("check length: %w", custom_errors.ErrWrongInputData)
	}

	return nil
}

func (req *ConfirmEmailSingUpHandlerReq) Validate() error {
	req.Hash = strings.Trim(req.Hash, " \n\t")

	if len(req.Login) != 16 || len(req.Hash) == 0 || len(strconv.Itoa(req.Code)) != 4 {
		return fmt.Errorf("check length: %w", custom_errors.ErrWrongInputData)
	}

	return nil
}

func (req *CreateUserHandlerReq) Validate() error {
	req.Password, req.HashFirst, req.HashSecond = strings.Trim(req.Password, " \n\t"), strings.Trim(req.HashFirst, " \n\t"), strings.Trim(req.HashSecond, " \n\t")

	if len(req.Login) != 16 || len(req.Password) == 0 || len(req.HashFirst) == 0 || len(req.HashSecond) == 0 {
		return fmt.Errorf("check length: %w", custom_errors.ErrWrongInputData)
	}

	return nil
}

func (req *AuthMiddlewareHandlerReq) Validate() error {
	req.Token = strings.Trim(req.Token, " \n\t")

	if len(req.Token) == 0 {
		return fmt.Errorf("check length: %w", custom_errors.ErrWrongInputData)
	}

	return nil
}

func (req *SendEmailRecoverHandlerReq) Validate() error {
	if len(req.Login) == 16 {
		return fmt.Errorf("check length: %w", custom_errors.ErrWrongInputData)
	}

	return nil
}

func (req *ConfirmEmailRecoverHandlerReq) Validate() error {
	req.Hash = strings.Trim(req.Hash, " \n\t")

	if len(req.Hash) == 0 || len(req.Login) != 16 || len(strconv.Itoa(int(req.Code))) != 4 {
		return fmt.Errorf("check length: %w", custom_errors.ErrWrongInputData)
	}

	return nil
}

func (req *ChangePasswordRecoverHandlerReq) Validate() error {
	req.Hash, req.Password = strings.Trim(req.Hash, " \n\t"), strings.Trim(req.Password, " \n\t")

	if len(req.Hash) == 0 || len(req.Password) == 0 || len(strconv.Itoa(req.Code)) != 8 || len(req.Login) != 16 {
		return fmt.Errorf("check length: %w", custom_errors.ErrWrongInputData)
	}

	return nil
}
