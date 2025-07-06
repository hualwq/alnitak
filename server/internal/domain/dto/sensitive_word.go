package dto

// AddSensitiveWordReq 添加敏感词请求
type AddSensitiveWordReq struct {
	Word string `json:"word" binding:"required"`
}

// DeleteSensitiveWordReq 删除敏感词请求
type DeleteSensitiveWordReq struct {
	ID uint `json:"id" binding:"required"`
}

// CheckSensitiveWordReq 检测敏感词请求
type CheckSensitiveWordReq struct {
	Text string `json:"text" binding:"required"`
}

// ReplaceSensitiveWordReq 替换敏感词请求
type ReplaceSensitiveWordReq struct {
	Text      string `json:"text" binding:"required"`
	ReplaceBy string `json:"replace_by"` // 替换字符，默认为*
}

// SensitiveWordListReq 敏感词列表请求
type SensitiveWordListReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Keyword  string `json:"keyword" form:"keyword"`
	Status   *uint  `json:"status" form:"status"`
}
