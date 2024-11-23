package params

import "strconv"

func ArticleLikeKey(articleID uint) string {
	return "like" + ":" + strconv.Itoa(int(articleID))
}

func FollowKey(userID string) string {
	return "follow" + ":" + userID
}
