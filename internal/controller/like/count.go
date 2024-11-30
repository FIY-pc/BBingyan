package like

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetLikeNum 获取文章赞数
func GetLikeNum(c echo.Context) error {
	// 获取文章id
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	count, err := model.GetLikeNum(articleId)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "get like num failed", err)
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
