package services

import (
	"go.uber.org/zap"

	"scrapper/internal/dto"
)

// tildaService реализация TildaService
type tildaService struct {
	logger *zap.Logger
}

// NewTildaService создает новый экземпляр TildaService
func NewTildaService(logger *zap.Logger) TildaService {
	return &tildaService{
		logger: logger,
	}
}

// DetectPlatform проверяет, соответствует ли страница Tilda
func (s *tildaService) DetectPlatform(html string) bool {
	panic("implement me")
}

// ParseHeader парсит шапку сайта Tilda
func (s *tildaService) ParseHeader(html string) (*dto.Block, error) {
	panic("implement me")
}

// ParseFooter парсит подвал сайта Tilda
func (s *tildaService) ParseFooter(html string) (*dto.Block, error) {
	panic("implement me")
}
