package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetArticleCommentCount(c echo.Context) error {
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	count, err := model.GetArticleCommentCount(uint(articleId))
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Get article comment count failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get comment count successfully",
		Data: map[string]interface{}{
			"articleId": articleId,
			"Count":     count,
		},
	})
}

func GetUserCommentCount(c echo.Context) error {
	var commentNum int64
	var err error
	userId, err := params.GetUserId(c)

	commentNum, err = model.GetUserCommentCount(userId)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Get user comment count failed", err)
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get comment count successfully",
		Data: map[string]interface{}{
			"commentNum": commentNum,
		},
	})
}
