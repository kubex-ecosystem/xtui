### `table_screen.go` Documentation

#### Package `components`

This package provides a `TableRenderer` for rendering tables in the terminal with customizable styles and dynamic behavior.

#### Imports

- `encoding/csv`, `encoding/json`, `encoding/xml`: For parsing different file formats.
- `fmt`: For formatted I/O.
- `github.com/atotto/clipboard`: For clipboard operations.
- `tea "github.com/charmbracelet/bubbletea"`: For building terminal user interfaces.
- `github.com/charmbracelet/lipgloss`: For styling terminal output.
- `github.com/charmbracelet/lipgloss/table`: For creating tables.
- `l "github.com/kubex-ecosystem/logz"`: For logging.
- `. "github.com/kubex-ecosystem/xtui/types"`: For custom types.
- `github.com/johnfercher/maroto/pkg/consts`, `p "github.com/johnfercher/maroto/pkg/pdf"`, `props "github.com/johnfercher/maroto/pkg/props"`: For PDF generation.
- `gopkg.in/yaml.v2`: For YAML parsing.
- `os`, `sort`, `strconv`, `strings`: Standard library packages.

#### Types

- **`TableRenderer`**: Struct for rendering tables with various properties like headers, rows, filters, sorting, pagination, etc.
- **`StyleFunc`**: Type definition for a function that returns a `lipgloss.Style` based on row, column, and cell value.

#### Functions

- **`NewTableRenderer`**: Creates a new `TableRenderer` with custom styles and an optional style function.
- **`(k *TableRenderer) Init`**: Initializes the table renderer.
- **`(k *TableRenderer) Update`**: Updates the table renderer based on user input.
- **`(k *TableRenderer) View`**: Returns the string representation of the table for rendering.
- **`(k *TableRenderer) GetHeaders`**: Returns the table headers.
- **`(k *TableRenderer) GetRows`**: Returns the table rows.
- **`(k *TableRenderer) GetArrayMap`**: Returns the table data as a map of arrays.
- **`(k *TableRenderer) GetHashMap`**: Returns the table data as a hash map.
- **`(k *TableRenderer) GetObjectMap`**: Returns the table data as a slice of maps.
- **`(k *TableRenderer) GetByteMap`**: Returns the table data as a map of byte slices.
- **`(k *TableRenderer) RowsNavigate`**: Navigates through the table rows.
- **`(k *TableRenderer) ApplyFilter`**: Applies a filter to the table rows.
- **`(k *TableRenderer) SortRows`**: Sorts the table rows.
- **`(k *TableRenderer) GetCurrentPageRows`**: Returns the rows for the current page.
- **`(k *TableRenderer) ExportToCSV`**: Exports the table data to a CSV file.
- **`(k *TableRenderer) ExportToYAML`**: Exports the table data to a YAML file.
- **`(k *TableRenderer) ExportToJSON`**: Exports the table data to a JSON file.
- **`(k *TableRenderer) ExportToXML`**: Exports the table data to an XML file.
- **`(k *TableRenderer) ExportToExcel`**: Placeholder for exporting the table data to an Excel file.
- **`(k *TableRenderer) ExportToPDF`**: Exports the table data to a PDF file.
- **`(k *TableRenderer) ExportToMarkdown`**: Placeholder for exporting the table data to a Markdown file.
- **`(k *TableRenderer) ToggleColumnVisibility`**: Toggles the visibility of table columns.

#### Execution Functions

- **`GetTableScreenCustom`**: Returns the table screen view with custom styles and an optional style function.
- **`NavigateAndExecuteTableCustom`**: Navigates and executes the table screen with custom styles and an optional style function.
- **`StartTableScreenCustom`**: Starts the table screen with custom styles and an optional style function.
- **`GetTableScreen`**: Returns the table screen view with custom styles.
- **`NavigateAndExecuteTable`**: Navigates and executes the table screen with custom styles.
- **`StartTableScreen`**: Starts the table screen with custom styles.
- **`StartTableScreenFromRenderer`**: Starts the table screen from an existing `TableRenderer`.

This documentation provides an overview of the `table_screen.go` file, its types, functions, and their purposes.