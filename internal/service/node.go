package service

import (
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/infrastructure"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/FIY-pc/BBingyan/internal/model"
)

// GetNodeByID 获取节点
func GetNodeByID(nodeId uint) (dto.NodeDTO, error) {
	var nodeDTO dto.NodeDTO
	node := model.Node{}
	if err := infrastructure.PostgresDb.Model(&model.Node{}).Where("id = ?", nodeId).First(&node).Error; err != nil {
		logger.Log.Error(nil, ModelError, "nodeID", nodeId, "error", err)
		return nodeDTO, err
	}
	// 构建返回值
	nodeDTO.ID = node.ID
	nodeDTO.Name = node.Name
	nodeDTO.Intro = node.Intro
	nodeDTO.Avatar = node.Avatar
	logger.Log.Info(nil, Success)
	return nodeDTO, nil
}

func CreateNode(nodeDTO dto.NodeDTO) error {
	// 构建新节点
	creates := model.Node{}
	if nodeDTO.Name != "" {
		creates.Name = nodeDTO.Name
	}
	if nodeDTO.Intro != "" {
		creates.Intro = nodeDTO.Intro
	}
	if nodeDTO.Avatar != "" {
		creates.Avatar = nodeDTO.Avatar
	}
	// 创建节点
	if err := infrastructure.PostgresDb.Model(&model.Node{}).Create(&creates).Error; err != nil {
		logger.Log.Error(nil, ModelError, "nodeID", nodeDTO.ID, "error", err)
		return err
	}
	logger.Log.Info(nil, Success)
	return nil
}
func UpdateNode(nodeDTO dto.NodeDTO) error {
	// 更新节点
	update := make(map[string]interface{})
	if nodeDTO.Name != "" {
		update["name"] = nodeDTO.Name
	}
	if nodeDTO.Intro != "" {
		update["intro"] = nodeDTO.Intro
	}
	if nodeDTO.Avatar != "" {
		update["avatar"] = nodeDTO.Avatar
	}
	if err := infrastructure.PostgresDb.Model(&model.Node{}).Where("id = ?", nodeDTO.ID).Updates(&update).Error; err != nil {
		logger.Log.Error(nil, ModelError, "nodeID", nodeDTO.ID, "error", err)
		return err
	}
	logger.Log.Info(nil, Success)
	return nil
}

// SortDeleteNode 软删除节点
func SortDeleteNode(nodeId uint) error {
	// 软删除节点
	if err := infrastructure.PostgresDb.Model(&model.Node{}).Where("id = ?", nodeId).Delete(&model.Node{}).Error; err != nil {
		logger.Log.Error(nil, ModelError, "nodeID", nodeId, "error", err)
		return err
	}
	logger.Log.Info(nil, Success)
	return nil
}

// HardDeleteNode 硬删除节点
func HardDeleteNode(nodeId uint) error {
	// 硬删除节点
	if err := infrastructure.PostgresDb.Model(&model.Node{}).Unscoped().Where("id = ?", nodeId).Delete(&model.Node{}).Error; err != nil {
		logger.Log.Error(nil, ModelError, "nodeID", nodeId, "error", err)
		return err
	}
	logger.Log.Info(nil, Success)
	return nil
}

// DeletePostUnderNode 删除节点下的所有帖子
func DeletePostUnderNode(nodeId uint) error {
	if err := infrastructure.PostgresDb.Model(&model.Post{}).Where("node_id = ?", nodeId).Delete(&model.Post{}).Error; err != nil {
		logger.Log.Error(nil, ModelError, "nodeID", nodeId, "error", err)
		return err
	}
	logger.Log.Info(nil, Success)
	return nil
}
