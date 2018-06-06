package utils

import (
	"github.com/ninjadotorg/handshake-bazzar/configs"
)

func CdnUrlFor(fileUrl string) string {
	if fileUrl == "" {
		return ""
	}
	result := ""
	result += configs.AppConf.CdnUrl + "/" + fileUrl
	return result
}

func CdnUrlFor2(filePath string, fileUrl string) string {
	if fileUrl == "" {
		return ""
	}
	result := ""
	result += configs.AppConf.CdnUrl + "/" + filePath + fileUrl
	return result
}
