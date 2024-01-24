package handler

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"lk-service/internal/service"
	"lk-service/pkg/custom_errors"
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
	case "/lk":
		h.GetUserInfoHandler(ctx)
	case "/lk/me":
		h.GetMyProfileHandler(ctx)
	case "/lk/change/info":
		h.ChangeUserInfoHandler(ctx)
	case "/lk/change/email":
		h.ChangeUserEmailConfirmHandler(ctx)
	case "/lk/change/email/confirm":
		h.ChangeUserEmailCodeHandler(ctx)
	case "/lk/change/password":
		h.ChangeUserPasswordEmailConfirmHandler(ctx)
	case "/lk/change/password/confirm":
		h.ChangeUserPasswordCodeHandler(ctx)
	case "/lk/change/avatar":
		h.ChangeUserProfilePictureLinkHandler(ctx)
	default:
		jsonErr, err := json.Marshal(custom_errors.ErrPageNotFound)
		if err != nil {
			logging.Info(err)
		}
		ctx.Write(jsonErr)
	}
}
