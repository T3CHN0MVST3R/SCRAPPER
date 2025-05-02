package services

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"

	"scrapper/internal/dto"
	"scrapper/internal/repos"
)

// parserService реализация ParserService
type parserService struct {
	logger           *zap.Logger
	repo             repos.ParserRepo
	wordpressService WordPressService
	tildaService     TildaService
	bitrixService    BitrixService
	html5Service     HTML5Service
}

// NewParserService создает новый экземпляр ParserService
func NewParserService(
	logger *zap.Logger,
	repo repos.ParserRepo,
	wordpressService WordPressService,
	tildaService TildaService,
	bitrixService BitrixService,
	html5Service HTML5Service,
) ParserService {
	return &parserService{
		logger:           logger,
		repo:             repo,
		wordpressService: wordpressService,
		tildaService:     tildaService,
		bitrixService:    bitrixService,
		html5Service:     html5Service,
	}
}

// ParseURL парсит URL и сохраняет результаты в базу данных
func (s *parserService) ParseURL(ctx context.Context, url string) (uuid.UUID, error) {
	// Создаем операцию в БД
	operationID, err := s.repo.CreateOperation(ctx, url)
	if err != nil {
		s.logger.Error("Failed to create operation", zap.Error(err))
		return uuid.Nil, err
	}

	// Обновляем статус операции
	err = s.repo.UpdateOperationStatus(ctx, operationID, dto.StatusProcessing)
	if err != nil {
		s.logger.Error("Failed to update operation status", zap.Error(err))
		return operationID, err
	}

	// Запускаем парсинг в отдельной горутине
	go func() {
		// Создаем новый контекст для горутины
		goCtx := context.Background()

		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath("C:/Program Files/Google/Chrome/Application/chrome.exe"), // путь к Chrome
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
		)

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		// Создаем контекст с таймаутом
		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		// Устанавливаем общий таймаут на выполнение
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var html string

		// Навигация и извлечение outer HTML
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(2*time.Second), // дать JS отработать
			chromedp.OuterHTML("html", &html),
		)

		if err != nil {
			s.logger.Error("Failed to render page with chromedp", zap.Error(err))
			s.repo.UpdateOperationStatus(goCtx, operationID, dto.StatusError)
		}

		// Определяем платформу сайта
		platform := s.DetectPlatform(html)

		// Парсим шапку и подвал в зависимости от платформы
		var headerBlock, footerBlock *dto.Block

		switch platform {
		case dto.PlatformWordPress:
			headerBlock, err = s.wordpressService.ParseHeader(html)
			if err != nil {
				s.logger.Error("Failed to parse WordPress header", zap.Error(err))
			}

			footerBlock, err = s.wordpressService.ParseFooter(html)
			if err != nil {
				s.logger.Error("Failed to parse WordPress footer", zap.Error(err))
			}
		case dto.PlatformTilda:
			headerBlock, err = s.tildaService.ParseHeader(html)
			if err != nil {
				s.logger.Error("Failed to parse Tilda header", zap.Error(err))
			}

			footerBlock, err = s.tildaService.ParseFooter(html)
			if err != nil {
				s.logger.Error("Failed to parse Tilda footer", zap.Error(err))
			}
		case dto.PlatformBitrix:
			headerBlock, err = s.bitrixService.ParseHeader(html)
			if err != nil {
				s.logger.Error("Failed to parse Bitrix header", zap.Error(err))
			}

			footerBlock, err = s.bitrixService.ParseFooter(html)
			if err != nil {
				s.logger.Error("Failed to parse Bitrix footer", zap.Error(err))
			}
		case dto.PlatformHTML5:
			var blocks []*dto.Block
			templates, err := s.repo.GetAllTemplates(platform)
			if err != nil {
				s.logger.Error("Failed to load templates:", zap.Error(err))
			}

			blocks, err = s.html5Service.ParseAndClassifyPage(html, templates)
			if err != nil {
				s.logger.Error("Failed to parse HTML5 content", zap.Error(err))
			}

			// Сохраняем найденные блоки в БД
			for _, block := range blocks {
				if block != nil {
					block.OperationID = operationID

					err = s.repo.SaveBlock(goCtx, block)
					if err != nil {
						s.logger.Error("Failed to save block", zap.Error(err), zap.Any("block", block))
					}
				}
			}
		}

		// Сохраняем найденные блоки в БД
		if headerBlock != nil {
			headerBlock.OperationID = operationID
			err = s.repo.SaveBlock(goCtx, headerBlock)
			if err != nil {
				s.logger.Error("Failed to save header block", zap.Error(err))
			}
		}

		if footerBlock != nil {
			footerBlock.OperationID = operationID
			err = s.repo.SaveBlock(goCtx, footerBlock)
			if err != nil {
				s.logger.Error("Failed to save footer block", zap.Error(err))
			}
		}

		// Обновляем статус операции
		err = s.repo.UpdateOperationStatus(goCtx, operationID, dto.StatusCompleted)
		if err != nil {
			s.logger.Error("Failed to update operation status", zap.Error(err))
		}
	}()

	return operationID, nil
}

// GetOperationResult получает результаты операции по ID
func (s *parserService) GetOperationResult(ctx context.Context, operationID uuid.UUID) (*dto.GetOperationResultResponse, error) {
	// Получаем операцию из БД
	operation, err := s.repo.GetOperationByID(ctx, operationID)
	if err != nil {
		s.logger.Error("Failed to get operation", zap.Error(err))
		return nil, err
	}

	// Получаем блоки операции
	blocks, err := s.repo.GetBlocksByOperationID(ctx, operationID)
	if err != nil {
		s.logger.Error("Failed to get blocks", zap.Error(err))
		return nil, err
	}

	// Формируем ответ
	response := &dto.GetOperationResultResponse{
		Operation: *operation,
		Blocks:    blocks,
	}

	return response, nil
}

// ExportOperation экспортирует результаты операции в файл
func (s *parserService) ExportOperation(ctx context.Context, operationID uuid.UUID, format string) ([]byte, string, error) {
	// Получаем результаты операции
	result, err := s.GetOperationResult(ctx, operationID)
	if err != nil {
		return nil, "", err
	}

	var filename string
	var content []byte

	// В зависимости от формата экспортируем результаты
	switch format {
	case "excel":
		// Создаем Excel-файл
		f := excelize.NewFile()

		// Устанавливаем заголовки для первого листа (Информация об операции)
		f.SetCellValue("Sheet1", "A1", "ID")
		f.SetCellValue("Sheet1", "B1", "URL")
		f.SetCellValue("Sheet1", "C1", "Status")
		f.SetCellValue("Sheet1", "D1", "Created At")
		f.SetCellValue("Sheet1", "E1", "Updated At")

		// Заполняем данные операции
		f.SetCellValue("Sheet1", "A2", result.Operation.ID.String())
		f.SetCellValue("Sheet1", "B2", result.Operation.URL)
		f.SetCellValue("Sheet1", "C2", result.Operation.Status)
		f.SetCellValue("Sheet1", "D2", result.Operation.CreatedAt.Format(time.RFC3339))
		f.SetCellValue("Sheet1", "E2", result.Operation.UpdatedAt.Format(time.RFC3339))

		// Создаем новый лист для блоков
		f.NewSheet("Blocks")

		// Устанавливаем заголовки для листа блоков
		f.SetCellValue("Blocks", "A1", "ID")
		f.SetCellValue("Blocks", "B1", "Type")
		f.SetCellValue("Blocks", "C1", "Platform")
		f.SetCellValue("Blocks", "D1", "Created At")
		f.SetCellValue("Blocks", "E1", "HTML")

		// Заполняем данные блоков
		for i, block := range result.Blocks {
			row := i + 2
			f.SetCellValue("Blocks", fmt.Sprintf("A%d", row), block.ID.String())
			f.SetCellValue("Blocks", fmt.Sprintf("B%d", row), block.BlockType)
			f.SetCellValue("Blocks", fmt.Sprintf("C%d", row), block.Platform)
			f.SetCellValue("Blocks", fmt.Sprintf("D%d", row), block.CreatedAt.Format(time.RFC3339))
			f.SetCellValue("Blocks", fmt.Sprintf("E%d", row), block.HTML)
		}

		// Сохраняем Excel-файл в буфер
		buffer, err := f.WriteToBuffer()
		if err != nil {
			return nil, "", fmt.Errorf("failed to write Excel file: %w", err)
		}

		content = buffer.Bytes()
		filename = fmt.Sprintf("operation_%s.xlsx", operationID.String())

	case "text":
		// Формируем текстовый отчет
		textContent := fmt.Sprintf("Operation ID: %s\n", result.Operation.ID.String())
		textContent += fmt.Sprintf("URL: %s\n", result.Operation.URL)
		textContent += fmt.Sprintf("Status: %s\n", result.Operation.Status)
		textContent += fmt.Sprintf("Created At: %s\n", result.Operation.CreatedAt.Format(time.RFC3339))
		textContent += fmt.Sprintf("Updated At: %s\n\n", result.Operation.UpdatedAt.Format(time.RFC3339))

		textContent += "Blocks:\n"
		for _, block := range result.Blocks {
			textContent += fmt.Sprintf("  ID: %s\n", block.ID.String())
			textContent += fmt.Sprintf("  Type: %s\n", block.BlockType)
			textContent += fmt.Sprintf("  Platform: %s\n", block.Platform)
			textContent += fmt.Sprintf("  Created At: %s\n", block.CreatedAt.Format(time.RFC3339))
			textContent += fmt.Sprintf("  HTML: %s\n\n", block.HTML)
		}

		content = []byte(textContent)
		filename = fmt.Sprintf("operation_%s.txt", operationID.String())

	default:
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}

	return content, filename, nil
}

// DetectPlatform определяет платформу сайта по HTML
func (s *parserService) DetectPlatform(html string) dto.Platform {

	if s.html5Service.DetectPlatform(html) {
		return dto.PlatformHTML5
	}

	if s.wordpressService.DetectPlatform(html) {
		return dto.PlatformWordPress
	}

	if s.tildaService.DetectPlatform(html) {
		return dto.PlatformTilda
	}

	if s.bitrixService.DetectPlatform(html) {
		return dto.PlatformBitrix
	}

	return dto.PlatformUnknown
}
