package controller

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
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
	// 将关注者id添加到被关注者id的follow set中
	addCmd := rdb.ZAdd(ctx, params.FollowKey(targetId), redis.Z{
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
	remCmd := rdb.ZRem(ctx, params.FollowKey(targetId), userId)
	if remCmd.Err() != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "cancel follow failed", nil)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "unfollow success",
		Data: nil,
	})
}

// IsFollowed 检查用户自身是否关注了目标用户
func IsFollowed(c echo.Context) error {
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
	// 判断是否关注目标用户
	var isMember bool
	cmd := rdb.ZScore(ctx, targetId, strconv.Itoa(int(userId)))
	if errors.Is(cmd.Err(), redis.Nil) {
		isMember = false
	} else if cmd.Err() == nil {
		isMember = true
	} else {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "judge isFollowed failed", nil)
	}
	// 返回关注状态
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Followed",
		Data: map[string]interface{}{
			"ismember": isMember,
		},
	})
}

// GetFollowerNum 若提供目标用户id,则获取该id用户的追随者数量,若未提供,则根据token获取自己的追随者数量
func GetFollowerNum(c echo.Context) error {
	var num int64
	var err error
	rdb := util.Rdb
	ctx := context.Background()

	// 获取目标查询用户的Id,未提供则查自己
	userId, err := params.GetUserId(c)
	if err != nil {
		return err
	}
	// 查询
	cmd := rdb.ZCard(ctx, params.FollowKey(strconv.Itoa(int(userId))))
	if cmd.Err() != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "get follower num failed", nil)
	}
	num, err = cmd.Result()
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "get follower num failed", nil)
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get follower num success",
		Data: map[string]interface{}{
			"userid": userId,
			"num":    num,
		},
	})
}

func ListFollower(c echo.Context) error {
	rdb := util.Rdb
	ctx := context.Background()
	// 获取目标查询用户的Id,未提供则查自己
	userId, err := params.GetUserId(c)
	if err != nil {
		return err
	}
	// 获取查询页数和条数,默认30条一页
	page, pageSize, err := params.GetPageParams(c)
	if err != nil {
		return err
	}
	// 获取总follower数
	strUserId := strconv.Itoa(int(userId))
	cardCmd := rdb.ZCard(ctx, params.FollowKey(strUserId))
	count, err := cardCmd.Result()
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "List follower failed", err)
	}
	// 从redis获取需要查询的id列表,分页已通过redis的zset类型处理,默认按照id升序排序
	cmd := rdb.ZRange(ctx, params.FollowKey(strUserId), int64((page-1)*pageSize), int64(page*pageSize))
	if cmd.Err() != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "List follower failed", err)
	}
	// 获取结果并进行类型转换
	StrIds, err := cmd.Result()
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "List follower failed", err)
	}
	uintIds, err := strIdsToUint(StrIds)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "List follower failed", err)
	}
	// 从postgres查询follower详细信息
	followers, err := model.GetUsersByIDs(uintIds)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "List follower failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get followers success",
		Data: map[string]interface{}{
			"followers": followers,
			"page":      page,
			"pageSize":  pageSize,
			"count":     count,
		},
	})
}

// strIdsToUint 用于把一系列id从[]string转换到[]uint类型
func strIdsToUint(strIds []string) ([]uint, error) {
	uintIds := make([]uint, len(strIds))
	for _, strId := range strIds {
		_, err := strconv.ParseUint(strId, 10, 32)
		if err != nil {
			return uintIds, err
		}
	}
	return uintIds, nil
}
