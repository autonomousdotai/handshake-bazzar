package utils

import (
	"fmt"
	"net/http"
	"mime/multipart"
	"log"
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/configs"
	"encoding/json"
	"io/ioutil"
	"errors"
	"bytes"
	"net/url"
)

type GSService struct {
}

func (gsService GSService) UploadFile(file string, sourceFile *multipart.File) error {
	result := make(map[string]interface{})
	buffer, err := ioutil.ReadAll(*sourceFile)
	if err != nil {
		fmt.Println("Read file error: ", err)
		return err
	}
	filePostBytes := bytes.NewReader(buffer)

	var urlReq *url.URL
	urlReq, err = url.Parse(configs.StorageServiceUrl)
	if err != nil {
		return err
	}
	urlReq.Path += "/"
	parameters := url.Values{}
	parameters.Add("file", file)
	urlReq.RawQuery = parameters.Encode()

	req, err := http.NewRequest("POST", urlReq.String(), filePostBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("http status code %d", resp.StatusCode))
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	if val, ok := result["status"]; ok {
		status := val.(int64)
		if status != 1 {
			return errors.New(fmt.Sprintf("response status %d", resp.StatusCode))
		}
	}
	return nil
}
