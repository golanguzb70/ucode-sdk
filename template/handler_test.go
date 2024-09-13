package function

// import (
// 	"encoding/json"
// 	"fmt"
// 	"testing"
// 	"time"

// 	ucodesdk "github.com/golanguzb70/ucode-sdk"
// 	"github.com/stretchr/testify/assert"
// )

// type TestHandlerI interface {
// 	GetAsserts() []Asserts
// 	GetBenchmarkRequest() Asserts
// }

// func NewAssert(f FunctionAssert) TestHandlerI {
// 	return f
// }

// func TestHandler(t *testing.T) {
// 	asserts := NewAssert(FunctionAssert{})
// 	if len(asserts.GetAsserts()) < 2 {
// 		t.Error("There should be more than 2 cases to check")
// 	}

// 	for _, a := range asserts.GetAsserts() {
// 		requestByte, err := json.Marshal(a.Request)
// 		assert.Nil(t, err)
// 		response := Handle(requestByte)
// 		status, err := ConvertResponse([]byte(response))
// 		assert.Nil(t, err)
// 		assert.Equal(t, a.Response.Status, status.Status)
// 	}
// }

// func BenchmarkHandler(b *testing.B) {
// 	if !IsHTTP {
// 		return
// 	}
// 	a := NewAssert(FunctionAssert{})
// 	var start time.Time

// 	for i := 0; i < b.N; i++ {
// 		reqByte, err := json.Marshal(a.GetBenchmarkRequest().Request)
// 		assert.Nil(b, err)

// 		start = time.Now()

// 		response := Handle(reqByte)

// 		resStatus, err := ConvertResponse([]byte(response))
// 		assert.Nil(b, err)
// 		assert.Equal(b, "done", resStatus.Status)

// 		if time.Since(start) > time.Millisecond*5000 {
// 			assert.Nil(b, fmt.Errorf("took more time than %d ms: %v", 500, time.Since(start)))
// 		}
// 	}
// }

// func ConvertResponse(data []byte) (ucodesdk.ResponseStatus, error) {
// 	response := ucodesdk.ResponseStatus{}
// 	err := json.Unmarshal(data, &response)

// 	return response, err
// }
