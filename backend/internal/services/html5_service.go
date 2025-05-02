package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"scrapper/internal/dto"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// html5Service реализация HTML5Service
type html5Service struct {
}

// NewHTML5Service создает новый экземпляр HTML5Service
func NewHTML5Service() HTML5Service {
	return &html5Service{}
}

// DetectPlatform проверяет, соответствует ли страница HTML5
func (s *html5Service) DetectPlatform(html string) bool {
	// Проверяем наличие характерных признаков WordPress
	htmlPatterns := []string{
		`<!DOCTYPE html>`,
		`<meta charset="UTF-8">`,
		`<html lang`,
	}

	for _, pattern := range htmlPatterns {
		if strings.Contains(html, pattern) {
			return true
		}
	}

	return false
}

// ParseHeader парсит шапку сайта Tilda
func (s *html5Service) ParseHeader(html string) (*dto.Block, error) {
	panic("implement me")
}

// ParseFooter парсит подвал сайта Tilda
func (s *html5Service) ParseFooter(html string) (*dto.Block, error) {
	panic("implement me")
}

func (s *html5Service) ParseAndClassifyPage(html string, templates []dto.BlockTemplate) ([]*dto.Block, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var blocks []*dto.Block

	// Найти header
	header := doc.Find("header").First()
	if header.Length() != 0 {
		blocks = append(blocks, &dto.Block{
			BlockType: dto.BlockTypeHeader,
			Platform:  dto.PlatformHTML5,
			Content: map[string]string{
				"template_name": "Шапка",
			},
			HTML: "",
		})
	}
	// Идем сверху вниз по логике документа
	for sibling := header.Next(); sibling.Length() > 0; sibling = sibling.Next() {
		node := sibling.Get(0)
		if node.Data == "footer" {
			blocks = append(blocks, &dto.Block{
				BlockType: dto.BlockTypeFooter,
				Platform:  dto.PlatformHTML5,
				Content: map[string]string{
					"template_name": "Футер",
				},
				HTML: "",
			})
			break // Достигли footer — останавливаемся
		}
		if node.Data == "section" || node.Data == "div" {
			var out bytes.Buffer
			walkVisible(sibling, &out)

			htmlContent := out.String()
			//htmlContent := out.String()

			matchedTemplate := matchBlock(htmlContent, templates)
			blockType := dto.BlockTypeContent // по умолчанию

			content := map[string]interface{}{}
			if matchedTemplate != nil {
				content["template_name"] = matchedTemplate.BlockType
			}

			blocks = append(blocks, &dto.Block{
				BlockType: blockType,
				Platform:  dto.PlatformHTML5,
				Content:   content,
				HTML:      "",
			})
		}
	}

	return blocks, nil
}

func matchBlock(html string, templates []dto.BlockTemplate) *dto.BlockTemplate {
	for _, template := range templates {
		var tagSequence map[string]interface{}
		if err := json.Unmarshal([]byte(template.HTMLTags), &tagSequence); err != nil {
			continue
		}

		matched := true

		for i := 1; ; i++ {
			key := fmt.Sprintf("step%d", i)
			raw, ok := tagSequence[key]
			if !ok {
				break
			}

			switch val := raw.(type) {
			case string:
				// Поддержка "ИЛИ"
				if strings.Contains(val, "|") {
					found := false
					for _, option := range strings.Split(val, "|") {
						if strings.Contains(html, option) {
							found = true
							break
						}
					}
					if !found {
						matched = false
						break
					}
				} else {
					if !strings.Contains(html, val) {
						matched = false
						break
					}
				}
			case []interface{}:
				found := false
				for _, item := range val {
					if tag, ok := item.(string); ok && strings.Contains(html, tag) {
						found = true
						break
					}
				}
				if !found {
					matched = false
					break
				}
			}
		}

		if matched {
			return &template
		}
	}

	return nil
}

// Обходит DOM, собирает HTML без скрытых элементов
func walkVisible(sel *goquery.Selection, buf *bytes.Buffer) {

	sel.Find("div").Each(func(i int, s *goquery.Selection) {
		// Пропускаем скрытые элементы и их родителей
		if isHidden(s) {
			s.Remove()
		}
	})

	html.Render(buf, sel.Get(0))

}

func isHidden(sel *goquery.Selection) bool {
	//html, _ := sel.Html()
	html, _ := sel.Html()
	if strings.Contains(html, "visibility:hidden") || strings.Contains(html, "visibility: hidden") {
		return true
	}

	// Проверка родителей
	hidden := false
	return hidden
}
