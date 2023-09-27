package services

import (
	"context"
	"fmt"
	"time"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/flytrap/telegram-bot/pkg/redis"
	"github.com/mitchellh/mapstructure"
)

type UserService interface {
	Check(ctx context.Context, userId int64) bool
	List(q string, page int64, size int64, ordering string, data interface{}) (n int64, err error)
	CreateOrUpdate(ctx context.Context, info map[string]interface{}) error
	GetOrCreate(ctx context.Context, info map[string]interface{}) error

	Update(ctx context.Context, info map[string]interface{}) error
	Create(ctx context.Context, info map[string]interface{}) error
	Delete(ctx context.Context, ids []int64) (err error)
}

func NewUserService(repo repositories.UserRepository, store *redis.Store) UserService {
	return &userServiceImp{repo: repo, store: store}
}

type userServiceImp struct {
	store *redis.Store
	repo  repositories.UserRepository
}

func (s *userServiceImp) Check(ctx context.Context, userId int64) bool {
	key := fmt.Sprintf("user:%d", userId)
	if s.store.IsExist(ctx, key) {
		return true
	}
	u, err := s.repo.Get(userId)
	if err == nil && u != nil {
		s.store.Set(ctx, key, u.Username, time.Second*60*60*3)
		return true
	}
	return false
}

func (s *userServiceImp) List(q string, page int64, size int64, ordering string, data interface{}) (n int64, err error) {
	n, err = s.repo.List(q, (page-1)*size, size, ordering, data)
	return
}

func (s *userServiceImp) Create(ctx context.Context, info map[string]interface{}) error {
	data := models.User{}
	err := mapstructure.Decode(info, &data)
	if err != nil {
		return err
	}
	err = s.repo.Create(&data)
	s.store.Set(ctx, fmt.Sprintf("user:%d", data.UserId), data.Username, time.Second*60*60*3)
	return err
}

func (s *userServiceImp) Update(ctx context.Context, info map[string]interface{}) error {
	if info["UserId"] == nil {
		return nil
	}
	userId := info["UserId"].(int64)

	return s.repo.Update(userId, info)
}

func (s *userServiceImp) Delete(ctx context.Context, ids []int64) (err error) {
	return s.repo.Delete(ids)
}

func (s *userServiceImp) CreateOrUpdate(ctx context.Context, info map[string]interface{}) error {
	userId := info["UserId"].(int64)
	t, _ := s.repo.Get(userId)
	if t != nil {
		return s.Update(ctx, info)
	}
	return s.Create(ctx, info)
}

func (s *userServiceImp) GetOrCreate(ctx context.Context, info map[string]interface{}) error {
	userId := info["UserId"].(int64)
	t, _ := s.repo.Get(userId)
	if t != nil {
		return nil
	}
	return s.Create(ctx, info)
}
