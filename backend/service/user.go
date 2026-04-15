package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/sanbei101/second/middleware"
	"github.com/sanbei101/second/model"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(phone, password, nickname string, role model.UserRole) (*model.User, string, error) {
	var exist model.User
	if err := s.db.Where("phone = ?", phone).First(&exist).Error; err == nil {
		log.Warn().Str("phone", phone).Msg("register failed: phone already exists")
		return nil, "", errors.New("phone already registered")
	}

	user := model.User{
		Phone:    phone,
		Password: password,
		Nickname: nickname,
		Role:     role,
		Avatar:   "https://img.yzcdn.cn/vant/cat.jpeg",
	}
	if err := s.db.Create(&user).Error; err != nil {
		log.Error().Err(err).Str("phone", phone).Msg("register failed: db error")
		return nil, "", err
	}

	log.Info().Uint("user_id", user.ID).Str("phone", phone).Msg("user registered")
	token, err := s.generateToken(&user)
	if err != nil {
		log.Error().Err(err).Uint("user_id", user.ID).Msg("generate token failed")
	}
	return &user, token, err
}

func (s *UserService) Login(phone, password string) (*model.User, string, error) {
	var user model.User
	if err := s.db.Where("phone = ? AND password = ?", phone, password).First(&user).Error; err != nil {
		log.Warn().Str("phone", phone).Msg("login failed: invalid credentials")
		return nil, "", errors.New("invalid phone or password")
	}

	log.Info().Uint("user_id", user.ID).Str("phone", phone).Msg("user logged in")
	token, err := s.generateToken(&user)
	if err != nil {
		log.Error().Err(err).Uint("user_id", user.ID).Msg("generate token failed")
	}
	return &user, token, err
}

func (s *UserService) WxLogin(openid string, role model.UserRole) (*model.User, string, error) {
	var user model.User
	err := s.db.Where("openid = ?", openid).First(&user).Error
	if err == nil {
		log.Info().Uint("user_id", user.ID).Str("openid", openid).Msg("wx login success")
		token, err := s.generateToken(&user)
		return &user, token, err
	}

	user = model.User{
		Openid:   openid,
		Nickname: "微信用户",
		Role:     role,
		Avatar:   "https://img.yzcdn.cn/vant/cat.jpeg",
	}
	if err := s.db.Create(&user).Error; err != nil {
		log.Error().Err(err).Str("openid", openid).Msg("wx register failed")
		return nil, "", err
	}

	log.Info().Uint("user_id", user.ID).Str("openid", openid).Msg("wx user registered")
	token, err := s.generateToken(&user)
	return &user, token, err
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		log.Warn().Err(err).Uint("user_id", id).Msg("get user failed")
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Update(id uint, nickname, avatar, phone string) error {
	updates := map[string]interface{}{}
	if nickname != "" {
		updates["nickname"] = nickname
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}
	if phone != "" {
		updates["phone"] = phone
	}
	if err := s.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Error().Err(err).Uint("user_id", id).Msg("update user failed")
		return err
	}
	log.Info().Uint("user_id", id).Msg("user updated")
	return nil
}

func (s *UserService) generateToken(user *model.User) (string, error) {
	claims := middleware.Claims{
		UserID: user.ID,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JWTSecret)
}
