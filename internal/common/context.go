package common

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AppContext struct {
	DB    *gorm.DB
	Cache *redis.Client
}
