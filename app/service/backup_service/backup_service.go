package backup_service

import (
	"archive/zip"
	"github.com/cnmade/bsmi-mail-kernel/pkg/common"
	"github.com/kardianos/osext"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)
func DoBackup() {
	BackupConfig := common.Config.BackupConfig
	if BackupConfig.BackupEnabled == 1 {


		KbDir, err := osext.ExecutableFolder()
		if err != nil {
			common.Sugar.Fatal(err)
			return
		}
		common.Sugar.Infof("KbDir: %+v", KbDir)

		loc, err := time.LoadLocation("Asia/Shanghai")
		fileName := time.Now().In(loc).Format("2006-01-02_15_04_05")

		if _, err := os.Stat(BackupConfig.BackupDir); os.IsNotExist(err) {
			err := os.MkdirAll(BackupConfig.BackupDir, 755)
			if err != nil {
				common.LogError(err)
				common.Sugar.Error("备份目录不存在，且无法创建")
				return
			}
		}
		targetFileName := BackupConfig.BackupDir + "/" + fileName + ".zip"
		err = zipit(KbDir, targetFileName)
		if err != nil {
			common.Sugar.Error(err)
		}
	}
}

func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
