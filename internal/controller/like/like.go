package like

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Like(c echo.Context) error {
	// 获取文章id
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	// 获取点赞对象ID
	userId, err := params.GetUserId(c)
	if err != nil {
		return err
	}
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
	// 获取点赞对象ID
	userId, err := params.GetUserId(c)
	if err != nil {
		return err
	}
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
