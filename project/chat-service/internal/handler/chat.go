package handler

import (
	"chat-service/internal/model"
	"chat-service/pkg/custom_errors"
	"chat-service/pkg/tool"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"strconv"
	"strings"
)

func (h *Handler) FullTextSearchHandler(ctx *fasthttp.RequestCtx) (res *model.FullTextSearchHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)

		jsonRes, _ := json.Marshal(resFin)

		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsGet() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	if string(ctx.QueryArgs().Peek("search")) == "" {
		return nil, nil
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, fmt.Errorf("handler: no authorization header: %w", custom_errors.ErrUnauthorized)
	}

	token := strings.Split(string(tokenGet), " ")[1]

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("bytrip"), nil
	})
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	req := model.FullTextSearchHandlerReq{Search: string(ctx.QueryArgs().Peek("search")), UID: int64(id)}

	res, err = h.services.FullTextSearchService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: full-text search service: %w", err)
	}

	return
}

func (h *Handler) GetChatByUsernameHandler(ctx *fasthttp.RequestCtx) (res *model.GetChatByUsernameHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)

		jsonRes, _ := json.Marshal(resFin)

		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.GetChatByUsernameHandlerReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, fmt.Errorf("handler: no authorization header: %w", custom_errors.ErrUnauthorized)
	}

	token := strings.Split(string(tokenGet), " ")[1]

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("bytrip"), nil
	})
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	req.UID = int64(id)

	res, err = h.services.GetChatByUsernameService(req)
	if err != nil {
		return nil, fmt.Errorf("get chat by username service: %w", err)
	}

	return
}

func (h *Handler) GetAllUserChatsHandler(ctx *fasthttp.RequestCtx) (res []*model.GetAllUserChatsHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)

		jsonRes, _ := json.Marshal(resFin)

		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsGet() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, fmt.Errorf("handler: no authorization header: %w", custom_errors.ErrUnauthorized)
	}

	token := strings.Split(string(tokenGet), " ")[1]

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("bytrip"), nil
	})
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	res, err = h.services.GetAllUserChatsService(model.GetAllUserChatsHandlerReq{ID: int64(id)})
	if err != nil {
		return nil, fmt.Errorf("handler: get all user chats service: %w", err)
	}

	return
}

func (h *Handler) GetChatHandler(ctx *fasthttp.RequestCtx) (res *model.GetChatHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)

		jsonRes, _ := json.Marshal(resFin)

		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsGet() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, fmt.Errorf("handler: no authorization header: %w", custom_errors.ErrUnauthorized)
	}

	token := strings.Split(string(tokenGet), " ")[1]

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("bytrip"), nil
	})
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	paramValue := string(ctx.QueryArgs().Peek("chat-id"))

	chatID, err := strconv.Atoi(paramValue)
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	res, err = h.services.GetChatService(model.GetChatHandlerReq{ChatID: int64(chatID), UID: int64(id)})
	if err != nil {
		return nil, fmt.Errorf("handler: get chat service: %w", err)
	}

	return
}

func (h *Handler) SendMessageHandler(ctx *fasthttp.RequestCtx) (res *model.SendMessageHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)

		jsonRes, _ := json.Marshal(resFin)

		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.SendMessageHandlerReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, fmt.Errorf("handler: no authorization header: %w", custom_errors.ErrUnauthorized)
	}

	token := strings.Split(string(tokenGet), " ")[1]

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("bytrip"), nil
	})
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	req.UID = int64(id)

	res, err = h.services.SendMessageService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: send message service: %w", err)
	}

	return
}

func (h *Handler) SendFileHandler(ctx *fasthttp.RequestCtx) (res *model.SendFilesHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)

		jsonRes, _ := json.Marshal(resFin)

		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: multipart form: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	var files []*multipart.FileHeader

	for k := range form.File {
		file, err := ctx.FormFile(k)
		if err != nil {
			return nil, fmt.Errorf(fmt.Errorf("handler: form file: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
		}
		files = append(files, file)
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, fmt.Errorf("handler: no authorization header: %w", custom_errors.ErrUnauthorized)
	}

	token := strings.Split(string(tokenGet), " ")[1]

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("bytrip"), nil
	})
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	res, err = h.services.SendFileService(model.SendFilesHandlerReq{
		Files:            files,
		ReceiverUsername: form.Value["receiver_username"][0],
		UID:              int64(id),
	})
	if err != nil {
		return nil, fmt.Errorf("handler: send file service: %w", err)
	}

	return
}

func (h *Handler) GetOnlineHandler(ctx *fasthttp.RequestCtx) (err error) {
	defer func() {
		resFin := tool.Response(err, nil)

		jsonRes, _ := json.Marshal(resFin)

		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsGet() {
		return fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return fmt.Errorf("handler: no authorization header: %w", custom_errors.ErrUnauthorized)
	}

	token := strings.Split(string(tokenGet), " ")[1]

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("bytrip"), nil
	})
	if err != nil {
		return fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrUnauthorized)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	err = h.services.SetOnlineService(model.SetOnlineHandlerReq{UID: int64(id)})
	if err != nil {
		return fmt.Errorf("handler: set online service: %w", err)
	}

	return
}

func (h *Handler) FullTextMessageSearchHandler(ctx *fasthttp.RequestCtx) (res *model.FullTextMessageSearchHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)

		jsonRes, _ := json.Marshal(resFin)

		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsGet() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	req := model.FullTextMessageSearchHandlerReq{Search: string(ctx.QueryArgs().Peek("search")), ChatID: string(ctx.QueryArgs().Peek("chat-id"))}

	res, err = h.services.FullTextMessageSearchService(req)
	if err != nil {
		return nil, fmt.Errorf("full test message search service: %w", err)
	}

	return res, nil
}
