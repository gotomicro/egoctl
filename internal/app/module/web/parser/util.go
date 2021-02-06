package parser

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gotomicro/egoctl/internal/pkg/utils"
	"github.com/gotomicro/egoctl/logger"
)

const (
	AnnotationOverwrite    = "@EgoctlOverwrite"
	AnnotationGenerateTime = "@EgoctlGenerateTime"
)

var CompareExcept = []string{AnnotationGenerateTime}

// write to file
func (c *RenderFile) write(filename string, buf []byte) (err error) {
	if utils.IsExist(filename) && !isNeedOverwrite(filename) {
		return
	}

	filePath := filepath.Dir(filename)
	err = createPath(filePath)
	if err != nil {
		err = errors.New("write create path " + err.Error())
		return
	}

	filePathBak := filePath + "/bak"
	err = createPath(filePathBak)
	if err != nil {
		err = errors.New("write create path bak " + err.Error())
		return
	}

	name := path.Base(filename)

	if utils.IsExist(filename) {
		bakName := fmt.Sprintf("%s/%s.%s.bak", filePathBak, filepath.Base(name), time.Now().Format("2006.01.02.15.04.05"))
		logger.Log.Infof("bak file '%s'", bakName)
		if err := os.Rename(filename, bakName); err != nil {
			err = errors.New("file is bak error, path is " + bakName)
			return err
		}
	}

	file, err := os.Create(filename)
	defer func() {
		err = file.Close()
		if err != nil {
			logger.Log.Fatalf("file close error, err %s", err)
		}
	}()
	if err != nil {
		err = errors.New("write create file " + err.Error())
		return
	}

	err = ioutil.WriteFile(filename, buf, 0644)
	if err != nil {
		err = errors.New("write write file " + err.Error())
		return
	}
	return
}

func isNeedOverwrite(fileName string) (flag bool) {
	seg := GetSeg(filepath.Ext(fileName))

	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()
	overwrite := ""
	var contentByte []byte
	contentByte, err = ioutil.ReadAll(f)
	if err != nil {
		return
	}
	for _, s := range strings.Split(string(contentByte), "\n") {
		s = strings.TrimSpace(strings.TrimPrefix(s, seg))
		if strings.HasPrefix(s, AnnotationOverwrite) {
			overwrite = strings.TrimSpace(s[len(AnnotationOverwrite):])
		}
	}
	if strings.ToLower(overwrite) == "yes" {
		flag = true
		return
	}
	return
}

// createPath 调用os.MkdirAll递归创建文件夹
func createPath(filePath string) error {
	if !utils.IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

func getPackagePath(projectPath string) (packagePath string) {
	f, err := os.Open(projectPath + "/go.mod")
	if err != nil {
		fmt.Println("getPackagePath", err)
		return
	}
	defer f.Close()
	var contentByte []byte
	contentByte, err = ioutil.ReadAll(f)
	if err != nil {
		return
	}
	for _, s := range strings.Split(string(contentByte), "\n") {
		packagePath = strings.TrimSpace(strings.TrimPrefix(s, "module"))
		return
	}
	return
}

func FileContentChange(org, new []byte, seg string) bool {
	if len(org) == 0 {
		return true
	}
	orgContent := GetFilterContent(string(org), seg)
	newContent := GetFilterContent(string(new), seg)
	orgMd5 := md5.Sum([]byte(orgContent))
	newMd5 := md5.Sum([]byte(newContent))
	if orgMd5 != newMd5 {
		return true
	}
	logger.Log.Infof("File has no change in the content")
	return false
}

func GetFilterContent(content string, seg string) string {
	res := ""
	for _, s := range strings.Split(content, "\n") {
		s = strings.TrimSpace(strings.TrimPrefix(s, seg))
		var have bool
		for _, except := range CompareExcept {
			if strings.HasPrefix(s, except) {
				have = true
			}
		}
		if !have {
			res += s
		}
	}
	return res
}

func GetSeg(ext string) string {
	switch ext {
	case ".sql":
		return "--"
	default:
		return "//"
	}
}
