package model

import (
	"time"
)

// SensitiveWord 敏感词模型
type SensitiveWord struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Word      string    `json:"word" gorm:"type:varchar(100);not null;uniqueIndex;comment:敏感词"`
	Status    uint      `json:"status" gorm:"type:tinyint;default:1;comment:状态 1:启用 0:禁用"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SensitiveWord) TableName() string {
	return "sensitive_words"
}
