package handler

import (
	"auth-service/internal/model"
	"auth-service/pkg/custom_errors"
	"auth-service/pkg/tool"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func (h *Handler) SendEmailSingUpHandler(ctx *fasthttp.RequestCtx) (res *model.SendEmailSingUpHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.SendEmailSingUpHandlerReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	req.IP = ctx.RemoteIP().String()

	if err = req.Validate(); err != nil {
		return nil, fmt.Errorf("handler: validation request: %w", err)
	}

	res, err = h.services.SendEmailSingUpService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: send email sing up service: %w", err)
	}

	return
}

func (h *Handler) ConfirmEmailSingUpHandler(ctx *fasthttp.RequestCtx) (res *model.ConfirmEmailSingUpHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.ConfirmEmailSingUpHandlerReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	req.IP = ctx.RemoteIP().String()

	if err = req.Validate(); err != nil {
		return nil, fmt.Errorf("handler: validation request: %w", err)
	}

	res, err = h.services.ConfirmEmailSingUpService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: confirm sing up email service: %v", err)
	}

	return
}

func (h *Handler) CreateUserHandler(ctx *fasthttp.RequestCtx) (res *model.CreateUserHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.CreateUserHandlerReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	req.IP = ctx.RemoteIP().String()

	if err = req.Validate(); err != nil {
		return nil, fmt.Errorf("handler: validation request: %w", err)
	}

	res, err = h.services.CreateUserService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: create user service: %v", err)
	}

	return
}

func (h *Handler) GenerateHashSingInHandler(ctx *fasthttp.RequestCtx) (res *model.GenerateHashHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.GenerateHashHandlerReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	if err = req.Validate(); err != nil {
		return nil, fmt.Errorf("handler: validation request: %w", err)
	}

	req.IP = ctx.RemoteIP().String()

	res, err = h.services.GenerateHashService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: generate hash service: %w", err)
	}

	return
}

func (h *Handler) CheckDFASingInHandler(ctx *fasthttp.RequestCtx) (res *model.CheckDFAHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.CheckDFAHandlerReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	if err = req.Validate(); err != nil {
		return nil, fmt.Errorf("handler: validation request: %w", err)
	}

	req.IP = ctx.RemoteIP().String()

	res, err = h.services.CheckDFAService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: check dfa service: %w", err)
	}

	return
}

func (h *Handler) ConfirmEmailSingInHandler(ctx *fasthttp.RequestCtx) (res *model.ConfirmEmailSingInHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.ConfirmEmailSingInHandlerReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	if err = req.Validate(); err != nil {
		return nil, fmt.Errorf("handler: validation request: %v", err)
	}

	res, err = h.services.ConfirmEmailSingInService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: confirm email service: %v", err)
	}

	return
}

func (h *Handler) SendEmailRecoverHandler(ctx *fasthttp.RequestCtx) (res *model.SendEmailRecoverHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req *model.SendEmailRecoverHandlerReq

	if err = json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	req.IP = ctx.RemoteIP().String()

	res, err = h.services.SendEmailRecoverService(req)
	if err != nil {
		return nil, fmt.Errorf("send email recover service: %w", err)
	}

	return
}

func (h *Handler) ConfirmEmailRecoverHandler(ctx *fasthttp.RequestCtx) (res *model.ConfirmEmailRecoverHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req *model.ConfirmEmailRecoverHandlerReq

	if err = json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	req.IP = ctx.RemoteIP().String()

	res, err = h.services.ConfirmEmailRecoverService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: confirm email recover service: %w", err)
	}

	return
}

func (h *Handler) ChangePasswordRecoverHandler(ctx *fasthttp.RequestCtx) (res *model.ChangePasswordRecoverHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req *model.ChangePasswordRecoverHandlerReq

	if err = json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	res, err = h.services.ChangePasswordRecoverService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: change password recover service: %w", err)
	}

	return

}

func (h *Handler) AuthMiddleWareHandler(ctx *fasthttp.RequestCtx) (res *model.AuthMiddleWareHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var tokenStr model.AuthMiddlewareHandlerReq

	if err = json.Unmarshal(ctx.Request.Body(), &tokenStr); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.key")), nil
	})
	if err != nil || !token.Valid || claims["name"] == nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: token invalid").Error()+": %w", custom_errors.ErrInternal)
	}

	res, err = h.services.AuthMiddleWareService(claims)
	if err != nil {
		return nil, fmt.Errorf("handler: user exist service: %v", err)
	}

	return
}
