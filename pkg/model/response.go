package model

import (
	"encoding/json"
	"time"
)

type ResponseData struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
	Time time.Time   `json:"time"`
}

func UnmarshalReponseData(data []byte) (ResponseData, error) {
	var r ResponseData
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ResponseData) MarshalResponseData() ([]byte, error) {
	return json.Marshal(r)
}
