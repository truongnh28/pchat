package service

import (
	"chat-app/config"
	"chat-app/internal/domain"
	"chat-app/pkg/utils/auth"
	"encoding/json"
	"sync"
)

var (
	jwtOnce     sync.Once
	jwtInstance auth.JWTAuth
)

func GetJWTInstance() auth.JWTAuth {
	jwtOnce.Do(func() {
		var (
			cfg = config.GetAppConfig().Authentication
		)
		jwtInstance = auth.NewJWTAuth(cfg.SecretKey, cfg.ExpiredTime, getInfoFromToken)
	})
	return jwtInstance
}

func getInfoFromToken(dt string) (interface{}, error) {
	var (
		acc = &domain.User{}
		err error
	)
	err = json.Unmarshal([]byte(dt), acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
