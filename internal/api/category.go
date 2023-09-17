package api

import (
	"context"

	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/flytrap/telegram-bot/pb/v1"
	"github.com/jinzhu/copier"
)

func NewCategoryApi(categoryService services.CategoryService) pb.CategoryServiceServer {
	return &CategoryApi{categoryService: categoryService}
}

type CategoryApi struct {
	categoryService services.CategoryService
	pb.UnimplementedCategoryServiceServer
}

func (s *CategoryApi) ListCategory(ctx context.Context, req *pb.QueryRequest) (*pb.QueryTagResp, error) {
	n, data, err := s.categoryService.List(req.Q, req.Page, req.Size, req.Order)
	if err != nil {
		return &pb.QueryTagResp{Ret: &pb.RetInfo{Status: false, Msg: err.Error()}}, err
	}
	results := []*pb.Tag{}
	for _, item := range data {
		i := pb.Tag{Id: int64(item["id"].(uint)), Name: item["name"].(string), Weight: item["weight"].(int32)}
		results = append(results, &i)
	}
	copier.Copy(&results, &data)
	return &pb.QueryTagResp{Ret: &pb.RetInfo{Status: true}, Data: results, Total: n}, nil
}
func (s *CategoryApi) CreateCategory(ctx context.Context, req *pb.Tag) (*pb.RetInfo, error) {
	err := s.categoryService.Create(req.Name, req.Weight)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}
func (s *CategoryApi) UpdateCategory(ctx context.Context, req *pb.Tag) (*pb.RetInfo, error) {
	err := s.categoryService.Update(uint(req.Id), req.Name, req.Weight)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}
func (s *CategoryApi) DeleteCategory(ctx context.Context, req *pb.DeleteIds) (*pb.RetInfo, error) {
	ids := []uint{}
	copier.Copy(&ids, &req.Ids)
	err := s.categoryService.Delete(ids)
	if err != nil {
		return &pb.RetInfo{Status: false, Msg: err.Error()}, err
	}
	return &pb.RetInfo{Status: true}, nil
}
