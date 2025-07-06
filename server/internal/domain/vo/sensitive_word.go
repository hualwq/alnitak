package vo

import "interastral-peace.com/alnitak/internal/domain/model"

// SensitiveWordVO 敏感词VO
type SensitiveWordVO struct {
	ID        uint   `json:"id"`
	Word      string `json:"word"`
	Status    uint   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// SensitiveWordListVO 敏感词列表VO
type SensitiveWordListVO struct {
	List     []SensitiveWordVO `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

// CheckSensitiveWordVO 检测敏感词VO
type CheckSensitiveWordVO struct {
	HasSensitive bool     `json:"has_sensitive"`
	Words        []string `json:"words"`
	FirstWord    string   `json:"first_word"`
}

// ReplaceSensitiveWordVO 替换敏感词VO
type ReplaceSensitiveWordVO struct {
	OriginalText string `json:"original_text"`
	ReplacedText string `json:"replaced_text"`
	ReplacedBy   string `json:"replaced_by"`
}

// ToSensitiveWordVO 转换为VO
func ToSensitiveWordVO(word model.SensitiveWord) SensitiveWordVO {
	return SensitiveWordVO{
		ID:        word.ID,
		Word:      word.Word,
		Status:    word.Status,
		CreatedAt: word.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: word.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToSensitiveWordVOList 转换为VO列表
func ToSensitiveWordVOList(words []model.SensitiveWord) []SensitiveWordVO {
	var voList []SensitiveWordVO
	for _, word := range words {
		voList = append(voList, ToSensitiveWordVO(word))
	}
	return voList
}
