package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// AddNodeAdmin 添加节点管理员
func AddNodeAdmin(c echo.Context) error {
	nodeId, err := params.GetNodeID(c)
	if err != nil {
		return err
	}

	rawAddAdmin := c.FormValue("user_id")
	addAdmin, err := strconv.Atoi(rawAddAdmin)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "Invalid user_id", err)
	}

	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	permission := claims.Permission
	// 权限认证,需节点管理员以上
	if permission < util.PermissionAdmin {
		if !model.IsNodeAdmin(nodeId, userId) {
			return params.CommonErrorGenerate(c, http.StatusUnauthorized, "permission not allowed", err)
		}
	}
	// 调用model
	err = model.AddNodeAdmin(nodeId, uint(addAdmin))
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "add admin failed", err)
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
	nodeId, err := params.GetNodeID(c)
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
			return params.CommonErrorGenerate(c, http.StatusBadRequest, "Invalid user_id", err)
		}
		if permission == util.PermissionAdmin {
			err = model.DeleteNodeAdmin(nodeId, userId)
			if err != nil {
				return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Delete node admin failed", err)
			}
		}
	} else {
		// 如果不指定,则是管理员自己主动辞职
		userId = claims.UserId
		// 检查是否为该节点的管理员
		if model.IsNodeAdmin(nodeId, userId) {
			err = model.DeleteNodeAdmin(nodeId, userId)
			if err != nil {
				return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Delete node admin failed", err)
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
	nodeId, err := params.GetNodeID(c)
	if err != nil {
		return err
	}
	// 调用model
	admins, err = model.ListNodeAdmin(nodeId)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "List node admin failed", err)
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
