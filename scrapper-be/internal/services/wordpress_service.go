package services

import (
	"regexp"
	"strings"

	"go.uber.org/zap"

	"scrapper/internal/dto"
)

// wordPressService реализация WordPressService
type wordPressService struct {
	logger *zap.Logger
}

// NewWordPressService создает новый экземпляр WordPressService
func NewWordPressService(logger *zap.Logger) WordPressService {
	return &wordPressService{
		logger: logger,
	}
}

// DetectPlatform проверяет, соответствует ли страница WordPress
func (s *wordPressService) DetectPlatform(html string) bool {
	// Проверяем наличие характерных признаков WordPress
	wpPatterns := []string{
		`wp-content`,
		`wp-includes`,
		`wp-json`,
		`<meta name="generator" content="WordPress`,
		`class="wordpress"`,
	}

	for _, pattern := range wpPatterns {
		if strings.Contains(html, pattern) {
			s.logger.Debug("WordPress pattern detected", zap.String("pattern", pattern))
			return true
		}
	}

	return false
}

// ParseHeader парсит шапку сайта WordPress
func (s *wordPressService) ParseHeader(html string) (*dto.Block, error) {
	// Создаем регулярное выражение для поиска шапки
	reHeader := regexp.MustCompile(`(?s)<header.*?>(.*?)</header>`)
	headerMatch := reHeader.FindStringSubmatch(html)

	// Если не найдено <header>, пробуем найти по id или классу
	if len(headerMatch) < 2 {
		reHeaderClass := regexp.MustCompile(`(?s)<div\s+(?:class=".*?header.*?".*?|id=".*?header.*?".*?)>(.*?)</div>`)
		headerMatch = reHeaderClass.FindStringSubmatch(html)
	}

	// Если все еще не найдено, пробуем найти .site-header
	if len(headerMatch) < 2 {
		reSiteHeader := regexp.MustCompile(`(?s)<div\s+class=".*?site-header.*?".*?>(.*?)</div>`)
		headerMatch = reSiteHeader.FindStringSubmatch(html)
	}

	// Если шапка не найдена, возвращаем nil
	if len(headerMatch) < 2 {
		s.logger.Warn("WordPress header not found")
		return nil, nil
	}

	headerHtml := headerMatch[0]

	// Парсим логотип
	reLogo := regexp.MustCompile(`(?s)<a.*?class=".*?logo.*?".*?>(.*?)</a>`)
	logoMatch := reLogo.FindStringSubmatch(headerHtml)

	var logoHtml string
	if len(logoMatch) >= 2 {
		logoHtml = logoMatch[0]
	}

	// Парсим навигационное меню
	reMenu := regexp.MustCompile(`(?s)<nav.*?>(.*?)</nav>`)
	menuMatch := reMenu.FindStringSubmatch(headerHtml)

	var menuHtml string
	if len(menuMatch) >= 2 {
		menuHtml = menuMatch[0]
	}

	// Создаем структуру для хранения содержимого шапки
	content := map[string]interface{}{
		"logo": logoHtml,
		"menu": menuHtml,
	}

	// Создаем блок
	block := &dto.Block{
		BlockType: dto.BlockTypeHeader,
		Platform:  dto.PlatformWordPress,
		Content:   content,
		HTML:      headerHtml,
	}

	s.logger.Debug("WordPress header parsed successfully")
	return block, nil
}

// ParseFooter парсит подвал сайта WordPress
func (s *wordPressService) ParseFooter(html string) (*dto.Block, error) {
	// Создаем регулярное выражение для поиска подвала
	reFooter := regexp.MustCompile(`(?s)<footer.*?>(.*?)</footer>`)
	footerMatch := reFooter.FindStringSubmatch(html)

	// Если не найдено <footer>, пробуем найти по id или классу
	if len(footerMatch) < 2 {
		reFooterClass := regexp.MustCompile(`(?s)<div\s+(?:class=".*?footer.*?".*?|id=".*?footer.*?".*?)>(.*?)</div>`)
		footerMatch = reFooterClass.FindStringSubmatch(html)
	}

	// Если все еще не найдено, пробуем найти .site-footer
	if len(footerMatch) < 2 {
		reSiteFooter := regexp.MustCompile(`(?s)<div\s+class=".*?site-footer.*?".*?>(.*?)</div>`)
		footerMatch = reSiteFooter.FindStringSubmatch(html)
	}

	// Если подвал не найден, возвращаем nil
	if len(footerMatch) < 2 {
		s.logger.Warn("WordPress footer not found")
		return nil, nil
	}

	footerHtml := footerMatch[0]

	// Парсим виджеты
	reWidgets := regexp.MustCompile(`(?s)<div\s+class=".*?widgets.*?".*?>(.*?)</div>`)
	widgetsMatch := reWidgets.FindStringSubmatch(footerHtml)

	var widgetsHtml string
	if len(widgetsMatch) >= 2 {
		widgetsHtml = widgetsMatch[0]
	}

	// Парсим копирайт
	reCopyright := regexp.MustCompile(`(?s)<div\s+class=".*?copyright.*?".*?>(.*?)</div>`)
	copyrightMatch := reCopyright.FindStringSubmatch(footerHtml)

	var copyrightHtml string
	if len(copyrightMatch) >= 2 {
		copyrightHtml = copyrightMatch[0]
	}

	// Создаем структуру для хранения содержимого подвала
	content := map[string]interface{}{
		"widgets":   widgetsHtml,
		"copyright": copyrightHtml,
	}

	// Создаем блок
	block := &dto.Block{
		BlockType: dto.BlockTypeFooter,
		Platform:  dto.PlatformWordPress,
		Content:   content,
		HTML:      footerHtml,
	}

	s.logger.Debug("WordPress footer parsed successfully")
	return block, nil
}
