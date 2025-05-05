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
	gl "github.com/faelmori/xtui/logger"
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

// TableRenderer is responsible for rendering tables in the terminal with customizable styles and dynamic behavior.
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

// StyleFunc defines a function that returns a lipgloss.Style based on row, column, and cell value.
type StyleFunc func(row, col int, cellValue string) lipgloss.Style

// NewTableRenderer creates a new TableRenderer with custom styles and an optional style function.
func NewTableRenderer(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color, styleFunc StyleFunc) *TableRenderer {
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

	if styleFunc == nil {
		styleFunc = func(row, col int, cellValue string) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}

			rowIndex := row - 1
			if rowIndex < 0 || rowIndex >= len(rows) {
				return baseStyle
			}
			if rows != nil && len(rows) < rowIndex {
				if rows[rowIndex] != nil && len(rows[rowIndex]) > 1 {
					if rows[rowIndex][1] == "Bug" {
						return selectedStyle
					}
				}
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
		}
	} else {
		styleFunc = func(row, col int, cellValue string) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			rowIndex := row - 1
			if rowIndex < 0 || rowIndex >= len(rows) {
				return baseStyle
			}
			return styleFunc(rowIndex, col, rows[rowIndex][col])
		}
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
			if row == 1 {
				return selectedStyle
			}
			rowIndex := row - 1
			if rowIndex < 0 || len(rows) <= rowIndex {
				return baseStyle
			}
			if rows != nil && rowIndex < len(rows) {
				if rows[rowIndex] != nil &&
					len(rows[rowIndex]) > 1 &&
					len(rows[rowIndex]) > col &&
					len(rows[rowIndex]) > 1 {
					if rows[rowIndex][1] == "Bug" {
						return selectedStyle
					}
				}
			}
			if col >= len(rows[rowIndex]) {
				return baseStyle
			}

			return styleFunc(row, col, rows[row][col])
			//return styleFunc(row, col, rows[rowIndex][col])
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

// Init initializes the table renderer.
func (k *TableRenderer) Init() tea.Cmd {
	return nil
}

// Update updates the table renderer based on user input.
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
	k.kTb.ClearRows()                             // Clear the table rows before adding new ones
	k.kTb = k.kTb.Rows(k.GetCurrentPageRows()...) // Update the table with the current rows
	return k, cmd
}

// View returns the string representation of the table for rendering.
func (k *TableRenderer) View() string {
	helpText := "\nShortcuts:\n" +
		"  - q, ctrl+c: Quit\n" +
		"  - enter: Copy selected row to clipboard\n" +
		"  - esc: Exit selection mode\n" +
		"  - backspace: Remove last character from filter\n" +
		"  - ctrl+o: Toggle sorting\n" +
		"  - right: Next page\n" +
		"  - left: Previous page\n" +
		"  - down: Select next row\n" +
		"  - up: Select previous row\n" +
		"  - ctrl+e: Export to CSV\n" +
		"  - ctrl+y: Export to YAML\n" +
		"  - ctrl+j: Export to JSON\n" +
		"  - ctrl+x: Export to XML\n" +
		"  - ctrl+l: Export to Excel\n" +
		"  - ctrl+p: Export to PDF\n" +
		"  - ctrl+m: Export to Markdown\n" +
		"  - ctrl+k: Toggle column visibility\n"

	toggleHelpText := "\nPress ctrl+h to show/hide shortcuts."

	if k.showHelp {
		return fmt.Sprintf("\nFilter: %s\n\n%s\nPage: %d/%d\n%s%s", k.filter, k.kTb.String(), k.page+1, (len(k.filteredRows)+k.pageSize-1)/k.pageSize, helpText, toggleHelpText)
	}
	return fmt.Sprintf("\nFilter: %s\n\n%s\nPage: %d/%d\n%s", k.filter, k.kTb.String(), k.page+1, (len(k.filteredRows)+k.pageSize-1)/k.pageSize, toggleHelpText)
}

// GetHeaders returns the table headers.
func (k *TableRenderer) GetHeaders() []string { return k.headers }

// GetRows returns the table rows.
func (k *TableRenderer) GetRows() [][]string { return k.rows }

// GetArrayMap returns the table data as a map of arrays.
func (k *TableRenderer) GetArrayMap() map[string][]string {
	m := make(map[string][]string)
	for _, row := range k.rows {
		m[row[0]] = row[1:]
	}
	return m
}

// GetHashMap returns the table data as a hash map.
func (k *TableRenderer) GetHashMap() map[string]string {
	m := make(map[string]string)
	for _, row := range k.rows {
		m[row[0]] = row[1]
	}
	return m
}

// GetObjectMap returns the table data as a slice of maps.
func (k *TableRenderer) GetObjectMap() []map[string]string {
	var m []map[string]string
	for _, row := range k.rows {
		m = append(m, map[string]string{row[0]: row[1]})
	}
	return m
}

// GetByteMap returns the table data as a map of byte slices.
func (k *TableRenderer) GetByteMap() map[string][]byte {
	m := make(map[string][]byte)
	for _, row := range k.rows {
		m[row[0]] = []byte(row[1])
	}
	return m
}

// RowsNavigate navigates through the table rows.
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

// ApplyFilter applies a filter to the table rows.
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

// SortRows sorts the table rows.
func (k *TableRenderer) SortRows() {
	sort.SliceStable(k.filteredRows, func(i, j int) bool {
		if k.sortAsc {
			return k.filteredRows[i][k.sortColumn] < k.filteredRows[j][k.sortColumn]
		}
		return k.filteredRows[i][k.sortColumn] > k.filteredRows[j][k.sortColumn]
	})
	k.kTb = k.kTb.Rows(k.GetCurrentPageRows()...)
}

// GetCurrentPageRows returns the rows for the current page.
func (k *TableRenderer) GetCurrentPageRows() [][]string {
	start := k.page * k.pageSize
	end := start + k.pageSize
	if end > len(k.filteredRows) {
		end = len(k.filteredRows)
	}
	return k.filteredRows[start:end]
}

// ExportToCSV exports the table data to a CSV file.
func (k *TableRenderer) ExportToCSV(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		gl.Log("error", "Error creating file: "+err.Error())
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if writerErr := writer.Write(k.headers); writerErr != nil {
		gl.Log("error", "Error writing headers to CSV:"+writerErr.Error())
		return
	}

	// Write rows
	for _, row := range k.filteredRows {
		if writerRowsErr := writer.Write(row); writerRowsErr != nil {
			gl.Log("error", "Error writing row to CSV: "+writerRowsErr.Error())
			return
		}
	}

	gl.Log("info", "Data exported to CSV:"+filename)
}

// ExportToYAML exports the table data to a YAML file.
func (k *TableRenderer) ExportToYAML(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		gl.Log("error", "Error creating file:"+err.Error())
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
		gl.Log("error", "Error writing data to YAML:"+err.Error())
		return
	}

	gl.Log("info", "Data exported to YAML:"+filename)
}

// ExportToJSON exports the table data to a JSON file.
func (k *TableRenderer) ExportToJSON(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		gl.Log("error", "Error creating file:"+err.Error())
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	data := k.GetObjectMap()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		gl.Log("error", "Error writing data to JSON:"+err.Error())
		return
	}
	gl.Log("info", "Data exported to JSON:"+filename)
}

// ExportToXML exports the table data to an XML file.
func (k *TableRenderer) ExportToXML(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		gl.Log("error", "Error creating file:"+err.Error())
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	data := k.GetObjectMap()
	encoder := xml.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		gl.Log("error", "Error writing data to XML:"+err.Error())
		return
	}
	gl.Log("info", "Data exported to XML:"+filename)
}

// ExportToExcel is a placeholder for exporting the table data to an Excel file.
func (k *TableRenderer) ExportToExcel(filename string) {
	// Implementation for exporting to Excel
}

// ExportToPDF exports the table data to a PDF file.
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
		gl.Log("error", "Could not save PDF: "+err.Error())
	}
}

// ExportToMarkdown exports the table data to a Markdown file.
func (k *TableRenderer) ExportToMarkdown(filename string) {
	// Implementation for exporting to Markdown
}

// ToggleColumnVisibility toggles the visibility of the columns in the table.
func (k *TableRenderer) ToggleColumnVisibility() {
	for header := range k.visibleCols {
		k.visibleCols[header] = !k.visibleCols[header]
	}
	k.kTb = k.kTb.Rows(k.GetCurrentPageRows()...)
}

// Execution functions

// GetTableScreenCustom returns the string representation of the table with custom styles and an optional style function.
func GetTableScreenCustom(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color, styleFunc StyleFunc) string {
	k := NewTableRenderer(tbHandler, customStyles, styleFunc)
	return k.View()
}

// NavigateAndExecuteTableCustom navigates and executes the table screen with custom styles and an optional style function.
func NavigateAndExecuteTableCustom(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color, styleFunc StyleFunc) error {
	k := NewTableRenderer(tbHandler, customStyles, styleFunc)

	prog := tea.NewProgram(k, tea.WithAltScreen())
	if _, err := prog.Run(); err != nil {
		gl.Log("error", "Error running table screen: "+err.Error())
		return nil
	}
	return nil
}

// StartTableScreenCustom starts the table screen with custom styles and an optional style function.
func StartTableScreenCustom(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color, styleFunc StyleFunc) error {
	k := NewTableRenderer(tbHandler, customStyles, styleFunc)

	prog := tea.NewProgram(k, tea.WithAltScreen())
	if _, err := prog.Run(); err != nil {
		gl.Log("error", "Error running table screen: "+err.Error())
		return nil
	}
	return nil
}

// GetTableScreen returns the string representation of the table with custom styles.
func GetTableScreen(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color) string {
	return GetTableScreenCustom(tbHandler, customStyles, nil)
}

// NavigateAndExecuteTable navigates and executes the table screen with custom styles.
func NavigateAndExecuteTable(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color) error {
	return NavigateAndExecuteTableCustom(tbHandler, customStyles, nil)
}

// StartTableScreen starts the table screen with custom styles.
func StartTableScreen(tbHandler TableDataHandler, customStyles map[string]lipgloss.Color) error {
	return StartTableScreenCustom(tbHandler, customStyles, nil)
}

// StartTableScreenFromRenderer starts the table screen from a given TableRenderer.
func StartTableScreenFromRenderer(k *TableRenderer) error {
	prog := tea.NewProgram(k, tea.WithAltScreen())
	if _, err := prog.Run(); err != nil {
		gl.Log("error", "Error running table screen: "+err.Error())
		return err
	}
	return nil
}
