package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

// NodeInfo 获取节点基本信息,可以使用ID或Name进行查询
func NodeInfo(c echo.Context) error {
	var resultNode model.Node
	// 获取ID
	nodeID, err := params.GetNodeID(c)

	// 若ID为空,获取name
	if nodeID == params.NodeIDIsEmpty {
		name := c.QueryParam("name")
		resultNode, err = model.GetNodeByName(name)
	} else {
		// 使用ID进行查询
		if nodeID == params.NodeIDIsInvalid {
			return err
		}
		resultNode, err = model.GetNodeById(nodeID)
	}

	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Get node failed", err)
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
	nodeID, err := params.GetNodeID(c)
	if err != nil {
		return err
	}
	// 获取原节点信息
	node, err = model.GetNodeById(nodeID)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Get node failed", err)
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
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Update node failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Update node successfully",
		Data: node,
	})
}

// DeleteNode 删除节点
func DeleteNode(c echo.Context) error {
	nodeID, err := params.GetNodeID(c)
	if err != nil {
		return err
	}
	err = model.DeleteNodeById(nodeID)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Delete node failed", err)
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
	nodeID, err := params.GetNodeID(c)
	if err != nil {
		return err
	}
	// 检查是否已经存在
	_, err = model.GetNodeById(nodeID)
	if err == nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "node existed", err)
	}
	// 获取其余信息
	name := c.FormValue("name")
	if name == "" {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "name can't be empty", err)
	}
	logo := c.FormValue("logo")
	if logo != "" {
		node.Logo = logo
	}
	// 创建节点
	err = model.CreateNode(node)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Create node failed", err)
	}
	// 将创建者设置为默认管理员
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	err = model.AddNodeAdmin(nodeID, userId)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "add admin failed", err)
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Create node successfully",
		Data: node,
	})
}
