package handler

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"lk-service/internal/model"
	"lk-service/pkg/custom_errors"
	"lk-service/pkg/jwtRegister"
	"lk-service/pkg/tool"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) GetMyProfileHandler(ctx *fasthttp.RequestCtx) (res *model.GetMyProfileHandlerRes, err error) {
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

	claims := model.JWTCustomClaims{}

	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("bytrip"), nil
	})
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	email, username, name, surname, patronymic, birthday, phone, imgLink := claims.Email,
		claims.Username, claims.Name, claims.Surname, claims.Patronymic, claims.Birthday, claims.Phone, claims.IMGLink

	res = &model.GetMyProfileHandlerRes{
		ID:                 int64(id),
		Email:              email,
		Username:           username,
		Name:               name,
		Surname:            surname,
		Patronymic:         patronymic,
		Birthday:           birthday,
		Phone:              phone,
		ProfilePictureLink: imgLink,
	}
	return
}

func (h *Handler) GetUserInfoHandler(ctx *fasthttp.RequestCtx) (res *model.GetUserInfoHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsGet() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	res, err = h.services.GetUserInfoService(model.GetUserInfoHandlerReq{Username: string(ctx.QueryArgs().Peek("username"))})
	if err != nil {
		return nil, fmt.Errorf("get user info service: %w", err)
	}

	return
}

func (h *Handler) ChangeUserInfoHandler(ctx *fasthttp.RequestCtx) (res *model.ChangeUserInfoHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.ChangeUserInfoHandlerReq

	if err = json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, fmt.Errorf("handler: no authorization header: %w", custom_errors.ErrUnauthorized)
	}

	token := strings.Split(string(tokenGet), " ")[1]

	claims := model.JWTCustomClaims{}

	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("bytrip"), nil
	})
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	//if evolutionOfFront == true {delete a lot of fucking if}
	if req.Name == claims.Name {
		req.Name = ""
	}
	if req.Surname == claims.Surname {
		req.Surname = ""
	}
	if req.Patronymic == claims.Patronymic {
		req.Patronymic = ""
	}
	if req.Username == claims.Username {
		req.Username = ""
	}
	if req.Phone == claims.Phone {
		req.Phone = ""
	}
	if req.Birthday == claims.Birthday {
		req.Birthday = ""
	}
	req.ID = int64(id)

	res, err = h.services.ChangeUserInfoService(req)
	if err != nil {
		return nil, fmt.Errorf("change user info service: %w", err)
	}

	jtiNew := rand.Intn(99999999-10000000+1) + 10000000
	claimsNew := model.JWTCustomClaims{
		Email:      res.Email,
		Username:   res.Username,
		Name:       res.Name,
		Surname:    res.Surname,
		Patronymic: res.Patronymic,
		Birthday:   res.Birthday,
		Phone:      res.Phone,
		IMGLink:    res.ProfilePictureLink,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    fmt.Sprintf("https://%s", viper.GetString("jwt.domain")),
			Audience:  jwt.ClaimStrings{strconv.Itoa(id)},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 120)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        strconv.Itoa(jtiNew),
		},
	}

	tokenNew, err := jwtRegister.GenerateToken(claimsNew)
	if err != nil {
		return nil, fmt.Errorf("generate token : %w", err)
	}

	res.Token = tokenNew

	return
}

func (h *Handler) ChangeUserEmailConfirmHandler(ctx *fasthttp.RequestCtx) (err error) {
	defer func() {
		resFin := tool.Response(err, nil)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.ChangeUserEmailConfirmHandlerReq

	if err = json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		return fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
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
		return fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	req.UID = int64(id)

	err = h.services.ChangeUserEmailConfirmService(req)
	if err != nil {
		return fmt.Errorf("change user email confirm service: %w", err)
	}

	return
}

func (h *Handler) ChangeUserEmailCodeHandler(ctx *fasthttp.RequestCtx) (res *model.ChangeUserEmailCodeHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.ChangeUserEmailCodeHandlerReq

	if err = json.Unmarshal(ctx.Request.Body(), &req); err != nil {
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
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	req.UID = int64(id)

	res, err = h.services.ChangeUserEmailCodeService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: change user email code service: %w", err)
	}

	jtiNew := rand.Intn(99999999-10000000+1) + 10000000
	claimsNew := model.JWTCustomClaims{
		Email:      res.Email,
		Username:   res.Username,
		Name:       res.Name,
		Surname:    res.Surname,
		Patronymic: res.Patronymic,
		Birthday:   res.Birthday,
		Phone:      res.Phone,
		IMGLink:    res.ProfilePictureLink,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    fmt.Sprintf("https://%s", viper.GetString("jwt.domain")),
			Audience:  jwt.ClaimStrings{strconv.Itoa(id)},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 120)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        strconv.Itoa(jtiNew),
		},
	}

	tokenNew, err := jwtRegister.GenerateToken(claimsNew)
	if err != nil {
		return nil, fmt.Errorf("generate token : %w", err)
	}

	res.Token = tokenNew

	return
}

func (h *Handler) ChangeUserPasswordEmailConfirmHandler(ctx *fasthttp.RequestCtx) (err error) {
	defer func() {
		resFin := tool.Response(err, nil)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.ChangeUserPasswordEmailConfirmHandlerReq

	if err = json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		return fmt.Errorf(fmt.Errorf("handler: unmarshal request: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
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
		return fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	req.UID = int64(id)

	err = h.services.ChangeUserPasswordEmailConfirmService(req)
	if err != nil {
		return fmt.Errorf("change user password email confirm service: %w", err)
	}

	return
}

func (h *Handler) ChangeUserPasswordCodeHandler(ctx *fasthttp.RequestCtx) (res *model.ChangeUserPasswordCodeHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	var req model.ChangeUserPasswordCodeHandlerReq

	if err = json.Unmarshal(ctx.Request.Body(), &req); err != nil {
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
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	req.UID = int64(id)

	res, err = h.services.ChangeUserPasswordCodeService(req)
	if err != nil {
		return nil, fmt.Errorf("handler: change user email code service: %w", err)
	}

	return
}

func (h *Handler) ChangeUserProfilePictureLinkHandler(ctx *fasthttp.RequestCtx) (res *model.ChangeUserProfilePictureLinkHandlerRes, err error) {
	defer func() {
		resFin := tool.Response(err, res)
		jsonRes, _ := json.Marshal(resFin)
		ctx.Write(jsonRes)
		ctx.Response.SetStatusCode(resFin.Code)
	}()

	if !ctx.IsPost() {
		return nil, fmt.Errorf("handler: check method: %w", custom_errors.ErrWrongMethod)
	}

	file, err := ctx.FormFile("profile_picture")
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("form file: %w", err).Error()+": %w", custom_errors.ErrWrongInputData)
	}

	/*form, err := ctx.MultipartForm()
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
	}*/

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
		return nil, fmt.Errorf(fmt.Errorf("handler: parse with claims: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	idStr, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: get audience: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("handler: atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	res, err = h.services.ChangeUserProfilePictureLinkService(model.ChangeUserProfilePictureLinkHandlerReq{UID: int64(id), File: file})
	if err != nil {
		return nil, fmt.Errorf("change user profile picture link service: %w", err)
	}

	jtiNew := rand.Intn(99999999-10000000+1) + 10000000
	claimsNew := model.JWTCustomClaims{
		Email:      res.Email,
		Username:   res.Username,
		Name:       res.Name,
		Surname:    res.Surname,
		Patronymic: res.Patronymic,
		Birthday:   res.Birthday,
		Phone:      res.Phone,
		IMGLink:    res.ProfilePictureLink,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    fmt.Sprintf("https://%s", viper.GetString("jwt.domain")),
			Audience:  jwt.ClaimStrings{strconv.Itoa(id)},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 120)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        strconv.Itoa(jtiNew),
		},
	}

	tokenNew, err := jwtRegister.GenerateToken(claimsNew)
	if err != nil {
		return nil, fmt.Errorf("generate token : %w", err)
	}

	res.Token = tokenNew

	return
}
