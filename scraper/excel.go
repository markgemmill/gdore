package scraper

import (
	"fmt"
	"time"

	"github.com/markgemmill/pathlib"
	"github.com/xuri/excelize/v2"
)

type ExcelDocument struct {
	path    pathlib.Path
	records []*Document
}

func (ex *ExcelDocument) AddDocument(doc *Document) {
	ex.records = append(ex.records, doc)
}

func (ex *ExcelDocument) Save() error {
	// rows := []string{}
	// for _, doc := range ex.records {
	// 	rows = append(rows, doc.Csv())
	// }
	// content := strings.Join(rows, "\n")
	//
	// err := ex.path.Write([]byte(content))
	// if err != nil {
	// 	return err
	// }
	return createExcelDoc(ex)
}

func NewExcelDocument(outputDir pathlib.Path) *ExcelDocument {
	timestamp := time.Now().Format("2006-01-02-150405")
	docPath := outputDir.Join(fmt.Sprintf("sobeys-documents-%s.xlsx", timestamp))
	doc := ExcelDocument{
		path: docPath,
	}
	return &doc
}

func cellId(letter string, number int) string {
	return fmt.Sprintf("%s%d", letter, number)
}

func createExcelDoc(doc *ExcelDocument) error {
	f := excelize.NewFile()
	defer func() {
		err := f.Close()
		if err != nil {
		}
	}()

	// write header
	colIndex := 0
	for _, header := range doc.records[0].Headers() {
		colIndex += 1
		colId, _ := excelize.ColumnNumberToName(colIndex)
		f.SetCellValue("Sheet1", cellId(colId, 1), header)
	}
	for _, header := range doc.records[0].Banners[0].Headers() {
		colIndex += 1
		colId, _ := excelize.ColumnNumberToName(colIndex)
		f.SetCellValue("Sheet1", cellId(colId, 1), header)
	}
	for _, header := range doc.records[0].Articles[0].Headers() {
		colIndex += 1
		colId, _ := excelize.ColumnNumberToName(colIndex)
		f.SetCellValue("Sheet1", cellId(colId, 1), header)
	}

	// write details
	rowIndex := 1
	for _, document := range doc.records {
		for _, row := range document.Table() {
			rowIndex += 1
			colIndex = 0
			for _, cell := range row {
				colIndex += 1
				colId, _ := excelize.ColumnNumberToName(colIndex)
				f.SetCellValue("Sheet1", cellId(colId, rowIndex), cell)
			}
		}
	}

	// save document
	err := f.SaveAs(doc.path.String())
	if err != nil {
		return err
	}

	return nil
}
