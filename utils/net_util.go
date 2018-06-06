package utils

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"
	"errors"
)

type NetUtil struct {
}


func (netUtil NetUtil) ParseRequest (req *http.Request, results interface{}) (error) {
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (netUtil NetUtil) CurlRequest (req *http.Request) ([]byte, error) {
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return make([]byte, 0), err
	}
	if resp.StatusCode != http.StatusOK {
		return make([]byte, 0), errors.New("HttpResponse Status Code = " + strconv.Itoa(resp.StatusCode))
	}
	return ioutil.ReadAll(resp.Body)
}

func (netUtil NetUtil) Curl (urlStr string) (error) {
	log.Println("curl -> " + urlStr)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(strconv.Itoa(res.StatusCode))
	}
	return nil
}

