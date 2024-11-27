package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	nodeIDIsEmpty   = 0
	nodeIDIsInvalid = 1
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
		resultNode, err = model.GetNodeById(nodeID)
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
	node, err = model.GetNodeById(nodeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get node failed",
			Error: err.Error(),
		})
	}
	// 更新信息
	name := c.FormValue("name")
	if name != "" {
		node.Name = name
	}
	logo := c.FormValue("logo")
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
	err = model.DeleteNodeById(nodeID)
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
	_, err = model.GetNodeById(nodeID)
	if err == nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Node already exists",
			Error: "",
		})
	}
	// 获取其余信息
	name := c.FormValue("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Name is empty",
			Error: "",
		})
	}
	logo := c.FormValue("logo")
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
	// 将创建者设置为默认管理员
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	err = model.AddNodeAdmin(nodeID, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Add node admin failed",
			Error: err.Error(),
		})
	}
	// 返回结果
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
	rawPage := c.QueryParam("page")
	rawPageSize := c.QueryParam("pageSize")
	rawSort := c.QueryParam("sort")
	// 获取page
	if rawPage == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "page param is empty",
			Error: "",
		})
	}
	page, err = strconv.Atoi(rawSort)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "page param is invalid",
			Error: err.Error(),
		})
	}
	// 获取pageSize
	if rawPageSize == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "page param is empty",
			Error: "",
		})
	}
	pageSize, err = strconv.Atoi(rawPageSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "page param is invalid",
			Error: err.Error(),
		})
	}
	// 获取排序方式
	if rawSort == "" {
		// 默认按时间排序
		sort = model.SortByTime
	} else {
		sort, err = strconv.Atoi(rawSort)
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

// AddNodeAdmin 添加节点管理员
func AddNodeAdmin(c echo.Context) error {
	nodeId, err := getNodeID(c)
	if err != nil {
		return err
	}

	rawAddAdmin := c.FormValue("user_id")
	addAdmin, err := strconv.Atoi(rawAddAdmin)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "add admin param is invalid",
			Error: err.Error(),
		})
	}

	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	permission := claims.Permission
	// 权限认证,需节点管理员以上
	if permission < util.PermissionAdmin {
		if !model.IsNodeAdmin(nodeId, userId) {
			return c.JSON(http.StatusUnauthorized, params.CommonErrorResp{
				Code:  http.StatusUnauthorized,
				Msg:   "permission denied",
				Error: "",
			})
		}
	}
	// 调用model
	err = model.AddNodeAdmin(nodeId, uint(addAdmin))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Add node admin failed",
			Error: err.Error(),
		})
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Add node admin successfully",
		Data: map[string]interface{}{
			"nodeId": nodeId,
			"userId": userId,
		},
	})
}

// DeleteNodeAdmin 删除节点管理员
func DeleteNodeAdmin(c echo.Context) error {
	var userId uint
	nodeId, err := getNodeID(c)
	if err != nil {
		return err
	}
	// 获取claims信息
	claims := c.Get("claims").(util.JwtClaims)
	permission := claims.Permission
	// 如果指定删除任意一位节点管理员,需要超级管理员权限
	if rawUserId := c.QueryParam("user_id"); rawUserId != "" {
		iUserId, err := strconv.Atoi(rawUserId)
		userId = uint(iUserId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
				Code:  http.StatusBadRequest,
				Msg:   "userId param is invalid",
				Error: err.Error(),
			})
		}
		if permission == util.PermissionAdmin {
			err = model.DeleteNodeAdmin(nodeId, userId)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
					Code:  http.StatusInternalServerError,
					Msg:   "Delete node admin failed",
					Error: err.Error(),
				})
			}
		}
	} else {
		// 如果不指定,则是管理员自己主动辞职
		userId = claims.UserId
		// 检查是否为该节点的管理员
		if model.IsNodeAdmin(nodeId, userId) {
			err = model.DeleteNodeAdmin(nodeId, userId)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
					Code:  http.StatusInternalServerError,
					Msg:   "Delete node admin failed",
					Error: err.Error(),
				})
			}
		}
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Delete node admin successfully",
		Data: nil,
	})
}

// ListNodeAdmin 列出节点所有管理员
func ListNodeAdmin(c echo.Context) error {
	var admins []model.User
	nodeId, err := getNodeID(c)
	if err != nil {
		return err
	}
	// 调用model
	admins, err = model.ListNodeAdmin(nodeId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "List node admin failed",
			Error: err.Error(),
		})
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "List node admin successfully",
		Data: map[string]interface{}{
			"adminList": admins,
		},
	})
}

// getNodeID 获取并转换nodeID
func getNodeID(c echo.Context) (uint, error) {
	rawNodeID := c.QueryParam("node_id")
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
