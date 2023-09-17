package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/flytrap/telegram-bot/pb/v1"
	"github.com/jinzhu/copier"
)

func NewTgBotApi(dataService services.DataService) pb.TgBotServiceServer {
	return &TgBotService{dataService: dataService}
}

type TgBotService struct {
	dataService services.DataService
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
				fmt.Println("Send error:", err)
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
				log.Fatalf("recv error:%v", err)
			}
			fmt.Printf("Recved :%v \n", req.Name)
			msgCh <- s.dataService.UpdateOrCreate(req.Code, req.Tid, req.Name, req.Desc, uint32(req.Number), req.Tags, req.Category, req.Lang)
		}
		close(msgCh)
	}()
	wg.Wait()
	return nil
}

func (s *TgBotService) SearchData(ctx context.Context, req *pb.DataSearchRequest) (*pb.QueryDataResp, error) {
	n, data, err := s.dataService.List(req.Q, req.Category, req.Lang, req.Page, req.Size, req.Order)
	if err != nil {
		return &pb.QueryDataResp{Ret: &pb.RetInfo{Status: false, Msg: err.Error()}}, err
	}
	results := []*pb.DataItem{}
	copier.Copy(&results, &data)
	return &pb.QueryDataResp{Ret: &pb.RetInfo{Status: true}, Data: results, Total: n}, nil
}

func (s *TgBotService) UpdateData(ctx context.Context, req *pb.DataItem) (*pb.RetInfo, error) {
	err := s.dataService.Update(req.Code, req.Tid, req.Name, req.Desc, req.Number, -1)
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