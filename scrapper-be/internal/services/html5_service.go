package services

import (
	"go.uber.org/zap"

	"scrapper/internal/dto"
)

// html5Service реализация HTML5Service
type html5Service struct {
	logger *zap.Logger
}

// NewHTML5Service создает новый экземпляр HTML5Service
func NewHTML5Service(logger *zap.Logger) HTML5Service {
	return &html5Service{
		logger: logger,
	}
}

// DetectPlatform проверяет, соответствует ли страница HTML5
func (s *html5Service) DetectPlatform(html string) bool {
	panic("implement me")
}

// ParseHeader парсит шапку сайта HTML5
func (s *html5Service) ParseHeader(html string) (*dto.Block, error) {
	panic("implement me")
}

// ParseFooter парсит подвал сайта HTML5
func (s *html5Service) ParseFooter(html string) (*dto.Block, error) {
	panic("implement me")
}
