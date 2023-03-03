package comm

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/helloh2o/lucky/cache"
	"github.com/helloh2o/lucky/log"
	"github.com/helloh2o/lucky/utils"
	"time"
	"xvid/entity"
)

const UserCacheTime = time.Hour * 3

func SetAppUserCache(token string, info *entity.AppUser) {
	if data, err := json.Marshal(info); err != nil {
		log.Error("SetGameUserInfo err %v", err)
		return
	} else {
		cache.RedisC.Set(context.Background(), token, string(data), UserCacheTime)
	}
}

// InvalidCache 删除缓存
func InvalidCache(token string) {
	n, err := cache.RedisC.Del(context.Background(), token).Result()
	if err != nil {
		log.Error("InvalidCache error %v, key:%s", err, token)
	}
	if n == 0 {
		log.Error("InvalidCache nothing , key:%s", token)
	}
}

// ExistedCache 是否存在redis key
func ExistedCache(key string) bool {
	n, _ := cache.RedisC.Exists(context.Background(), key).Result()
	return n == 1
}

func GetAppUser(token string) (*entity.AppUser, func()) {
	var info *entity.AppUser
GetSyncLock:
	release, ok, wait := utils.RDLockOpTimeout("lock:"+token, time.Minute)
	if !ok {
		<-wait
		log.Debug("token:%s, get app user info wait done, re-try to get sync lock.", token)
		goto GetSyncLock
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	val, err := cache.RedisC.Get(timeoutCtx, token).Result()
	if err != nil {
		if err.Error() != redis.Nil.Error() {
			log.Error("redis error %v", err)
		}
		return nil, release
	}
	err = json.Unmarshal([]byte(val), info)
	if err != nil {
		log.Error("Unmarshal GameUserInfo error %v", err)
		return nil, release
	}
	// reset last op time
	info.OpTime = time.Now()
	return info, func() {
		// rewrite
		data, err := json.Marshal(info)
		// update cache
		if err == nil {
			duration := time.Hour * 3
			cache.RedisC.Set(context.Background(), token, string(data), duration)
		}
		// after logic
		if info.AfterLogic != nil {
			info.AfterLogic()
		}
		// release user token
		release()
	}
}
