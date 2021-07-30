package admincontroller

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common/vo"
	"github.com/disintegration/imaging"
	"github.com/flosch/pongo2/v4"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	gonanoid "github.com/matoous/go-nanoid"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"strings"
	"time"
)

const MaxImgWidth = 750

// AdminLoginForm is the login form for Admin
type AdminLoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func AdminPermissionCheck(c *gin.Context) (err error, username interface{}, isAdmin interface{}) {
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

// Export
func ExportCtr(c *gin.Context) {
	err, _, _ := AdminPermissionCheck(c)
	if err != nil {
		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	var blogDataList []model.Article

	result := common.NewDb.
		Order("aid desc").
		Find(&blogDataList)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		common.ShowMessage(c, &vo.Msg{
			Msg: "文章不存在",
		})
		return
	}

	c.JSON(http.StatusOK, blogDataList)
}

func Files(c *gin.Context) {

	err, username, isAdmin := AdminPermissionCheck(c)
	if err != nil {
		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	objectLists := make([]string, 0)
	s, err := awsSession.NewSession(&aws.Config{
		Region: aws.String(common.Config.ObjectStorage.Aws_region),
		Credentials: credentials.NewStaticCredentials(
			common.Config.ObjectStorage.Aws_access_key_id,
			common.Config.ObjectStorage.Aws_secret_access_key,
			"",
		),
	})
	if err != nil {
		common.LogError(err)
		common.ShowUMessage(c, &vo.Umsg{err.Error(), "/"})
		return
	}
	s3o := s3.New(s)
	params := &s3.ListObjectsInput{
		Bucket: aws.String(common.Config.ObjectStorage.Aws_bucket),
	}
	resp, err := s3o.ListObjects(params)
	if err != nil {
		common.LogError(err)
	} else {
		for _, key := range resp.Contents {
			if strings.Contains(*key.Key, ".") {
				objectLists = append(objectLists, *key.Key)
				fmt.Println(*key.Key)
			}
		}
	}
	c.HTML(200, "admin/files.html",
		common.Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        common.Config.Site_name,
			"siteDescription": common.Config.Site_description,
			"cdnurl":          common.Config.ObjectStorage.Cdn_url,
			"username":        username.(string),
			"isAdmin":         isAdmin.(string),
		}))
	return
}
func FileUpload(c *gin.Context) {

	err, _, _ := AdminPermissionCheck(c)
	if err != nil {

		common.LogError(err)
		c.Redirect(301, "/admin/login")
		return
	}
	if common.Config.ObjectStorageType == 1 {

		preFileName, done := UploadByLocalStorage(c)
		if done {
			return
		}
		c.JSON(200, gin.H{"location": preFileName})
	} else {

		preFileName, done := UploadByAwsS3(c)
		if done {
			return
		}
		c.JSON(200, gin.H{"location": common.Config.ObjectStorage.Cdn_url + "/" + preFileName})
	}
}

func UploadByLocalStorage(c *gin.Context) (string, bool) {

	file, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
		common.LogError(err)
		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		common.LogError(err)
		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}
	prefix := time.Now().In(loc).Format("2006/01/02")

	//newFileName := time.Now().UnixNano()
	newFileName, err := gonanoid.Nanoid(20)
	if err != nil {

		common.LogError(err)
		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}

	body, err := ioutil.ReadAll(file)

	kind, _ := filetype.Match(body)
	if kind == filetype.Unknown {
		common.LogInfo("未知文件类型")

		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}

	common.LogInfoF("filetype was: %s", kind.Extension)

	iiReadFile, err := fileHeader.Open()
	if err != nil {
		common.LogError(err)

		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}
	imageInfo, _, err := image.DecodeConfig(iiReadFile)

	if err != nil {
		common.LogError(err)

		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}

	if imageInfo.Width > MaxImgWidth {

		reReadFile, err := fileHeader.Open()
		src, err := imaging.Decode(reReadFile)
		if err != nil {
			common.LogError(err)

			c.JSON(http.StatusBadRequest, "上传失败")
			return "", true
		}

		src = imaging.Resize(src, MaxImgWidth, 0, imaging.Lanczos)
		buf := new(bytes.Buffer)
		imgFormat, err := imaging.FormatFromExtension(kind.Extension)
		if err != nil {
			common.LogError(err)
			c.JSON(http.StatusBadRequest, "获取图片类型失败，上传失败")
			return "", true
		}
		_ = imaging.Encode(buf, src, imgFormat)

		body = buf.Bytes()
	}

	preFileName := fmt.Sprintf("/oss/%s/%s", prefix, newFileName +"."+kind.Extension)
	writeToFileName := "./vol" + preFileName

	common.Sugar.Infof("The writeToFileName: %+v", writeToFileName)
	targetDirectory := filepath.Dir(writeToFileName)
	if _, err := os.Stat(targetDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(targetDirectory, 755)
		if err != nil {
			common.LogError(err)
			c.JSON(http.StatusBadRequest, "上传失败")
			return "", true
		}
	}
	err = ioutil.WriteFile(writeToFileName, body, 0755)

	if err != nil {
		common.LogError(err)
		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}
	return preFileName, false
}

func UploadByAwsS3(c *gin.Context) (string, bool) {
	s, err := awsSession.NewSession(&aws.Config{
		Region:   aws.String(common.Config.ObjectStorage.Aws_region),
		Endpoint: aws.String("https://s3.us-west-001.backblazeb2.com"),
		Credentials: credentials.NewStaticCredentials(
			common.Config.ObjectStorage.Aws_access_key_id,
			common.Config.ObjectStorage.Aws_secret_access_key,
			"",
		),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		common.LogError(err)
		c.JSON(http.StatusBadRequest, "初始化失败，上传失败")
		return "", true
	}
	s3o := s3.New(s)

	file, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
		common.LogError(err)
		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		common.LogError(err)
		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}
	prefix := time.Now().In(loc).Format("2006/01/02")

	//newFileName := time.Now().UnixNano()
	newFileName, err := gonanoid.Nanoid(20)
	if err != nil {

		common.LogError(err)
		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}
	body, err := ioutil.ReadAll(file)

	kind, _ := filetype.Match(body)
	if kind == filetype.Unknown {
		common.LogInfo("未知文件类型")

		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}

	common.LogInfoF("filetype was: %s", kind.Extension)

	iiReadFile, err := fileHeader.Open()
	if err != nil {
		common.LogError(err)

		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}
	imageInfo, _, _ := image.DecodeConfig(iiReadFile)

	if imageInfo.Width > MaxImgWidth {

		reReadFile, err := fileHeader.Open()
		src, err := imaging.Decode(reReadFile)
		if err != nil {
			common.LogError(err)

			c.JSON(http.StatusBadRequest, "上传失败")
			return "", true
		}

		src = imaging.Resize(src, MaxImgWidth, 0, imaging.Lanczos)
		buf := new(bytes.Buffer)
		imgFormat, err := imaging.FormatFromExtension(kind.Extension)
		if err != nil {
			common.LogError(err)
			c.JSON(http.StatusBadRequest, "获取图片类型失败，上传失败")
			return "", true
		}
		_ = imaging.Encode(buf, src, imgFormat)

		body = buf.Bytes()
	}

	preFileName := fmt.Sprintf("%s/%s", prefix, newFileName +"."+kind.Extension)
	params := &s3.PutObjectInput{
		Bucket:      aws.String(common.Config.ObjectStorage.Aws_bucket),
		Key:         aws.String(preFileName),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(fileHeader.Header.Get("content-type")),
	}
	_, err = s3o.PutObject(params)
	if err != nil {
		common.LogError(err)
		c.JSON(http.StatusBadRequest, "上传失败")
		return "", true
	}
	return preFileName, false
}
