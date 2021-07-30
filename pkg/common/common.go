package common

import (
	"database/sql"
	"github.com/cnmade/bsmi-mail-kernel/app/orm/model"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common/vo"
	"github.com/flosch/pongo2/v4"
	"github.com/gin-gonic/gin"
	"github.com/grokify/html-strip-tags-go"
	"github.com/kataras/hcaptcha"
	"github.com/naoina/toml"
	"github.com/ztrue/tracerr"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

/**
 * Logging error
 */
func LogError(err error) {
	if err != nil {
		Sugar.Error(tracerr.Sprint(tracerr.Wrap(err)))
	}
}

/**
 * Logging info
 */
func LogInfo(msg string) {
	if msg != "" {
		Sugar.Info(msg)
	}
}

func LogInfoF(msg string, v interface{}) {
	if msg != "" {
		Sugar.Infof(msg, v)
	}
}



/**
 * close rows defer
 */
func CloseRowsDefer(rows *sql.Rows) {
	_ = rows.Close()
}

/*
* ShowMessage with template
 */
func ShowMessage(c *gin.Context, m *vo.Msg) {

	c.HTML(200, "message-traditional.html",
		Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        Config.Site_name,
			"siteDescription": Config.Site_description,
			"message":         m.Msg,
		}))
	return
}

func ShowUMessage(c *gin.Context, m *vo.Umsg) {

	c.HTML(200, "message-traditional.html",
		Pongo2ContextWithVersion(pongo2.Context{
			"siteName":        Config.Site_name,
			"siteDescription": Config.Site_description,
			"message":         m.Msg,
			"url":             m.Url,
		}))
	return
}

func GetMinutes() string {
	return time.Now().Format("200601021504")
}

func GetNewDb(config *vo.AppConfig) *gorm.DB {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      false,         // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open(config.Dbdsn), &gorm.Config{

		Logger: newLogger,
	})
	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(&model.Category{})

	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(&model.Tag{})

	if err != nil {
		panic(err.Error())
	}


	err = db.AutoMigrate(&model.Article{})

	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(&model.ArticleHistory{})

	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(&model.TwoAuth{})

	if err != nil {
		panic(err.Error())
	}


	err = db.AutoMigrate(&model.FailBan{})

	if err != nil {
		panic(err.Error())
	}


	return db
}

func GetConfig() *vo.AppConfig {
	_cm := "GetConfig@pkg/common/common"
	//TODO load config from cmd line argument
	f, err := os.Open("./vol/config.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config vo.AppConfig
	if err := toml.Unmarshal(buf, &config); err != nil {
		Sugar.Infof(_cm + " error: %+v", err)
	}
	return &config
}



var (
	Config    *vo.AppConfig
	NewDb     *gorm.DB
	Logger, _ = zap.NewProduction()
	Sugar *zap.SugaredLogger
	BsmiKbVersion string
	HCaptchClient *hcaptcha.Client
)

func InitApp() {
	Config = GetConfig()
//	gin.SetMode(Config.SrvMode)
	gin.SetMode(gin.DebugMode)
	NewDb = GetNewDb(Config)
	defer Logger.Sync()
	Sugar = Logger.Sugar()

	HCaptchClient = hcaptcha.New(Config.HCaptchaSecretKey)
}

func OutPutHtml( c *gin.Context, s string) {
	c.Header("Content-Type", "text/html;charset=UTF-8")
	c.String(200, "%s", s)
	return
}
func OutPutText( c *gin.Context, s string) {
	c.Header("Content-Type", "text/plain;charset=UTF-8")
	c.String(200, "%s", s)
	return
}
/**
 * 截取指定长度的字符串，中文
 */
func SubCutContent(content string, length int) string {
	if len(content) <= length {
		return content
	}

	content = strip.StripTags(content)
	content = strings.TrimSpace(content)
	content = strings.Replace(content, "<!DOCTYPE html>", "", 1)
	content = strings.Replace(content, "&nbsp;", "", 1)

	tmpContent := []rune(content)

	rawLen := len(tmpContent)

	if length > rawLen {
		return content
	}

	return string(tmpContent[0:length])
}

