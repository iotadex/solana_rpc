package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HttpRequest(url string, method string, postParams []byte, headers map[string]string) ([]byte, error) {
	httpClient := &http.Client{}
	var reader io.Reader
	if len(postParams) > 0 {
		reader = strings.NewReader(string(postParams))
		if headers == nil {
			headers = map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
		}
	} else {
		reader = nil
	}

	//构建request
	request, err := http.NewRequest(method, url, reader)
	if nil != err {
		return nil, fmt.Errorf("NewRequest error. %v", err)
	}

	//添加header
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	// 发出请求
	response, err := httpClient.Do(request)
	if nil != err {
		return nil, fmt.Errorf("do the request error. %v", err)
	}

	defer response.Body.Close()

	// 解析响应内容
	body, err := io.ReadAll(response.Body)
	if nil != err {
		return nil, fmt.Errorf("readAll response.Body error. %v", err)
	}

	return body, nil
}

func HttpJsonPost(url string, postParams interface{}) ([]byte, error) {
	httpClient := http.Client{}

	dataByte, err := json.Marshal(postParams)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(dataByte)
	request, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if nil != err {
		return nil, fmt.Errorf("NewRequest error. %v", err)
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := httpClient.Do(request)
	if nil != err {
		return nil, fmt.Errorf("do the request error. %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if nil != err {
		return nil, fmt.Errorf("readAll response.Body error. %v", err)
	}

	return body, nil
}

func HttpGetWithHeader(url string, headers map[string]string) ([]byte, error) {
	return HttpRequest(url, "GET", nil, headers)
}
