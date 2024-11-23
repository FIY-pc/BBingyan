package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

func Follow(c echo.Context) error {
	rdb := util.Rdb
	ctx := context.Background()
	claims := c.Get("claims").(util.JwtClaims)
	// 关注者id
	userId := claims.UserId
	// 被关注者id
	targetId := c.QueryParam("id")
	if targetId == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "target userid missing",
			Error: "",
		})
	}

	// 将关注者id添加到被关注者id的follow set中
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

func Unfollow(c echo.Context) error {
	rdb := util.Rdb
	ctx := context.Background()
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	targetId := c.QueryParam("id")
	if targetId == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "target userid missing",
			Error: "",
		})
	}
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

func IsFollowed(c echo.Context) error {
	rdb := util.Rdb
	ctx := context.Background()
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	targetId := c.QueryParam("id")
	if targetId == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "target userid missing",
			Error: "",
		})
	}
	cmd := rdb.SIsMember(ctx, targetId, userId)
	isMember, err := cmd.Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "judge isFollowed failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Followed",
		Data: map[string]interface{}{
			"ismember": isMember,
		},
	})
}

func MyFollowerNum(c echo.Context) error {
	rdb := util.Rdb
	ctx := context.Background()
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	idStr := strconv.Itoa(int(userId))
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
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get follower num success",
		Data: map[string]interface{}{
			"num": num,
		},
	})
}

func GetFollowerNum(c echo.Context) error {
	rdb := util.Rdb
	ctx := context.Background()
	userId := c.QueryParam("id")
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
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get follower num success",
		Data: map[string]interface{}{
			"userid": userId,
			"num":    num,
		},
	})
}
