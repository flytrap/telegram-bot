package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/flytrap/telegram-bot/internal/serializers"
	"github.com/flytrap/telegram-bot/pkg/redis"
	"github.com/sirupsen/logrus"
)

type GroupSettingSerivce interface {
	GetSetting(ctx context.Context, code string) (*serializers.GroupSetting, error)
	SetWelcome(ctx context.Context, code string, item serializers.Welcome) error
	SetNotRobot(ctx context.Context, code string, item serializers.NotRobot) error
	GetWelcome(ctx context.Context, code string) (*serializers.Welcome, error)
	GetNotRobot(ctx context.Context, code string) (*serializers.NotRobot, error)
	Init(ctx context.Context, code string) error
}

func NewGroupSettingService(repo repositories.GroupSettingRepository, store *redis.Store) GroupSettingSerivce {
	return &groupSettingServiceImp{repo: repo, store: store}
}

type groupSettingServiceImp struct {
	repo  repositories.GroupSettingRepository
	store *redis.Store
}

func (s *groupSettingServiceImp) GetSetting(ctx context.Context, code string) (*serializers.GroupSetting, error) {
	res, err := s.getCache(ctx, code)
	if res != nil {
		return res, err
	}
	gs, err := s.repo.Get(code)
	if err != nil {
		s.Init(ctx, code)
		return nil, err
	}
	result := serializers.GroupSetting{NotRobot: serializers.NotRobot{IsOpen: gs.NotRobot, Timeout: gs.RobotTimeout},
		Welcome: serializers.Welcome{IsOpen: gs.Welcome, Desc: gs.WelcomeDesc, Pinned: gs.WelcomePinned, Template: gs.WelcomeTemplate, KillMe: gs.WelcomeKillme},
	}
	return &result, s.setCache(ctx, code, &result)
}

func (s *groupSettingServiceImp) GetWelcome(ctx context.Context, code string) (*serializers.Welcome, error) {
	gs, err := s.GetSetting(ctx, code)
	if err != nil {
		return nil, err
	}
	return &gs.Welcome, nil
}

func (s *groupSettingServiceImp) GetNotRobot(ctx context.Context, code string) (*serializers.NotRobot, error) {
	gs, err := s.GetSetting(ctx, code)
	if err != nil {
		return nil, err
	}
	return &gs.NotRobot, nil
}

func (s *groupSettingServiceImp) Init(ctx context.Context, code string) error {
	gs, _ := s.repo.Get(code)
	if gs != nil {
		result := serializers.GroupSetting{NotRobot: serializers.NotRobot{IsOpen: gs.NotRobot, Timeout: gs.RobotTimeout},
			Welcome: serializers.Welcome{IsOpen: gs.Welcome, Desc: gs.WelcomeDesc, Pinned: gs.WelcomePinned, Template: gs.WelcomeTemplate, KillMe: gs.WelcomeKillme},
		}
		return s.setCache(ctx, code, &result)
	}
	gs = &models.GroupSetting{Code: code}
	return s.repo.Create(gs)
}

func (s *groupSettingServiceImp) SetWelcome(ctx context.Context, code string, info serializers.Welcome) error {
	gs, err := s.getCache(ctx, code)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	gs.Welcome = info
	err = s.setCache(ctx, code, gs)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	data := map[string]interface{}{"welcome": info.IsOpen, "welcome_desc": info.Desc, "welcome_pinned": info.Pinned, "welcome_killme": info.KillMe, "welcome_template": info.Template}
	return s.repo.Update(code, data)
}

func (s *groupSettingServiceImp) SetNotRobot(ctx context.Context, code string, info serializers.NotRobot) error {
	gs, err := s.getCache(ctx, code)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	gs.NotRobot = info
	err = s.setCache(ctx, code, gs)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	return s.repo.Update(code, map[string]interface{}{"not_robot": info.IsOpen, "robot_timeout": info.Timeout})
}

func (s *groupSettingServiceImp) setCache(ctx context.Context, code string, info *serializers.GroupSetting) error {
	key := fmt.Sprintf("groupsetting:%s", code)
	res, err := json.Marshal(info)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	return s.store.Set(ctx, key, res, time.Hour*24)
}

func (s *groupSettingServiceImp) getCache(ctx context.Context, code string) (info *serializers.GroupSetting, err error) {
	key := fmt.Sprintf("groupsetting:%s", code)
	res := s.store.Get(ctx, key)
	if res != nil {
		err = json.Unmarshal([]byte(res.(string)), &info)
	}
	return
}

// func (s *groupSettingServiceImp) clearCache(ctx context.Context, code string) error {
// 	return s.store.Delete(ctx, fmt.Sprintf("groupsetting:%s", code))
// }
