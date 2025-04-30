package services

import (
	"go.uber.org/zap"

	"scrapper/internal/dto"
)

// bitrixService реализация BitrixService
type bitrixService struct {
	logger *zap.Logger
}

// NewBitrixService создает новый экземпляр BitrixService
func NewBitrixService(logger *zap.Logger) BitrixService {
	return &bitrixService{
		logger: logger,
	}
}

// DetectPlatform проверяет, соответствует ли страница Bitrix
func (s *bitrixService) DetectPlatform(html string) bool {
	panic("implement me")
}

// ParseHeader парсит шапку сайта Bitrix
func (s *bitrixService) ParseHeader(html string) (*dto.Block, error) {
	panic("implement me")
}

// ParseFooter парсит подвал сайта Bitrix
func (s *bitrixService) ParseFooter(html string) (*dto.Block, error) {
	panic("implement me")
}
