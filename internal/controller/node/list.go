package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// ListArticleFromNode 分页查询同一节点下的文章,提供多种排序方式
func ListArticleFromNode(c echo.Context) error {
	nodeID, err := params.GetNodeID(c)
	if err != nil {
		return err
	}
	var page, pageSize, sort int
	rawPage := c.QueryParam("page")
	rawPageSize := c.QueryParam("pageSize")
	rawSort := c.QueryParam("sort")
	// 获取page
	if rawPage == "" {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "page can't be empty", err)
	}
	page, err = strconv.Atoi(rawSort)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "Invalid page params", err)
	}
	// 获取pageSize
	if rawPageSize == "" {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "page can't be empty", err)
	}
	pageSize, err = strconv.Atoi(rawPageSize)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "Invalid pageSize params", err)
	}
	// 获取排序方式
	if rawSort == "" {
		// 默认按时间排序
		sort = model.SortByTime
	} else {
		sort, err = strconv.Atoi(rawSort)
		if err != nil {
			return params.CommonErrorGenerate(c, http.StatusBadRequest, "Invalid sort params", err)
		}
	}
	// 执行查询
	articleList, err := model.ListArticleFromNode(nodeID, page, pageSize, sort)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "List article failed", err)
	}
	// 查询节点文章数
	count, err := model.CountArticleFromNode(nodeID)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Count article failed", err)
	}
	// 返回结果
	return c.JSON(http.StatusOK, map[string]interface{}{
		"Code":  http.StatusOK,
		"Msg":   "List article successfully",
		"Count": count,
		"Data": map[string]interface{}{
			"nodeID":      nodeID,
			"articleList": articleList,
		},
	})
}
