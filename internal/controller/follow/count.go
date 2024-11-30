package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model/modelParams"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

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
	cmd := rdb.ZCard(ctx, modelParams.FollowKey(strconv.Itoa(int(userId))))
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
