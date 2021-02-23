package request

import (
	"fmt"
	"log"
	"os"
	"path"
)

var DefaultClient *HttpClient
var logger *log.Logger

var (
	kJsonParseErrFilePath = "jsonParseFailSrcData.log"
	kReqLinkLogPath       = "requestLink.log"
)

func setLogFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, 0666)
	if err != nil {
		return wrapError("setlogfolder:", err)
	}
	kJsonParseErrFilePath = path.Join(folderPath, kJsonParseErrFilePath)
	kReqLinkLogPath = path.Join(folderPath, kReqLinkLogPath)
	return nil
}

func init() {
	DefaultClient = NewClient()
	err := setLogFolder("./log")
	if err != nil {
		log.Panic(err)
	}
	f, err := os.OpenFile(kReqLinkLogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Panic("创建requsetLog失败", err)
	}
	logger = log.New(f, "[request.Log]", 0)
}

func wrapError(newErr string, wrappedErr error) error {
	return fmt.Errorf("%s%s", newErr, wrappedErr.Error())
}
