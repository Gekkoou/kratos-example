package data

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	pb "kratos-example/api/bff/v1"
	pbUser "kratos-example/api/user/v1"
	"kratos-example/app/bff/internal/biz"
	"kratos-example/app/bff/internal/pkg/jwt"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Login(ctx context.Context, req *pb.LoginRequest) (*biz.User, error) {
	res, err := r.data.uc.GetUserByName(ctx, &pbUser.GetUserByNameRequest{Name: req.Name, Password: req.Password})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:    res.Id,
		Name:  res.Name,
		Phone: res.Phone,
		Gold:  res.Gold,
	}, nil
}

func (r *userRepo) GetJwtToken(baseClaims jwt.BaseClaims) (string, int32, error) {
	j := jwt.NewJWT(r.data.conf.Jwt.Server.SigningKey)
	claims := j.CreateClaims(baseClaims)
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", 0, errors.New("获取token失败")
	}
	return token, int32(claims.RegisteredClaims.ExpiresAt.Unix()), nil
}

func (r *userRepo) Register(ctx context.Context, req *pb.CreateUserRequest) (*biz.User, error) {
	res, err := r.data.uc.CreateUser(ctx, &pbUser.CreateUserRequest{
		Name:     req.Name,
		Password: req.Password,
		Phone:    req.Phone,
	})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:    res.Id,
		Name:  res.Name,
		Phone: res.Phone,
		Gold:  res.Gold,
	}, nil
}

func (r *userRepo) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*biz.User, error) {
	res, err := r.data.uc.GetUser(ctx, &pbUser.GetUserRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:    res.Id,
		Name:  res.Name,
		Phone: res.Phone,
		Gold:  res.Gold,
	}, nil
}

func (r *userRepo) GetUserList(ctx context.Context, req *pb.GetUserListRequest) ([]*biz.User, int64, error) {
	res, err := r.data.uc.GetUserList(ctx, &pbUser.GetUserListRequest{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, 0, err
	}
	var userList []*biz.User
	for _, user := range res.List {
		userList = append(userList, &biz.User{
			Id:    user.Id,
			Name:  user.Name,
			Phone: user.Phone,
			Gold:  user.Gold,
		})
	}
	return userList, res.Total, nil
}

func (r *userRepo) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (bool, error) {
	res, err := r.data.uc.DeleteUser(ctx, &pbUser.DeleteRequest{Id: req.Id})
	if err != nil {
		return false, err
	}
	return res.Bool, nil
}

func (r *userRepo) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (bool, error) {
	res, err := r.data.uc.ChangePassword(ctx, &pbUser.ChangePasswordRequest{Id: req.Id, Password: req.Password, NewPassword: req.NewPassword})
	if err != nil {
		return false, err
	}
	return res.Bool, nil
}
