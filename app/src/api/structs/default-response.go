package structs

type DefaultResponse struct {
    Code    int16       `json:"code"`
    Message []string    `json:"message"`
    Status  string      `json:"status"`
}
