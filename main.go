package ucodesdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UcodeApis interface {
	GetList(arg *Argument) (GetListClientApiResponse, Response, error)
}

type object struct {
	config Config
}

func New(cfg Config) UcodeApis {
	return &object{
		config: cfg,
	}
}

func (o *object) GetList(arg *Argument) (GetListClientApiResponse, Response, error) {
	var (
		response      Response
		getListObject GetListClientApiResponse
		url           = fmt.Sprintf("%s/v1/object/get-list/%s?from-ofs=%t", o.config.BaseURL, o.config.TableSlug, arg.DisableFaas)
	)

	getListResponseInByte, err := DoRequest(url, "POST", arg.Request, o.config.AppId)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	err = json.Unmarshal(getListResponseInByte, &getListObject)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, errors.New("invalid response")
	}

	return getListObject, response, nil
}

func DoRequest(url string, method string, body interface{}, appId string) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	request.Header.Add("authorization", "API-KEY")
	request.Header.Add("X-API-KEY", appId)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	return respByte, err
}
