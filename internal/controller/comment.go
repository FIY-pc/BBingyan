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
	rawarticleId := c.QueryParam("article_id")
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
	// 获取内容
	content := c.FormValue("content")
	// 调用model
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
	rawCommentId := c.QueryParam("comment_id")
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
	if permission == util.PermissionUser {
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
	if permission == util.PermissionUser {
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

// CommentGetById 使用id查询单条评论
func CommentGetById(c echo.Context) error {
	rawCommentId := c.QueryParam("comment_id")
	commentId, err := strconv.Atoi(rawCommentId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid commentID",
			Error: "",
		})
	}
	comment, err := model.GetCommentByID(uint(commentId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Comment does not exist",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get comment successfully",
		Data: comment,
	})
}

// CommentList 获取文章评论,进行分页处理,关键参数:page,pageSize,articleId
func CommentList(c echo.Context) error {
	// 非空检查
	rawArticleId := c.QueryParam("article_id")
	rawPage := c.QueryParam("page")
	rawPageSize := c.QueryParam("pageSize")
	if rawPage != "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "The page param missing",
			Error: "",
		})
	}
	if rawPageSize != "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "The pageSize param missing",
			Error: "",
		})
	}
	if rawArticleId != "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "The articleId param missing",
			Error: "",
		})
	}
	// 转换类型
	articleId, err := strconv.Atoi(rawArticleId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid articleId",
			Error: "",
		})
	}
	page, err := strconv.Atoi(rawPage)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid page",
			Error: "",
		})
	}
	pageSize, err := strconv.Atoi(rawPageSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid pageSize",
			Error: "",
		})
	}
	// 获取当前页数的评论
	list, err := model.GetCommentByPage(uint(articleId), page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get comment failed",
			Error: err.Error(),
		})
	}
	// 获取当前文章总评论数
	count, err := model.GetArticleCommentCount(uint(articleId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get comment count failed",
			Error: err.Error(),
		})
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

func GetArticleCommentCount(c echo.Context) error {
	rawArticleId := c.QueryParam("article_id")
	articleId, err := strconv.Atoi(rawArticleId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid articleId",
			Error: "",
		})
	}
	count, err := model.GetArticleCommentCount(uint(articleId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get comment count failed",
			Error: err.Error(),
		})
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
	rawUserId := c.QueryParam("user_id")
	// 若提供userId,则查询该用户,否则查询自身评论数
	if rawUserId == "" {
		claims := c.Get("claims").(util.JwtClaims)
		userId := claims.UserId
		commentNum, err = model.GetUserCommentCount(userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
				Code:  http.StatusInternalServerError,
				Msg:   "Get comment count failed",
				Error: err.Error(),
			})
		}
	} else {
		var userId int
		userId, err = strconv.Atoi(rawUserId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
				Code:  http.StatusBadRequest,
				Msg:   "Invalid userId",
				Error: "",
			})
		}
		commentNum, err = model.GetUserCommentCount(uint(userId))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
				Code:  http.StatusInternalServerError,
				Msg:   "Get comment count failed",
				Error: err.Error(),
			})
		}
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
