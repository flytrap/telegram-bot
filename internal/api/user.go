package api

import (
	"context"

	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/flytrap/telegram-bot/pb/v1"
	"github.com/mitchellh/mapstructure"
)

func NewUserApi(userService services.UserService) pb.UserServiceServer {
	return &UserApi{userService: userService}
}

type UserApi struct {
	userService services.UserService
	pb.UnimplementedUserServiceServer
}

func (s *UserApi) ListUser(ctx context.Context, req *pb.QueryRequest) (*pb.QueryUserResp, error) {
	results := []*pb.BotUser{}
	n, err := s.userService.List(req.Q, req.Page, req.Size, req.Order, &results)
	if err != nil {
		return &pb.QueryUserResp{Ret: &pb.RetInfo{Status: false, Msg: err.Error()}}, err
	}
	return &pb.QueryUserResp{Ret: &pb.RetInfo{Status: true}, Data: results, Total: n}, nil
}

func (s *UserApi) CreateUse(ctx context.Context, req *pb.BotUser) (*pb.RetInfo, error) {
	info := map[string]interface{}{}
	mapstructure.Decode(&req, info)
	err := s.userService.Create(info)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}

func (s *UserApi) UpdateUser(ctx context.Context, req *pb.BotUser) (*pb.RetInfo, error) {
	info := map[string]interface{}{}
	mapstructure.Decode(&req, info)
	err := s.userService.Update(info)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}

func (s *UserApi) DeleteUser(ctx context.Context, req *pb.DeleteIds) (*pb.RetInfo, error) {
	err := s.userService.Delete(req.Ids)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}
