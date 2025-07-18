package services

import (
	"bytes"
	"fmt"
	"io"

	// THAY ĐỔI: Sử dụng đường dẫn import chính xác từ Github
	"baliance.com/gooxml/document"
	"github.com/ledongthuc/pdf"
)

// ExtractTextFromPDF extracts all text from a PDF file given as an io.Reader
func ExtractTextFromPDF(r io.Reader) (string, error) {
    // Giữ nguyên logic cho PDF
    data, err := io.ReadAll(r)
    if err != nil {
        return "", fmt.Errorf("failed to read PDF data: %w", err)
    }
    reader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
    if err != nil {
        return "", fmt.Errorf("failed to create PDF reader: %w", err)
    }

	var buf bytes.Buffer
	numPages := reader.NumPage()
	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		content, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %w", i, err)
		}
		buf.WriteString(content)
	}
	return buf.String(), nil
}

// ExtractTextFromDOCX extracts all text from a DOCX file given as an io.Reader
func ExtractTextFromDOCX(r io.Reader) (string, error) {
	// Logic sử dụng thư viện mới vẫn giữ nguyên
	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to read DOCX data into buffer: %w", err)
	}
	reader := bytes.NewReader(data)

	doc, err := document.Read(reader, int64(reader.Len()))
	if err != nil {
		return "", fmt.Errorf("failed to open DOCX with new library: %w", err)
	}

	var buf bytes.Buffer
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			buf.WriteString(run.Text())
		}
		buf.WriteString("\n")
	}
	return buf.String(), nil
}