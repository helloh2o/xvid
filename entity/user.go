package entity

import "time"

type Model struct {
	Id int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
}
type AppUser struct {
	Model
	UserUid  string `json:"user_uid" gorm:"unique"`
	OpenId   string `json:"open_id" gorm:"index"`
	Channel  string `json:"channel"  gorm:"index"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Level    int    `json:"level"`
	AdTimes  int64  `json:"ad_times"`
	// 水印次数
	LeftTimesWT int64     `json:"left_times_wt"`
	RegTime     time.Time `json:"reg_time"`
	// 最后活动时间
	OpTime     time.Time `json:"op_time"`
	Token      string    `json:"token"   gorm:"index"`
	AfterLogic func()    `json:"-" gorm:"-"`
}
