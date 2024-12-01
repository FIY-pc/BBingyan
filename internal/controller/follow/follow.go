package controller

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model/modelParams"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

// score 计算关注列表中用户的分数,暂时以id为准,后续可以根据用户活动、用户等级等设计更复杂的算法
func score(userId uint) float64 {
	return float64(userId)
}

// Follow 关注
func Follow(c echo.Context) error {
	rdb := util.Rdb
	ctx := context.Background()
	claims := c.Get("claims").(util.JwtClaims)
	// 关注者id
	userId := claims.UserId
	// 被关注者id
	targetId := c.QueryParam("userId")
	if targetId == "" {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "params missing", errors.New("target userId missing"))
	}

	// 检查是否已关注
	cmd := rdb.SIsMember(ctx, targetId, userId)
	isMember, err := cmd.Result()
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "judge isFollowed failed", err)
	}
	if isMember {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "you had followed this user already", nil)
	}
	// 将关注者id添加到被关注者id的follow sort set中
	addCmd := rdb.ZAdd(ctx, modelParams.FollowKey(targetId), redis.Z{
		Member: userId,
		Score:  score(userId),
	})

	if addCmd.Err() != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "follow failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "follow success",
		Data: nil,
	})
}

// Unfollow 取消关注
func Unfollow(c echo.Context) error {
	rdb := util.Rdb
	ctx := context.Background()
	// 从claims中获取关注者id
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	// 被关注者id
	targetId := c.QueryParam("userId")
	if targetId == "" {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "params missing", errors.New("target userId missing"))
	}
	// 检查是否已关注
	cmd := rdb.ZScore(ctx, targetId, strconv.Itoa(int(userId)))
	if !errors.Is(cmd.Err(), redis.Nil) {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "already cancel", nil)
	}
	// 取消关注
	remCmd := rdb.ZRem(ctx, modelParams.FollowKey(targetId), userId)
	if remCmd.Err() != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "cancel follow failed", nil)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "unfollow success",
		Data: nil,
	})
}
