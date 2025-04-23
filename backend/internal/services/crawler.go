package services

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// crawlerService реализация CrawlerService
type crawlerService struct {
	userAgent      string
	maxDepth       int
	allowedDomains []string
	client         *http.Client
	visitedURLs    map[string]bool
	mutex          sync.Mutex
}

// NewCrawlerService создает новый экземпляр CrawlerService
func NewCrawlerService(allowedDomains []string) CrawlerService {
	return &crawlerService{
		userAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
		maxDepth:       2,
		allowedDomains: allowedDomains,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		visitedURLs: make(map[string]bool),
	}
}

// CrawlURL обходит URL и собирает ссылки
func (s *crawlerService) CrawlURL(ctx context.Context, url string, maxDepth int) ([]string, error) {
	// TODO: Реализовать обход URL и сбор ссылок.
	// Алгоритм:
	// 1. Проверить, разрешен ли домен для обхода
	// 2. Создать пустой список для хранения найденных URL
	// 3. Обойти сайт, начиная с указанного URL, до указанной глубины
	// 4. Собрать все уникальные URL и добавить их в список
	// 5. Вернуть список найденных URL
	//
	// Подробный алгоритм работы:
	// - Используйте рекурсивную функцию для обхода URL
	// - На каждом уровне рекурсии:
	//   1. Проверьте, не превышена ли максимальная глубина
	//   2. Проверьте, не посещался ли URL ранее
	//   3. Отправьте HTTP запрос
	//   4. Извлеките все ссылки из HTML
	//   5. Для каждой ссылки:
	//      a. Нормализуйте URL (абсолютный путь, удаление якорей)
	//      b. Проверьте, относится ли URL к разрешенному домену
	//      c. Рекурсивно обойдите URL, уменьшив максимальную глубину на 1
	//
	// Пример реализации:
	// func (s *crawlerService) crawlRecursive(ctx context.Context, urlStr string, depth int, results []string) ([]string, error) {
	//     // Проверяем глубину
	//     if depth <= 0 {
	//         return results, nil
	//     }
	//
	//     // Нормализуем URL
	//     u, err := url.Parse(urlStr)
	//     if err != nil {
	//         return results, err
	//     }
	//
	//     // Проверяем, не посещали ли URL ранее
	//     s.mutex.Lock()
	//     if s.visitedURLs[u.String()] {
	//         s.mutex.Unlock()
	//         return results, nil
	//     }
	//     s.visitedURLs[u.String()] = true
	//     s.mutex.Unlock()
	//
	//     // Отправляем HTTP запрос
	//     req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	//     if err != nil {
	//         return results, err
	//     }
	//     req.Header.Set("User-Agent", s.userAgent)
	//
	//     resp, err := s.client.Do(req)
	//     if err != nil {
	//         return results, err
	//     }
	//     defer resp.Body.Close()
	//
	//     // Проверяем статус
	//     if resp.StatusCode != http.StatusOK {
	//         return results, fmt.Errorf("status code: %d", resp.StatusCode)
	//     }
	//
	//     // Добавляем URL в результаты
	//     results = append(results, u.String())
	//
	//     // Парсим HTML и извлекаем ссылки
	//     doc, err := html.Parse(resp.Body)
	//     if err != nil {
	//         return results, err
	//     }
	//
	//     // Извлекаем ссылки и рекурсивно обходим их
	//     links := extractLinks(doc, u)
	//     for _, link := range links {
	//         if s.isAllowedDomain(link) {
	//             results, err = s.crawlRecursive(ctx, link, depth-1, results)
	//             if err != nil {
	//                 // Обрабатываем ошибку, но продолжаем обход
	//             }
	//         }
	//     }
	//
	//     return results, nil
	// }

	panic("implement me: разработайте алгоритм обхода URL и сбора ссылок")
}

// IsAllowedDomain проверяет, разрешен ли домен для обхода
func (s *crawlerService) IsAllowedDomain(urlStr string) bool {
	// TODO: Реализовать проверку разрешенного домена.
	// Алгоритм:
	// 1. Распарсить URL для получения домена
	// 2. Сравнить домен с списком разрешенных доменов
	// 3. Вернуть true, если домен разрешен, иначе false
	//
	// Пример реализации:
	// u, err := url.Parse(urlStr)
	// if err != nil {
	//     return false
	// }
	//
	// host := u.Hostname()
	//
	// // Проверяем, находится ли домен в списке разрешенных
	// for _, allowedDomain := range s.allowedDomains {
	//     if host == allowedDomain || strings.HasSuffix(host, "."+allowedDomain) {
	//         return true
	//     }
	// }
	//
	// return false

	panic("implement me: разработайте алгоритм проверки разрешенного домена")
}

// SetUserAgent устанавливает User-Agent для запросов
func (s *crawlerService) SetUserAgent(userAgent string) {
	// TODO: Реализовать установку User-Agent.
	// s.userAgent = userAgent

	panic("implement me: установите User-Agent для HTTP запросов")
}

// SetMaxDepth устанавливает максимальную глубину обхода
func (s *crawlerService) SetMaxDepth(depth int) {
	// TODO: Реализовать установку максимальной глубины.
	// if depth < 1 {
	//     s.maxDepth = 1
	// } else {
	//     s.maxDepth = depth
	// }

	panic("implement me: установите максимальную глубину обхода")
}

// Дополнительные вспомогательные функции, которые можно реализовать:

// extractLinks извлекает все ссылки из HTML документа
// func extractLinks(n *html.Node, baseURL *url.URL) []string {
//     var links []string
//
//     var traverse func(*html.Node)
//     traverse = func(n *html.Node) {
//         if n.Type == html.ElementNode && n.Data == "a" {
//             for _, attr := range n.Attr {
//                 if attr.Key == "href" {
//                     link, err := url.Parse(attr.Val)
//                     if err != nil {
//                         continue
//                     }
//
//                     // Преобразуем относительный URL в абсолютный
//                     resolvedURL := baseURL.ResolveReference(link)
//
//                     // Удаляем фрагмент (якорь)
//                     resolvedURL.Fragment = ""
//
//                     // Проверяем, что это HTTP или HTTPS URL
//                     if resolvedURL.Scheme == "http" || resolvedURL.Scheme == "https" {
//                         links = append(links, resolvedURL.String())
//                     }
//                 }
//             }
//         }
//
//         for c := n.FirstChild; c != nil; c = c.NextSibling {
//             traverse(c)
//         }
//     }
//
//     traverse(n)
//     return links
// }
