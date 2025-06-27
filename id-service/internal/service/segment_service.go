package service

import (
	"database/sql"
	"fmt"
	"id-service/internal/db"
	"id-service/internal/model"
)

// GetAndUpdateSegment 获取并更新号段
// 通过MySQL的原子性UPDATE，保障多实例下号段申请不重复
func GetAndUpdateSegment(bizTag string) (*model.Segment, error) {
	tx, err := db.DB.Beginx() // 启动事务，保障操作原子性
	if err != nil {
		return nil, fmt.Errorf("开启事务失败: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // 遇到panic，事务回滚
		}
	}()

	// 查询当前号段信息，锁定行防止并发冲突
	var seg model.Segment
	query := `SELECT biz_tag, max_id, step FROM id_segments WHERE biz_tag = ? FOR UPDATE`
	if err := tx.Get(&seg, query, bizTag); err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("业务标识[%s]不存在", bizTag)
		}
		return nil, fmt.Errorf("查询号段信息失败: %v", err)
	}

	// 计算新号段的最大ID
	newMaxID := seg.MaxID + int64(seg.Step)

	// 更新号段表
	updateSQL := `UPDATE id_segments SET max_id = ? WHERE biz_tag = ?`
	if _, err := tx.Exec(updateSQL, newMaxID, bizTag); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新号段失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("事务提交失败: %v", err)
	}

	// 返回本次号段信息
	seg.MaxID = newMaxID
	return &seg, nil
}
