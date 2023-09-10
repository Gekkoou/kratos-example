package service

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-example/app/bff/internal/biz"
	"kratos-example/app/bff/internal/pkg/jwt"

	pb "kratos-example/api/bff/v1"
)

type UserService struct {
	pb.UnimplementedUserServer

	user *biz.UserUsecase
	log  *log.Helper
}

func NewUserService(user *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{user: user, log: log.NewHelper(logger)}
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	res, err := s.user.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	token, expiresAt, err := s.user.GetJwtToken(jwt.BaseClaims{
		Id:   res.Id,
		Name: res.Name,
	})
	if err != nil {
		return nil, errors.New("获取token失败")
	}
	return &pb.LoginReply{
		Code:    0,
		Message: "",
		Data: &pb.LoginReply_LoginReplyData{
			UserInfo: &pb.UserInfoData{
				Id:    res.Id,
				Name:  res.Name,
				Phone: res.Phone,
				Gold:  res.Gold,
			},
			Token:     token,
			ExpiresAt: expiresAt,
		},
	}, nil
}

func (s *UserService) Register(ctx context.Context, req *pb.CreateUserRequest) (*pb.GetUserInfoReply, error) {
	res, err := s.user.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserInfoReply{
		Code:    0,
		Message: "",
		Data: &pb.GetUserInfoReply_GetUserInfoReplyData{
			UserInfo: &pb.UserInfoData{
				Id:    res.Id,
				Name:  res.Name,
				Phone: res.Phone,
				Gold:  res.Gold,
			},
		},
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	res, err := s.user.GetUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserInfoReply{
		Code:    0,
		Message: "",
		Data: &pb.GetUserInfoReply_GetUserInfoReplyData{
			UserInfo: &pb.UserInfoData{
				Id:    res.Id,
				Name:  res.Name,
				Phone: res.Phone,
				Gold:  res.Gold,
			},
		},
	}, nil
}

func (s *UserService) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListReply, error) {
	res, total, err := s.user.GetUserList(ctx, req)
	if err != nil {
		return nil, err
	}
	var userList []*pb.UserInfoData
	for _, user := range res {
		userList = append(userList, &pb.UserInfoData{
			Id:    user.Id,
			Name:  user.Name,
			Phone: user.Phone,
			Gold:  user.Gold,
		})
	}
	return &pb.GetUserListReply{
		Code:    0,
		Message: "",
		Data: &pb.GetUserListReply_GetUserListReplyData{
			List:  userList,
			Total: total,
		},
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (*pb.Response, error) {
	if jwt.GetUserID(ctx) != req.Id {
		return nil, errors.New("操作失败")
	}
	ok, err := s.user.DeleteUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("注销失败")
	}
	return &pb.Response{
		Code:    0,
		Message: "注销成功",
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.Response, error) {
	if jwt.GetUserID(ctx) != req.Id {
		return nil, errors.New("操作失败")
	}
	ok, err := s.user.ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("修改失败")
	}
	return &pb.Response{
		Code:    0,
		Message: "修改成功",
	}, nil
}
