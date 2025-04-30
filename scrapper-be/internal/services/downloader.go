package services

import (
	"context"
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// downloaderService реализация DownloaderService
type downloaderService struct {
	logger        *zap.Logger
	parserService ParserService
}

// NewDownloaderService создает новый экземпляр DownloaderService
func NewDownloaderService(logger *zap.Logger, parserService ParserService) DownloaderService {
	return &downloaderService{
		logger:        logger,
		parserService: parserService,
	}
}

// DownloadByOperationID загружает файлы по ID операции и сохраняет их в указанный путь
func (s *downloaderService) DownloadByOperationID(ctx context.Context, operationID uuid.UUID, path string) error {
	panic("implement me")
}

// GetAvailableFormats возвращает список доступных форматов для загрузки
func (s *downloaderService) GetAvailableFormats() []string {
	return []string{"excel", "text"}
}

// DownloadByOperationIDWithFormat загружает файлы в указанном формате
func (s *downloaderService) DownloadByOperationIDWithFormat(ctx context.Context, operationID uuid.UUID, format string, path string) error {
	panic("implement me")
}

// DownloadByOperationID загружает файлы по ID операции
func (s *downloaderService) DownloadByOperationID(ctx context.Context, operationID uuid.UUID, format string) ([]byte, string, error) {
	// Получаем результаты операции через парсер-сервис
	result, err := s.parserService.GetOperationResult(ctx, operationID)
	if (err != nil) {
		return nil, "", err
	}

	var buffer bytes.Buffer
	var filename string

	switch format {
	case "excel":
		data, name, err := s.generateExcel(result)
		if err != nil {
			return nil, "", err
		}
		buffer.Write(data)
		filename = name

	case "pdf":
		data, name, err := s.generatePDF(result)
		if err != nil {
			return nil, "", err
		}
		buffer.Write(data)
		filename = name

	case "text":
		data, name, err := s.generateText(result)
		if err != nil {
			return nil, "", err
		}
		buffer.Write(data)
		filename = name

	default:
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}

	return buffer.Bytes(), filename, nil
}

// generatePDF создает PDF-документ с результатами
func (s *downloaderService) generatePDF(result *dto.GetOperationResultResponse) ([]byte, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	// Заголовок
	pdf.Cell(190, 10, fmt.Sprintf("Operation Results (ID: %s)", result.Operation.ID))
	pdf.Ln(10)

	// Информация об операции
	pdf.Cell(190, 10, fmt.Sprintf("URL: %s", result.Operation.URL))
	pdf.Ln(10)
	pdf.Cell(190, 10, fmt.Sprintf("Status: %s", result.Operation.Status))
	pdf.Ln(10)

	// Блоки
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 10, "Blocks:")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	for _, block := range result.Blocks {
		pdf.Cell(190, 10, fmt.Sprintf("Type: %s, Platform: %s", block.BlockType, block.Platform))
		pdf.Ln(10)
		pdf.MultiCell(190, 10, block.HTML, "", "", false)
		pdf.Ln(5)
	}

	// Используем буферизированный вывод
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, "", fmt.Errorf("failed to write PDF: %w", err)
	}

	filename := fmt.Sprintf("operation_%s.pdf", result.Operation.ID)
	return buf.Bytes(), filename, nil
}

func (s *downloaderService) generateText(result *dto.GetOperationResultResponse) ([]byte, string, error) {
	var buffer bytes.Buffer

	// Пишем информацию об операции
	buffer.WriteString(fmt.Sprintf("Operation ID: %s\n", result.Operation.ID))
	buffer.WriteString(fmt.Sprintf("URL: %s\n", result.Operation.URL))
	buffer.WriteString(fmt.Sprintf("Status: %s\n", result.Operation.Status))
	buffer.WriteString(fmt.Sprintf("Created: %s\n\n", result.Operation.CreatedAt))

	// Пишем информацию о блоках
	buffer.WriteString("Blocks:\n")
	for _, block := range result.Blocks {
		buffer.WriteString(fmt.Sprintf("\nType: %s\n", block.BlockType))
		buffer.WriteString(fmt.Sprintf("Platform: %s\n", block.Platform))
		buffer.WriteString("Content:\n")
		buffer.WriteString(block.HTML)
		buffer.WriteString("\n---\n")
	}

	filename := fmt.Sprintf("operation_%s.txt", result.Operation.ID)
	return buffer.Bytes(), filename, nil
}
