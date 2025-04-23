package services

import (
	"regexp"
	"strings"

	"scrapper/internal/dto"
)

// wordPressService реализация WordPressService
type wordPressService struct {
}

// NewWordPressService создает новый экземпляр WordPressService
func NewWordPressService() WordPressService {
	return &wordPressService{}
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
		`/wp-admin/`,
		`/wp-login.php`,
	}

	for _, pattern := range wpPatterns {
		if strings.Contains(html, pattern) {
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
		return nil, nil
	}

	headerHtml := headerMatch[0]

	// Парсим логотип
	reLogo := regexp.MustCompile(`(?s)<a.*?class=".*?logo.*?".*?>(.*?)</a>`)
	logoMatch := reLogo.FindStringSubmatch(headerHtml)

	// Если не найден логотип по классу logo, пробуем найти по alt или title с упоминанием logo
	if len(logoMatch) < 2 {
		reLogo = regexp.MustCompile(`(?s)<a.*?><img.*?(?:alt|title)=".*?logo.*?".*?></a>`)
		logoMatch = reLogo.FindStringSubmatch(headerHtml)
	}

	// Если все еще не найден, ищем просто первое изображение в шапке
	if len(logoMatch) < 2 {
		reLogo = regexp.MustCompile(`(?s)<a.*?><img.*?src=".*?".*?></a>`)
		logoMatch = reLogo.FindStringSubmatch(headerHtml)
	}

	var logoHtml string
	if len(logoMatch) >= 1 {
		logoHtml = logoMatch[0]
	}

	// Парсим навигационное меню
	reMenu := regexp.MustCompile(`(?s)<nav.*?>(.*?)</nav>`)
	menuMatch := reMenu.FindStringSubmatch(headerHtml)

	// Если не найдено <nav>, ищем по классу 'menu' или 'navigation'
	if len(menuMatch) < 2 {
		reMenu = regexp.MustCompile(`(?s)<(?:div|ul)\s+class=".*?(?:menu|navigation).*?".*?>(.*?)</(?:div|ul)>`)
		menuMatch = reMenu.FindStringSubmatch(headerHtml)
	}

	var menuHtml string
	if len(menuMatch) >= 1 {
		menuHtml = menuMatch[0]
	}

	// Парсим контактную информацию (если есть)
	reContacts := regexp.MustCompile(`(?s)<div\s+class=".*?(?:contact|phone|email|address).*?".*?>(.*?)</div>`)
	contactsMatch := reContacts.FindStringSubmatch(headerHtml)

	var contactsHtml string
	if len(contactsMatch) >= 1 {
		contactsHtml = contactsMatch[0]
	}

	// Парсим поиск (если есть)
	reSearch := regexp.MustCompile(`(?s)<form.*?(?:class=".*?search.*?".*?|id=".*?search.*?".*?)>(.*?)</form>`)
	searchMatch := reSearch.FindStringSubmatch(headerHtml)

	var searchHtml string
	if len(searchMatch) >= 1 {
		searchHtml = searchMatch[0]
	}

	// Создаем структуру для хранения содержимого шапки
	content := map[string]interface{}{
		"logo":     logoHtml,
		"menu":     menuHtml,
		"contacts": contactsHtml,
		"search":   searchHtml,
	}

	// Создаем блок
	block := &dto.Block{
		BlockType: dto.BlockTypeHeader,
		Platform:  dto.PlatformWordPress,
		Content:   content,
		HTML:      headerHtml,
	}

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
		return nil, nil
	}

	footerHtml := footerMatch[0]

	// Парсим виджеты
	reWidgets := regexp.MustCompile(`(?s)<div\s+class=".*?widgets.*?".*?>(.*?)</div>`)
	widgetsMatch := reWidgets.FindStringSubmatch(footerHtml)

	// Если не найдены виджеты по классу widgets, ищем по классу sidebar
	if len(widgetsMatch) < 2 {
		reWidgets = regexp.MustCompile(`(?s)<div\s+class=".*?sidebar.*?".*?>(.*?)</div>`)
		widgetsMatch = reWidgets.FindStringSubmatch(footerHtml)
	}

	var widgetsHtml string
	if len(widgetsMatch) >= 1 {
		widgetsHtml = widgetsMatch[0]
	}

	// Парсим копирайт
	reCopyright := regexp.MustCompile(`(?s)<div\s+class=".*?copyright.*?".*?>(.*?)</div>`)
	copyrightMatch := reCopyright.FindStringSubmatch(footerHtml)

	// Если не найден копирайт по классу, ищем по тексту
	if len(copyrightMatch) < 2 {
		reCopyright = regexp.MustCompile(`(?s)(?:©|&copy;).*?\d{4}`)
		copyrightMatch = reCopyright.FindStringSubmatch(footerHtml)
	}

	var copyrightHtml string
	if len(copyrightMatch) >= 1 {
		copyrightHtml = copyrightMatch[0]
	}

	// Парсим социальные сети
	reSocial := regexp.MustCompile(`(?s)<div\s+class=".*?social.*?".*?>(.*?)</div>`)
	socialMatch := reSocial.FindStringSubmatch(footerHtml)

	// Если не найдены соцсети по классу social, ищем по иконкам
	if len(socialMatch) < 2 {
		reSocial = regexp.MustCompile(`(?s)<ul\s+class=".*?(?:social|socials).*?".*?>(.*?)</ul>`)
		socialMatch = reSocial.FindStringSubmatch(footerHtml)
	}

	var socialHtml string
	if len(socialMatch) >= 1 {
		socialHtml = socialMatch[0]
	}

	// Парсим ссылки навигации в подвале
	reNavLinks := regexp.MustCompile(`(?s)<nav\s+class=".*?footer-navigation.*?".*?>(.*?)</nav>`)
	navLinksMatch := reNavLinks.FindStringSubmatch(footerHtml)

	// Если не найдены ссылки навигации по классу, ищем по тегам ul/li в подвале
	if len(navLinksMatch) < 2 {
		reNavLinks = regexp.MustCompile(`(?s)<ul\s+class=".*?footer-menu.*?".*?>(.*?)</ul>`)
		navLinksMatch = reNavLinks.FindStringSubmatch(footerHtml)
	}

	var navLinksHtml string
	if len(navLinksMatch) >= 1 {
		navLinksHtml = navLinksMatch[0]
	}

	// Создаем структуру для хранения содержимого подвала
	content := map[string]interface{}{
		"widgets":   widgetsHtml,
		"copyright": copyrightHtml,
		"social":    socialHtml,
		"nav_links": navLinksHtml,
	}

	// Создаем блок
	block := &dto.Block{
		BlockType: dto.BlockTypeFooter,
		Platform:  dto.PlatformWordPress,
		Content:   content,
		HTML:      footerHtml,
	}

	return block, nil
}
