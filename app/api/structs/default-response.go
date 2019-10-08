package structs

type DefaultResponse struct {
    Code    int16       `json:"code"`
    Message []string    `json:"message"`
    Status  string      `json:"status"`
}

type Video struct {
    Name    string  `json:"name"`
    Id      string  `json:"id"`
    List    string  `json:"list"`
}

type VideoResponse struct {
    Code    int16       `json:"code"`
    Message []Video     `json:"message"`
    Status  string      `json:"status"` 
} 
