package handler

import (
	"chat-service/internal/service"
	"chat-service/pkg/custom_errors"
	"encoding/json"
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
	case "/chat/search":
		h.FullTextSearchHandler(ctx)
	case "/chat/search/get":
		h.GetChatByUsernameHandler(ctx)
	case "/chat/get-all":
		h.GetAllUserChatsHandler(ctx)
	case "/chat/get":
		h.GetChatHandler(ctx)
	case "/chat/send/message":
		h.SendMessageHandler(ctx)
	case "/chat/send/file":
		h.SendFileHandler(ctx)
	case "/online":
		h.GetOnlineHandler(ctx)
	case "/chat/message/search":
		h.FullTextMessageSearchHandler(ctx)
	default:
		jsonErr, err := json.Marshal(custom_errors.ErrPageNotFound)
		if err != nil {
			logging.Info(err)
		}
		ctx.Write(jsonErr)
	}
}
