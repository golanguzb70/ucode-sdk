package ucodesdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cast"
)

type UcodeApis interface {
	/*
		GetList is function that get list of objects from specific table using filter.
		This method works slower because it gets all the information
		about the table, fields and view.
		default_value:
			page = 1
			limit = 10
	*/
	GetList(arg *Argument) (GetListClientApiResponse, Response, error)
	/*
		GetSingle is function that get one object with all the information of fields, formulas, views and relations.
		It is better to use GetSlim for better performance

		guid="your_guid"
	*/
	GetSingle(arg *Argument) (ClientApiResponse, Response, error)
	/*
		GetListSlim is function that get list of objects from specific table using filter.
		This method works much lighter than GetList because it doesn't get all information about the table, fields and view.
		default_value:
			page = 1
			limit = 10
	*/
	GetListSlim(arg *Argument) (GetListClientApiResponse, Response, error)
	/*
		GetSingleSlim is function that get one object with its fields.
		It is light and fast to use.

		guid="your_guid"
	*/
	GetSingleSlim(arg *Argument) (ClientApiResponse, Response, error)
	/*
		CreateObject is a function that creates new object.

	*/
	CreateObject(arg *Argument) (Datas, Response, error)
	/*
		UpdateObject is a function that updates specific object
	*/
	UpdateObject(arg *Argument) (ClientApiUpdateResponse, Response, error)
	/*
		MultipleUpdate is a function that updates multiple objects at once
	*/
	MultipleUpdate(arg *Argument) (ClientApiMultipleUpdateResponse, Response, error)
	/*
		Delete is a function that is used to delete one object
		map[guid]="actual_guid"
	*/
	Delete(arg *Argument) (Response, error)
}

type object struct {
	config Config
}

func New(cfg Config) UcodeApis {
	return &object{
		config: cfg,
	}
}

func (o *object) CreateObject(arg *Argument) (Datas, Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		createdObject Datas
		url           = fmt.Sprintf("%s/v1/object/%s?from-ofs=%t", o.config.BaseURL, o.config.TableSlug, arg.DisableFaas)
	)

	createObjectResponseInByte, err := DoRequest(url, "POST", arg.Request, o.config.AppId)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return Datas{}, response, err
	}

	err = json.Unmarshal(createObjectResponseInByte, &createdObject)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while unmarshalling create object", "error": err.Error()}
		response.Status = "error"
		return Datas{}, response, err
	}

	return createdObject, response, nil
}

func (o *object) UpdateObject(arg *Argument) (ClientApiUpdateResponse, Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		updateObject ClientApiUpdateResponse
		url          = fmt.Sprintf("%s/v1/object/%s?from-ofs=%t", o.config.BaseURL, o.config.TableSlug, arg.DisableFaas)
	)

	updateObjectResponseInByte, err := DoRequest(url, "PUT", arg.Request, o.config.AppId)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while updating object", "error": err.Error()}
		response.Status = "error"
		return ClientApiUpdateResponse{}, response, err
	}

	err = json.Unmarshal(updateObjectResponseInByte, &updateObject)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while unmarshalling update object", "error": err.Error()}
		response.Status = "error"
		return ClientApiUpdateResponse{}, response, err
	}

	return updateObject, response, nil
}

func (o *object) MultipleUpdate(arg *Argument) (ClientApiMultipleUpdateResponse, Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		multipleUpdateObject ClientApiMultipleUpdateResponse
		url                  = fmt.Sprintf("%s/v1/object/multiple-update/%s?from-ofs=%t", o.config.BaseURL, o.config.TableSlug, arg.DisableFaas)
	)
	multipleUpdateObjectsResponseInByte, err := DoRequest(url, "PUT", arg.Request, o.config.AppId)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while multiple updating objects", "error": err.Error()}
		response.Status = "error"
		return ClientApiMultipleUpdateResponse{}, response, errors.New("error")
	}

	err = json.Unmarshal(multipleUpdateObjectsResponseInByte, &multipleUpdateObject)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while unmarshalling multiple update objects", "error": err.Error()}
		response.Status = "error"
		return ClientApiMultipleUpdateResponse{}, response, errors.New("error")
	}

	return ClientApiMultipleUpdateResponse{}, response, nil
}

func (o *object) GetList(arg *Argument) (GetListClientApiResponse, Response, error) {
	var (
		response      Response
		getListObject GetListClientApiResponse
		url           = fmt.Sprintf("%s/v1/object/get-list/%s?from-ofs=%t", o.config.BaseURL, o.config.TableSlug, arg.DisableFaas)
	)

	page := cast.ToInt(arg.Request.Data["page"])
	limit := cast.ToInt(arg.Request.Data["limit"])
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	arg.Request.Data["offset"] = (page - 1) * limit
	arg.Request.Data["limit"] = limit

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

func (o *object) GetListSlim(arg *Argument) (GetListClientApiResponse, Response, error) {
	var (
		response Response
		listSlim GetListClientApiResponse
		url      = fmt.Sprintf("%s/v1/object-slim/get-list/%s?from-ofs=%t", o.config.BaseURL, o.config.TableSlug, arg.DisableFaas)
	)

	reqObject, err := json.Marshal(arg.Request.Data)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while marshalling request getting list slim object", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	page := cast.ToInt(arg.Request.Data["page"])
	limit := cast.ToInt(arg.Request.Data["limit"])
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	url = fmt.Sprintf("%s&data=%s&offset=%d&limit=%d", url, string(reqObject), (page-1)*limit, limit)

	getListResponseInByte, err := DoRequest(url, "GET", nil, o.config.AppId)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	err = json.Unmarshal(getListResponseInByte, &listSlim)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, errors.New("invalid response")
	}

	return listSlim, response, nil
}

func (o *object) GetSingle(arg *Argument) (ClientApiResponse, Response, error) {
	var (
		response  Response
		getObject ClientApiResponse
		url       = fmt.Sprintf("%s/v1/object/%s/%s?from-ofs=%t", o.config.BaseURL, o.config.TableSlug, cast.ToString(arg.Request.Data["guid"]), arg.DisableFaas)
	)

	resByte, err := DoRequest(url, "GET", nil, o.config.AppId)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	err = json.Unmarshal(resByte, &getObject)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	return getObject, response, nil
}

func (o *object) GetSingleSlim(arg *Argument) (ClientApiResponse, Response, error) {
	var (
		response  Response
		getObject ClientApiResponse
		url       = fmt.Sprintf("%s/v1/object-slim/%s/%s?from-ofs=%t", o.config.BaseURL, o.config.TableSlug, cast.ToString(arg.Request.Data["guid"]), arg.DisableFaas)
	)

	resByte, err := DoRequest(url, "GET", nil, o.config.AppId)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	err = json.Unmarshal(resByte, &getObject)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while unmarshalling to object", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	return getObject, response, nil
}

func (o *object) Delete(arg *Argument) (Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		url = fmt.Sprintf("%s/v1/object/%s/%s?from-ofs=%t", o.config.BaseURL, o.config.TableSlug, cast.ToString(arg.Request.Data["guid"]), arg.DisableFaas)
	)

	_, err := DoRequest(url, "DELETE", nil, o.config.AppId)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while deleting object", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
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
