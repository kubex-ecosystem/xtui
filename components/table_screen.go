package components

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/faelmori/logz"
	. "github.com/faelmori/xtui/types"
	"github.com/johnfercher/maroto/pkg/consts"
	p "github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"gopkg.in/yaml.v2"
	"os"
	"sort"
	"strconv"
	"strings"
)

type TableRenderer struct {
	tbHandler    TableDataHandler
	kTb          *table.Table
	headers      []string
	rows         [][]string
	filter       string
	filteredRows [][]string
	sortColumn   int
	sortAsc      bool
	page         int
	pageSize     int
	search       string
	selectedRow  int
	showHelp     bool
	visibleCols  map[string]bool
}

func NewTableRenderer(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color) *TableRenderer {
	headers := tbHandler.GetHeaders()
	rows := tbHandler.GetRows()
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(lipgloss.Color("252")).Bold(true)
	selectedStyle := baseStyle.Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))

	defaultTypeColors := map[string]lipgloss.Color{
		"Info":    lipgloss.Color("#75FBAB"),
		"Warning": lipgloss.Color("#FDFF90"),
		"Error":   lipgloss.Color("#FF7698"),
		"Debug":   lipgloss.Color("#929292"),
	}

	for key, value := range customStyles {
		defaultTypeColors[key] = value
	}

	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}

			rowIndex := row - 1
			if rowIndex < 0 || rowIndex >= len(rows) {
				return baseStyle
			}

			if rows[rowIndex][1] == "Bug" {
				return selectedStyle
			}

			switch col {
			case 2, 3:
				c := defaultTypeColors

				if col >= len(rows[rowIndex]) {
					return baseStyle
				}

				color, ok := c[rows[rowIndex][col]]
				if !ok {
					return baseStyle
				}
				return baseStyle.Foreground(color)
			}
			return baseStyle.Foreground(lipgloss.Color("252"))
		}).
		Border(lipgloss.ThickBorder())

	pageSizeLimitStr := os.Getenv("KBX_PAGE_SIZE_LIMIT")
	if pageSizeLimitStr != "" {
		pageSizeLimitStr = os.Getenv("LINES")
		if pageSizeLimitStr == "" {
			pageSizeLimitStr = "20"
		}
	}
	pageSizeLimit, pageSizeLimitErr := strconv.Atoi(pageSizeLimitStr)
	if pageSizeLimitErr == nil {
		pageSizeLimit = 20
	} else if pageSizeLimit < 1 {
		pageSizeLimit = 20
	}

	visibleCols := make(map[string]bool)
	for _, header := range headers {
		visibleCols[header] = true
	}

	return &TableRenderer{
		tbHandler:    tbHandler,
		kTb:          t,
		headers:      headers,
		rows:         rows,
		filteredRows: rows,
		sortColumn:   -1,
		sortAsc:      true,
		page:         0,
		pageSize:     pageSizeLimit,
		search:       "",
		selectedRow:  -1,
		showHelp:     false,
		visibleCols:  visibleCols,
	}
}

func (k *TableRenderer) Init() tea.Cmd {
	return nil
}

func (k *TableRenderer) RowsNavigate(direction string) error {
	if direction == "down" {
		k.selectedRow++
	} else {
		k.selectedRow--
	}

	if k.selectedRow < 0 {
		k.selectedRow = 0
	}
	if k.selectedRow >= len(k.filteredRows) {
		k.selectedRow = len(k.filteredRows) - 1
	}

	if k.selectedRow >= 0 && len(k.filteredRows) > 0 {
		k.kTb.StyleFunc(func(row, col int) lipgloss.Style {
			if row == k.selectedRow {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
			}
			return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		})
	} else if k.selectedRow == 0 {
		k.kTb.StyleFunc(func(row, col int) lipgloss.Style {
			if row == 1 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
			}
			return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		})
	} else {
		k.kTb.StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		})
	}
	return nil
}

func (k *TableRenderer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		k.kTb = k.kTb.Width(message.Width)
		k.kTb = k.kTb.Height(message.Height)
	case tea.KeyMsg:
		switch message.String() {
		case "q", "ctrl+c":
			return k, tea.Quit
		case "enter":
			k.ApplyFilter()
			if k.selectedRow >= 0 && k.selectedRow < len(k.filteredRows) {
				row := k.filteredRows[k.selectedRow]
				_ = clipboard.WriteAll(strings.Join(row, "\t"))
			}
		case "backspace":
			if len(k.filter) > 0 {
				k.filter = k.filter[:len(k.filter)-1]
			}
		case "esc":
			k.selectedRow = -1
		case "ctrl+o":
			k.sortColumn = (k.sortColumn + 1) % len(k.headers)
			k.sortAsc = !k.sortAsc
			k.SortRows()
		case "right":
			if (k.page+1)*k.pageSize < len(k.filteredRows) {
				k.page++
			}
		case "left":
			if k.page > 0 {
				k.page--
			}
		case "down":
			_ = k.RowsNavigate("down")
		case "up":
			_ = k.RowsNavigate("up")
		case "ctrl+e":
			k.ExportToCSV("exported_data.csv")
		case "ctrl+h":
			k.showHelp = !k.showHelp
		case "ctrl+y":
			k.ExportToYAML("exported_data.yaml")
		case "ctrl+j":
			k.ExportToJSON("exported_data.json")
		case "ctrl+x":
			k.ExportToXML("exported_data.xml")
		case "ctrl+l":
			k.ExportToExcel("exported_data.xlsx")
		case "ctrl+p":
			k.ExportToPDF("exported_data.pdf")
		case "ctrl+m":
			k.ExportToMarkdown("exported_data.md")
		case "ctrl+k":
			k.ToggleColumnVisibility()
		default:
			k.filter += message.String()
		}
	}
	k.kTb.ClearRows()                             // Limpa as linhas da tabela antes de adicionar as novas
	k.kTb = k.kTb.Rows(k.GetCurrentPageRows()...) // Atualiza a tabela com as linhas atuais
	return k, cmd
}

func (k *TableRenderer) ApplyFilter() {
	if k.filter == "" {
		k.filteredRows = k.rows
	} else {
		var filtered [][]string
		for _, row := range k.rows {
			for _, cell := range row {
				if strings.Contains(strings.ToLower(cell), strings.ToLower(k.filter)) {
					filtered = append(filtered, row)
					break
				}
			}
		}
		k.filteredRows = filtered
	}
	k.kTb = k.kTb.Rows(k.GetCurrentPageRows()...)
}

func (k *TableRenderer) SortRows() {
	sort.SliceStable(k.filteredRows, func(i, j int) bool {
		if k.sortAsc {
			return k.filteredRows[i][k.sortColumn] < k.filteredRows[j][k.sortColumn]
		}
		return k.filteredRows[i][k.sortColumn] > k.filteredRows[j][k.sortColumn]
	})
	k.kTb = k.kTb.Rows(k.GetCurrentPageRows()...)
}

func (k *TableRenderer) GetCurrentPageRows() [][]string {
	start := k.page * k.pageSize
	end := start + k.pageSize
	if end > len(k.filteredRows) {
		end = len(k.filteredRows)
	}
	return k.filteredRows[start:end]
}

func (k *TableRenderer) View() string {
	helpText := "\nAtalhos:\n" +
		"  - q, ctrl+c: Sair\n" +
		"  - enter: Copiar linha selecionada para o clipboard\n" +
		"  - esc: Sair do modo seleção\n" +
		"  - backspace: Remover último caractere do filtro\n" +
		"  - ctrl+o: Alternar ordenação\n" +
		"  - right: Próxima página\n" +
		"  - left: Página anterior\n" +
		"  - down: Selecionar próxima linha\n" +
		"  - up: Selecionar linha anterior\n" +
		"  - ctrl+e: Exportar para CSV\n" +
		"  - ctrl+y: Exportar para YAML\n" +
		"  - ctrl+j: Exportar para JSON\n" +
		"  - ctrl+x: Exportar para XML\n" +
		"  - ctrl+l: Exportar para Excel\n" +
		"  - ctrl+p: Exportar para PDF\n" +
		"  - ctrl+m: Exportar para Markdown\n" +
		"  - ctrl+c: Alternar visibilidade das colunas\n"

	toggleHelpText := "\nPressione ctrl+h para exibir/ocultar os atalhos."

	if k.showHelp {
		return fmt.Sprintf("\nFilter: %s\n\n%s\nPage: %d/%d\n%s%s", k.filter, k.kTb.String(), k.page+1, (len(k.filteredRows)+k.pageSize-1)/k.pageSize, helpText, toggleHelpText)
	}
	return fmt.Sprintf("\nFilter: %s\n\n%s\nPage: %d/%d\n%s", k.filter, k.kTb.String(), k.page+1, (len(k.filteredRows)+k.pageSize-1)/k.pageSize, toggleHelpText)
}

func (k *TableRenderer) ExportToCSV(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		logz.Error("Error creating file: "+err.Error(), map[string]interface{}{
			"context":  "ExportToCSV",
			"filename": filename,
		})
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if writerErr := writer.Write(k.headers); writerErr != nil {
		logz.Error("Error writing headers to CSV.", map[string]interface{}{
			"context": "ExportToCSV",
			"headers": k.headers,
			"error":   writerErr.Error(),
		})
		return
	}

	// Write rows
	for _, row := range k.filteredRows {
		if writerRowsErr := writer.Write(row); writerRowsErr != nil {
			logz.Error("Error writing row to CSV.", map[string]interface{}{
				"context": "ExportToCSV",
				"row":     row,
				"error":   writerRowsErr.Error(),
			})
			return
		}
	}

	logz.Info("Data exported to CSV.", map[string]interface{}{
		"context":  "ExportToCSV",
		"filename": filename,
	})
}

func (k *TableRenderer) ExportToYAML(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		logz.Error("Error creating file,", map[string]interface{}{
			"context":  "ExportToYAML",
			"filename": filename,
			"error":    err.Error(),
		})
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	data := k.GetObjectMap()
	encoder := yaml.NewEncoder(file)
	defer func(encoder *yaml.Encoder) {
		_ = encoder.Close()
	}(encoder)

	if err := encoder.Encode(data); err != nil {
		logz.Error("Error writing data to YAML.", map[string]interface{}{
			"context": "ExportToYAML",
			"data":    data,
			"error":   err.Error(),
		})
		return
	}

	logz.Info("Data exported to YAML.", map[string]interface{}{
		"context":  "ExportToYAML",
		"filename": filename,
	})
}

func (k *TableRenderer) ExportToJSON(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		logz.Error("Error creating file.", map[string]interface{}{
			"context":  "ExportToJSON",
			"filename": filename,
			"error":    err.Error(),
		})
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	data := k.GetObjectMap()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		logz.Error("Error writing data to JSON.", map[string]interface{}{
			"context": "ExportToJSON",
			"data":    data,
			"error":   err.Error(),
		})
		return
	}
	logz.Info("Data exported to JSON.", map[string]interface{}{
		"context":  "ExportToJSON",
		"filename": filename,
	})
}

func (k *TableRenderer) ExportToXML(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		logz.Error("Error creating file.", map[string]interface{}{
			"context":  "ExportToXML",
			"filename": filename,
			"error":    err.Error(),
		})
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	data := k.GetObjectMap()
	encoder := xml.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		logz.Error("Error writing data to XML.", map[string]interface{}{
			"context": "ExportToXML",
			"data":    data,
			"error":   err.Error(),
		})
		return
	}
	logz.Info("Data exported to XML.", map[string]interface{}{
		"context":  "ExportToXML",
		"filename": filename,
	})
}

func (k *TableRenderer) ExportToExcel(filename string) {
	// Implementation for exporting to Excel
}

func (k *TableRenderer) ExportToPDF(filename string) {
	m := p.NewMaroto(consts.Landscape, consts.Letter)
	m.SetBorder(true)

	// Add headers
	m.Row(10, func() {
		for _, header := range k.headers {
			w := uint(12 / len(k.headers))
			m.Col(w, func() {
				m.Text(header, props.Text{Align: consts.Center, Style: consts.Bold})
			})
		}
	})

	// Add rows
	for _, row := range k.filteredRows {
		m.Row(10, func() {
			for _, cell := range row {
				w := uint(12 / len(row))
				m.Col(w, func() {
					m.Text(cell, props.Text{Align: consts.Left})
				})
			}
		})
	}

	// Save the PDF
	err := m.OutputFileAndClose(filename)
	if err != nil {
		logz.Error("Could not save PDF: "+err.Error(), map[string]interface{}{
			"context":  "ExportToPDF",
			"filename": filename,
		})
	}
}

func (k *TableRenderer) ExportToMarkdown(filename string) {
	// Implementation for exporting to Markdown
}

func (k *TableRenderer) ToggleColumnVisibility() {
	for header := range k.visibleCols {
		k.visibleCols[header] = !k.visibleCols[header]
	}
	k.kTb = k.kTb.Rows(k.GetCurrentPageRows()...)
}

func (k *TableRenderer) GetHeaders() []string { return k.headers }

func (k *TableRenderer) GetRows() [][]string { return k.rows }

func (k *TableRenderer) GetArrayMap() map[string][]string {
	m := make(map[string][]string)
	for _, row := range k.rows {
		m[row[0]] = row[1:]
	}
	return m
}

func (k *TableRenderer) GetHashMap() map[string]string {
	m := make(map[string]string)
	for _, row := range k.rows {
		m[row[0]] = row[1]
	}
	return m
}

func (k *TableRenderer) GetObjectMap() []map[string]string {
	var m []map[string]string
	for _, row := range k.rows {
		m = append(m, map[string]string{row[0]: row[1]})
	}
	return m
}

func (k *TableRenderer) GetByteMap() map[string][]byte {
	m := make(map[string][]byte)
	for _, row := range k.rows {
		m[row[0]] = []byte(row[1])
	}
	return m
}

func GetTableScreen(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color) string {
	k := NewTableRenderer(tbHandler, customStyles)
	return k.View()
}

func StartTableScreen(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color) error {
	k := NewTableRenderer(tbHandler, customStyles)

	p := tea.NewProgram(k, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logz.Error("Error running table screen: "+err.Error(), map[string]interface{}{
			"context":   "StartTableScreen",
			"tbHandler": tbHandler,
		})
		return nil
	}
	return nil
}

func StartTableScreenFromRenderer(k *TableRenderer) error {
	p := tea.NewProgram(k, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logz.Error("Error running table screen: "+err.Error(), map[string]interface{}{
			"context": "StartTableScreenFromRenderer",
		})
		return err
	}
	return nil
}

func NavigateAndExecuteTable(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color) error {
	k := NewTableRenderer(tbHandler, customStyles)

	p := tea.NewProgram(k, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logz.Error("Error running table screen: "+err.Error(), map[string]interface{}{
			"context":   "NavigateAndExecuteTable",
			"tbHandler": tbHandler,
		})
		return nil
	}
	return nil
}
