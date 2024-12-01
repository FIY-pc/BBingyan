package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/controller/permission"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ArticleInfo 获取文章
func ArticleInfo(c echo.Context) error {
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	article, err := model.GetArticleByID(articleId)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "article info failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Info article success",
		Data: article,
	})
}

// ArticleCreate 创建文章
func ArticleCreate(c echo.Context) error {
	var err error
	article := model.Article{}
	// 进行信息更换
	if title := c.FormValue("title"); title != "" {
		article.Title = title
	}
	if content := c.FormValue("content"); content != "" {
		article.Content.Text = content
	}
	// 绑定作者ID
	userId, _, err := params.GetClaimsInfo(c)
	if err != nil {
		return err
	}
	article.UserID = userId

	err = model.CreateArticle(article)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "create article failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Create article success",
		Data: nil,
	})
}

// ArticleUpdate 更新文章
func ArticleUpdate(c echo.Context) error {
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	article, err := model.GetArticleByID(articleId)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "article update failed", err)
	}
	// 检查权限
	if !permission.ArticlePermissionCheck(c, articleId) {
		return params.CommonErrorGenerate(c, http.StatusUnauthorized, "permission check failed", err)
	}

	// 获取其余参数
	if title := c.FormValue("title"); title != "" {
		article.Title = title
	}
	if content := c.FormValue("content"); content != "" {
		article.Content.Text = content
	}
	// 更新文章
	err = model.UpdateArticle(*article)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "article update failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Update article success",
		Data: nil,
	})
}

// ArticleDelete 删除文章
func ArticleDelete(c echo.Context) error {
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	// 权限检查
	if !permission.ArticlePermissionCheck(c, articleId) {
		return params.CommonErrorGenerate(c, http.StatusUnauthorized, "permission check failed", err)
	}
	// 删除文章
	err = model.DeleteArticleByID(articleId)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "article delete failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Delete article success",
		Data: nil,
	})
}
