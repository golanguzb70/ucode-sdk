package ucodesdk

type (
	Request struct {
		Data Data `json:"data"`
	}

	RequestUpdate struct {
		Data map[string]interface{} `json:"data"`
	}

	Argument struct {
		AppId         string        `json:"app_id"`
		TableSlug     string        `json:"table_slug"`
		Request       Request       `json:"request"`
		DisableFaas   bool          `json:"disable_faas"`
		RequestUpdate RequestUpdate `json:"update_request"`
	}

	ArgumentWithPegination struct {
		AppId         string        `json:"app_id"`
		TableSlug     string        `json:"table_slug"`
		Request       Request       `json:"request"`
		DisableFaas   bool          `json:"disable_faas"`
		RequestUpdate RequestUpdate `json:"update_request"`
		Limit         int           `json:"limit"`
		Page          int           `json:"page"`
	}

	Data struct {
		AppId      string                 `json:"app_id"`
		Method     string                 `json:"method"`
		ObjectData map[string]interface{} `json:"object_data"`
		ObjectIds  []string               `json:"object_ids"`
		TableSlug  string                 `json:"table_slug"`
		UserId     string                 `json:"user_id"`
	}
)

// Response structures
type (
	// Create function response body >>>>> CREATE
	Datas struct {
		Data struct {
			Data struct {
				Data map[string]interface{} `json:"data"`
			} `json:"data"`
		} `json:"data"`
	}

	// ClientApiResponse This is get single api response >>>>> GET_SINGLE_BY_ID, GET_SLIM_BY_ID
	ClientApiResponse struct {
		Data ClientApiData `json:"data"`
	}

	ClientApiData struct {
		Data ClientApiResp `json:"data"`
	}

	ClientApiResp struct {
		Response map[string]interface{} `json:"response"`
	}

	Response struct {
		Status string                 `json:"status"`
		Error  string                 `json:"error"`
		Data   map[string]interface{} `json:"data"`
	}

	// GetListClientApiResponse This is get list api response >>>>> GET_LIST, GET_LIST_SLIM
	GetListClientApiResponse struct {
		Data GetListClientApiData `json:"data"`
	}

	GetListClientApiData struct {
		Data GetListClientApiResp `json:"data"`
	}

	GetListClientApiResp struct {
		Response []map[string]interface{} `json:"response"`
	}
	// GetListAggregationClientApiResponse  This is get list aggregation response
	GetListAggregationClientApiResponse struct {
		Data struct {
			Data struct {
				Data []map[string]interface{} `json:"data"`
			} `json:"data"`
		} `json:"data"`
	}

	// ClientApiUpdateResponse This is single update api response >>>>> UPDATE
	ClientApiUpdateResponse struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			TableSlug string                 `json:"table_slug"`
			Data      map[string]interface{} `json:"data"`
		} `json:"data"`
	}

	// ClientApiMultipleUpdateResponse This is multiple update api response >>>>> MULTIPLE_UPDATE
	ClientApiMultipleUpdateResponse struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			Data struct {
				Objects []map[string]interface{} `json:"objects"`
			} `json:"data"`
		} `json:"data"`
	}

	ResponseStatus struct {
		Status string `json:"status"`
	}
)
