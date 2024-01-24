package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"lk-service/internal/api_db/reindexer_db"
	"lk-service/internal/model"
	"lk-service/pkg/custom_errors"
	"lk-service/pkg/email"
	"lk-service/pkg/files"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type LKServiceImpl struct {
	userApiDB      reindexer_db.UserApi
	generalApiDB   reindexer_db.GeneralApi
	emailCodeApiDB reindexer_db.EmailCodeApi
}

func NewLKService(userApiDB reindexer_db.UserApi, generalApiDB reindexer_db.GeneralApi, emailCodeApiDB reindexer_db.EmailCodeApi) *LKServiceImpl {
	return &LKServiceImpl{
		userApiDB:      userApiDB,
		generalApiDB:   generalApiDB,
		emailCodeApiDB: emailCodeApiDB,
	}
}

func (s *LKServiceImpl) GetMyProfileService(req model.GetMyProfileHandlerReq) (*model.GetMyProfileHandlerRes, error) {
	user, err := s.userApiDB.FindUserByIDDB(req.UID)
	if err != nil {
		return nil, fmt.Errorf("find user by ID DB: %w", err)
	}

	record, err := s.generalApiDB.GetRecordByUIDDB(user.ID)
	if err != nil {
		return nil, fmt.Errorf("get record by UID DB: %w", err)
	}

	phoneStr := strconv.Itoa(int(user.Phone))

	return &model.GetMyProfileHandlerRes{
		ID:                 user.ID,
		Email:              user.Email,
		Username:           user.Username,
		Name:               user.Name,
		Surname:            user.Surname,
		Patronymic:         user.Patronymic,
		Birthday:           birthdayResParse(user.Birthday),
		Phone:              "+7(" + phoneStr[1:4] + ")" + phoneStr[4:7] + "-" + phoneStr[7:9] + "-" + phoneStr[9:],
		ProfilePictureLink: record.IMGLink}, nil
}

func (s *LKServiceImpl) GetUserInfoService(req model.GetUserInfoHandlerReq) (*model.GetUserInfoHandlerRes, error) {
	user, err := s.userApiDB.FindUserByUsernameDB(req.Username)
	if err != nil {
		return nil, fmt.Errorf("find user by username DB: %w", err)
	}

	record, err := s.generalApiDB.GetRecordByUIDDB(user.ID)
	if err != nil {
		return nil, fmt.Errorf("get record by UID DB: %w", err)
	}

	phoneStr := strconv.Itoa(int(user.Phone))

	return &model.GetUserInfoHandlerRes{
		ID:                 user.ID,
		Email:              user.Email,
		Username:           user.Username,
		Name:               user.Name,
		Surname:            user.Surname,
		Patronymic:         user.Patronymic,
		Birthday:           birthdayResParse(user.Birthday),
		Phone:              "+7(" + phoneStr[1:4] + ")" + phoneStr[4:7] + "-" + phoneStr[7:9] + "-" + phoneStr[9:],
		ProfilePictureLink: record.IMGLink}, nil
}

func (s *LKServiceImpl) ChangeUserInfoService(req model.ChangeUserInfoHandlerReq) (*model.ChangeUserInfoHandlerRes, error) {
	var phone int64
	var err error
	var birthday int64
	if req.Phone != "" {
		phone, err = getIntPhone(req.Phone)
		if err != nil {
			return nil, fmt.Errorf("get int phone: %w", err)
		}
	}

	if req.Username != "" {
		found, err := s.userApiDB.CheckUserByUsernameDB(req.Username)
		if err != nil {
			return nil, fmt.Errorf("check user by username DB: %w", err)
		}

		if found {
			return nil, fmt.Errorf(fmt.Errorf("user with this username alreay exist").Error()+": %w", custom_errors.ErrWrongInputData)
		}
	}

	if phone != 0 {
		found, err := s.userApiDB.CheckUserByPhoneDB(phone)
		if err != nil {
			return nil, fmt.Errorf("check user by username DB: %w", err)
		}

		if found {
			return nil, fmt.Errorf(fmt.Errorf("user with this phone alreay exist").Error()+": %w", custom_errors.ErrWrongInputData)
		}
	}

	if req.Birthday != "" {
		birthday, err = birthdayReqParse(req.Birthday)
		if err != nil {
			return nil, fmt.Errorf("birthday req parse: %w", err)
		}
	}

	if err = s.userApiDB.UpdateUserInfoDB(req, phone, birthday); err != nil {
		return nil, fmt.Errorf("update user info DB: %w", err)
	}

	user, err := s.userApiDB.FindUserByIDDB(req.ID)
	if err != nil {
		return nil, fmt.Errorf("find user by ID DB: %w", err)
	}

	record, err := s.generalApiDB.GetRecordByUIDDB(user.ID)
	if err != nil {
		return nil, fmt.Errorf("get record by UID DB: %w", err)
	}

	phoneStr := strconv.Itoa(int(user.Phone))

	return &model.ChangeUserInfoHandlerRes{
		ID:                 user.ID,
		Email:              user.Email,
		Username:           user.Username,
		Name:               user.Name,
		Surname:            user.Surname,
		Patronymic:         user.Patronymic,
		Birthday:           birthdayResParse(user.Birthday),
		Phone:              "+7(" + phoneStr[1:4] + ")" + phoneStr[4:7] + "-" + phoneStr[7:9] + "-" + phoneStr[9:],
		ProfilePictureLink: record.IMGLink}, nil
}

func (s *LKServiceImpl) ChangeUserEmailConfirmService(req model.ChangeUserEmailConfirmHandlerReq) error {
	found, err := s.userApiDB.CheckUserByEmailDB(req.Email)
	if err != nil {
		return fmt.Errorf("check user by email DB: %w", err)
	}

	if found {
		return fmt.Errorf(fmt.Errorf("user with this email alreay exist").Error()+": %w", custom_errors.ErrWrongInputData)
	}

	code := rand.Intn(9999-1000+1) + 1000

	err = email.SendMail(req.Email, code)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	err = s.emailCodeApiDB.CreateRecordEmailCodeDB(req.UID, int64(code), req.Email)
	if err != nil {
		return fmt.Errorf("create record email code DB: %w", err)
	}

	return nil
}

func (s *LKServiceImpl) ChangeUserEmailCodeService(req model.ChangeUserEmailCodeHandlerReq) (*model.ChangeUserEmailCodeHandlerRes, error) {
	codeRecord, err := s.emailCodeApiDB.GetRecordEmailCodeDB(req.UID)
	if err != nil {
		return nil, fmt.Errorf("get record email code DB: %w", err)
	}

	if int64(req.Code) != codeRecord.Code {
		return nil, fmt.Errorf("wrong code: %w", custom_errors.ErrWrongInputData)
	}

	user, err := s.userApiDB.UpdateUserEmailDB(req.UID, codeRecord.Email)
	if err != nil {
		return nil, fmt.Errorf("update user email DB: %w", err)
	}

	phoneStr := strconv.Itoa(int(user.Phone))

	return &model.ChangeUserEmailCodeHandlerRes{
		ID:         user.ID,
		Email:      user.Email,
		Username:   user.Username,
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: user.Patronymic,
		Birthday:   birthdayResParse(user.Birthday),
		Phone:      "+7(" + phoneStr[1:4] + ")" + phoneStr[4:7] + "-" + phoneStr[7:9] + "-" + phoneStr[9:]}, err
}

func (s *LKServiceImpl) ChangeUserPasswordEmailConfirmService(req model.ChangeUserPasswordEmailConfirmHandlerReq) error {
	user, err := s.userApiDB.FindUserByIDDB(req.UID)
	if err != nil {
		return fmt.Errorf("find user by id DB: %w", err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.PasswordNew), 10)
	if err != nil {
		return fmt.Errorf(fmt.Errorf("generate from password: %w", err).Error()+": %w", custom_errors.ErrInternal)
	}

	code := rand.Intn(9999-1000+1) + 1000

	err = email.SendMail(user.Email, code)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	err = s.emailCodeApiDB.CreateRecordEmailCodeWithPasswordDB(req.UID, int64(code), user.Email, string(password))
	if err != nil {
		return fmt.Errorf("create record email code with password DB: %w", err)
	}

	return nil
}

func (s *LKServiceImpl) ChangeUserPasswordCodeService(req model.ChangeUserPasswordCodeHandlerReq) (*model.ChangeUserPasswordCodeHandlerRes, error) {
	codeRecord, err := s.emailCodeApiDB.GetRecordEmailCodeDB(req.UID)
	if err != nil {
		return nil, fmt.Errorf("get record email code DB: %w", err)
	}

	if int64(req.Code) != codeRecord.Code {
		return nil, fmt.Errorf("wrong code: %w", custom_errors.ErrWrongInputData)
	}

	user, err := s.userApiDB.UpdateUserPasswordDB(req.UID, codeRecord.Password)
	if err != nil {
		return nil, fmt.Errorf("update user email DB: %w", err)
	}

	phoneStr := strconv.Itoa(int(user.Phone))

	return &model.ChangeUserPasswordCodeHandlerRes{
		ID:         user.ID,
		Email:      user.Email,
		Username:   user.Username,
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: user.Patronymic,
		Birthday:   birthdayResParse(user.Birthday),
		Phone:      "+7(" + phoneStr[1:4] + ")" + phoneStr[4:7] + "-" + phoneStr[7:9] + "-" + phoneStr[9:]}, err
}

func (s *LKServiceImpl) ChangeUserProfilePictureLinkService(req model.ChangeUserProfilePictureLinkHandlerReq) (*model.ChangeUserProfilePictureLinkHandlerRes, error) {
	user, err := s.userApiDB.FindUserByIDDB(req.UID)
	if err != nil {
		return nil, fmt.Errorf("find user by id DB: %w", err)
	}
	record, err := s.generalApiDB.GetRecordByUIDDB(req.UID)
	if err != nil {
		return nil, fmt.Errorf("get record by uid DB: %w", err)
	}

	fileLink, err := files.DownloadAvatars(record.UID, req.File)
	if err != nil || fileLink == "" {
		return nil, fmt.Errorf("download avatars: %w", err)
	}

	err = s.generalApiDB.ChangeProfilePictureLinkDB(fileLink, record.ID)
	if err != nil {
		return nil, fmt.Errorf("change profile picture link DB: %w", err)
	}

	phoneStr := strconv.Itoa(int(user.Phone))

	return &model.ChangeUserProfilePictureLinkHandlerRes{
		ID:                 user.ID,
		Email:              user.Email,
		Username:           user.Username,
		Name:               user.Name,
		Surname:            user.Surname,
		Patronymic:         user.Patronymic,
		Birthday:           birthdayResParse(user.Birthday),
		Phone:              "+7(" + phoneStr[1:4] + ")" + phoneStr[4:7] + "-" + phoneStr[7:9] + "-" + phoneStr[9:],
		ProfilePictureLink: fileLink}, nil
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
