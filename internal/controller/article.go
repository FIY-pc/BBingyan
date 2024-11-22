package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func ArticleInfo(c echo.Context) error {
	rawArticleId := c.QueryParam("article_id")
	ArticleId, err := strconv.Atoi(rawArticleId)
	article, err := model.GetArticleByID(uint(ArticleId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get article failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Info article success",
		Data: article,
	})
}

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
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	article.UserID = userId

	err = model.CreateArticle(article)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Create article failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Create article success",
		Data: nil,
	})
}

func ArticleUpdate(c echo.Context) error {
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Parse id failed",
			Error: err.Error(),
		})
	}
	article, err := model.GetArticleByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Parse article failed",
			Error: err.Error(),
		})
	}

	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	Permission := claims.Permission
	// 查询权限，若为管理员以下，则检查是否为文章作者
	if Permission < model.PermissionAdmin {
		article, err := model.GetArticleByID(uint(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
				Code:  http.StatusInternalServerError,
				Msg:   "Get article failed",
				Error: err.Error(),
			})
		}
		if article.UserID != userId {
			return c.JSON(http.StatusUnauthorized, params.CommonErrorResp{
				Code:  http.StatusUnauthorized,
				Msg:   "Unauthorized",
				Error: "",
			})
		}
	}

	if title := c.QueryParam("title"); title != "" {
		article.Title = title
	}
	if content := c.QueryParam("content"); content != "" {
		article.Content.Text = content
	}
	err = model.UpdateArticle(*article)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Update article failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Update article success",
		Data: nil,
	})
}

func ArticleDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Parse id failed",
			Error: err.Error(),
		})
	}
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	Permission := claims.Permission
	// 查询权限，若为管理员以下，则检查是否为文章作者
	if Permission < model.PermissionAdmin {
		article, err := model.GetArticleByID(uint(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
				Code:  http.StatusInternalServerError,
				Msg:   "Get article failed",
				Error: err.Error(),
			})
		}
		if article.UserID != userId {
			return c.JSON(http.StatusUnauthorized, params.CommonErrorResp{
				Code:  http.StatusUnauthorized,
				Msg:   "Unauthorized",
				Error: "",
			})
		}
	}

	err = model.DeleteArticleByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Delete article failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Delete article success",
		Data: nil,
	})
}
