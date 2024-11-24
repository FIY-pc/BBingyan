package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func Like(c echo.Context) error {
	// 获取文章id
	rawArticleId := c.QueryParam("articleId")
	if rawArticleId == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "articleId missing",
			Error: "",
		})
	}
	articleId, err := strconv.Atoi(rawArticleId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid articleId",
			Error: err.Error(),
		})
	}
	// 获取点赞用户Id
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId

	// 调用model
	model.Like(uint(articleId), userId)
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
	rawArticleId := c.QueryParam("articleId")
	if rawArticleId == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "articleId missing",
			Error: "",
		})
	}
	articleId, err := strconv.Atoi(rawArticleId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid articleId",
			Error: err.Error(),
		})
	}
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	model.Unlike(uint(articleId), userId)
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "OK",
		Data: map[string]interface{}{
			"articleId": articleId,
			"userId":    userId,
		},
	})
}

func GetLikeNum(c echo.Context) error {
	rawArticleId := c.QueryParam("articleId")
	if rawArticleId == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "articleId missing",
			Error: "",
		})
	}
	articleId, err := strconv.Atoi(rawArticleId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Invalid articleId",
			Error: err.Error(),
		})
	}
	count, err := model.GetLikeNum(uint(articleId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get like Num failed",
			Error: err.Error(),
		})
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
