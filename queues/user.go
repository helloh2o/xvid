package queues

import (
	"context"
	"encoding/json"
	"github.com/helloh2o/lucky/cache"
	"github.com/helloh2o/lucky/log"
	"github.com/helloh2o/lucky/utils"
	"github.com/helloh2o/lucky/xdb"
	"xvid/entity"
)

var AppUserQ *utils.LazyQueue

func init() {
	var err error
	// 单节点5万人同时在线有数据变化
	AppUserQ, err = utils.NewLazyQueue(30, 50000, func(key interface{}) error {
		if appUserToken, ok := key.(string); ok {
			var newUserInfo entity.AppUser
			val, err := cache.RedisC.Get(context.Background(), appUserToken).Result()
			if err == nil {
				// find in redis
				err = json.Unmarshal([]byte(val), &newUserInfo)
				if err != nil {
					log.Error("Lazy queue unmarshal app user error %v", err)
				} else {
					if err := xdb.QqsDB().Save(&newUserInfo).Error; err != nil {
						log.Error("Lazy queue save app user error %v", err)
					} else {
						log.Debug("Lazy queue saved app user:: %+v", newUserInfo)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	go AppUserQ.Run()
}
