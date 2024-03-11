package service

import repository "github.com/codersidprogrammer/gonotif/app/notification/repositories"

type notificationBucket struct {
	notificationRepository repository.NotificationBucketRepository
}

// Create implements NotificationService.

type NotificationBucketService interface {
	Create(d repository.NotificationBucket) (repository.NotificationBucket, error)
	Get(id string) (repository.NotificationBucket, error)
	GetAll() ([]repository.NotificationBucket, error)
	GetAllWithPrefix(prefix string) ([]repository.NotificationBucket, error)
}

func NewNotificationBucketService() NotificationBucketService {
	return &notificationBucket{
		notificationRepository: repository.NewBucketRepo(),
	}
}

func (n *notificationBucket) Create(d repository.NotificationBucket) (repository.NotificationBucket, error) {
	nb, err := n.notificationRepository.Create(d)
	return nb, err
}

func (n *notificationBucket) Get(id string) (repository.NotificationBucket, error) {
	nb, err := n.notificationRepository.Get(id)
	return nb, err
}

func (n *notificationBucket) GetAll() ([]repository.NotificationBucket, error) {
	nb, err := n.notificationRepository.GetAll()
	return nb, err
}

func (n *notificationBucket) GetAllWithPrefix(prefix string) ([]repository.NotificationBucket, error) {
	nb, err := n.notificationRepository.GetAllPrefixed(prefix)
	return nb, err
}
