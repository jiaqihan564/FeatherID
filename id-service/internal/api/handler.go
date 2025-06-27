package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"id-service/internal/service"
	"id-service/pkg/logger"
)

// response 定义接口统一返回结构
type response struct {
	Code int         `json:"code"`           // 错误码，0表示成功
	Msg  string      `json:"msg"`            // 提示消息
	Data interface{} `json:"data,omitempty"` // 返回数据
}

// Handler 结构，持有ID生成器引用
type Handler struct {
	Generator *service.Generator
}

// NewHandler 构造函数
func NewHandler(gen *service.Generator) *Handler {
	return &Handler{Generator: gen}
}

// GetIDHandler 单个ID获取接口
func (h *Handler) GetIDHandler(w http.ResponseWriter, r *http.Request) {
	// 解析业务标识
	bizTag := r.URL.Query().Get("biz_tag")
	if bizTag == "" {
		writeJSON(w, response{Code: 1, Msg: "biz_tag参数不能为空"})
		return
	}

	// 调用生成器获取ID
	id, err := h.Generator.GetID(bizTag)
	if err != nil {
		writeJSON(w, response{Code: 2, Msg: err.Error()})
		return
	}

	writeJSON(w, response{Code: 0, Msg: "success", Data: id})
}

// GetBatchIDHandler 批量ID获取接口
func (h *Handler) GetBatchIDHandler(w http.ResponseWriter, r *http.Request) {
	bizTag := r.URL.Query().Get("biz_tag")
	countStr := r.URL.Query().Get("count")

	if bizTag == "" {
		writeJSON(w, response{Code: 1, Msg: "biz_tag参数不能为空"})
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 || count > 1000 {
		writeJSON(w, response{Code: 3, Msg: "count参数无效，范围1-1000"})
		return
	}

	ids := make([]int64, 0, count)
	for i := 0; i < count; i++ {
		id, err := h.Generator.GetID(bizTag)
		if err != nil {
			writeJSON(w, response{Code: 2, Msg: err.Error()})
			return
		}
		ids = append(ids, id)
	}

	writeJSON(w, response{Code: 0, Msg: "success", Data: ids})
}

// writeJSON 统一返回JSON格式
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// SetLogLevelHandler 动态调整日志级别接口
func (h *Handler) SetLogLevelHandler(w http.ResponseWriter, r *http.Request) {
	level := r.URL.Query().Get("level")
	if level == "" {
		writeJSON(w, response{Code: 1, Msg: "参数level不能为空"})
		return
	}

	logger.SetLevel(level)
	writeJSON(w, response{Code: 0, Msg: "日志级别调整成功"})
}
