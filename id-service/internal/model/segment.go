package model

// Segment 代表号段的数据库结构映射
type Segment struct {
	BizTag     string `db:"biz_tag"`     // 业务标识
	MaxID      int64  `db:"max_id"`      // 当前最大ID值
	Step       int    `db:"step"`        // 号段步长
	UpdateTime string `db:"update_time"` // 更新时间，字符串即可满足展示需求
}
