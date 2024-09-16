package ucodesdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UcodeApis interface {
	/*
		CreateObject is a function that creates new object.

		Works for [Mongo, Postgres]
	*/
	CreateObject(arg *Argument) (Datas, Response, error)
	/*
		GetList is function that get list of objects from specific table using filter.
		This method works slower because it gets all the information
		about the table, fields and view.
		default_value:
			page = 1
			limit = 10

		Works for [Mongo, Postgres]
	*/
	GetList(arg *ArgumentWithPegination) (GetListClientApiResponse, Response, error)
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

		Works for [Mongo, Postgres]
	*/
	GetListSlim(arg *ArgumentWithPegination) (GetListClientApiResponse, Response, error)
	/*
		GetSingleSlim is function that get one object with its fields.
		It is light and fast to use.

		guid="your_guid"

		Works for [Mongo, Postgres]
	*/
	GetSingleSlim(arg *Argument) (ClientApiResponse, Response, error)
	/*
		GetListAggregation is function that get list of objects with its fields (not include relational data)
		from specific table using filter which you give.
		This method works faster because it does not get all the information
		You should write filter in mongoDB and this function works for only mongoDB.

		pipelines=[]map[string]interface{}{} as filter

		Works for [Mongo]
	*/
	GetListAggregation(arg *Argument) (GetListAggregationClientApiResponse, Response, error)
	/*
		UpdateObject is a function that updates specific object

		Works for [Mongo, Postgres]
	*/
	UpdateObject(arg *Argument) (ClientApiUpdateResponse, Response, error)
	/*
		MultipleUpdate is a function that updates multiple objects at once

		Works for [Mongo, Postgres]
	*/
	MultipleUpdate(arg *Argument) (ClientApiMultipleUpdateResponse, Response, error)
	/*
		Delete is a function that is used to delete one object
		map[guid]="actual_guid"

		Works for [Mongo, Postgres]
	*/
	Delete(arg *Argument) (Response, error)
	/*
		MultipleDelete is a function that is used to delete multiple objects
		map[ids]=[list of ids]

		Works for [Mongo, Postgres]
	*/
	MultipleDelete(arg *Argument) (Response, error)
	/*
		AppendManyToMany is a function that is used to append to a field which referenced many-to-many

		"table_from": "table_name",		// main table
		"table_to":   "table_name",		// relation table
		"id_from":    "table_id", 		// main table id
		"id_to":      "table_id",		// relation table id

		Works for [Mongo]
	*/
	AppendManyToMany(arg *Argument) (Response, error)
	/*
		AppendManyToMany is a function that is used to delete from a field which referenced many-to-many

		"table_from": "table_name",		// main table
		"table_to":   "table_name",		// relation table
		"id_from":    "table_id", 		// main table id
		"id_to":      "table_id",		// relation table id

		Works for [Mongo]
	*/
	DeleteManyToMany(arg *Argument) (Response, error)

	Config() *Config

	DoRequest(url string, method string, body interface{}, headers map[string]string) ([]byte, error)
}

type object struct {
	config *Config
}

func New(cfg *Config) UcodeApis {
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
		url           = fmt.Sprintf("%s/v2/items/%s?from-ofs=%t", o.config.BaseURL, arg.TableSlug, arg.DisableFaas)
	)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	createObjectResponseInByte, err := o.DoRequest(url, "POST", arg.Request, header)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(createObjectResponseInByte), "message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return Datas{}, response, err
	}

	err = json.Unmarshal(createObjectResponseInByte, &createdObject)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(createObjectResponseInByte), "message": "Error while unmarshalling create object", "error": err.Error()}
		response.Status = "error"
		return Datas{}, response, err
	}

	return createdObject, response, nil
}

func (o *object) GetList(arg *ArgumentWithPegination) (GetListClientApiResponse, Response, error) {
	var (
		response      = Response{Status: "done"}
		getListObject GetListClientApiResponse
		url           = fmt.Sprintf("%s/v2/object/get-list/%s?from-ofs=%t", o.config.BaseURL, arg.TableSlug, arg.DisableFaas)
		page          int
		limit         int
	)

	if arg.Page > 0 {
		page = arg.Page
	}

	if arg.Limit > 0 {
		limit = arg.Limit
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	arg.Request.Data["offset"] = (page - 1) * limit
	arg.Request.Data["limit"] = limit

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	getListResponseInByte, err := o.DoRequest(url, "POST", arg.Request, header)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(getListResponseInByte), "message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	err = json.Unmarshal(getListResponseInByte, &getListObject)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(getListResponseInByte), "message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	return getListObject, response, nil
}

func (o *object) GetListSlim(arg *ArgumentWithPegination) (GetListClientApiResponse, Response, error) {
	var (
		response    = Response{Status: "done"}
		listSlim    GetListClientApiResponse
		url         = fmt.Sprintf("%s/v2/object-slim/get-list/%s?from-ofs=%t", o.config.BaseURL, arg.TableSlug, arg.DisableFaas)
		page, limit int
	)

	reqObject, err := json.Marshal(arg.Request.Data)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while marshalling request getting list slim object", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	if arg.Page > 0 {
		page = arg.Page
	}

	if arg.Limit > 0 {
		limit = arg.Limit
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	url = fmt.Sprintf("%s&data=%s&offset=%d&limit=%d", url, string(reqObject), (page-1)*limit, limit)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	getListResponseInByte, err := o.DoRequest(url, "GET", nil, header)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(getListResponseInByte), "message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	err = json.Unmarshal(getListResponseInByte, &listSlim)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(getListResponseInByte), "message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return GetListClientApiResponse{}, response, err
	}

	return listSlim, response, nil
}

func (o *object) GetSingle(arg *Argument) (ClientApiResponse, Response, error) {
	var (
		response  = Response{Status: "done"}
		getObject ClientApiResponse
		url       = fmt.Sprintf("%s/v2/items/%s/%v?from-ofs=%t", o.config.BaseURL, arg.TableSlug, arg.Request.Data["guid"], arg.DisableFaas)
	)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	resByte, err := o.DoRequest(url, "GET", nil, header)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(resByte), "message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	err = json.Unmarshal(resByte, &getObject)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(resByte), "message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	return getObject, response, nil
}

func (o *object) GetSingleSlim(arg *Argument) (ClientApiResponse, Response, error) {
	var (
		response  = Response{Status: "done"}
		getObject ClientApiResponse
		url       = fmt.Sprintf("%s/v1/object-slim/%s/%v?from-ofs=%t", o.config.BaseURL, arg.TableSlug, arg.Request.Data["guid"], arg.DisableFaas)
	)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	resByte, err := o.DoRequest(url, "GET", nil, header)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(resByte), "message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	err = json.Unmarshal(resByte, &getObject)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(resByte), "message": "Error while unmarshalling to object", "error": err.Error()}
		response.Status = "error"
		return ClientApiResponse{}, response, err
	}

	return getObject, response, nil
}

func (o *object) GetListAggregation(arg *Argument) (GetListAggregationClientApiResponse, Response, error) {
	var (
		response           = Response{Status: "done"}
		getListAggregation GetListAggregationClientApiResponse
		url                = fmt.Sprintf("%s/v2/items/%s/aggregation", o.config.BaseURL, arg.TableSlug)
	)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	getListAggregationResponseInByte, err := o.DoRequest(url, "POST", arg.Request, header)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(getListAggregationResponseInByte), "message": "Can't sent request", "error": err.Error()}
		response.Status = "error"
		return GetListAggregationClientApiResponse{}, response, err
	}

	err = json.Unmarshal(getListAggregationResponseInByte, &getListAggregation)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(getListAggregationResponseInByte), "message": "Error while unmarshalling get list object", "error": err.Error()}
		response.Status = "error"
		return GetListAggregationClientApiResponse{}, response, err
	}

	return getListAggregation, response, nil
}

func (o *object) UpdateObject(arg *Argument) (ClientApiUpdateResponse, Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		updateObject ClientApiUpdateResponse
		url          = fmt.Sprintf("%s/v2/items/%s?from-ofs=%t", o.config.BaseURL, arg.TableSlug, arg.DisableFaas)
	)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	updateObjectResponseInByte, err := o.DoRequest(url, "PUT", arg.Request, header)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(updateObjectResponseInByte), "message": "Error while updating object", "error": err.Error()}
		response.Status = "error"
		return ClientApiUpdateResponse{}, response, err
	}

	err = json.Unmarshal(updateObjectResponseInByte, &updateObject)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(updateObjectResponseInByte), "message": "Error while unmarshalling update object", "error": err.Error()}
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
		url                  = fmt.Sprintf("%s/v1/object/multiple-update/%s?from-ofs=%t", o.config.BaseURL, arg.TableSlug, arg.DisableFaas)
	)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	multipleUpdateObjectsResponseInByte, err := o.DoRequest(url, "PUT", arg.Request, header)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(multipleUpdateObjectsResponseInByte), "message": "Error while multiple updating objects", "error": err.Error()}
		response.Status = "error"
		return ClientApiMultipleUpdateResponse{}, response, err
	}

	err = json.Unmarshal(multipleUpdateObjectsResponseInByte, &multipleUpdateObject)
	if err != nil {
		response.Data = map[string]interface{}{"description": string(multipleUpdateObjectsResponseInByte), "message": "Error while unmarshalling multiple update objects", "error": err.Error()}
		response.Status = "error"
		return ClientApiMultipleUpdateResponse{}, response, err
	}

	return multipleUpdateObject, response, nil
}

func (o *object) Delete(arg *Argument) (Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		url = fmt.Sprintf("%s/v2/items/%s/%v?from-ofs=%t", o.config.BaseURL, arg.TableSlug, arg.Request.Data["guid"], arg.DisableFaas)
	)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	_, err := o.DoRequest(url, "DELETE", Request{Data: map[string]interface{}{}}, header)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while deleting object", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
}

func (o *object) MultipleDelete(arg *Argument) (Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		url = fmt.Sprintf("%s/v1/object/%s/?from-ofs=%t", o.config.BaseURL, arg.TableSlug, arg.DisableFaas)
	)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	_, err := o.DoRequest(url, "DELETE", arg.Request.Data, header)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while deleting objects", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
}

func (o *object) AppendManyToMany(arg *Argument) (Response, error) {
	var (
		response = Response{Status: "done"}
		url      = fmt.Sprintf("%s/v2/items/many-to-many?from-ofs=%t", o.config.BaseURL, arg.DisableFaas)
	)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	_, err := o.DoRequest(url, "PUT", arg.Request.Data, header)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while appending many-to-many object", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
}

func (o *object) DeleteManyToMany(arg *Argument) (Response, error) {
	var (
		response = Response{Status: "done"}
		url      = fmt.Sprintf("%s/v2/items/many-to-many?from-ofs=%t", o.config.BaseURL, arg.DisableFaas)
	)

	var appId = o.config.AppId
	if arg.AppId != "" {
		appId = arg.AppId
	}

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	_, err := o.DoRequest(url, "DELETE", arg.Request.Data, header)
	if err != nil {
		response.Data = map[string]interface{}{"message": "Error while deleting many-to-many object", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
}

/*
DoRequest is a function to send http request easily
It gets url, method, body, app_id(for ucode purpose) as paramters

Returns body of the response as array of bytes and error
*/
func (o *object) DoRequest(url string, method string, body interface{}, headers map[string]string) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	if o.config.RequestTimeout > 0 {
		client.Timeout = o.config.RequestTimeout
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Add headers from the map
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	return respByte, err
}

func (o *object) Config() *Config {
	return o.config
}
