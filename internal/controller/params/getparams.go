package params

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// GetClaimsInfo 获取claims中的userId,permission
func GetClaimsInfo(c echo.Context) (uint, int, error) {
	claims, ok := c.Get("claims").(util.JwtClaims)
	if !ok {
		return 0, 0, c.JSON(http.StatusInternalServerError, CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "GetClaimsInfo error",
			Error: "",
		})
	}
	userId := claims.UserId
	permission := claims.Permission
	return userId, permission, nil
}

// GetUserId 获取UserId,参数未提供则获取自身ID
func GetUserId(c echo.Context) (uint, error) {
	var userId uint
	rawUserId := c.QueryParam("userId")
	if rawUserId == "" {
		claims := c.Get("claims").(util.JwtClaims)
		userId = claims.UserId
	} else {
		iUserId, err := strconv.Atoi(rawUserId)
		if err != nil {
			return 0, CommonErrorGenerate(c, http.StatusBadRequest, "Invalid userId", nil)
		}
		userId = uint(iUserId)
	}
	return userId, nil
}

func GetArticleId(c echo.Context) (uint, error) {
	rawArticleId := c.QueryParam("articleId")
	if rawArticleId == "" {
		return 0, c.JSON(http.StatusBadRequest, CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "articleId missing",
			Error: "",
		})
	}
	articleId, err := strconv.Atoi(rawArticleId)
	if err != nil {
		return 0, c.JSON(http.StatusBadRequest, CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid articleId",
			Error: err.Error(),
		})
	}
	return uint(articleId), nil
}

// 定义GetNodeID错误类型,用于一些特殊判断逻辑
const (
	NodeIDIsEmpty   = 0
	NodeIDIsInvalid = 1
)

// GetNodeID 获取并转换nodeID
func GetNodeID(c echo.Context) (uint, error) {
	rawNodeID := c.QueryParam("node_id")
	if rawNodeID == "" {
		return NodeIDIsEmpty, c.JSON(http.StatusBadRequest, CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "NodeID is empty",
			Error: "",
		})
	}
	nodeID, err := strconv.Atoi(rawNodeID)
	if err != nil {
		return NodeIDIsInvalid, c.JSON(http.StatusBadRequest, CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "NodeID is invalid",
			Error: err.Error(),
		})
	}
	return uint(nodeID), nil
}

// GetPageParams 获取分页查询必需参数page,pageSize
func GetPageParams(c echo.Context) (int, int, error) {
	var page, pageSize int
	rawPage := c.QueryParam("page")
	rawPageSize := c.QueryParam("pageSize")
	if rawPage == "" {
		return 0, 0, CommonErrorGenerate(c, http.StatusBadRequest, "params missing", errors.New("page missing"))
	}
	page, err := strconv.Atoi(rawPage)
	if err != nil {
		return 0, 0, CommonErrorGenerate(c, http.StatusBadRequest, "Invalid page param", nil)
	}

	if rawPageSize == "" {
		pageSize = 30
	}
	pageSize, err = strconv.Atoi(rawPageSize)
	if err != nil {
		return 0, 0, CommonErrorGenerate(c, http.StatusBadRequest, "Invalid pageSize param", err)
	}
	return page, pageSize, nil
}

// GetCommentId 获取评论id
func GetCommentId(c echo.Context) (uint, error) {
	rawCommentId := c.QueryParam("comment_id")
	if rawCommentId == "" {
		return 0, c.JSON(http.StatusBadRequest, CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "commentID missing",
			Error: "",
		})
	}
	commentId, err := strconv.Atoi(rawCommentId)
	if err != nil {
		return 0, c.JSON(http.StatusBadRequest, CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid commentID",
			Error: "",
		})
	}
	return uint(commentId), nil
}
