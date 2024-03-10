package database

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/storage/etcd"
)

var EtcdConnection *etcd.Storage
var lock = &sync.Mutex{}

func DoConnectEtcd() (*etcd.Storage, error) {
	port, err := strconv.Atoi(os.Getenv("ETCD_PORT"))
	utils.ExitIfErr(err, "Failed to get etcd port")

	host := os.Getenv("ETCD_HOST")

	url := fmt.Sprintf("%s:%d", host, port)
	options := &etcd.Config{
		Endpoints:   []string{url},
		DialTimeout: 10 * time.Second,
	}

	if EtcdConnection == nil {
		lock.Lock()
		defer lock.Unlock()

		if EtcdConnection == nil {
			log.Info("Creating new etcd connection")
			EtcdConnection = etcd.New(*options)
		} else {
			log.Info("Using existing etcd connection")
		}
	}

	log.Info("etcd conneection established")
	return EtcdConnection, nil
}
