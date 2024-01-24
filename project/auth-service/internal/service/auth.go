package service

import (
	"auth-service/internal/api_db/redis_db"
	"auth-service/internal/api_db/reindexer_db"
	"auth-service/internal/model"
	"auth-service/pkg/custom_errors"
	"auth-service/pkg/email"
	"auth-service/pkg/jwtRegister"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
	mrand "math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type AuthServiceImpl struct {
	userApiDB      reindexer_db.UserApi
	generalApiDB   reindexer_db.GeneralApi
	authRedisApiDB redis_db.RedisApi
}

func NewAuthService(userApiDB reindexer_db.UserApi, generalApiDB reindexer_db.GeneralApi, authRedisApiDB redis_db.RedisApi) *AuthServiceImpl {
	return &AuthServiceImpl{
		userApiDB:      userApiDB,
		generalApiDB:   generalApiDB,
		authRedisApiDB: authRedisApiDB,
	}
}
func (s *AuthServiceImpl) SendEmailSingUpService(req model.SendEmailSingUpHandlerReq) (*model.SendEmailSingUpHandlerRes, error) {
	phone, err := getIntPhone(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("get int phone: %w", err)
	}

	found, err := s.userApiDB.CheckUserByPhoneDB(phone)
	if err != nil {
		return nil, fmt.Errorf("check user by phone DB: %w", err)
	}

	if found {
		return nil, fmt.Errorf("check user by phone DB: %w", custom_errors.ErrUserExist)
	}

	found, err = s.userApiDB.CheckUsernameDB(req.Username)
	if err != nil {
		return nil, fmt.Errorf("check user by username DB: %w", err)
	}

	if found {
		return nil, fmt.Errorf("check user by username DB: %w", custom_errors.ErrUserExist)
	}

	authCode := mrand.Intn(9999-1000+1) + 1000
	salt := make([]byte, 256)

	if _, err = rand.Read(salt); err != nil {
		return nil, fmt.Errorf(fmt.Errorf("fill salt array: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}
	hashVer := sha3.Sum512([]byte(fmt.Sprintf("%d%s", phone, hex.EncodeToString(salt))))

	res := &model.SendEmailSingUpHandlerRes{
		Login: req.Phone,
		Hash:  fmt.Sprintf("%x", hashVer),
	}

	redisAuthReq := model.CreateRedisAuthSingUpDB7Req{
		Key:      fmt.Sprintf("%d_-_%x", phone, hashVer),
		AuthCode: authCode,
		Salt:     hex.EncodeToString(salt),
		IP:       req.IP,
	}

	birthday, err := birthdayReqParse(req.Birthday)
	if err != nil {
		return nil, fmt.Errorf("birthday req parse: %w", err)
	}

	redisUserReq := model.CreateRedisUserSingUpDB8Req{
		Key:        fmt.Sprintf("%d_-_%x", phone, hashVer),
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
		Email:      req.Email,
		Username:   req.Username,
		Birthday:   birthday,
		Phone:      phone,
	}

	err = s.authRedisApiDB.CreateRedisAuthSingUpDB7(redisAuthReq)
	if err != nil {
		return nil, fmt.Errorf("create redis auth sing up DB: %w", err)
	}

	err = s.authRedisApiDB.CreateRedisUserSingUpDB8(redisUserReq)
	if err != nil {
		return nil, fmt.Errorf("create redis user sing up DB: %w", err)
	}

	err = email.SendMail(req.Email, authCode)

	if err != nil {
		return nil, fmt.Errorf("send mail: %w", err)
	}

	return res, nil
}

func (s *AuthServiceImpl) ConfirmEmailSingUpService(req model.ConfirmEmailSingUpHandlerReq) (*model.ConfirmEmailSingUpHandlerRes, error) {
	phone, err := getIntPhone(req.Login)
	if err != nil {
		return nil, fmt.Errorf("get int phone: %w", err)
	}

	values, err := s.authRedisApiDB.GetRedisAuthSingUpDB7(fmt.Sprintf("%d_-_%s", phone, req.Hash))
	if err != nil {
		return nil, fmt.Errorf("get redis user sing up DB: %w", err)
	}

	if values["auth_code"] == strconv.Itoa(req.Code) {
		hashVer := sha3.Sum512([]byte(fmt.Sprintf("%d%s", phone, values["salt"])))
		if fmt.Sprintf("%x", hashVer) == req.Hash && req.IP == values["ip"] {
			salt := make([]byte, 256)

			_, err = rand.Read(salt)
			if err != nil {
				return nil, fmt.Errorf(fmt.Errorf("fill salt array: %w", err).Error()+": %w", custom_errors.ErrInternal)
			}

			hashVerSecond := sha3.Sum512([]byte(fmt.Sprintf("%d%s", phone, hex.EncodeToString(salt))))
			res := &model.ConfirmEmailSingUpHandlerRes{Login: req.Login, Hash: fmt.Sprintf("%x", hashVerSecond)}

			err = s.authRedisApiDB.CreateRedisConfirmEmailDB7(model.CreateRedisConfirmEmailDB7Req{
				Key:  fmt.Sprintf("%d_-_%x", phone, hashVerSecond),
				Salt: hex.EncodeToString(salt),
				IP:   req.IP,
			})
			if err != nil {
				return nil, fmt.Errorf("create redis confirm email DB: %w", err)
			}

			return res, nil
		}
	}
	return nil, fmt.Errorf("wrong code: %w", custom_errors.ErrUnauthorized)
}

func (s *AuthServiceImpl) CreateUserService(req model.CreateUserHandlerReq) (*model.CreateUserHandlerRes, error) {
	phone, err := getIntPhone(req.Login)
	if err != nil {
		return nil, fmt.Errorf("get int phone: %w", err)
	}

	valuesSecondHashAuth, err := s.authRedisApiDB.GetRedisAuthSingUpDB7(fmt.Sprintf("%d_-_%s", phone, req.HashSecond))
	if err != nil {
		return nil, fmt.Errorf("get redis auth sing up DB: %w", err)
	}

	secondHashVer := sha3.Sum512([]byte(fmt.Sprintf("%d%s", phone, valuesSecondHashAuth["salt"])))

	if fmt.Sprintf("%x", secondHashVer) == req.HashSecond && req.IP == valuesSecondHashAuth["ip"] {
		valuesFirstHashUser, err := s.authRedisApiDB.GetRedisUserSingUpDB8(fmt.Sprintf("%d_-_%s", phone, req.HashFirst))
		if err != nil {
			return nil, fmt.Errorf("get redis user sing up DB: %w", err)
		}

		password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			return nil, fmt.Errorf(fmt.Errorf("generate from password: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}

		phoneNew, err := strconv.Atoi(valuesFirstHashUser["phone"])
		if err != nil {
			return nil, fmt.Errorf(fmt.Errorf("atoi phone: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}

		birthday, err := strconv.Atoi(valuesFirstHashUser["birthday"])
		if err != nil {
			return nil, fmt.Errorf(fmt.Errorf("atoi birthday: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}

		newUser := model.CreateUserDBReq{
			Email:      valuesFirstHashUser["email"],
			Phone:      int64(phoneNew),
			Username:   valuesFirstHashUser["username"],
			Name:       valuesFirstHashUser["name"],
			Surname:    valuesFirstHashUser["surname"],
			Patronymic: valuesFirstHashUser["patronymic"],
			Birthday:   int64(birthday),
			Password:   string(password),
		}

		err = s.userApiDB.CreateUserDB(newUser)
		if err != nil {
			return nil, fmt.Errorf("create user DB: %w", err)
		}

		user, err := s.userApiDB.GetUserByPhoneDB(int64(phoneNew))
		if err != nil {
			return nil, fmt.Errorf("get user by email DB: %w", err)
		}

		newGeneral := model.CreateGeneralDBReq{
			UID: user.ID,
			IP:  req.IP,
		}

		err = s.generalApiDB.CreateGeneralDB(newGeneral)
		if err != nil {
			return nil, fmt.Errorf("create general DB: %w", err)
		}

		general, err := s.generalApiDB.GetRecordByUIDDB(user.ID)
		if err != nil {
			return nil, fmt.Errorf("get record by UID DB: %w", err)
		}

		arr, err := s.authRedisApiDB.GetJtiArrayDB(phone)
		if err != nil {
			return nil, fmt.Errorf("get jti array DB: %w", err)
		}

		if len(arr) > 5 {
			err = s.authRedisApiDB.RemoveLastJtiRecordDB(phone)
			if err != nil {
				return nil, fmt.Errorf("remove last jti record DB: %w", err)
			}
		}

		jtiNew := mrand.Intn(99999999-10000000+1) + 10000000
		redisRes := s.authRedisApiDB.CheckExistingJtiDB(phone, jtiNew)

		if redisRes == 0 || redisRes == -1 {
			err = s.authRedisApiDB.AddJtiToArrayDB(phone, jtiNew)
			if err != nil {
				return nil, fmt.Errorf("add jti to array: %w", err)
			}
		} else {
			err = s.authRedisApiDB.RemoveNewJtiDB(phone, jtiNew)
			if err != nil {
				return nil, fmt.Errorf("remove new jti: %w", err)
			}
			jtiNew = mrand.Intn(99999999-10000000+1) + 10000000
			for i := 0; s.authRedisApiDB.CheckExistingJtiDB(phone, jtiNew) != -1 && i < 10; i++ {
				err = s.authRedisApiDB.RemoveNewJtiDB(phone, jtiNew)
				if err != nil {
					return nil, fmt.Errorf("remove new jti: %w", err)
				}
				jtiNew = mrand.Intn(99999999-10000000+1) + 10000000
			}
			err = s.authRedisApiDB.AddJtiToArrayDB(phone, jtiNew)
			if err != nil {
				return nil, fmt.Errorf("add jti to array: %w", err)
			}
		}

		claims := model.JWTCustomClaims{
			Email:      user.Email,
			Username:   user.Username,
			Name:       user.Name,
			Surname:    user.Surname,
			Patronymic: user.Patronymic,
			Birthday:   birthdayResParse(user.Birthday),
			Phone:      req.Login,
			IMGLink:    general.IMGLink,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    fmt.Sprintf("https://%s", viper.GetString("jwt.domain")),
				Audience:  jwt.ClaimStrings{strconv.Itoa(int(user.ID))},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 120)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        strconv.Itoa(jtiNew),
			},
		}

		token, err := jwtRegister.GenerateToken(claims)
		res := &model.CreateUserHandlerRes{
			ID:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
			DFA:     general.DFA,
			Token:   token,
		}

		return res, nil
	}

	return nil, fmt.Errorf("hash or ip is invalid: %w", custom_errors.ErrUnauthorized)
}

func (s *AuthServiceImpl) GenerateHashService(req model.GenerateHashHandlerReq) (*model.GenerateHashHandlerRes, error) {
	phone, err := getIntPhone(req.Login)
	if err != nil {
		return nil, fmt.Errorf("get int phone: %w", err)
	}

	salt := make([]byte, 256)

	user, err := s.userApiDB.GetUserByPhoneDB(phone)
	if err != nil {
		return nil, fmt.Errorf("get user by phone DB: %w", err)
	}

	var trial bool

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(strings.Trim(req.Password, " ")))
	if err == nil {
		trial = true
	}

	general, err := s.generalApiDB.GetRecordByUIDDB(user.ID)
	if err != nil {
		return nil, fmt.Errorf("get record by UID DB: %w", err)
	}

	_, err = rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("service singIn salt generating: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	hashVer := sha3.Sum512([]byte(fmt.Sprintf("%d%s%s", user.Phone, user.Password, hex.EncodeToString(salt))))

	loss := &model.GenerateHashHandlerRes{
		Login: req.Login,
		Hash:  fmt.Sprintf("%x", hashVer),
		DFA:   general.DFA,
	}

	if !general.ActiveStatus || !trial {
		return nil, fmt.Errorf("check status and pass: %w", custom_errors.ErrUnauthorized)
	}

	reqRedis := &model.CreateRedisUserSingInDB7Req{
		Key:      fmt.Sprintf("%d__%s", user.Phone, hex.EncodeToString(hashVer[:])),
		Password: trial,
		Salt:     hex.EncodeToString(salt),
		ID:       general.ID,
		IP:       req.IP,
		DFA:      general.DFA,
	}

	err = s.authRedisApiDB.CreateRedisUserSingInDB7(reqRedis)
	if err != nil {
		return nil, fmt.Errorf("createRedisUserDB: %w", err)
	}

	return loss, nil
}

func (s *AuthServiceImpl) CheckDFAService(req model.CheckDFAHandlerReq) (*model.CheckDFAHandlerRes, error) {
	phone, err := getIntPhone(req.Login)
	if err != nil {
		return nil, fmt.Errorf("get int phone: %w", err)
	}

	var res *model.CheckDFAHandlerRes

	values, err := s.authRedisApiDB.GetRedisUserSingInDB7(fmt.Sprintf("%d__%s", phone, req.Hash))
	if err != nil {
		return nil, fmt.Errorf("redis get user data by key DB7: %w", err)
	}

	user, err := s.userApiDB.GetUserByPhoneDB(phone)
	if err != nil {
		return nil, fmt.Errorf("get user by phone DB: %w", err)
	}

	general, err := s.generalApiDB.GetRecordByUIDDB(user.ID)
	if err != nil {
		return nil, fmt.Errorf("general get record by UID DB: %w", err)
	}

	if general.DFA == 0 {
		arr, err := s.authRedisApiDB.GetJtiArrayDB(phone)
		if err != nil {
			return nil, fmt.Errorf("get jti array DB: %w", err)
		}

		if len(arr) > 5 {
			err = s.authRedisApiDB.RemoveLastJtiRecordDB(phone)
			if err != nil {
				return nil, fmt.Errorf("remove last jti record DB: %w", err)
			}
		}

		jtiNew := mrand.Intn(99999999-10000000+1) + 10000000
		redisRes := s.authRedisApiDB.CheckExistingJtiDB(phone, jtiNew)

		if redisRes == 0 || redisRes == -1 {
			err = s.authRedisApiDB.AddJtiToArrayDB(phone, jtiNew)
			if err != nil {
				return nil, fmt.Errorf("add jti to array: %w", err)
			}
		} else {
			err = s.authRedisApiDB.RemoveNewJtiDB(phone, jtiNew)
			if err != nil {
				return nil, fmt.Errorf("remove new jti: %w", err)
			}
			jtiNew = mrand.Intn(99999999-10000000+1) + 10000000
			for i := 0; s.authRedisApiDB.CheckExistingJtiDB(phone, jtiNew) != -1 && i < 10; i++ {
				err = s.authRedisApiDB.RemoveNewJtiDB(phone, jtiNew)
				if err != nil {
					return nil, fmt.Errorf("remove new jti: %w", err)
				}
				jtiNew = mrand.Intn(99999999-10000000+1) + 10000000
			}
			err = s.authRedisApiDB.AddJtiToArrayDB(phone, jtiNew)
			if err != nil {
				return nil, fmt.Errorf("add jti to array: %w", err)
			}
		}

		claims := model.JWTCustomClaims{
			Email:      user.Email,
			Username:   user.Username,
			Name:       user.Name,
			Surname:    user.Surname,
			Patronymic: user.Patronymic,
			Birthday:   birthdayResParse(user.Birthday),
			Phone:      req.Login,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    fmt.Sprintf("https://%s", viper.GetString("jwt.domain")),
				Audience:  jwt.ClaimStrings{strconv.Itoa(int(user.ID))},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 120)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        strconv.Itoa(jtiNew),
			},
		}

		token, err := jwtRegister.GenerateToken(claims)
		if err != nil {
			return nil, fmt.Errorf("sing in second: %w", err)
		}

		res = &model.CheckDFAHandlerRes{
			ID:         user.ID,
			Name:       user.Name,
			Surname:    user.Surname,
			Patronymic: user.Patronymic,
			Email:      user.Email,
			Phone:      req.Login,
			DFA:        general.DFA,
			Token:      token,
		}

		return res, nil
	}

	hashVer := sha3.Sum512([]byte(fmt.Sprintf("%d%s%s", phone, user.Password, values["salt"])))

	if fmt.Sprintf("%x", hashVer) == req.Hash && req.IP == values["ip"] {
		authCode := mrand.Intn(9999-1000+1) + 1000
		err = s.authRedisApiDB.CreateRedisAuthSingInDB8(&model.CreateRedisAuthSingInDB8Req{
			Key:      fmt.Sprintf("%d__%s", phone, hex.EncodeToString(hashVer[:])),
			AuthCode: authCode,
			ID:       general.ID,
			IP:       req.IP,
			DFA:      general.DFA,
		})
		if err != nil {
			return nil, fmt.Errorf("create redis auth data: %w", err)
		}

		err = email.SendMail(user.Email, authCode)
		if err != nil {
			return nil, fmt.Errorf("email: %w", err)
		}

		return res, nil
	}

	return nil, fmt.Errorf("check hash or ip: %w", custom_errors.ErrUnauthorized)
}

func (s *AuthServiceImpl) ConfirmEmailSingInService(req model.ConfirmEmailSingInHandlerReq) (*model.ConfirmEmailSingInHandlerRes, error) {
	phone, err := getIntPhone(req.Login)
	if err != nil {
		return nil, fmt.Errorf("get int phone: %w", err)
	}

	code := strconv.Itoa(req.Code)

	values, err := s.authRedisApiDB.GetRedisAuthSingInDB8(fmt.Sprintf("%d__%s", phone, req.Hash))
	if err != nil {
		return nil, fmt.Errorf("get auth date DB: %w", err)
	}

	if values["auth_code"] == code {
		arr, err := s.authRedisApiDB.GetJtiArrayDB(phone)
		if err != nil {
			return nil, fmt.Errorf("get jti array DB: %w", err)
		}

		if len(arr) > 5 {
			err = s.authRedisApiDB.RemoveLastJtiRecordDB(phone)
			if err != nil {
				return nil, fmt.Errorf("remove last jti record DB: %w", err)
			}
		}

		jtiNew := mrand.Intn(99999999-10000000+1) + 10000000
		redisRes := s.authRedisApiDB.CheckExistingJtiDB(phone, jtiNew)

		if redisRes == 0 || redisRes == -1 {
			err = s.authRedisApiDB.AddJtiToArrayDB(phone, jtiNew)
			if err != nil {
				return nil, fmt.Errorf("add jti to array: %w", err)
			}
		} else {
			err = s.authRedisApiDB.RemoveNewJtiDB(phone, jtiNew)
			if err != nil {
				return nil, fmt.Errorf("remove new jti: %w", err)
			}
			jtiNew = mrand.Intn(99999999-10000000+1) + 10000000
			for i := 0; s.authRedisApiDB.CheckExistingJtiDB(phone, jtiNew) != -1 && i < 10; i++ {
				err = s.authRedisApiDB.RemoveNewJtiDB(phone, jtiNew)
				if err != nil {
					return nil, fmt.Errorf("remove new jti: %w", err)
				}
				jtiNew = mrand.Intn(99999999-10000000+1) + 10000000
			}
			err = s.authRedisApiDB.AddJtiToArrayDB(phone, jtiNew)
			if err != nil {
				return nil, fmt.Errorf("add jti to array: %w", err)
			}
		}

		user, err := s.userApiDB.GetUserByPhoneDB(phone)
		if err != nil {
			return nil, fmt.Errorf("get user by phone DB: %w", err)
		}

		general, err := s.generalApiDB.GetRecordByUIDDB(user.ID)
		if err != nil {
			return nil, fmt.Errorf("get record by UID: %w", err)
		}

		claims := model.JWTCustomClaims{
			Email:      user.Email,
			Username:   user.Username,
			Name:       user.Name,
			Surname:    user.Surname,
			Patronymic: user.Patronymic,
			Phone:      req.Login,
			Birthday:   birthdayResParse(user.Birthday),
			IMGLink:    general.IMGLink,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    fmt.Sprintf("https://%s", viper.GetString("jwt.domain")),
				Audience:  jwt.ClaimStrings{strconv.Itoa(int(user.ID))},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 120)),
				NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        strconv.Itoa(jtiNew),
			},
		}

		token, err := jwtRegister.GenerateToken(claims)
		if err != nil {
			return nil, fmt.Errorf("sing in second: %w", err)
		}

		res := &model.ConfirmEmailSingInHandlerRes{
			ID:         user.ID,
			Name:       user.Name,
			Surname:    user.Surname,
			Patronymic: user.Patronymic,
			Email:      user.Email,
			Phone:      req.Login,
			DFA:        general.DFA,
			Token:      token,
		}

		return res, nil
	}

	return nil, fmt.Errorf("wrong data: %w", custom_errors.ErrUnauthorized)
}

func (s *AuthServiceImpl) SendEmailRecoverService(req *model.SendEmailRecoverHandlerReq) (*model.SendEmailRecoverHandlerRes, error) {
	phone, err := getIntPhone(req.Login)
	if err != nil {
		return nil, fmt.Errorf("get int phone: %w", err)
	}

	salt := make([]byte, 256)

	_, err = rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("salt generating: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	user, err := s.userApiDB.GetUserByPhoneDB(phone)
	if err != nil {
		return nil, fmt.Errorf("get user by phone DB: %w", err)
	}

	general, err := s.generalApiDB.GetRecordByUIDDB(user.ID)
	if err != nil {
		return nil, fmt.Errorf("get record by UID DB: %w", err)
	}

	authCode := mrand.Intn(9999-1000+1) + 1000

	hashVer := sha3.Sum512([]byte(fmt.Sprintf("%d%s", phone, hex.EncodeToString(salt))))

	res := &model.SendEmailRecoverHandlerRes{
		Login: req.Login,
		Hash:  fmt.Sprintf("%x", hashVer),
		DFA:   general.DFA,
	}

	err = s.authRedisApiDB.CreateRedisGeneralRecoveryDB7(model.CreateRedisGeneralRecoveryDB7Req{
		Key:          fmt.Sprintf("%d_recovery_%s", phone, fmt.Sprintf("%x", hashVer)),
		AuthCode:     authCode,
		ActiveStatus: general.ActiveStatus,
		Salt:         hex.EncodeToString(salt),
		IP:           req.IP,
		ID:           user.ID,
		DFA:          general.DFA,
	})
	if err != nil {
		return nil, fmt.Errorf("create redis general recovery DB7: %w", err)
	}

	err = email.SendMail(user.Email, authCode)
	if err != nil {
		return nil, fmt.Errorf("send email: %w", err)
	}

	return res, nil
}

func (s *AuthServiceImpl) ConfirmEmailRecoverService(req *model.ConfirmEmailRecoverHandlerReq) (*model.ConfirmEmailRecoverHandlerRes, error) {
	phone, err := getIntPhone(req.Login)
	if err != nil {
		return nil, fmt.Errorf("get int phone: %w", err)
	}

	values, err := s.authRedisApiDB.GetRedisGeneralRecoveryDB7(fmt.Sprintf("%d_recovery_%s", phone, req.Hash))
	if err != nil {
		return nil, fmt.Errorf("get redis general recovery DB7: %w", err)
	}

	if values["auth_code"] == strconv.Itoa(int(req.Code)) && req.IP == values["ip"] && values["active_status"] == "1" {
		salt := make([]byte, 256)

		_, err := rand.Read(salt)
		if err != nil {
			return nil, fmt.Errorf(fmt.Errorf("salt generating: %w", err).Error()+": %w", custom_errors.ErrInternal)
		}

		hashVer := sha3.Sum512([]byte(fmt.Sprintf("%s%s", req.Hash, hex.EncodeToString(salt))))
		authCode := mrand.Intn(99999999-10000000+1) + 10000000

		err = s.authRedisApiDB.CreateRedisConfirmRecoveryDB8(model.CreateRedisConfirmRecoveryDB8Req{
			Key:      fmt.Sprintf("%d_recovery_password_%s", phone, fmt.Sprintf("%x", hashVer)),
			AuthCode: authCode,
			ID:       values["id"],
		})

		res := &model.ConfirmEmailRecoverHandlerRes{Code: authCode, Hash: fmt.Sprintf("%x", hashVer)}

		return res, nil
	}

	return nil, fmt.Errorf("wrong data: %w", custom_errors.ErrUnauthorized)
}

func (s *AuthServiceImpl) ChangePasswordRecoverService(req *model.ChangePasswordRecoverHandlerReq) (*model.ChangePasswordRecoverHandlerRes, error) {
	phone, err := getIntPhone(req.Login)
	if err != nil {
		return nil, fmt.Errorf("get int phone: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("generate hash from password: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	values, err := s.authRedisApiDB.GetRedisConfirmRecoveryDB8(fmt.Sprintf("%d_recovery_password_%s", phone, req.Hash))
	if err != nil {
		return nil, fmt.Errorf("get redis confirm recovery DB8: %w", err)
	}

	idInt, err := strconv.Atoi(values["id"])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("id atoi: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	if values["auth_code"] == strconv.Itoa(req.Code) {
		err = s.userApiDB.ChangePassDB(int64(idInt), string(passwordHash))
		if err != nil {
			return nil, fmt.Errorf("change pass DB: %w", err)
		}

		err = s.generalApiDB.ChangePassTimeDB(int64(idInt))
		if err != nil {
			return nil, fmt.Errorf("change pass time DB: %w", err)
		}

		return &model.ChangePasswordRecoverHandlerRes{Message: "Password successfully changed!"}, nil
	}

	return nil, fmt.Errorf("wrong data: %w", custom_errors.ErrUnauthorized)
}

func (s *AuthServiceImpl) AuthMiddleWareService(claims jwt.MapClaims) (*model.AuthMiddleWareHandlerRes, error) {
	id, err := claims.GetAudience()
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("get audience: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	idInt, err := strconv.Atoi(id[0])
	if err != nil {
		return nil, fmt.Errorf(fmt.Errorf("atoi id: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	found, err := s.userApiDB.CheckUserByIDDB(int64(idInt))
	if err != nil {
		return nil, fmt.Errorf("check user by ID DB: %w", err)
	} else if !found {
		return nil, fmt.Errorf("check user by ID DB: %w", custom_errors.ErrNotFound)
	}

	res := &model.AuthMiddleWareHandlerRes{
		ID:         int64(idInt),
		Name:       claims["name"].(string),
		Surname:    claims["surname"].(string),
		Patronymic: claims["patronymic"].(string),
		Birthday:   claims["birthday"].(string),
	}

	return res, nil
}

func getIntPhone(str string) (int64, error) {
	strRune := []rune(str)
	var phone []rune
	for _, elem := range strRune {
		if unicode.IsDigit(elem) {
			phone = append(phone, elem)
		}
	}

	phoneInt, err := strconv.Atoi(string(phone))
	if err != nil {
		return 0, fmt.Errorf(fmt.Errorf("atoi: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return int64(phoneInt), nil
}

func birthdayResParse(intBirthday int64) string {
	newDate := strings.Split(strings.Split(time.Unix(intBirthday, 0).String(), " ")[0], "-")

	return newDate[2] + "." + newDate[1] + "." + newDate[0]
}

func birthdayReqParse(str string) (int64, error) {
	birthday, err := time.Parse("02.01.2006", str)
	if err != nil {
		return 0, fmt.Errorf(fmt.Errorf("atoi: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	return birthday.Unix(), nil
}
