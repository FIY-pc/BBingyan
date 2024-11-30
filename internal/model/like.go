package model

import (
	"github.com/FIY-pc/BBingyan/internal/model/modelParams"
	"github.com/FIY-pc/BBingyan/internal/util"
	"golang.org/x/net/context"
)

func Like(articleID uint, UserID uint) {
	ctx := context.Background()
	rdb := util.Rdb
	rdb.SAdd(ctx, modelParams.ArticleLikeKey(articleID), UserID)
}

func Unlike(articleID uint, UserID uint) {
	ctx := context.Background()
	rdb := util.Rdb
	rdb.SRem(ctx, modelParams.ArticleLikeKey(articleID), UserID)
}

func GetLikeNum(articleID uint) (int64, error) {
	ctx := context.Background()
	rdb := util.Rdb
	num, err := rdb.SCard(ctx, modelParams.ArticleLikeKey(articleID)).Result()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func IsUserLikeArticle(articleID uint, userID uint) bool {
	ctx := context.Background()
	rdb := util.Rdb
	cmd := rdb.SIsMember(ctx, modelParams.ArticleLikeKey(articleID), userID)
	isMember, err := cmd.Result()
	if err != nil {
		return false
	}
	return isMember
}
