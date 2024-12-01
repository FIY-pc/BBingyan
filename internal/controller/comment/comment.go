package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CommentCreate 评论
func CommentCreate(c echo.Context) error {
	articleId, err := params.GetArticleId(c)
	if err != nil {
		return err
	}
	// 获取claims,读取用户id
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	// 获取内容
	content := c.FormValue("content")
	// 调用model
	comment := model.Comment{
		UserID:    userId,
		ArticleID: articleId,
		Content:   content,
	}
	err = model.CreateComment(comment)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Comment create failed", err)
	}
	// 成功创建
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Create comment successfully",
		Data: nil,
	})
}

// CommentDelete 删评,做了权限区分,user仅能删自己的评论,admin可删所有评论
func CommentDelete(c echo.Context) error {
	// 获取评论id
	commentId, err := params.GetCommentId(c)
	// 获取对应评论
	comment, err := model.GetCommentByID(uint(commentId))
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "delete comment failed", err)
	}
	// 获取claims,读取权限与用户id
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	permission := claims.Permission
	// 用户权限仅能删自己的评论
	if permission == util.PermissionUser {
		if comment.UserID != userId {
			return params.CommonErrorGenerate(c, http.StatusBadRequest, "Permission not allowed", err)
		}
		err = model.DeleteCommentByID(uint(commentId))
		if err != nil {
			return params.CommonErrorGenerate(c, http.StatusInternalServerError, "delete comment failed", err)
		}
	}
	// 管理员能删所有人的评论
	if permission == util.PermissionUser {
		err = model.DeleteCommentByID(uint(commentId))
		if err != nil {
			return params.CommonErrorGenerate(c, http.StatusInternalServerError, "delete comment failed", err)
		}
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Delete comment successfully",
		Data: nil,
	})
}

// CommentGetById 使用id查询单条评论
func CommentGetById(c echo.Context) error {
	commentId, err := params.GetCommentId(c)
	comment, err := model.GetCommentByID(uint(commentId))
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "get comment failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get comment successfully",
		Data: comment,
	})
}
