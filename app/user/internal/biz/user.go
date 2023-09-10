package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"

	"github.com/go-kratos/kratos/v2/log"
	pb "kratos-example/api/user/v1"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(pb.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type User struct {
	Id        uint64                `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Name      string                `gorm:"column:name;type:varchar(50);comment:用户名;NOT NULL" json:"name"`
	Password  string                `gorm:"column:password;type:varchar(100);comment:密码;NOT NULL" json:"password"`
	Phone     string                `gorm:"column:phone;type:varchar(50);comment:手机号;NOT NULL" json:"phone"`
	Gold      int32                 `gorm:"column:gold;type:int(11);comment:金币;NOT NULL;default:0" json:"gold"`
	CreatedAt int32                 `gorm:"column:created_at;type:int(11) unsigned;comment:创建时间;NOT NULL;autoCreateTime" json:"created_at"`
	UpdatedAt int32                 `gorm:"column:updated_at;type:int(11) unsigned;comment:更新时间;NOT NULL;autoUpdateTime" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;type:int(11) unsigned;comment:删除时间;NOT NULL;default:0" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	fmt.Println("AfterUpdate User", u)
	return
}

type UserRepo interface {
	GetUser(context.Context, *pb.GetUserRequest) (*User, error)
	GetUserList(context.Context, *pb.GetUserListRequest) ([]*User, int64, error)
	GetUserByName(context.Context, *pb.GetUserByNameRequest) (*User, error)
	CreateUser(context.Context, *pb.CreateUserRequest) (*User, error)
	DeleteUser(context.Context, *pb.DeleteRequest) (bool, error)
	ChangePassword(context.Context, *pb.ChangePasswordRequest) (bool, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (user *UserUsecase) GetUser(ctx context.Context, req *pb.GetUserRequest) (*User, error) {
	return user.repo.GetUser(ctx, req)
}

func (user *UserUsecase) GetUserList(ctx context.Context, req *pb.GetUserListRequest) ([]*User, int64, error) {
	return user.repo.GetUserList(ctx, req)
}

func (user *UserUsecase) GetUserByName(ctx context.Context, req *pb.GetUserByNameRequest) (*User, error) {
	return user.repo.GetUserByName(ctx, req)
}

func (user *UserUsecase) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*User, error) {
	return user.repo.CreateUser(ctx, req)
}

func (user *UserUsecase) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (bool, error) {
	return user.repo.DeleteUser(ctx, req)
}

func (user *UserUsecase) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (bool, error) {
	return user.repo.ChangePassword(ctx, req)
}
