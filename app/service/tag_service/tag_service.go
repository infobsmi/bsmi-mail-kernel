package tag_service

import (
	"errors"
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"gorm.io/gorm"
)

func GetOrNewTagId(tagName string) int64 {
	var tag model.Tag
	result := common.NewDb.
		Where("name = ?", tagName).
		First(&tag)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		tag.Name = tagName
		 result = common.NewDb.Create(&tag)
		 if result.Error != nil {
		 	common.Sugar.Info("Failed to create tag")
		 }
	}
	return tag.ID
}

func BatchGetTagName(tagIds []int64) []model.Tag{
	var tags []model.Tag
	if len(tagIds) <= 0 {
		return tags
	}
	result := common.NewDb.
		Find(&tags, tagIds)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {

		if result.Error != nil {
			common.Sugar.Info("Failed to get tag")
		}
		return nil
	} else {
		common.Sugar.Infof("Tags: %v", tags)
		return tags
	}
}

func CountWithTagId(tagId int64) int64 {
	type Result struct {
		Cnt int64
	}

	var result Result
	common.NewDb.Raw("select count(*) as cnt from bk_article, json_each(bk_article.tag_ids) where json_each.value = ?", tagId).
		Find(&result);
	return result.Cnt;
}

func RefreshCountOfArticle(tagId int64) bool {
	var tag model.Tag
	result := common.NewDb.First(&tag, tagId)

	if result.Error == nil {
		tag.TotalNums = CountWithTagId(tagId)
		common.Sugar.Infof(" totalNums now: %d", tag.TotalNums)
		common.NewDb.Save(&tag)
	} else {
		common.Sugar.Infof("can not find tag with tagId: %s", tagId)
	}
	return true
}