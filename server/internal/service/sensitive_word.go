package service

import (
	"errors"
	"strings"

	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/domain/vo"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/internal/initialize"
	"interastral-peace.com/alnitak/utils"
)

var SensitiveFilter *initialize.Filter

// InitSensitiveFilter 初始化敏感词过滤器
func InitSensitiveFilter() error {
	SensitiveFilter = initialize.New()

	// 从数据库加载敏感词
	var words []model.SensitiveWord
	if err := global.Mysql.Where("status = ?", 1).Find(&words).Error; err != nil {
		return err
	}

	var wordList []string
	for _, word := range words {
		wordList = append(wordList, word.Word)
	}

	SensitiveFilter.AddWord(wordList...)
	return nil
}

// LoadSensitiveWordsFromFile 从文件加载敏感词到数据库
func LoadSensitiveWordsFromFile(filePath string) error {
	// 读取文件内容
	content, err := utils.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 按行分割
	lines := strings.Split(content, "\n")

	// 批量插入数据库
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检查是否已存在
		var count int64
		global.Mysql.Model(&model.SensitiveWord{}).Where("word = ?", line).Count(&count)
		if count > 0 {
			continue
		}

		// 插入新敏感词
		sensitiveWord := model.SensitiveWord{
			Word:   line,
			Status: 1,
		}

		if err := global.Mysql.Create(&sensitiveWord).Error; err != nil {
			utils.ErrorLog("插入敏感词失败", "sensitive_word", err.Error())
		}
	}

	return nil
}

// AddSensitiveWord 添加敏感词
func AddSensitiveWord(req dto.AddSensitiveWordReq) error {
	// 检查是否已存在
	var count int64
	global.Mysql.Model(&model.SensitiveWord{}).Where("word = ?", req.Word).Count(&count)
	if count > 0 {
		return errors.New("敏感词已存在")
	}

	// 创建敏感词
	sensitiveWord := model.SensitiveWord{
		Word:   req.Word,
		Status: 1,
	}

	if err := global.Mysql.Create(&sensitiveWord).Error; err != nil {
		return err
	}

	// 更新过滤器
	SensitiveFilter.AddWord(req.Word)

	return nil
}

// DeleteSensitiveWord 删除敏感词
func DeleteSensitiveWord(req dto.DeleteSensitiveWordReq) error {
	var sensitiveWord model.SensitiveWord
	if err := global.Mysql.First(&sensitiveWord, req.ID).Error; err != nil {
		return errors.New("敏感词不存在")
	}

	if err := global.Mysql.Delete(&sensitiveWord).Error; err != nil {
		return err
	}

	// 重新初始化过滤器
	return InitSensitiveFilter()
}

// GetSensitiveWordList 获取敏感词列表
func GetSensitiveWordList(req dto.SensitiveWordListReq) (vo.SensitiveWordListVO, error) {
	var words []model.SensitiveWord
	var total int64

	query := global.Mysql.Model(&model.SensitiveWord{})

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("word LIKE ?", "%"+req.Keyword+"%")
	}

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&words).Error; err != nil {
		return vo.SensitiveWordListVO{}, err
	}

	return vo.SensitiveWordListVO{
		List:     vo.ToSensitiveWordVOList(words),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// CheckSensitiveWord 检测敏感词
func CheckSensitiveWord(req dto.CheckSensitiveWordReq) vo.CheckSensitiveWordVO {
	hasSensitive, firstWord := SensitiveFilter.FindIn(req.Text)
	allWords := SensitiveFilter.FindAll(req.Text)

	return vo.CheckSensitiveWordVO{
		HasSensitive: hasSensitive,
		Words:        allWords,
		FirstWord:    firstWord,
	}
}

// ReplaceSensitiveWord 替换敏感词
func ReplaceSensitiveWord(req dto.ReplaceSensitiveWordReq) vo.ReplaceSensitiveWordVO {
	replaceBy := "*"
	if req.ReplaceBy != "" {
		replaceBy = req.ReplaceBy
	}

	replacedText := SensitiveFilter.Replace(req.Text, []rune(replaceBy)[0])

	return vo.ReplaceSensitiveWordVO{
		OriginalText: req.Text,
		ReplacedText: replacedText,
		ReplacedBy:   replaceBy,
	}
}
