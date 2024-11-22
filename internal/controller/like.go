package controller

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func like(articleID uint, UserID uint) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	rdb.SAdd(ctx, params.ArticleLikeKey(articleID), UserID)
}

func Unlike(articleID uint, UserID uint) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	rdb.SRem(ctx, params.ArticleLikeKey(articleID))
}

func GetLikeNum(articleID uint) (int64, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	num, err := rdb.SCard(ctx, params.ArticleLikeKey(articleID)).Result()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func isUserLikeArticle(articleID uint, userID uint) bool {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	cmd := rdb.SIsMember(ctx, params.ArticleLikeKey(articleID), userID)
	isMember, err := cmd.Result()
	if err != nil {
		return false
	}
	return isMember
}
