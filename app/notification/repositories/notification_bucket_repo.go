package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/codersidprogrammer/gonotif/platform/database"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/storage/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type bucketRepo struct {
	etcd *etcd.Storage
	key  string
}

func NewBucketRepo() NotificationBucketRepository {
	return &bucketRepo{
		etcd: database.EtcdConnection,
		key:  "/notification/bucket",
	}
}

func NewBucketRepoWithoutConn(etcd *etcd.Storage) NotificationBucketRepository {
	return &bucketRepo{
		etcd: etcd,
		key:  "/notification/bucket",
	}
}

// Create implements NotificationBucketRepository.
func (b *bucketRepo) Create(d NotificationBucket) (NotificationBucket, error) {
	d.ID = utils.GenerateRandomString(24)
	k := fmt.Sprintf("%s/%s", b.key, d.ID)

	bucket, err := d.MarshalNotificationBucket()
	if err != nil {
		return NotificationBucket{}, err
	}

	_, err2 := b.etcd.Conn().Put(context.Background(), k, string(bucket))
	if err2 != nil {
		return NotificationBucket{}, err
	}

	return d, nil
}

func (b *bucketRepo) CreateWithLease(d NotificationBucket, leaseDuration time.Duration) (NotificationBucket, error) {
	d.ID = utils.GenerateRandomString(24)
	k := fmt.Sprintf("%s/%s", b.key, d.ID)

	bucket, err := d.MarshalNotificationBucket()
	if err != nil {
		return NotificationBucket{}, err
	}

	err2 := b.etcd.Set(k, bucket, leaseDuration)
	if err2 != nil {
		return NotificationBucket{}, err
	}

	return d, nil
}

// Delete implements NotificationBucket
func (b *bucketRepo) Delete(id string) error {
	err := b.etcd.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// Get implements NotificationBucketRepository.
func (b *bucketRepo) Get(id string) (NotificationBucket, error) {
	k := fmt.Sprintf("%s/%s", b.key, id)

	b2, err := b.etcd.Get(k)
	if err != nil {
		return NotificationBucket{}, err
	}

	nb, err := UnmarshalNotificationBucket(b2)
	if err != nil {
		return NotificationBucket{}, err
	}

	return nb, nil
}

// Get implements NotificationBucketRepository.
func (b *bucketRepo) GetAll() ([]NotificationBucket, error) {
	b2, err := b.etcd.Conn().Get(context.Background(), b.key, clientv3.WithPrefix())
	if err != nil {
		return []NotificationBucket{}, err
	}

	var result []NotificationBucket
	for _, resp := range b2.Kvs {
		log.Info(string(resp.Value))
		nb, err := UnmarshalNotificationBucket(resp.Value)
		if err != nil {
			return []NotificationBucket{}, err
		}
		result = append(result, nb)
	}

	if err != nil {
		return []NotificationBucket{}, err
	}

	log.Info(result)
	return result, nil
}

func (b *bucketRepo) GetAllPrefixed(prefix string) ([]NotificationBucket, error) {
	b2, err := b.etcd.Conn().Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return []NotificationBucket{}, err
	}

	var result []NotificationBucket
	for _, resp := range b2.Kvs {
		nb, err := UnmarshalNotificationBucket(resp.Value)
		if err != nil {
			return []NotificationBucket{}, err
		}
		result = append(result, nb)
	}

	if err != nil {
		return []NotificationBucket{}, err
	}

	return result, nil
}

// Update implements NotificationBucketRepository.
func (b *bucketRepo) Update(d NotificationBucket) (NotificationBucket, error) {

	if d.ID == "" {
		return NotificationBucket{}, fmt.Errorf("id cannot be empty")
	}

	k := fmt.Sprintf("%s/%s", b.key, d.ID)

	bucket, err := d.MarshalNotificationBucket()
	if err != nil {
		return NotificationBucket{}, err
	}

	_, err2 := b.etcd.Conn().Put(context.Background(), k, string(bucket))
	if err2 != nil {
		return NotificationBucket{}, err
	}

	return d, nil
}
