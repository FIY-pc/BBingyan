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
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "target userid missing",
			Error: "",
		})
	}

	// 检查是否已关注
	cmd := rdb.SIsMember(ctx, targetId, userId)
	isMember, err := cmd.Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "judge isFollowed failed",
			Error: err.Error(),
		})
	}
	if isMember {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "you had followed this user already",
			Error: "",
		})
	}
	// 将关注者id添加到被关注者id的follow set中
	addCmd := rdb.ZAdd(ctx, params.FollowKey(targetId), redis.Z{
		Member: userId,
		Score:  score(userId),
	})

	if addCmd.Err() != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Add follow failed",
			Error: addCmd.Err().Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "follow success",
		Data: nil,
	})
}

// score 计算关注列表中用户的分数,暂时以id为准,后续可以根据用户活动、用户等级等设计更复杂的算法
func score(userId uint) float64 {
	return float64(userId)
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
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "target userid missing",
			Error: "",
		})
	}
	// 检查是否已关注
	cmd := rdb.ZScore(ctx, targetId, strconv.Itoa(int(userId)))
	if !errors.Is(cmd.Err(), redis.Nil) {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "you had already cancel this follow",
			Error: "",
		})
	}
	// 取消关注
	remCmd := rdb.ZRem(ctx, params.FollowKey(targetId), userId)
	if remCmd.Err() != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "cancel follow failed",
			Error: remCmd.Err().Error(),
		})
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
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "target userid missing",
			Error: "",
		})
	}
	// 判断是否关注目标用户
	var isMember bool
	cmd := rdb.ZScore(ctx, targetId, strconv.Itoa(int(userId)))
	if errors.Is(cmd.Err(), redis.Nil) {
		isMember = false
	} else if cmd.Err() == nil {
		isMember = true
	} else {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "judge IsFollowed failed",
			Error: cmd.Err().Error(),
		})
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
	var resultId string
	rdb := util.Rdb
	ctx := context.Background()
	// 获取查询目标id
	if userId := c.QueryParam("userId"); userId != "" {
		cmd := rdb.ZCard(ctx, params.FollowKey(userId))
		if cmd.Err() != nil {
			return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
				Code:  http.StatusInternalServerError,
				Msg:   "Get follower num failed",
				Error: cmd.Err().Error(),
			})
		}
		num, err = cmd.Result()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
				Code:  http.StatusInternalServerError,
				Msg:   "Get follower num failed",
				Error: err.Error(),
			})
		}
		resultId = userId
	} else {
		// 从claims获取用户id
		claims := c.Get("claims").(util.JwtClaims)
		rawMyId := claims.UserId
		myId := strconv.Itoa(int(rawMyId))
		cmd := rdb.ZCard(ctx, params.FollowKey(myId))
		if cmd.Err() != nil {
			return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
				Code:  http.StatusInternalServerError,
				Msg:   "Get follower num failed",
				Error: cmd.Err().Error(),
			})
		}
		num, err = cmd.Result()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
				Code:  http.StatusInternalServerError,
				Msg:   "Get follower num failed",
				Error: err.Error(),
			})
		}
		resultId = myId
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get follower num success",
		Data: map[string]interface{}{
			"userid": resultId,
			"num":    num,
		},
	})
}

func ListFollower(c echo.Context) error {
	var err error
	var userId uint
	rdb := util.Rdb
	ctx := context.Background()

	// 获取目标查询用户的Id,未提供则查自己
	rawUserId := c.QueryParam("userId")
	if rawUserId == "" {
		claims := c.Get("claims").(util.JwtClaims)
		userId = claims.UserId
	} else {
		iUserId, err := strconv.Atoi(rawUserId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
				Code:  http.StatusBadRequest,
				Msg:   "Invalid user id",
				Error: err.Error(),
			})
		}
		userId = uint(iUserId)
	}
	// 获取查询页数和条数,默认30条一页
	var page, pageSize int
	rawPage := c.QueryParam("page")
	rawPageSize := c.QueryParam("pageSize")
	if rawPage == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "page param missing",
			Error: rawPage,
		})
	}
	page, err = strconv.Atoi(rawPage)
	if rawPageSize == "" {
		pageSize = 30
	}
	pageSize, err = strconv.Atoi(rawPageSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid page size",
			Error: err.Error(),
		})
	}
	strUserId := strconv.Itoa(int(userId))
	// 获取总follower数
	cardCmd := rdb.ZCard(ctx, params.FollowKey(strUserId))
	count, err := cardCmd.Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get follower num failed",
			Error: err.Error(),
		})
	}
	// 从redis获取需要查询的id列表,分页已通过redis的zset类型处理,默认按照id升序排序
	cmd := rdb.ZRange(ctx, params.FollowKey(strUserId), int64((page-1)*pageSize), int64(page*pageSize))
	if cmd.Err() != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get follower ids failed",
			Error: cmd.Err().Error(),
		})
	}
	StrIds, err := cmd.Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get follower ids failed",
			Error: err.Error(),
		})
	}
	uintIds, err := strIdsToUint(StrIds)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "strIdsToUint failed",
			Error: err.Error(),
		})
	}
	// 从postgres查询follower详细信息
	followers, err := model.GetUsersByIDs(uintIds)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code: http.StatusInternalServerError,
			Msg:  "Get followers failed",
		})
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
