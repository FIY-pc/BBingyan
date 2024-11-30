package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/model/modelParams"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

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
	cardCmd := rdb.ZCard(ctx, modelParams.FollowKey(strUserId))
	count, err := cardCmd.Result()
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "List follower failed", err)
	}
	// 从redis获取需要查询的id列表,分页已通过redis的zset类型处理,默认按照id升序排序
	cmd := rdb.ZRange(ctx, modelParams.FollowKey(strUserId), int64((page-1)*pageSize), int64(page*pageSize))
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
