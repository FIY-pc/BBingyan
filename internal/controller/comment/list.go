package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CommentList 获取文章评论,进行分页处理,关键参数:page,pageSize,articleId
func CommentList(c echo.Context) error {
	// 获取参数
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	page, pageSize, err := params.GetPageParams(c)
	if err != nil {
		return err
	}
	// 获取当前页数的评论
	list, err := model.GetCommentByPage(uint(articleId), page, pageSize)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "List comments failed", err)
	}
	// 获取当前文章总评论数
	count, err := model.GetArticleCommentCount(uint(articleId))
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "get comments count failed", err)
	}
	// 返回结果
	return c.JSON(http.StatusOK, map[string]interface{}{
		"Code":  http.StatusOK,
		"Msg":   "Get comment successfully",
		"Count": count,
		"Data": map[string]interface{}{
			"articleId":   articleId,
			"commentList": list,
		},
	})
}
