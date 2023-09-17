package api

import (
	"context"

	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/flytrap/telegram-bot/pb/v1"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
)

func NewAdApi(adService services.AdService) pb.AdServiceServer {
	return &AdApi{adService: adService}
}

type AdApi struct {
	adService services.AdService
	pb.UnimplementedAdServiceServer
}

func (s *AdApi) ListAd(ctx context.Context, req *pb.AdFilter) (*pb.QueryAdResp, error) {
	n, data, err := s.adService.List(req.Q, req.Page, req.Size, req.Order)
	if err != nil {
		return &pb.QueryAdResp{Ret: &pb.RetInfo{Status: false, Msg: err.Error()}}, err
	}
	results := []*pb.Ad{}
	for _, item := range data {
		i := pb.Ad{Keyword: item["keyword"].(string), Type: int32(item["type"].(int8)), Code: item["code"].(string), Desc: item["desc"].(string), Number: item["num"].(uint32),
			IsShowAd: item["is_show_ad"].(bool), CategoryId: int64(item["category_id"].(uint)), Category: item["category"].(string), Name: item["name"].(string),
			Global: int32(item["global"].(int8)), Language: item["language"].(string), Id: int64(item["id"].(uint))}
		results = append(results, &i)
	}
	return &pb.QueryAdResp{Ret: &pb.RetInfo{Status: true}, Data: results, Total: n}, nil
}
func (s *AdApi) CreateAd(ctx context.Context, req *pb.Ad) (*pb.RetInfo, error) {
	info := map[string]interface{}{}
	mapstructure.Decode(&req, info)
	err := s.adService.Create(info)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}
func (s *AdApi) UpdateAd(ctx context.Context, req *pb.Ad) (*pb.RetInfo, error) {
	info := map[string]interface{}{}
	mapstructure.Decode(&req, info)
	err := s.adService.Update(uint(req.Id), info)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}
func (s *AdApi) DeleteAd(ctx context.Context, req *pb.DeleteIds) (*pb.RetInfo, error) {
	ids := []uint{}
	copier.Copy(&ids, &req.Ids)
	err := s.adService.Delete(ids)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}
