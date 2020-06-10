package api

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type FailedResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
