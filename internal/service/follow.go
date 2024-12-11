package service

import (
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/utils/logger"
	"github.com/redis/go-redis/v9"
	Rcontext "golang.org/x/net/context"
	"strconv"
)

func score(followID uint) float64 {
	return float64(followID)
}

func followKey(targetUID uint) string {
	return strconv.Itoa(int(targetUID)) + ":followers"
}

func Follow(targetUID uint, followerUID uint) error {
	ctx := Rcontext.Background()
	model.Rdb.ZAdd(ctx, followKey(targetUID), redis.Z{
		Score:  score(followerUID),
		Member: followerUID,
	})
	logger.Log.Info(nil, "Follow success", "followerUID", followerUID, "targetUID", targetUID)
	return nil
}

func UnFollow(targetUID uint, followerUID uint) error {
	ctx := Rcontext.Background()
	model.Rdb.ZRem(ctx, followKey(targetUID), followerUID)
	logger.Log.Info(nil, "UnFollow success", "followerUID", followerUID, "targetUID", targetUID)
	return nil
}
