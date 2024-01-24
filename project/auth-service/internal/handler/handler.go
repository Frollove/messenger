package handler

import (
	"auth-service/internal/service"
	"auth-service/pkg/custom_errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"resenje.org/logging"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	logging.Info(fmt.Sprintf("%s %s", string(ctx.Path()), string(ctx.Method())))
	switch string(ctx.Path()) {
	case "/sing-up/email/send":
		h.SendEmailSingUpHandler(ctx)
	case "/sing-up/email/confirm":
		h.ConfirmEmailSingUpHandler(ctx)
	case "/sing-up/token":
		h.CreateUserHandler(ctx)
	case "/sing-in/hash":
		h.GenerateHashSingInHandler(ctx)
	case "/sing-in/dfa":
		h.CheckDFASingInHandler(ctx)
	case "/sing-in/token":
		h.ConfirmEmailSingInHandler(ctx)
	case "/recover/email/send":
		h.SendEmailRecoverHandler(ctx)
	case "/recover/email/confirm":
		h.ConfirmEmailRecoverHandler(ctx)
	case "/recover/password":
		h.ChangePasswordRecoverHandler(ctx)
	case "/auth":
		h.AuthMiddleWareHandler(ctx)
	default:
		ctx.Error(custom_errors.ErrPageNotFound.Message, custom_errors.ErrPageNotFound.Code)
	}
}
