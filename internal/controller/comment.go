package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// CommentCreate 评论
func CommentCreate(c echo.Context) error {
	// 获取并转换文章id,顺便检查id是否有效
	rawarticleId := c.QueryParam("id")
	articleId, err := strconv.Atoi(rawarticleId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "The articleID isn't a valid number",
			Error: err.Error(),
		})
	}
	// 获取claims,读取用户id
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	// 调用model
	content := c.FormValue("content")
	comment := model.Comment{
		UserID:    userId,
		ArticleID: uint(articleId),
		Content:   content,
	}
	err = model.CreateComment(comment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Create comment failed",
			Error: err.Error(),
		})
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
	rawCommentId := c.QueryParam("commentId")
	if rawCommentId == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "commentID missing",
			Error: "",
		})
	}
	commentId, err := strconv.Atoi(rawCommentId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid commentID",
			Error: "",
		})
	}
	// 获取对应评论
	comment, err := model.GetCommentByID(uint(commentId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Comment does not exist",
			Error: err.Error(),
		})
	}
	// 获取claims,读取权限与用户id
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	permission := claims.Permission
	// 用户权限仅能删自己的评论
	if permission == model.PermissionUser {
		if comment.UserID != userId {
			return c.JSON(http.StatusUnauthorized, params.CommonErrorResp{
				Code:  http.StatusUnauthorized,
				Msg:   "You are not authorized to delete this comment",
				Error: "",
			})
		}
		err = model.DeleteCommentByID(uint(commentId))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
				Code:  http.StatusInternalServerError,
				Msg:   "Delete comment failed",
				Error: err.Error(),
			})
		}
	}
	// 管理员能删所有人的评论
	if permission == model.PermissionUser {
		err = model.DeleteCommentByID(uint(commentId))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
				Code:  http.StatusInternalServerError,
				Msg:   "Delete comment failed",
				Error: err.Error(),
			})
		}
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Delete comment successfully",
		Data: nil,
	})
}

// CommentList 获取文章所有评论
func CommentList(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "Not implemented")
}
