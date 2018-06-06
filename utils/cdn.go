package utils

import (
	"github.com/ninjadotorg/handshake-bazzar/configs"
)

func CdnUrlFor(fileUrl string) string {
	if fileUrl == "" {
		return ""
	}
	result := ""
	if configs.CdnHttps == true {
		result += "https://"
	} else {
		result += "http://"
	}
	result += configs.CdnDomain + "/" + fileUrl
	return result
}

func CdnUrlFor2(filePath string, fileUrl string) string {
	if fileUrl == "" {
		return ""
	}
	result := ""
	if configs.CdnHttps == true {
		result += "https://"
	} else {
		result += "http://"
	}
	result += configs.CdnDomain + "/" + filePath + fileUrl
	return result
}
