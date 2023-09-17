package app

import (
	"fmt"
	"net"
	"os"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/pb/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

func InitGrpcServer(ps pb.TgBotServiceServer, tgAdService pb.AdServiceServer, tgTagService pb.TagServiceServer, tgCategoryService pb.CategoryServiceServer, tgUserService pb.UserServiceServer) *GrpcServer {
	return &GrpcServer{config: config.C.ServerConfig, tgBotService: ps, tgAdService: tgAdService, tgTagService: tgTagService, tgCategoryService: tgCategoryService, tgUserService: tgUserService}
}

type GrpcServer struct {
	config            config.ServerConfig
	tgBotService      pb.TgBotServiceServer
	tgAdService       pb.AdServiceServer
	tgTagService      pb.TagServiceServer
	tgCategoryService pb.CategoryServiceServer
	tgUserService     pb.UserServiceServer
}

func (s *GrpcServer) Run() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.GrpcPort)
	logrus.Info(addr)
	listener, err := net.Listen(s.config.GrpcProtocol, addr)

	if err != nil {
		logrus.Error(err)
		return err
	}

	opts := []grpc.ServerOption{}

	if len(s.config.Cert) > 0 && len(s.config.Key) > 0 {
		c, err := credentials.NewServerTLSFromFile(s.config.Cert, s.config.Key)
		if err != nil {
			logrus.Warn(err)
		} else {
			opts = append(opts, grpc.Creds(c))
		}
	}

	grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(grpcLog)

	srv := grpc.NewServer(opts...)

	pb.RegisterTgBotServiceServer(srv, s.tgBotService)
	pb.RegisterAdServiceServer(srv, s.tgAdService)
	pb.RegisterCategoryServiceServer(srv, s.tgCategoryService)
	pb.RegisterTagServiceServer(srv, s.tgTagService)
	pb.RegisterUserServiceServer(srv, s.tgUserService)

	reflection.Register(srv)

	if err := srv.Serve(listener); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
