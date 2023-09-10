package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "kratos-example/api/user/v1"
	"kratos-example/app/user/internal/biz"
)

type UserService struct {
	pb.UnimplementedUserServer

	user *biz.UserUsecase
	log  *log.Helper
}

func NewUserService(user *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{user: user, log: log.NewHelper(logger)}
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	res, err := s.user.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserReply{
		Id:        res.Id,
		Name:      res.Name,
		Phone:     res.Phone,
		Gold:      res.Gold,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (s *UserService) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListReply, error) {
	res, total, err := s.user.GetUserList(ctx, req)
	if err != nil {
		return nil, err
	}
	var userList []*pb.GetUserReply
	for _, user := range res {
		userList = append(userList, &pb.GetUserReply{
			Id:        user.Id,
			Name:      user.Name,
			Phone:     user.Phone,
			Gold:      user.Gold,
			CreatedAt: user.CreatedAt,
		})
	}
	return &pb.GetUserListReply{
		List:  userList,
		Total: total,
	}, nil
}

func (s *UserService) GetUserByName(ctx context.Context, req *pb.GetUserByNameRequest) (*pb.GetUserReply, error) {
	res, err := s.user.GetUserByName(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserReply{
		Id:        res.Id,
		Name:      res.Name,
		Phone:     res.Phone,
		Gold:      res.Gold,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.GetUserReply, error) {
	res, err := s.user.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserReply{
		Id:        res.Id,
		Name:      res.Name,
		Phone:     res.Phone,
		Gold:      res.Gold,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteUserReply, error) {
	res, err := s.user.DeleteUser(ctx, req)
	return &pb.DeleteUserReply{Bool: res}, err
}

func (s *UserService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordReply, error) {
	res, err := s.user.ChangePassword(ctx, req)
	return &pb.ChangePasswordReply{Bool: res}, err
}
