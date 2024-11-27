package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Like(c echo.Context) error {
	// 获取文章id
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	// 获取点赞用户Id
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId

	// 调用model
	model.Like(articleId, userId)
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "OK",
		Data: map[string]interface{}{
			"articleId": articleId,
			"userId":    userId,
		},
	})
}

func Unlike(c echo.Context) error {
	// 获取文章id
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	// 获取取消点赞对象ID
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	// 执行操作
	model.Unlike(articleId, userId)
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "OK",
		Data: map[string]interface{}{
			"articleId": articleId,
			"userId":    userId,
		},
	})
}

func GetLikeNum(c echo.Context) error {
	// 获取文章id
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	count, err := model.GetLikeNum(articleId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get like Num failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "OK",
		Data: map[string]interface{}{
			"articleId": articleId,
			"count":     count,
		},
	})
}
