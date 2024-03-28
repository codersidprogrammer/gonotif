package dto

import "encoding/json"

func UnmarshalUserActiveSession(data []byte) (UserActiveSession, error) {
	var r UserActiveSession
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UserActiveSession) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type UserActiveSession struct {
	Table []Table `json:"table"`
	Type  string  `json:"type"`
}

type Table struct {
	ClientID   string `json:"client_id"`
	IsOnline   bool   `json:"is_online"`
	Mountpoint string `json:"mountpoint"`
	PeerHost   string `json:"peer_host"`
	PeerPort   int64  `json:"peer_port"`
	User       string `json:"user"`
}
