package router

// Response represents the data to be sent back in the http response body
type Response struct {
	Info map[string]interface{} `json:"info"`
}
