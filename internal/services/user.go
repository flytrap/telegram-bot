package services

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/mitchellh/mapstructure"
)

type UserService interface {
	List(q string, page int64, size int64) (data []map[string]interface{}, err error)
	CreateOrUpdate(info map[string]interface{}) error
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &UserServiceImp{repo: repo}
}

type UserServiceImp struct {
	repo repositories.UserRepository
}

func (s *UserServiceImp) List(q string, page int64, size int64) (data []map[string]interface{}, err error) {
	result, err := s.repo.GetMany(q, (page-1)*size, size)
	if err != nil {
		return nil, err
	}
	for _, item := range result {
		data = append(data, map[string]interface{}{"Username": item.Username, "userID": item.UserID, "FirstName": item.FirstName, "LastName": item.LastName,
			"LanguageCode": item.LanguageCode, "Lang": item.Lang, "IsBot": item.IsBot, "IsPremium": item.IsPremium})
	}
	return
}

func (s *UserServiceImp) Update(info map[string]interface{}) error {
	userId := info["userId"].(int64)
	t := &models.User{}
	err := mapstructure.Decode(info, t)
	if err != nil {
		return err
	}

	return s.repo.Update(userId, info)
}

func (s *UserServiceImp) CreateOrUpdate(info map[string]interface{}) error {
	userId := info["userID"].(int64)
	t, _ := s.repo.Get(userId)
	if t != nil {
		return s.Update(info)
	}
	t = &models.User{}
	err := mapstructure.Decode(info, t)
	if err != nil {
		return err
	}
	return s.repo.Create(t)
}
