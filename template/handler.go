package function

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	sdk "github.com/golanguzb70/ucode-sdk"
)

var (
	baseUrl      = "https://api.admin.u-code.io"
	functionName = ""
)

/*
Answer below questions before starting the function.

When the function invoked?
  - table_slug -> AFTER | BEFORE | HTTP -> CREATE | UPDATE | MULTIPLE_UPDATE | DELETE | APPEND_MANY2MANY | DELETE_MANY2MANY

What does it do?
- Explain the purpose of the function.(O'zbekcha yozilsa ham bo'ladi.)
*/
// func main() {
// 	data := `{"data":{"app_id":"P-CgtoLQxIfoXuz081FuZCenSJbUSMCjOf","object_data":{"test_id":"41574168-4d2f-481a-8c6f-bc60be37e674"}}}`
// 	resp := Handle([]byte(data))
// 	fmt.Println(resp)
// }

// Handle a serverless request
func Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			request       sdk.Request
			response      sdk.Response
			errorResponse sdk.ResponseError
			ucodeApi      = sdk.New(&sdk.Config{BaseURL: baseUrl, FunctionName: functionName})
			returnError   = func(errorResponse sdk.ResponseError) string {
				response = sdk.Response{
					Status: "error",
					Data:   map[string]interface{}{"message": errorResponse.ClientErrorMessage, "error": errorResponse.ErrorMessage, "description": errorResponse.Description},
				}
				marshaledResponse, _ := json.Marshal(response)
				return string(marshaledResponse)
			}
		)
		// set timeout for request
		ucodeApi.Config().RequestTimeout = time.Duration(30 * time.Second)

		// set app_id from .env file
		ucodeApi.Config().AppId = ""

		requestByte, err := io.ReadAll(r.Body)
		if err != nil {
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(requestByte, &request)
		if err != nil {
			errorResponse.ClientErrorMessage = "Error on unmarshal request"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
			return
		}

		// // --------------------------CreateObject------------------------------
		// createHousesRequest := map[string]interface{}{
		// 	"name":       "house_1",
		// 	"price":      15000,
		// 	"room_count": 5,
		// }

		// _, _, err = ucodeApi.CreateObject(&sdk.Argument{DisableFaas: true, TableSlug: "houses", Request: sdk.Request{Data: createHousesRequest}})
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "error on creating new hourse"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }

		// // --------------------------GetList------------------------------
		// var (
		// 	getListRequest = sdk.Request{Data: map[string]interface{}{
		// 		"id":    "f4cca3a7-a0d3-4d7a-b34c-df7d3bf3905c",
		// 		"price": 15000,
		// 	}}
		// )
		// ExistObject, _, err := ucodeApi.GetList(&sdk.ArgumentWithPegination{TableSlug: "houses", Request: getListRequest})
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "Error on useing GetList method"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }
		// response.Data = map[string]interface{}{"result": ExistObject}

		// // --------------------------GetSingle------------------------------
		// house, response, err := ucodeApi.GetSingle(&sdk.Argument{
		// 	TableSlug:   "houses",
		// 	Request:     sdk.Request{Data: map[string]interface{}{"guid": "bf4061fc-f73d-4ecf-bb80-28b9b8a84e13"}},
		// 	DisableFaas: true,
		// })
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "Error on getting single"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }
		// response.Data = map[string]interface{}{"result": house}

		// // --------------------------GetListSlim------------------------------
		// getDemoCourseReq := sdk.Request{Data: map[string]interface{}{"with_relations": true, "selected_relations": []string{"room"}}}
		// getCoursesResp, response, err := ucodeApi.GetListSlim(&sdk.ArgumentWithPegination{
		// 	TableSlug: "houses",
		// 	Request:   getDemoCourseReq,
		// 	Limit:     10,
		// 	Page:      1,
		// })
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "Error on get list"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }
		// response.Data = map[string]interface{}{"result": getCoursesResp}

		// // --------------------------GetSingleSlim------------------------------
		// var id = "e208ae1e-349d-4613-89a1-87c148fa034f"

		// getCourseRequest := sdk.Request{Data: map[string]interface{}{"guid": id}}
		// courseResponse, response, err := ucodeApi.GetSingleSlim(&sdk.Argument{DisableFaas: true, TableSlug: "houses", Request: getCourseRequest})
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "Error on get-single course"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }
		// response.Data = map[string]interface{}{"result": courseResponse}

		// // --------------------------GetListAggregation FOR MongoDB------------------------------
		// getListAggregationPipeline := []map[string]interface{}{
		// 	{"$match": map[string]interface{}{
		// 		"price": map[string]interface{}{
		// 			"$exists": true,
		// 			"$eq":     1000,
		// 		},
		// 	}},
		// }
		// getListAggregationList, response, err := ucodeApi.GetListAggregation(&sdk.Argument{
		// 	TableSlug: "houses",
		// 	Request: sdk.Request{
		// 		Data: map[string]interface{}{"pipelines": getListAggregationPipeline},
		// 	},
		// 	DisableFaas: true,
		// })
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "error on GetListAggregation"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }
		// response.Data = map[string]interface{}{"result": getListAggregationList}

		// // --------------------------UpdateObject------------------------------
		// updateStudent := sdk.Request{
		// 	Data: map[string]interface{}{
		// 		"guid": "a3211725-c7e4-4e34-9375-dafef86e01c3",
		// 		// "name":       "house_13",
		// 		// "price":      15000,
		// 		"room_count": 10,
		// 	},
		// }
		// _, response, err = ucodeApi.UpdateObject(&sdk.Argument{
		// 	TableSlug:   "houses",
		// 	Request:     updateStudent,
		// 	DisableFaas: true,
		// })
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "error on UpdateObject"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }

		// // --------------------------MultipleUpdate------------------------------
		// //Get list demo courses
		// getMultipleUpdateReq := sdk.Request{Data: map[string]interface{}{"with_relations": true, "selected_relations": []string{"room"}}}
		// getMultipleUpdateResp, response, err := ucodeApi.GetListSlim(&sdk.ArgumentWithPegination{
		// 	TableSlug:   "houses",
		// 	Request:     getMultipleUpdateReq,
		// 	DisableFaas: true,
		// })
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "Error on GetListSlim"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }

		// var (
		// 	multipleUpdateRequest = []map[string]interface{}{}
		// )

		// for _, house := range getMultipleUpdateResp.Data.Data.Response {

		// 	multipleUpdateRequest = append(multipleUpdateRequest, map[string]interface{}{
		// 		"guid":       cast.ToString(house["guid"]),
		// 		"room_count": 10,
		// 	})
		// }
		// fmt.Println(multipleUpdateRequest)

		// _, response, err = ucodeApi.MultipleUpdate(&sdk.Argument{
		// 	DisableFaas: true,
		// 	TableSlug:   "houses",
		// 	Request:     sdk.Request{Data: map[string]interface{}{"objects": multipleUpdateRequest}},
		// })
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "Error on MultipleUpdate"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }

		// // --------------------------Delete------------------------------
		// var idDelete = "bf4061fc-f73d-4ecf-bb80-28b9b8a84e13"
		// // var idDelete = "f052e20d-c964-4b42-ace6-038165ed097f"
		// deleteStudentRequest := sdk.Request{Data: map[string]interface{}{"guid": idDelete}}
		// response, err = ucodeApi.Delete(&sdk.Argument{DisableFaas: true, TableSlug: "houses", Request: deleteStudentRequest})
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "Error while Delete"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }

		// // --------------------------MultipleDelete------------------------------
		// var idDelete = []string{"298be83f-3ee3-4fb4-b794-b2f125f4bb21", "e08c83ce-c8e5-4bf9-8511-34a8f747fe26"}
		// deleteStudentRequest := sdk.Request{Data: map[string]interface{}{"ids": idDelete}}
		// response, err = ucodeApi.MultipleDelete(&sdk.Argument{DisableFaas: true, TableSlug: "houses", Request: deleteStudentRequest})
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "Error while Delete"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }

		// // --------------------------AppendManyToMany------------------------------
		// var roomIds = []string{"3844e3c5-4f44-4b25-b01a-7b3b7973c595", "907098ad-3ccf-4011-894d-77b04873d1b1"}

		// appendManyToManyRequest := sdk.Request{
		// 	Data: map[string]interface{}{
		// 		"table_from": "houses",                               // main table
		// 		"table_to":   "room",                                 // relation table
		// 		"id_from":    "85853aa0-71b8-41ad-8fab-d432b2ca53ce", // main table id
		// 		"id_to":      roomIds,                                // relation table id
		// 	},
		// }
		// _, err = ucodeApi.AppendManyToMany(&sdk.Argument{TableSlug: "room", Request: appendManyToManyRequest})
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "Error while AppendManyToMany"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }

		// // --------------------------DeleteManyToMany------------------------------
		// var roomIds = []string{"3844e3c5-4f44-4b25-b01a-7b3b7973c595", "907098ad-3ccf-4011-894d-77b04873d1b1"}

		// appendManyToManyRequest := sdk.Request{
		// 	Data: map[string]interface{}{
		// 		"table_from": "houses",                               // main table
		// 		"table_to":   "room",                                 // relation table
		// 		"id_from":    "85853aa0-71b8-41ad-8fab-d432b2ca53ce", // main table id
		// 		"id_to":      roomIds,                                // relation table id
		// 	},
		// }
		// _, err = ucodeApi.DeleteManyToMany(&sdk.Argument{TableSlug: "room", Request: appendManyToManyRequest})
		// if err != nil {
		// 	errorResponse.Description = response.Data["description"]
		// 	errorResponse.ClientErrorMessage = "Error while AppendManyToMany"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
		// 	return
		// }

		response.Status = "done"
		handleResponse(w, response, http.StatusOK)
	}
}

func handleResponse(w http.ResponseWriter, body interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	bodyByte, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`
			{
				"error": "Error marshalling response"
			}
		`))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(bodyByte)
}
