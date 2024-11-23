package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
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
	addcmd := rdb.SAdd(ctx, params.FollowKey(targetId), userId)
	if addcmd.Err() != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Add follow failed",
			Error: addcmd.Err().Error(),
		})
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
	if !isMember {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "you had already cancel this follow",
			Error: "",
		})
	}
	// 取消关注
	remcmd := rdb.SRem(ctx, params.FollowKey(targetId), userId)
	if remcmd.Err() != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "cancel follow failed",
			Error: remcmd.Err().Error(),
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
	cmd := rdb.SIsMember(ctx, targetId, userId)
	isMember, err := cmd.Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "judge isFollowed failed",
			Error: err.Error(),
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

// MyFollowerNum 仅通过token信息获取用户自身的追随者数量
func MyFollowerNum(c echo.Context) error {
	rdb := util.Rdb
	ctx := context.Background()
	// 从claims中获取自身id
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	idStr := strconv.Itoa(int(userId))
	// 统计数量
	cmd := rdb.SCard(ctx, params.FollowKey(idStr))
	if cmd.Err() != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "judge isFollowed failed",
			Error: cmd.Err().Error(),
		})
	}
	num, err := cmd.Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "judge isFollowed failed",
			Error: err.Error(),
		})
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get follower num success",
		Data: map[string]interface{}{
			"num": num,
		},
	})
}

// GetFollowerNum 获取任何人的追随者数量,需要目标id
func GetFollowerNum(c echo.Context) error {
	rdb := util.Rdb
	ctx := context.Background()
	// 获取查询目标id
	userId := c.QueryParam("userId")
	// 查询数量
	cmd := rdb.SCard(ctx, params.FollowKey(userId))
	if cmd.Err() != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get follower num failed",
			Error: cmd.Err().Error(),
		})
	}
	num, err := cmd.Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get follower num failed",
			Error: err.Error(),
		})
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
