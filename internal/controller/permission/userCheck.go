package permission

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
)

// UserPermissionCheck 检查是否有权限对用户进行敏感操作
func UserPermissionCheck(c echo.Context, userId uint) bool {
	claimsId, permission, err := params.GetClaimsInfo(c)
	if err != nil {
		return false
	}
	if userId != claimsId && permission < util.PermissionAdmin {
		return false
	}
	return true
}
