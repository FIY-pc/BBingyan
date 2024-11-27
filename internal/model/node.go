package model

import (
	"gorm.io/gorm"
	"sort"
)

type Node struct {
	ID      uint      `json:"id" gorm:"primary_key"`
	Name    string    `json:"name" gorm:"unique"`
	Logo    string    `json:"logo"`
	Article []Article `json:"article" gorm:"foreignKey:NodeID;-"` // forbid preload
	// TODO 节点简介
	Users []User `gorm:"many2many:user_nodes;" json:"users"` // 关联用户表
}

// UserNode 表示用户对节点管理权限的关联表
type UserNode struct {
	UserID uint `gorm:"primaryKey"`
	NodeID uint `gorm:"primaryKey"`
}

const (
	SortByTitle      = 1
	SortByLikeNum    = 2
	SortByCommentNum = 3
	SortByTime       = 4
)

func InitNode(DB *gorm.DB) {
	if err := DB.AutoMigrate(&Node{}); err != nil {
		panic(err)
	}
	if err := DB.AutoMigrate(&UserNode{}); err != nil {
		panic(err)
	}
}

func CreateNode(node Node) error {
	result := postgresDb.Create(&node)
	return result.Error
}

func DeleteNodeById(nodeId uint) error {
	result := postgresDb.Where("id", nodeId).Delete(&Node{})
	return result.Error
}

func UpdateNode(node Node) error {
	result := postgresDb.Where("id", node.ID).Updates(node)
	return result.Error
}

func GetNodeById(nodeId uint) (Node, error) {
	var node Node
	result := postgresDb.Where("id", nodeId).First(&node)
	return node, result.Error
}

func GetNodeByName(name string) (Node, error) {
	var node Node
	result := postgresDb.Where("name=?", name).First(&node)
	return node, result.Error
}

func ListArticleFromNode(nodeId uint, page int, pageSize int, sortMethod int) ([]Article, error) {
	var articles []Article
	var result *gorm.DB
	if sortMethod == SortByTitle {
		result = postgresDb.
			Offset((page-1)*pageSize).
			Limit(pageSize).
			Where("node_id", nodeId).
			Order("title").
			Find(&articles)
	}
	if sortMethod == SortByLikeNum {
		result = postgresDb.
			Offset((page-1)*pageSize).
			Limit(pageSize).
			Where("node_id", nodeId).
			Find(&articles)
		sort.SliceStable(articles, func(i, j int) bool {
			num1, _ := GetLikeNum(articles[i].ID)
			num2, _ := GetLikeNum(articles[j].ID)
			return num1 < num2
		})
	}
	if sortMethod == SortByCommentNum {
		result = postgresDb.Table("articles").
			Select("articles.*, COUNT(comments.id) AS comment_count").
			Joins("LEFT JOIN comments ON comments.article_id = articles.id").
			Group("articles.id").
			Order("comment_count DESC").
			Offset((page-1)*pageSize).
			Limit(pageSize).
			Where("node_id", nodeId).
			Find(&articles)
	}
	if sortMethod == SortByTime {
		result = postgresDb.
			Offset((page-1)*pageSize).
			Limit(pageSize).
			Where("node_id", nodeId).
			Order("created_at DESC").
			Find(&articles)
	}
	if result != nil {
		return articles, result.Error
	}
	return articles, nil
}

func CountArticleFromNode(nodeId uint) (int64, error) {
	var count int64
	result := postgresDb.
		Model(&Article{}).
		Where("node_id", nodeId).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func AddNodeAdmin(nodeId, userId uint) error {
	link := UserNode{UserID: userId, NodeID: nodeId}
	if err := postgresDb.Create(&link).Error; err != nil {
		return err
	}
	return nil
}

func DeleteNodeAdmin(nodeId, userId uint) error {
	if err := postgresDb.Delete(&UserNode{UserID: userId, NodeID: nodeId}).Error; err != nil {
		return err
	}
	return nil
}

func ListNodeAdmin(nodeId uint) ([]User, error) {
	var Admin []User
	var result *gorm.DB
	result =
		postgresDb.Model(&User{}).
			Select("users.*").
			Joins("LEFT JOIN user_nodes ON users.id = user_nodes.user_id").
			Where("user_nodes.node_id = ?", nodeId).
			Find(&Admin)
	if result.Error != nil {
		return nil, result.Error
	}
	return Admin, nil
}

func IsNodeAdmin(nodeId, userId uint) bool {
	result := postgresDb.
		Where("node_id = ? AND user_id = ?", nodeId, userId).
		First(&UserNode{})
	if result.Error != nil {
		return false
	}
	return true
}
