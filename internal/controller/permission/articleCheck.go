package permission

import (
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
)

// ArticlePermissionCheck 检查是否有权限对文章进行编辑操作,仅管理员或文章作者有权操作
func ArticlePermissionCheck(c echo.Context, articleId uint) bool {
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	Permission := claims.Permission
	if Permission < util.PermissionAdmin {
		article, err := model.GetArticleByID(articleId)
		if err != nil {
			return false
		}
		if article.UserID != userId {
			return false
		}
	}
	return true
}
