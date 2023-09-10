package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	pb "kratos-example/api/user/v1"
	"kratos-example/app/user/internal/biz"
	"kratos-example/pkg/utils"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data  *Data
	log   *log.Helper
	group singleflight.Group
}

var (
	userCacheKey = func(key string) string {
		return "userInfo:" + key
	}
	userCacheTTL = 600
)

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) GetUser(ctx context.Context, req *pb.GetUserRequest) (*biz.User, error) {
	key := strconv.FormatUint(req.Id, 10)
	// 防缓存击穿
	user, err, _ := r.group.Do(key, func() (interface{}, error) {
		return r.DoGetUser(req.Id)
	})
	return user.(*biz.User), err
}

func (r *userRepo) DoGetUser(uid uint64) (user *biz.User, err error) {
	key := userCacheKey(fmt.Sprintf("%d", uid))
	json, err := r.data.redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		if err = r.data.db.Where("id = ?", uid).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 防缓存穿透
				r.data.redis.Set(context.Background(), key, "NA", time.Duration(userCacheTTL)*time.Second).Result()
				return user, nil
			}
			return
		}
		r.SetUserCache(user)
	} else if err == nil {
		if json == "NA" {
			return &biz.User{}, nil
		}
		err = sonic.UnmarshalString(json, &user)
	}
	return
}

func (r *userRepo) SetUserCache(user *biz.User) {
	key := userCacheKey(fmt.Sprintf("%d", user.Id))
	json, err := sonic.MarshalString(user)
	if err == nil {
		r.data.redis.Set(context.Background(), key, json, time.Duration(userCacheTTL)*time.Second).Result()
	}
}

func (r *userRepo) DelUserCache(uid uint64) {
	key := userCacheKey(fmt.Sprintf("%d", uid))
	r.data.redis.Del(context.Background(), key).Result()
}

func (r *userRepo) GetUserList(ctx context.Context, req *pb.GetUserListRequest) ([]*biz.User, int64, error) {
	var userList []*biz.User
	var total int64
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := r.data.db.Model(&biz.User{})
	err := db.Count(&total).Error
	if err != nil {
		return []*biz.User{}, 0, err
	}
	err = db.Limit(int(limit)).Offset(int(offset)).Find(&userList).Error
	return userList, total, err
}

func (r *userRepo) GetUserByName(ctx context.Context, req *pb.GetUserByNameRequest) (*biz.User, error) {
	var u biz.User
	if errors.Is(r.data.db.Where("name = ?", req.Name).First(&u).Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("登陆失败")
	}
	if ok := utils.BcryptCheck(req.Password, u.Password); !ok {
		return nil, errors.New("密码错误")
	}
	return &u, nil
}

func (r *userRepo) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*biz.User, error) {
	var u biz.User
	if !errors.Is(r.data.db.Where("name = ?", req.Name).First(&u).Error, gorm.ErrRecordNotFound) {
		return &u, errors.New("用户名已存在")
	}
	var user = biz.User{
		Name:     req.Name,
		Password: utils.BcryptHash(req.Password),
		Phone:    req.Phone,
	}
	err := r.data.db.Create(&user).Error
	r.DelUserCache(user.Id)
	return &user, err
}

func (r *userRepo) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (bool, error) {
	var u biz.User
	err := r.data.db.Where("id = ?", req.Id).Delete(&u).Error
	if err != nil {
		return false, err
	}
	r.DelUserCache(req.Id)
	return true, nil
}

func (r *userRepo) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (bool, error) {
	var u biz.User
	if err := r.data.db.Where("id = ?", req.Id).First(&u).Error; err != nil {
		return false, errors.New("用户不存在")
	}
	if ok := utils.BcryptCheck(req.Password, u.Password); !ok {
		return false, errors.New("原密码错误")
	}
	u.Password = utils.BcryptHash(req.NewPassword)
	if err := r.data.db.Save(&u).Error; err != nil {
		return false, err
	}
	r.DelUserCache(req.Id)
	return true, nil
}
