package cli

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	c "github.com/kubex-ecosystem/xtui/components"
	gl "github.com/kubex-ecosystem/xtui/logger"
	t "github.com/kubex-ecosystem/xtui/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func ViewsCmdsList() []*cobra.Command {
	tableCmd := tableViewCmd()

	return []*cobra.Command{
		tableCmd,
	}
}

func tableViewCmd() *cobra.Command {
	var jsonFile, xmlFile, yamlFile, csvFile string
	var delimiter, quote, comment string

	cmd := &cobra.Command{
		Use:     "table",
		Aliases: []string{"tb", "t"},
		Annotations: GetDescriptions(
			[]string{
				"Table view for any command",
				"Table view screen, interactive mode, for any command with flags",
			},
			false,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			var inputData [][]string
			var inputDataErr error

			// Check if input is from a pipe
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				pipeInput, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				inputData, inputDataErr = parseCSV(pipeInput, delimiter, quote, comment)
				if inputDataErr != nil {
					return inputDataErr
				}
			} else {
				// Process file inputs
				if csvFile != "" {
					data, err := os.ReadFile(csvFile)
					if err != nil {
						return err
					}
					inputData, inputDataErr = parseCSV(data, delimiter, quote, comment)
					if inputDataErr != nil {
						return inputDataErr
					}
				} else if jsonFile != "" {
					data, err := os.ReadFile(jsonFile)
					if err != nil {
						return err
					}
					mapper := t.NewMapperTypeWithObject[[][]string](&inputData, "/tmp")
					resB, err := mapper.Deserialize(data, "json")
					if err != nil {
						return err
					}
					inputData = resB
				} else if xmlFile != "" {
					data, err := os.ReadFile(xmlFile)
					if err != nil {
						return err
					}
					inputData, inputDataErr = parseXML(data)
					if inputDataErr != nil {
						return inputDataErr
					}
				} else if yamlFile != "" {
					data, err := os.ReadFile(yamlFile)
					if err != nil {
						return err
					}
					inputData, inputDataErr = parseYAML(data)
					if inputDataErr != nil {
						return inputDataErr
					}
				} else if len(args) > 0 {
					inputData, inputDataErr = parseArgs(args)
					if inputDataErr != nil {
						return inputDataErr
					}
				}
			}

			customStyles := map[string]lipgloss.Color{
				"Info":    lipgloss.Color("#75FBAB"),
				"Warning": lipgloss.Color("#FDFF90"),
				"Error":   lipgloss.Color("#FF7698"),
				"Debug":   lipgloss.Color("#929292"),
			}

			headers := inputData[0]
			rows := inputData[1:]

			tbC := c.NewTableRenderer(&t.TableHandler{Headers: headers, Rows: rows}, customStyles, nil)

			return c.StartTableScreenFromRenderer(tbC)
		},
	}

	cmd.Flags().StringVarP(&jsonFile, "json", "j", "", "Input JSON file")
	cmd.Flags().StringVarP(&xmlFile, "xml", "x", "", "Input XML file")
	cmd.Flags().StringVarP(&yamlFile, "yaml", "y", "", "Input YAML file")
	cmd.Flags().StringVarP(&csvFile, "csv", "c", "", "Input CSV file")
	cmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "CSV delimiter")
	cmd.Flags().StringVarP(&quote, "quote", "q", "\"", "CSV quote")
	cmd.Flags().StringVarP(&comment, "comment", "m", "#", "CSV comment")

	return cmd
}

func parseCSV(data []byte, delimiter, quote, comment string) ([][]string, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = []rune(delimiter)[0]
	reader.LazyQuotes = true
	reader.Comment = []rune(comment)[0]
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func parseXML(data []byte) ([][]string, error) {
	var result map[string]string
	err := xml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	var records [][]string
	for key, value := range result {
		records = append(records, []string{key, value})
	}
	return records, nil
}

func parseYAML(data []byte) ([][]string, error) {
	var result map[string]string
	err := yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	var records [][]string
	for key, value := range result {
		records = append(records, []string{key, value})
	}
	return records, nil
}

func parseArgs(args []string) ([][]string, error) {
	var records [][]string
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid argument: %s", arg)
		}
		records = append(records, parts)
	}
	return records, nil
}

func TestTableViewCmd(t *testing.T) {
	cmd := tableViewCmd()
	if cmd.Use != "table" {
		t.Errorf("expected 'table', got '%s'", cmd.Use)
	}
	if cmd.Short != "Table view for any command" {
		t.Errorf("expected 'Table view for any command', got '%s'", cmd.Short)
	}
	if cmd.Long != "Table view screen, interactive mode, for any command with flags" {
		t.Errorf("expected 'Table view screen, interactive mode, for any command with flags', got '%s'", cmd.Long)
	}
}

func NavigateAndExecuteViewCommand(cmd *cobra.Command, args []string) error {
	// Detect command and its flags
	commandName := cmd.Name()
	flags := cmd.Flags()

	// Display command selection and flag definition in a table view
	tableConfig := createTableConfig(commandName, flags)
	customStyles := map[string]lipgloss.Color{
		"Info":    lipgloss.Color("#75FBAB"),
		"Warning": lipgloss.Color("#FDFF90"),
		"Error":   lipgloss.Color("#FF7698"),
		"Debug":   lipgloss.Color("#929292"),
	}
	if err := c.StartTableScreen(tableConfig, customStyles); err != nil {
		return err
	}

	var tableValues = make(map[string]string)
	for _, header := range tableConfig.GetHeaders() {
		tableValues[header] = ""
	}
	for _, row := range tableConfig.GetRows() {
		var i = 0
		for key, _ := range tableValues {
			tableValues[key] = row[i]
			i++
		}
	}
	// Set flag values based on table input
	flags.VisitAll(func(flag *pflag.Flag) {
		if value, ok := tableValues[flag.Name]; ok {
			if err := flag.Value.Set(value); err != nil {
				gl.Log("fatal", err.Error())
				return
			}
		}
	})

	// Execute the command
	return cmd.Execute()
}

func createTableConfig(commandName string, flags *pflag.FlagSet) *c.TableRenderer {
	var tableFields *c.TableRenderer
	var tableHeaders []string
	var tableRows [][]string

	flags.VisitAll(func(flag *pflag.Flag) {
		tableHeaders = append(tableHeaders, flag.Name)
		tableRows = append(tableRows, []string{flag.Name, flag.Usage})
	})

	tableFields = c.NewTableRenderer(&t.TableHandler{
		Headers: tableHeaders,
		Rows:    tableRows,
	}, make(map[string]lipgloss.Color), nil)

	return tableFields
}
