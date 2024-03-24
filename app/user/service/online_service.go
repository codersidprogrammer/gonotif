package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/codersidprogrammer/gonotif/platform/transport"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"
)

type ActiveUser struct {
	ClientId   string `json:"client_id"`
	Mountpoint string `json:"mountpoint"`
}

type ActiveUserMonitor struct {
	Event string      `json:"event"`
	User  *ActiveUser `json:"user"`
}

type onlineUser struct {
	repository *redis.Client
	ctx        context.Context
}

type OnlineUserService interface {
	SetOnlineUser(u *ActiveUser) (*ActiveUser, error)
	GetOnlineUser(key string) ([]ActiveUser, error)
	DeleteOnlineUser(u *ActiveUser) error
}

func NewOnlineUserService() OnlineUserService {
	return &onlineUser{
		repository: transport.RedisClient,
		ctx:        context.Background(),
	}
}

// DeleteOnlineUser implements OnlineUserService.
func (user *onlineUser) DeleteOnlineUser(u *ActiveUser) error {
	project, err := utils.GetItemFromSplitText(u.ClientId, "_", 0)
	if err != nil {
		return nil
	}

	b, err2 := json.Marshal(u)
	if err2 != nil {
		return err2
	}

	isMember, err3 := user.repository.SIsMember(user.ctx, fmt.Sprintf("push:%s:active", project), string(b)).Result()
	if err3 != nil {
		return err3
	}

	if isMember {
		if _, err := user.repository.SRem(user.ctx, fmt.Sprintf("push:%s:active", project), string(b)).Result(); err != nil {
			return err
		}
	}

	return nil
}

// GetOnlineUser implements OnlineUserService.
func (user *onlineUser) GetOnlineUser(key string) ([]ActiveUser, error) {
	result, err := user.repository.SMembers(user.ctx, key).Result()
	if err != nil {
		log.Error("Failed to get online user key: %s, error: %v", key, err)
		return nil, err
	}

	var aus []ActiveUser
	for _, res := range result {
		var _json ActiveUser
		if err := json.Unmarshal([]byte(res), &_json); err != nil {
			log.Error("Error unmarshaling active user, error ", err)
			break
		}
		aus = append(aus, _json)
	}

	return aus, nil
}

// SetOnlineUser implements OnlineUserService.
func (user *onlineUser) SetOnlineUser(u *ActiveUser) (*ActiveUser, error) {
	project, err := utils.GetItemFromSplitText(u.ClientId, "_", 0)
	if err != nil {
		return nil, err
	}

	b, err2 := json.Marshal(u)
	if err2 != nil {
		return nil, err2
	}

	if _, err := user.repository.SAdd(user.ctx, fmt.Sprintf("push:%s:active", project), string(b)).Result(); err != nil {
		log.Errorf("Failed to set as active user member, error: %v", err)
		return nil, err
	}
	return u, nil
}
