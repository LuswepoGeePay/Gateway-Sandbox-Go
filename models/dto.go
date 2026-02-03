package models

type GetRequest struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
	SearchQuery string `json:"search_query,omitempty"`
}

func (p *GetRequest) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.PageSize == 0 {
		p.PageSize = 10
	}
}

type AuthorizationRequest struct {
	GetRequest
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

type StatusRequest struct {
	GetRequest
	Status string `json:"status,omitempty"`
}

type APIRequest struct {
	GetRequest
	ProjectID string `json:"project_id,omitempty"`
}

type APIResponsesRequest struct {
	GetRequest
	APIID string `json:"api_id,omitempty"`
}

type APIParametersRequest struct {
	GetRequest
	APIID string `json:"api_id,omitempty"`
}

type APIHeadersRequest struct {
	GetRequest
	APIID string `json:"api_id,omitempty"`
}
type GetRepliesRequest struct {
	GetRequest
	ReviewID string `json:"reviewID,omitempty"`
}

type CallbackPayload struct {
	Code    int                 `json:"code"`
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    CallbackPayloadData `json:"data"`
}

type CallbackPayloadData struct {
	TransactionReference string `json:"transaction_reference"`
	ExternalReference    string `json:"external_reference"`
	Customer             string `json:"customer"`
	Amount               string `json:"amount"`
}
