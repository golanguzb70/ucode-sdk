package ucodesdk

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/spf13/cast"
)

var (
	baseUrl      = "https://api.admin.u-code.io"
	functionName = ""
)

func TestEndToEnd(t *testing.T) {
	var (
		response      Response
		errorResponse ResponseError
		ucodeApi      = New(&Config{BaseURL: baseUrl, FunctionName: functionName})
		returnError   = func(errorResponse ResponseError) string {
			response = Response{
				Status: "error",
				Data:   map[string]interface{}{"message": errorResponse.ClientErrorMessage, "error": errorResponse.ErrorMessage, "description": errorResponse.Description},
			}
			marshaledResponse, _ := json.Marshal(response)
			return string(marshaledResponse)
		}
		houses []map[string]interface{}
		rooms  []map[string]interface{}
	)
	// Test Config
	t.Run("Config", func(t *testing.T) {
		ucodeApi.Config().RequestTimeout = time.Duration(30 * time.Second)
		ucodeApi.Config().SetBaseUrl(baseUrl)
		if err := ucodeApi.Config().SetAppId(); err != nil {
			t.Errorf("Error setting app_id: %v", err)
		}
	})

	// set base url
	ucodeApi.Config().SetBaseUrl(baseUrl)

	// set app_id from .env file
	if err := ucodeApi.Config().SetAppId(); err != nil {
		errorResponse.ClientErrorMessage = "Error on setting app_id from .env file"
		errorResponse.ErrorMessage = err.Error()
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}

	// --------------------------CreateObject------------------------------
	// create houses
	createHousesRequest := map[string]interface{}{
		"name":       "house",
		"price":      15000,
		"room_count": 5,
	}

	for i := 0; i < 2; i++ {
		_, _, err := ucodeApi.CreateObject(&Argument{DisableFaas: true, TableSlug: "houses", Request: Request{Data: createHousesRequest}})
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "error on creating new hourse"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
		}
	}

	// check error case
	_, _, err := ucodeApi.CreateObject(&Argument{DisableFaas: true, TableSlug: "houses", Request: Request{}})
	if err == nil {
		errorResponse.Description = response.Data["description"]
		errorResponse.ClientErrorMessage = "error: request not given but work"
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}

	// create rooms
	createRoomRequest := map[string]interface{}{
		"name": "room",
	}

	for i := 0; i < 4; i++ {
		_, _, err := ucodeApi.CreateObject(&Argument{DisableFaas: true, TableSlug: "room", Request: Request{Data: createRoomRequest}})
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "error on creating new hourse"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
		}
	}

	// --------------------------GetList------------------------------
	// getting houses
	ExistObject, _, err := ucodeApi.GetList(&ArgumentWithPegination{
		TableSlug:   "houses",
		Request:     Request{Data: map[string]interface{}{}},
		DisableFaas: true,
		Limit:       2,
		Page:        1,
	})
	if err != nil {
		errorResponse.Description = response.Data["description"]
		errorResponse.ClientErrorMessage = "Error on useing GetList method"
		errorResponse.ErrorMessage = err.Error()
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}
	response.Data = map[string]interface{}{"result": ExistObject}
	houses = ExistObject.Data.Data.Response

	// Test with invalid parameters
	_, _, err = ucodeApi.GetList(&ArgumentWithPegination{
		TableSlug: "invalid_table",
		Request:   Request{Data: map[string]interface{}{}},
		Limit:     -1,
		Page:      -1,
	})
	if err == nil {
		t.Error("Expected error for invalid parameters, got nil")
	}

	// --------------------------GetListSlim------------------------------
	getListSlimReq := Request{Data: map[string]interface{}{}}
	getListSlim, response, err := ucodeApi.GetListSlim(&ArgumentWithPegination{
		TableSlug: "room",
		Request:   getListSlimReq,
		Limit:     4,
		Page:      1,
	})
	if err != nil {
		errorResponse.Description = response.Data["description"]
		errorResponse.ClientErrorMessage = "Error on GetListSlim"
		errorResponse.ErrorMessage = err.Error()
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}
	response.Data = map[string]interface{}{"result": getListSlim}
	rooms = getListSlim.Data.Data.Response

	// Test with invalid parameters
	_, _, err = ucodeApi.GetListSlim(&ArgumentWithPegination{
		TableSlug: "invalid_table",
		Request:   getListSlimReq,
		Limit:     -1,
		Page:      -1,
	})
	if err == nil {
		t.Error("Expected error for invalid parameters, got nil")
	}

	// --------------------------UpdateObject------------------------------
	// update first house
	if len(houses) < 2 {
		t.Errorf("error houses count = %d\nExpected count = 2", len(houses))
	}
	updateStudent := Request{
		Data: map[string]interface{}{
			"guid":       houses[0]["guid"],
			"room_count": 10,
		},
	}
	_, response, err = ucodeApi.UpdateObject(&Argument{
		DisableFaas: true,
		TableSlug:   "houses",
		Request:     updateStudent,
	})
	if err != nil {
		errorResponse.Description = response.Data["description"]
		errorResponse.ClientErrorMessage = "error on UpdateObject"
		errorResponse.ErrorMessage = err.Error()
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}

	// Test with invalid parameters
	_, _, err = ucodeApi.UpdateObject(&Argument{
		DisableFaas: true,
		TableSlug:   "invalid_table",
		Request:     Request{Data: map[string]interface{}{"guid": "invalid_guid"}},
	})
	if err == nil {
		t.Error("Expected error for invalid parameters, got nil")
	}

	// --------------------------GetSingle------------------------------
	// get the house info
	houseInfo, response, err := ucodeApi.GetSingle(&Argument{
		TableSlug:   "houses",
		Request:     Request{Data: map[string]interface{}{"guid": cast.ToString(houses[0]["guid"])}},
		DisableFaas: true,
	})
	if err != nil {
		errorResponse.Description = response.Data["description"]
		errorResponse.ClientErrorMessage = "Error on getting single"
		errorResponse.ErrorMessage = err.Error()
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}
	response.Data = map[string]interface{}{"result": houseInfo}

	// Test with invalid parameters
	_, _, err = ucodeApi.GetSingle(&Argument{
		TableSlug:   "invalid_table",
		Request:     Request{Data: map[string]interface{}{"guid": "invalid_guid"}},
		DisableFaas: true,
	})
	if err == nil {
		t.Error("Expected error for invalid parameters, got nil")
	}

	// --------------------------MultipleUpdate------------------------------
	var (
		multipleUpdateRequest = []map[string]interface{}{}
	)

	for _, house := range houses {
		multipleUpdateRequest = append(multipleUpdateRequest, map[string]interface{}{
			"guid":       cast.ToString(house["guid"]),
			"room_count": 15,
		})
	}

	_, response, err = ucodeApi.MultipleUpdate(&Argument{
		DisableFaas: true,
		TableSlug:   "houses",
		Request:     Request{Data: map[string]interface{}{"objects": multipleUpdateRequest}},
	})
	if err != nil {
		errorResponse.Description = response.Data["description"]
		errorResponse.ClientErrorMessage = "Error on MultipleUpdate"
		errorResponse.ErrorMessage = err.Error()
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}

	// Test with invalid parameters
	_, _, err = ucodeApi.MultipleUpdate(&Argument{
		DisableFaas: true,
		TableSlug:   "",
		Request:     Request{Data: map[string]interface{}{"objects": []map[string]interface{}{}}},
	})
	if err == nil {
		t.Error("Expected error for invalid parameters, got nil")
	}

	// // --------------------------GetListAggregation FOR MongoDB------------------------------
	// getListAggregationPipeline := []map[string]interface{}{
	// 	{"$match": map[string]interface{}{
	// 		"price": map[string]interface{}{
	// 			"$exists": true,
	// 			"$eq":     15000,
	// 		},
	// 	}},
	// }
	// getListAggregationList, response, err := ucodeApi.GetListAggregation(&Argument{
	// 	TableSlug: "houses",
	// 	Request: Request{
	// 		Data: map[string]interface{}{"pipelines": getListAggregationPipeline},
	// 	},
	// 	DisableFaas: true,
	// })
	// if err != nil {
	// 	errorResponse.Description = response.Data["description"]
	// 	errorResponse.ClientErrorMessage = "error on GetListAggregation"
	// 	errorResponse.ErrorMessage = err.Error()
	// 	errorResponse.StatusCode = http.StatusInternalServerError
	// 	t.Error(returnError(errorResponse))
	// }
	// response.Data = map[string]interface{}{"result": getListAggregationList}

	// --------------------------AppendManyToMany------------------------------
	for i := 0; i < 2; i++ {
		var roomIds = []string{cast.ToString(rooms[i]["guid"]), cast.ToString(rooms[i+1]["guid"]),
			cast.ToString(rooms[i+1]["guid"])}

		appendManyToManyRequest := Request{
			Data: map[string]interface{}{
				"table_from": "houses",                         // main table
				"table_to":   "room",                           // relation table
				"id_from":    cast.ToString(houses[i]["guid"]), // main table id
				"id_to":      roomIds,                          // relation table id
			},
		}
		_, err = ucodeApi.AppendManyToMany(&Argument{TableSlug: "houses", Request: appendManyToManyRequest})
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error while AppendManyToMany"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
		}
	}

	// --------------------------GetSingleSlim------------------------------
	var id = cast.ToString(rooms[0]["guid"])

	getCourseRequest := Request{Data: map[string]interface{}{"guid": id}}
	courseResponse, response, err := ucodeApi.GetSingleSlim(&Argument{DisableFaas: true, TableSlug: "room", Request: getCourseRequest})
	if err != nil {
		errorResponse.Description = response.Data["description"]
		errorResponse.ClientErrorMessage = "Error on get-single course"
		errorResponse.ErrorMessage = err.Error()
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}
	response.Data = map[string]interface{}{"result": courseResponse}

	// Test with invalid parameters
	_, _, err = ucodeApi.GetSingleSlim(&Argument{
		DisableFaas: true,
		TableSlug:   "invalid_table",
		Request:     Request{Data: map[string]interface{}{"guid": "invalid_guid"}},
	})
	if err == nil {
		t.Error("Expected error for invalid parameters, got nil")
	}

	// --------------------------DeleteManyToMany------------------------------
	for i := 0; i < 2; i++ {
		var houseId = cast.ToString(houses[i]["guid"])

		deleteManyToManyRequest := Request{
			Data: map[string]interface{}{
				"table_from": "houses",   // main table
				"table_to":   "room",     // relation table
				"id_from":    houseId,    // main table id
				"id_to":      []string{}, // relation table id
			},
		}
		_, err = ucodeApi.DeleteManyToMany(&Argument{TableSlug: "houses", Request: deleteManyToManyRequest})
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error while AppendManyToMany"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
		}
	}

	// --------------------------Delete------------------------------
	DeleteRequest := Request{Data: map[string]interface{}{"guid": cast.ToString(houses[0]["guid"])}}
	response, err = ucodeApi.Delete(&Argument{DisableFaas: true, TableSlug: "houses", Request: DeleteRequest})
	if err != nil {
		errorResponse.Description = response.Data["description"]
		errorResponse.ClientErrorMessage = "Error while Delete"
		errorResponse.ErrorMessage = err.Error()
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}

	// --------------------------MultipleDelete------------------------------
	// deleting from houses
	var (
		idMultipleDeleteHouses = []string{}
	)
	for _, val := range houses {
		idMultipleDeleteHouses = append(idMultipleDeleteHouses, cast.ToString(val["guid"]))
	}

	MultipleDeleteHouses := Request{Data: map[string]interface{}{"ids": idMultipleDeleteHouses}}
	response, err = ucodeApi.MultipleDelete(&Argument{DisableFaas: true, TableSlug: "houses", Request: MultipleDeleteHouses})
	if err != nil {
		errorResponse.Description = response.Data["description"]
		errorResponse.ClientErrorMessage = "Error while Delete"
		errorResponse.ErrorMessage = err.Error()
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}

	// --------------------------MultipleDelete------------------------------
	// deleting from rooms
	var (
		idMultipleDeleteRoom = []string{}
	)
	for _, val := range rooms {
		idMultipleDeleteRoom = append(idMultipleDeleteRoom, cast.ToString(val["guid"]))
	}

	MultipleDeleteRooms := Request{Data: map[string]interface{}{"ids": idMultipleDeleteRoom}}
	response, err = ucodeApi.MultipleDelete(&Argument{DisableFaas: true, TableSlug: "room", Request: MultipleDeleteRooms})
	if err != nil {
		errorResponse.Description = response.Data["description"]
		errorResponse.ClientErrorMessage = "Error while Delete"
		errorResponse.ErrorMessage = err.Error()
		errorResponse.StatusCode = http.StatusInternalServerError
		t.Error(returnError(errorResponse))
	}
}

func TestDoRequest(t *testing.T) {
	ucodeApi := New(&Config{BaseURL: baseUrl, FunctionName: functionName})

	// Test successful request
	_, err := ucodeApi.DoRequest(baseUrl+"/test", "GET", nil, "test_app_id", nil)
	if err != nil {
		t.Errorf("Error on DoRequest: %v", err)
	}

	// Test with invalid URL
	_, err = ucodeApi.DoRequest("invalid-url", "GET", nil, "test_app_id", nil)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}

	// Test with custom headers
	customHeaders := map[string]string{
		"Custom-Header": "TestValue",
	}
	_, err = ucodeApi.DoRequest(baseUrl+"/test", "GET", nil, "test_app_id", customHeaders)
	if err != nil {
		t.Errorf("Error on DoRequest with custom headers: %v", err)
	}

	// Test with request timeout
	ucodeApi.Config().RequestTimeout = time.Duration(1 * time.Nanosecond)
	_, err = ucodeApi.DoRequest(baseUrl+"/test", "GET", nil, "test_app_id", nil)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestConfigMethods(t *testing.T) {
	// Create a new Config object for testing
	cfg := &Config{BaseURL: baseUrl, FunctionName: functionName}
	ucodeApi := New(cfg)

	// Test SetBaseUrl
	newBaseURL := "https://new.api.example.com"
	ucodeApi.Config().SetBaseUrl(newBaseURL)
	if ucodeApi.Config().BaseURL != newBaseURL {
		t.Errorf("SetBaseUrl failed, expected %s, got %s", newBaseURL, ucodeApi.Config().BaseURL)
	}

	// Backup the current value of APP_ID to restore it after the test
	originalAppID, exists := os.LookupEnv("APP_ID")

	// Ensure we clean up the environment variable at the end
	defer func() {
		if exists {
			os.Setenv("APP_ID", originalAppID)
		} else {
			os.Unsetenv("APP_ID")
		}
	}()

	// Test SetAppId with APP_ID set
	os.Setenv("APP_ID", "test_app_id")
	err := ucodeApi.Config().SetAppId()
	if err != nil {
		t.Errorf("SetAppId failed: %v", err)
	}
	if ucodeApi.Config().appId != "test_app_id" {
		t.Errorf("SetAppId failed, expected test_app_id, got %s", ucodeApi.Config().appId)
	}

	// Test SetAppId with APP_ID unset
	os.Setenv("APP_ID", "") // set the APP_ID environment variable
	err = ucodeApi.Config().SetAppId()
	if err == nil {
		t.Error("Expected error for missing APP_ID environment variable, got nil")
	}
}
