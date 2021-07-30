package fail_ban_service

import (
	"errors"
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"gorm.io/gorm"
	"strings"
)

func CheckBan(clientIp string) bool {
	var banRecord model.FailBan
	err := common.NewDb.First(&banRecord, "ip = ? ", CleanIp(clientIp))
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return false;
	}

	if banRecord.ID <= 0 {
		return false
	}
	if banRecord.Count >= 3 {
		return true
	}
	return false
}

func BanIp(clientIp string)  {
	var banRecord model.FailBan
	cleanIp := CleanIp(clientIp)
	banRecord.Ip = cleanIp

	err := common.NewDb.First(&banRecord, "ip = ? ", cleanIp)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {

		common.NewDb.Create(&banRecord)
	} else {
		banRecord.Count = banRecord.Count + 1;
		common.NewDb.Updates(&banRecord)
	}
}

func ClearBan(arg1 string) {
	common.NewDb.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.FailBan{})

}

func CleanIp(clientIp string) string {
	if strings.ContainsAny(clientIp, ":") {
		sidx := strings.Index(clientIp, ":")
		return clientIp[:sidx]
	}
	return clientIp
}