package service

import repository "github.com/codersidprogrammer/gonotif/app/notification/repositories"

type notification struct {
	notificationRepository repository.NotificationBucketRepository
}

// Create implements NotificationService.

type NotificationService interface {
	Create(d repository.NotificationBucket) (repository.NotificationBucket, error)
	Get(id string) (repository.NotificationBucket, error)
	GetAll() ([]repository.NotificationBucket, error)
	GetAllWithPrefix(prefix string) ([]repository.NotificationBucket, error)
}

func NewNotificationService() NotificationService {
	return &notification{
		notificationRepository: repository.NewBucketRepo(),
	}
}

func (n *notification) Create(d repository.NotificationBucket) (repository.NotificationBucket, error) {
	nb, err := n.notificationRepository.Create(d)
	return nb, err
}

func (n *notification) Get(id string) (repository.NotificationBucket, error) {
	nb, err := n.notificationRepository.Get(id)
	return nb, err
}

func (n *notification) GetAll() ([]repository.NotificationBucket, error) {
	nb, err := n.notificationRepository.GetAll()
	return nb, err
}

func (n *notification) GetAllWithPrefix(prefix string) ([]repository.NotificationBucket, error) {
	nb, err := n.notificationRepository.GetAllPrefixed(prefix)
	return nb, err
}
