package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	pb "kratos-example/api/bff/v1"
	"kratos-example/app/bff/internal/pkg/jwt"
)

type User struct {
	Id    uint64
	Name  string
	Phone string
	Gold  int32
}

type UserRepo interface {
	Login(context.Context, *pb.LoginRequest) (*User, error)
	Register(context.Context, *pb.CreateUserRequest) (*User, error)
	GetUserInfo(context.Context, *pb.GetUserInfoRequest) (*User, error)
	GetUserList(context.Context, *pb.GetUserListRequest) ([]*User, int64, error)
	DeleteUser(context.Context, *pb.DeleteRequest) (bool, error)
	ChangePassword(context.Context, *pb.ChangePasswordRequest) (bool, error)
	GetJwtToken(jwt.BaseClaims) (string, int32, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (user *UserUsecase) Login(ctx context.Context, req *pb.LoginRequest) (*User, error) {
	return user.repo.Login(ctx, req)
}

func (user *UserUsecase) Register(ctx context.Context, req *pb.CreateUserRequest) (*User, error) {
	return user.repo.Register(ctx, req)
}

func (user *UserUsecase) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*User, error) {
	return user.repo.GetUserInfo(ctx, req)
}
func (user *UserUsecase) GetUserList(ctx context.Context, req *pb.GetUserListRequest) ([]*User, int64, error) {
	return user.repo.GetUserList(ctx, req)
}

func (user *UserUsecase) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (bool, error) {
	return user.repo.DeleteUser(ctx, req)
}

func (user *UserUsecase) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (bool, error) {
	return user.repo.ChangePassword(ctx, req)
}

func (user *UserUsecase) GetJwtToken(baseClaims jwt.BaseClaims) (string, int32, error) {
	return user.repo.GetJwtToken(baseClaims)
}
