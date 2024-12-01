package controller

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

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
