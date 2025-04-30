package services

import (
	"context"
	
	"go.uber.org/zap"
)

// crawlerService реализация CrawlerService
type crawlerService struct {
	logger         *zap.Logger
	userAgent      string
	maxDepth       int
	allowedDomains []string
}

// NewCrawlerService создает новый экземпляр CrawlerService
func NewCrawlerService(logger *zap.Logger, allowedDomains []string) CrawlerService {
	return &crawlerService{
		logger:         logger,
		userAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
		maxDepth:       2,
		allowedDomains: allowedDomains,
	}
}

// CrawlURL обходит URL и собирает ссылки
func (s *crawlerService) CrawlURL(ctx context.Context, url string, maxDepth int) ([]string, error) {
	panic("implement me")
}

// IsAllowedDomain проверяет, разрешен ли домен для обхода
func (s *crawlerService) IsAllowedDomain(url string) bool {
	panic("implement me")
}

// SetUserAgent устанавливает User-Agent для запросов
func (s *crawlerService) SetUserAgent(userAgent string) {
	s.userAgent = userAgent
}

// SetMaxDepth устанавливает максимальную глубину обхода
func (s *crawlerService) SetMaxDepth(depth int) {
	s.maxDepth = depth
}
