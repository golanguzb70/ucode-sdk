package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	function "github.com/golanguzb70/ucode-sdk/template"
)

func main() {
	// Open and read the JSON file
	jsonFile, err := os.Open("template/request.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Parse the JSON into a struct
	var requestData struct {
		Data json.RawMessage `json:"data"`
	}
	err = json.Unmarshal(byteValue, &requestData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Create a new request
	req, err := http.NewRequest("", "", bytes.NewBuffer(requestData.Data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := function.Handle()

	// Invoke the handler
	handler.ServeHTTP(rr, req)

	// Print the response
	fmt.Printf("Status: %d\n", rr.Code)
	fmt.Printf("Body: %s\n", rr.Body.String())
}
