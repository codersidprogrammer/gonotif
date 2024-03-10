package repository

import (
	"encoding/json"
	"time"
)

type NotificationType int

const (
	Push NotificationType = iota
	Email
)

type NotificationBucket struct {
	ID          string   `json:"id"`
	ProjectID   string   `json:"projectId"`
	ProjectName string   `json:"projectName"`
	UserID      string   `json:"userId"`
	Topics      []string `json:"topics"`
}

func UnmarshalNotificationBucket(data []byte) (NotificationBucket, error) {
	var r NotificationBucket
	err := json.Unmarshal(data, &r)
	return r, err
}

func UnmarshalAllNotificationBucket(data []byte) ([]NotificationBucket, error) {
	var r []NotificationBucket
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NotificationBucket) MarshalNotificationBucket() ([]byte, error) {
	return json.Marshal(r)
}

type NotificationBucketRepository interface {
	Create(d NotificationBucket) (NotificationBucket, error)
	CreateWithLease(d NotificationBucket, leaseDuration time.Duration) (NotificationBucket, error)
	Update(d NotificationBucket) (NotificationBucket, error)
	Delete(id string) error
	Get(id string) (NotificationBucket, error)
	GetAll() ([]NotificationBucket, error)
	GetAllPrefixed(prefix string) ([]NotificationBucket, error)
}
