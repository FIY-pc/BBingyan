package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	nodeIDIsEmpty   = -1
	nodeIDIsInvalid = -2
)

// NodeInfo 获取节点基本信息,可以使用ID或Name进行查询
func NodeInfo(c echo.Context) error {
	var resultNode model.Node
	// 获取ID
	nodeID, err := getNodeID(c)

	// 若ID为空,获取name
	if nodeID == nodeIDIsEmpty {
		name := c.QueryParam("name")
		resultNode, err = model.GetNodeByName(name)
	} else {
		// 使用ID进行查询
		if nodeID == nodeIDIsInvalid {
			return err
		}
		resultNode, err = model.GetNodeById(uint(nodeID))
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get node failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Get node successfully",
		Data: resultNode,
	})
}

// UpdateNode 更新节点基本信息
func UpdateNode(c echo.Context) error {
	var node model.Node
	// 获取ID
	nodeID, err := getNodeID(c)
	if err != nil {
		return err
	}
	// 获取原节点信息
	node, err = model.GetNodeById(uint(nodeID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get node failed",
			Error: err.Error(),
		})
	}
	// 更新信息
	name := c.QueryParam("name")
	if name != "" {
		node.Name = name
	}
	logo := c.QueryParam("logo")
	if logo != "" {
		node.Logo = logo
	}
	// 调用model
	err = model.UpdateNode(node)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Update node failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Update node successfully",
		Data: node,
	})
}

// DeleteNode 删除节点
func DeleteNode(c echo.Context) error {
	nodeID, err := getNodeID(c)
	if err != nil {
		return err
	}
	err = model.DeleteNodeById(uint(nodeID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Delete node failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Delete node successfully",
		Data: nil,
	})
}

func CreateNode(c echo.Context) error {
	var node model.Node
	// 获取ID
	nodeID, err := getNodeID(c)
	if err != nil {
		return err
	}
	// 检查是否已经存在
	_, err = model.GetNodeById(uint(nodeID))
	if err == nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Node already exists",
			Error: "",
		})
	}
	// 获取其余信息
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Name is empty",
			Error: "",
		})
	}
	logo := c.QueryParam("logo")
	if logo != "" {
		node.Logo = logo
	}
	// 创建节点
	err = model.CreateNode(node)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Create node failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Create node successfully",
		Data: node,
	})
}

// ListArticleFromNode 分页查询同一节点下的文章,提供多种排序方式
func ListArticleFromNode(c echo.Context) error {
	nodeID, err := getNodeID(c)
	if err != nil {
		return err
	}
	var page, pageSize, sort int
	rawpage := c.QueryParam("page")
	rawpageSize := c.QueryParam("pageSize")
	rawsort := c.QueryParam("sort")
	// 获取page
	if rawpage == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "page param is empty",
			Error: "",
		})
	}
	page, err = strconv.Atoi(rawsort)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "page param is invalid",
			Error: err.Error(),
		})
	}
	// 获取pageSize
	if rawpageSize == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "page param is empty",
			Error: "",
		})
	}
	pageSize, err = strconv.Atoi(rawpageSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "page param is invalid",
			Error: err.Error(),
		})
	}
	// 获取排序方式
	if rawsort == "" {
		sort = model.SortByTitle
	} else {
		sort, err = strconv.Atoi(rawsort)
		if err != nil {
			return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
				Code:  http.StatusBadRequest,
				Msg:   "sort param is invalid",
				Error: err.Error(),
			})
		}
	}
	// 执行查询
	articleList, err := model.ListArticleFromNode(nodeID, page, pageSize, sort)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "List article failed",
			Error: err.Error(),
		})
	}
	// 查询节点文章数
	count, err := model.CountArticleFromNode(nodeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Count article failed",
			Error: err.Error(),
		})
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

// getNodeID 获取并转换nodeID
func getNodeID(c echo.Context) (uint, error) {
	rawNodeID := c.QueryParam("nodeId")
	if rawNodeID == "" {
		return nodeIDIsEmpty, c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "NodeID is empty",
			Error: "",
		})
	}
	nodeID, err := strconv.Atoi(rawNodeID)
	if err != nil {
		return nodeIDIsInvalid, c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "NodeID is invalid",
			Error: err.Error(),
		})
	}
	return uint(nodeID), nil
}
