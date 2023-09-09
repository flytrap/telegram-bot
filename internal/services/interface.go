package services

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/flytrap/telegram-bot/pb/v1"
)

func NewTgBotService(groupService DataService) pb.TgBotServiceServer {
	return &TgBotService{groupService: groupService}
}

type TgBotService struct {
	groupService DataService
	pb.UnimplementedTgBotServiceServer
}

func (s *TgBotService) ImportGroup(stream pb.TgBotService_ImportDataServer) error {
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
			err := stream.Send(&pb.ImportResponse{Status: v == nil, Msg: msg})
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
			msgCh <- s.groupService.UpdateOrCreate(req.Code, req.Tid, req.Name, req.Desc, uint32(req.Number), req.Tags, req.Category)
		}
		close(msgCh)
	}()
	wg.Wait()
	return nil
}
