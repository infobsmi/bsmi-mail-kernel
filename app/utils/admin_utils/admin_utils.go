package admin_utils

import (
	"errors"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func  AdminPermissionCheck(c *gin.Context) (err error, username interface{}, isAdmin interface{}) {
	session := sessions.Default(c)
	username = session.Get("username")
	isAdmin = session.Get("isAdmin")
	if username == nil {
		common.Sugar.Infof("username was nil")
		return errors.New("需要登录"), nil, nil
	}
	common.Sugar.Infof("username was: %s", username.(string))
	if isAdmin == nil {
		return errors.New("需要管理员权限"), nil, nil
	}
	if isAdmin != "yes" {
		return errors.New("需要管理员权限！"), nil, nil
	}
	return nil, username, isAdmin
}
