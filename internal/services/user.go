package services

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/mitchellh/mapstructure"
)

type UserService interface {
	List(q string, page int64, size int64, ordering string) (n int64, data []map[string]interface{}, err error)
	CreateOrUpdate(info map[string]interface{}) error
	GetOrCreate(info map[string]interface{}) error

	Update(info map[string]interface{}) error
	Create(info map[string]interface{}) error
	Delete(ids []int64) (err error)
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &UserServiceImp{repo: repo}
}

type UserServiceImp struct {
	repo repositories.UserRepository
}

func (s *UserServiceImp) List(q string, page int64, size int64, ordering string) (n int64, data []map[string]interface{}, err error) {
	n, result, err := s.repo.List(q, (page-1)*size, size, ordering)
	if err != nil {
		return
	}
	for _, item := range result {
		data = append(data, map[string]interface{}{"Username": item.Username, "userId": item.UserId, "FirstName": item.FirstName, "LastName": item.LastName,
			"LanguageCode": item.LanguageCode, "Lang": item.Lang, "IsBot": item.IsBot, "IsPremium": item.IsPremium})
	}
	return
}

func (s *UserServiceImp) Create(info map[string]interface{}) error {
	data := models.User{}
	err := mapstructure.Decode(info, &data)
	if err != nil {
		return err
	}
	return s.repo.Create(&data)
}

func (s *UserServiceImp) Update(info map[string]interface{}) error {
	if info["userID"] == nil {
		return nil
	}
	userId := info["userID"].(int64)
	t := &models.User{}
	err := mapstructure.Decode(info, t)
	if err != nil {
		return err
	}

	return s.repo.Update(userId, info)
}

func (s *UserServiceImp) Delete(ids []int64) (err error) {
	return s.repo.Delete(ids)
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

func (s *UserServiceImp) GetOrCreate(info map[string]interface{}) error {
	userId := info["userID"].(int64)
	t, _ := s.repo.Get(userId)
	if t != nil {
		return nil
	}
	t = &models.User{}
	err := mapstructure.Decode(info, t)
	if err != nil {
		return err
	}
	return s.repo.Create(t)
}
