package service

import (
	"fmt"
	"sync"
)

// SegmentBuffer 缓存结构，存储本地号段信息
type SegmentBuffer struct {
	BizTag    string     // 业务标识
	CurrentID int64      // 当前已发出的ID
	MaxID     int64      // 当前号段的最大ID
	Step      int        // 号段步长
	mu        sync.Mutex // 并发锁，保障线程安全
}

// Generator ID生成器核心结构，管理所有业务的号段缓存
type Generator struct {
	cache map[string]*SegmentBuffer // 业务标识对应的号段缓存
	mu    sync.RWMutex              // 读写锁保护cache并发安全
}

// NewGenerator 初始化ID生成器
func NewGenerator() *Generator {
	return &Generator{
		cache: make(map[string]*SegmentBuffer),
	}
}

// GetID 获取bizTag业务的下一个唯一ID
func (g *Generator) GetID(bizTag string) (int64, error) {
	// 读锁获取缓存
	g.mu.RLock()
	buffer, exists := g.cache[bizTag]
	g.mu.RUnlock()

	// 如果缓存不存在，加载号段
	if !exists {
		if err := g.loadSegment(bizTag); err != nil {
			return 0, fmt.Errorf("加载号段失败: %v", err)
		}
		// 重新读取缓存
		g.mu.RLock()
		buffer = g.cache[bizTag]
		g.mu.RUnlock()
	}

	// 号段操作需要加锁保证并发安全
	buffer.mu.Lock()
	defer buffer.mu.Unlock()

	// 判断是否还有剩余ID
	if buffer.CurrentID < buffer.MaxID {
		buffer.CurrentID++
		return buffer.CurrentID, nil
	}

	// 号段用完，申请新号段
	seg, err := GetAndUpdateSegment(bizTag)
	if err != nil {
		return 0, fmt.Errorf("申请新号段失败: %v", err)
	}

	// 更新本地缓存号段
	buffer.CurrentID = seg.MaxID - int64(seg.Step) + 1
	buffer.MaxID = seg.MaxID
	buffer.Step = seg.Step

	buffer.CurrentID++
	return buffer.CurrentID, nil
}

// loadSegment 初始化或刷新bizTag业务的号段缓存
func (g *Generator) loadSegment(bizTag string) error {
	seg, err := GetAndUpdateSegment(bizTag)
	if err != nil {
		return fmt.Errorf("数据库获取号段失败: %v", err)
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.cache[bizTag] = &SegmentBuffer{
		BizTag:    bizTag,
		CurrentID: seg.MaxID - int64(seg.Step) + 1,
		MaxID:     seg.MaxID,
		Step:      seg.Step,
	}
	return nil
}
