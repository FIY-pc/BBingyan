package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/util"
	"golang.org/x/net/context"
)

func like(articleID uint, UserID uint) {
	ctx := context.Background()
	rdb := util.Rdb
	rdb.SAdd(ctx, params.ArticleLikeKey(articleID), UserID)
}

func Unlike(articleID uint, UserID uint) {
	ctx := context.Background()
	rdb := util.Rdb
	rdb.SRem(ctx, params.ArticleLikeKey(articleID))
}

func GetLikeNum(articleID uint) (int64, error) {
	ctx := context.Background()
	rdb := util.Rdb
	num, err := rdb.SCard(ctx, params.ArticleLikeKey(articleID)).Result()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func isUserLikeArticle(articleID uint, userID uint) bool {
	ctx := context.Background()
	rdb := util.Rdb
	cmd := rdb.SIsMember(ctx, params.ArticleLikeKey(articleID), userID)
	isMember, err := cmd.Result()
	if err != nil {
		return false
	}
	return isMember
}
