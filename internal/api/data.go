package api

import (
	"context"
	"io"
	"sync"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/internal/serializers"
	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/flytrap/telegram-bot/pb/v1"
	"github.com/flytrap/telegram-bot/pkg/human"
	"github.com/sirupsen/logrus"
)

func NewTgBotApi(dataService services.DataService, categoryService services.CategoryService, indexService services.IndexMangerService) pb.TgBotServiceServer {
	return &TgBotService{dataService: dataService, categoryService: categoryService, indexService: indexService}
}

type TgBotService struct {
	dataService     services.DataService
	categoryService services.CategoryService
	indexService    services.IndexMangerService
	pb.UnimplementedTgBotServiceServer
}

func (s *TgBotService) ImportData(stream pb.TgBotService_ImportDataServer) error {
	var (
		wg    sync.WaitGroup
		msgCh = make(chan error)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range msgCh {
			msg := ""
			if v != nil {
				msg = v.Error()
			}
			err := stream.Send(&pb.RetInfo{Status: v == nil, Msg: msg})
			if err != nil {
				logrus.Println("Send error:", err)
				continue
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				logrus.Fatalf("recv error:%v", err)
			}
			info, err := human.Decode(req.Detail)
			if err != nil {
				logrus.Warning(err)
				continue
			}

			info["tags"] = req.Tags
			info["category"] = req.Category
			err = s.dataService.UpdateOrCreate(req.Detail.Code, info)
			if err != nil {
				logrus.Warning(err)
				msgCh <- err
				continue
			}
			if config.C.Bot.UseIndex {
				info["category"] = req.Category
				info["code"] = req.Detail.Code
				dc := serializers.DataCache{}
				err = human.Encode(info, &dc)
				if err != nil {
					logrus.Warning(err)
					msgCh <- err
					continue
				}
				data := map[string]*serializers.DataCache{req.Detail.Code: &dc}
				msgCh <- s.indexService.AddItems(context.Background(), s.indexService.IndexName(config.C.Index.Language), data)
			} else {
				msgCh <- nil
			}
		}
		close(msgCh)
	}()
	wg.Wait()
	return nil
}

func (s *TgBotService) SearchData(ctx context.Context, req *pb.DataSearchRequest) (*pb.QueryDataResp, error) {
	items := []*pb.DataItem{}
	n, err := s.dataService.List(req.Q, req.Category, req.Lang, req.Page, req.Size, req.Order, &items)
	if err != nil {
		return &pb.QueryDataResp{Ret: &pb.RetInfo{Status: false, Msg: err.Error()}}, err
	}
	results := []*pb.DataInfo{}
	for _, item := range items {
		c, _ := s.categoryService.GetName(uint(item.Category))
		i := pb.DataInfo{Detail: item, Category: c}
		results = append(results, &i)
	}
	return &pb.QueryDataResp{Ret: &pb.RetInfo{Status: true}, Data: results, Total: n}, nil
}

func (s *TgBotService) UpdateData(ctx context.Context, req *pb.DataItem) (*pb.RetInfo, error) {
	info, err := human.Decode(req)
	if err != nil {
		logrus.Warning(err)
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	err = s.dataService.Update(req.Code, info)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}

func (s *TgBotService) DeleteData(ctx context.Context, req *pb.DeleteCodes) (*pb.RetInfo, error) {
	err := s.dataService.Delete(req.Codes)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}
