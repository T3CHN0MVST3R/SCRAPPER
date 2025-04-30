package services

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"

	"scrapper/internal/dto"
)

// BitrixService интерфейс для работы с элементами Bitrix
type BitrixService interface {
	DetectPlatform(html string) bool
	ParseHeader(html string) (*dto.Block, error)
	ParseFooter(html string) (*dto.Block, error)
}

// bitrixService реализация BitrixService
type bitrixService struct {
	logger *zap.Logger
	config BitrixConfig
}

// BitrixConfig конфигурация селекторов для Bitrix
type BitrixConfig struct {
	HeaderSelectors struct {
		Container    []string `yaml:"container"`
		Logo         []string `yaml:"logo"`
		Menu         []string `yaml:"menu"`
		Search       []string `yaml:"search"`
		Phones       []string `yaml:"phones"`
		Cart         []string `yaml:"cart"`
		Auth         []string `yaml:"auth"`
	} `yaml:"header"`

	FooterSelectors struct {
		Container    []string `yaml:"container"`
		Copyright    []string `yaml:"copyright"`
		Menu         []string `yaml:"menu"`
		Contacts     []string `yaml:"contacts"`
		Social       []string `yaml:"social"`
		Developer    []string `yaml:"developer"`
	} `yaml:"footer"`
}

// NewBitrixService создает новый экземпляр BitrixService
func NewBitrixService(logger *zap.Logger, config BitrixConfig) BitrixService {
	return &bitrixService{
		logger: logger,
		config: config,
	}
}

// DetectPlatform проверяет, соответствует ли страница Bitrix
func (s *bitrixService) DetectPlatform(html string) bool {
	bitrixPatterns := []string{
		`bitrix/js`,
		`bitrix/templates`,
		`<meta name="generator" content="Bitrix`,
		`BX\.`,
		`b24-widget`,
		`class="bx-`,
		`id="bx_`,
		`<!-- Bitrix`,
		`1C-Bitrix`,
		`/bitrix/`,
	}

	for _, pattern := range bitrixPatterns {
		if strings.Contains(html, pattern) {
			s.logger.Debug("Bitrix pattern detected", 
				zap.String("pattern", pattern))
			return true
		}
	}

	return false
}

// detectBitrixVersion определяет версию Bitrix
func (s *bitrixService) detectBitrixVersion(html string) string {
	versionPatterns := map[string]string{
		"modern": `BX24\.|b24-`,
		"legacy": `bitrix\/js\/main\/core\/core`,
		"old":    `bitrix\/components\/bitrix`,
	}

	for version, pattern := range versionPatterns {
		if matched, _ := regexp.MatchString(pattern, html); matched {
			s.logger.Debug("Detected Bitrix version", 
				zap.String("version", version))
			return version
		}
	}

	return "unknown"
}

// ParseHeader парсит шапку сайта Bitrix
func (s *bitrixService) ParseHeader(html string) (*dto.Block, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		s.logger.Error("Failed to parse HTML", zap.Error(err))
		return nil, err
	}

	header := &dto.Block{
		Type:       "bitrix_header",
		Components: make(map[string]interface{}),
		Version:    s.detectBitrixVersion(html),
	}

	// Поиск контейнера шапки
	var headerContainer *goquery.Selection
	for _, selector := range s.config.HeaderSelectors.Container {
		if doc.Find(selector).Length() > 0 {
			headerContainer = doc.Find(selector).First()
			break
		}
	}

	if headerContainer == nil {
		s.logger.Warn("No header containers found")
		return header, nil
	}

	// Парсинг логотипа
	s.parseLogo(headerContainer, header)

	// Парсинг меню
	s.parseMenu(headerContainer, header)

	// Парсинг поиска
	s.parseSearch(headerContainer, header)

	// Парсинг телефонов
	s.parsePhones(headerContainer, header)

	// Парсинг корзины
	s.parseCart(headerContainer, header)

	// Парсинг авторизации
	s.parseAuth(headerContainer, header)

	// Валидация результата
	if !s.validateHeader(header) {
		s.logger.Warn("Header validation failed")
	}

	return header, nil
}

// parseLogo парсит логотип
func (s *bitrixService) parseLogo(container *goquery.Selection, header *dto.Block) {
	for _, selector := range s.config.HeaderSelectors.Logo {
		logo := container.Find(selector).First()
		if logo.Length() > 0 {
			if src, exists := logo.Find("img").Attr("src"); exists {
				header.Components["logo"] = src
				break
			} else if href, exists := logo.Attr("href"); exists {
				header.Components["logo_link"] = href
				break
			}
		}
	}
}

// parseMenu парсит меню
func (s *bitrixService) parseMenu(container *goquery.Selection, header *dto.Block) {
	var menuItems []string
	for _, selector := range s.config.HeaderSelectors.Menu {
		container.Find(selector).Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(i int, item *goquery.Selection) {
				if text := strings.TrimSpace(item.Text()); text != "" {
					menuItems = append(menuItems, text)
				}
			})
		})
		if len(menuItems) > 0 {
			header.Components["menu"] = menuItems
			break
		}
	}
}

// parseSearch парсит поиск
func (s *bitrixService) parseSearch(container *goquery.Selection, header *dto.Block) {
	for _, selector := range s.config.HeaderSelectors.Search {
		if container.Find(selector).Length() > 0 {
			header.Components["search"] = true
			break
		}
	}
}

// parsePhones парсит телефоны
func (s *bitrixService) parsePhones(container *goquery.Selection, header *dto.Block) {
	var phones []string
	for _, selector := range s.config.HeaderSelectors.Phones {
		container.Find(selector).Each(func(i int, s *goquery.Selection) {
			phone := strings.TrimSpace(s.Text())
			if s.isValidPhone(phone) {
				phones = append(phones, phone)
			}
		})
		if len(phones) > 0 {
			header.Components["phones"] = phones
			break
		}
	}
}

// parseCart парсит корзину
func (s *bitrixService) parseCart(container *goquery.Selection, header *dto.Block) {
	for _, selector := range s.config.HeaderSelectors.Cart {
		if container.Find(selector).Length() > 0 {
			header.Components["cart"] = true
			break
		}
	}
}

// parseAuth парсит авторизацию
func (s *bitrixService) parseAuth(container *goquery.Selection, header *dto.Block) {
	for _, selector := range s.config.HeaderSelectors.Auth {
		if container.Find(selector).Length() > 0 {
			header.Components["auth"] = true
			break
		}
	}
}

// ParseFooter парсит подвал сайта Bitrix
func (s *bitrixService) ParseFooter(html string) (*dto.Block, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		s.logger.Error("Failed to parse HTML", zap.Error(err))
		return nil, err
	}

	footer := &dto.Block{
		Type:       "bitrix_footer",
		Components: make(map[string]interface{}),
		Version:    s.detectBitrixVersion(html),
	}

	// Поиск контейнера подвала
	var footerContainer *goquery.Selection
	for _, selector := range s.config.FooterSelectors.Container {
		if doc.Find(selector).Length() > 0 {
			footerContainer = doc.Find(selector).First()
			break
		}
	}

	if footerContainer == nil {
		s.logger.Warn("No footer containers found")
		return footer, nil
	}

	// Парсинг копирайта
	s.parseCopyright(footerContainer, footer)

	// Парсинг меню
	s.parseFooterMenu(footerContainer, footer)

	// Парсинг контактов
	s.parseContacts(footerContainer, footer)

	// Парсинг соцсетей
	s.parseSocial(footerContainer, footer)

	// Парсинг разработчика
	s.parseDeveloper(footerContainer, footer)

	// Валидация результата
	if !s.validateFooter(footer) {
		s.logger.Warn("Footer validation failed")
	}

	return footer, nil
}

// parseCopyright парсит копирайт
func (s *bitrixService) parseCopyright(container *goquery.Selection, footer *dto.Block) {
	for _, selector := range s.config.FooterSelectors.Copyright {
		if elem := container.Find(selector).First(); elem.Length() > 0 {
			footer.Components["copyright"] = strings.TrimSpace(elem.Text())
			break
		}
	}
}

// parseFooterMenu парсит меню в подвале
func (s *bitrixService) parseFooterMenu(container *goquery.Selection, footer *dto.Block) {
	var menuItems []string
	for _, selector := range s.config.FooterSelectors.Menu {
		container.Find(selector).Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(i int, item *goquery.Selection) {
				if text := strings.TrimSpace(item.Text()); text != "" {
					menuItems = append(menuItems, text)
				}
			})
		})
		if len(menuItems) > 0 {
			footer.Components["menu"] = menuItems
			break
		}
	}
}

// parseContacts парсит контакты
func (s *bitrixService) parseContacts(container *goquery.Selection, footer *dto.Block) {
	var contacts []string
	for _, selector := range s.config.FooterSelectors.Contacts {
		container.Find(selector).Each(func(i int, s *goquery.Selection) {
			s.Find("p, div").Each(func(i int, item *goquery.Selection) {
				text := strings.TrimSpace(item.Text())
				if text != "" {
					contacts = append(contacts, text)
				}
			})
		})
		if len(contacts) > 0 {
			footer.Components["contacts"] = contacts
			break
		}
	}
}

// parseSocial парсит соцсети
func (s *bitrixService) parseSocial(container *goquery.Selection, footer *dto.Block) {
	var socialLinks []string
	for _, selector := range s.config.FooterSelectors.Social {
		container.Find(selector).Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(i int, item *goquery.Selection) {
				if href, exists := item.Attr("href"); exists {
					socialLinks = append(socialLinks, href)
				}
			})
		})
		if len(socialLinks) > 0 {
			footer.Components["social"] = socialLinks
			break
		}
	}
}

// parseDeveloper парсит ссылку на разработчика
func (s *bitrixService) parseDeveloper(container *goquery.Selection, footer *dto.Block) {
	for _, selector := range s.config.FooterSelectors.Developer {
		if elem := container.Find(selector).First(); elem.Length() > 0 {
			if href, exists := elem.Attr("href"); exists {
				footer.Components["developer"] = href
				break
			}
		}
	}
}

// validateHeader проверяет валидность шапки
func (s *bitrixService) validateHeader(header *dto.Block) bool {
	required := []string{"logo", "menu"}
	for _, field := range required {
		if _, exists := header.Components[field]; !exists {
			s.logger.Warn("Header validation failed - missing field", 
				zap.String("field", field))
			return false
		}
	}
	return true
}

// validateFooter проверяет валидность подвала
func (s *bitrixService) validateFooter(footer *dto.Block) bool {
	if _, exists := footer.Components["copyright"]; !exists {
		s.logger.Warn("Footer validation failed - missing copyright")
		return false
	}
	return true
}

// isValidPhone проверяет валидность телефонного номера
func (s *bitrixService) isValidPhone(phone string) bool {
	cleaned := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")
	return len(cleaned) >= 5 && regexp.MustCompile(`[\d]{5,}`).MatchString(cleaned)
}
