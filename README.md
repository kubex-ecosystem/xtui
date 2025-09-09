# ![XTui Banner](docs/assets/top_banner.png)

**A high-performance, easy-to-use terminal user interface (TUI) library for Go, enabling developers to build interactive and visually appealing terminal applications with minimal effort.**

---

![Go Version](https://img.shields.io/github/go-mod/go-version/kubex-ecosystem/xtui)
![License](https://img.shields.io/github/license/kubex-ecosystem/xtui)
![Build Status](https://img.shields.io/github/actions/workflow/status/kubex-ecosystem/xtui/release.yml)

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [CLI Examples](#cli-examples)
- [Module Examples](#module-examples)
- [Hotkeys](#hotkeys)
- [Form Handling](#form-handling)
- [Data Export](#data-export)
- [Command Navigation Functionalities](#command-navigation-functionalities)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Introduction

**xtui** is a high-performance, easy-to-use terminal user interface (TUI) library for Go. It enables developers to build interactive and visually appealing terminal applications with minimal effort while maintaining flexibility and performance.

## Features

- **Intuitive API** – Simplifies the creation of rich terminal interfaces.
- **Customizable Styles** – Tailor UI components to your specific needs with custom styles and configurations.
- **Interactive Form Handling** – Manage form inputs with validation, password protection, and navigation.
- **Data Filtering, Sorting, and Navigation** – Built-in support for table operations.
- **Keyboard Shortcuts** – Provides an efficient user experience with predefined hotkeys.
- **Paginated Views** – Allows smooth navigation through large datasets.
- **Multi-format Export** – Export data to CSV, YAML, JSON, and XML formats.
- **Error Logging** – Integrated with the **logz** library for error tracking and debugging.

## Installation

To install **xtui**, run the following command:

```sh
go get github.com/kubex-ecosystem/xtui
```

## Usage

Here’s a quick example demonstrating how to use **xtui** for displaying tables:

```go
package main

import (
 "github.com/kubex-ecosystem/xtui"
 "github.com/kubex-ecosystem/xtui/types"
    "github.com/charmbracelet/lipgloss"
)

func main() {
    config := types.FormConfig{
        Fields: []types.Field{
            {Name: "ID", Placeholder: "Unique ID"},
            {Name: "Name", Placeholder: "User Name"},
        },
    }

    customStyles := map[string]lipgloss.Color{
        "Info":    lipgloss.Color("#75FBAB"),
        "Warning": lipgloss.Color("#FDFF90"),
    }

    if err := xtui.StartTableScreen(config, customStyles); err != nil {
        panic(err)
    }
}
```

For form-based interactions:

```go
package main

import (
 "github.com/kubex-ecosystem/xtui"
 "github.com/kubex-ecosystem/xtui/types"
)

func main() {
    config := types.Config{
        Title: "User Registration",
        Fields: types.FormFields{
            Inputs: []types.InputField{
                {Ph: "Name", Tp: "text", Req: true, Err: "Name is required!"},
                {Ph: "Password", Tp: "password", Req: true, Err: "Password is required!"},
            },
        },
    }

    result, err := xtui.ShowForm(config)
    if err != nil {
        panic(err)
    }
    println("Form submitted:", result)
}
```

### Command Navigation Functionalities

The `xtui` module provides several command navigation functionalities to enhance the user experience. These functionalities include `NavigateAndExecuteCommand`, `NavigateAndExecuteFormCommand`, and `NavigateAndExecuteViewCommand`.

#### NavigateAndExecuteCommand

The `NavigateAndExecuteCommand` function handles command navigation and execution. It detects commands and their flags, displays command selection and flag definition in a form, sets flag values based on form input, and executes the command.

Example:

```go
package main

import (
 "github.com/kubex-ecosystem/xtui/cmd/cli"
    "github.com/spf13/cobra"
)

func main() {
    cmd := &cobra.Command{
        Use: "example-command",
        RunE: func(cmd *cobra.Command, args []string) error {
            return cli.NavigateAndExecuteCommand(cmd, args)
        },
    }

    if err := cmd.Execute(); err != nil {
        panic(err)
    }
}
```

#### NavigateAndExecuteFormCommand

The `NavigateAndExecuteFormCommand` function handles form-based navigation and execution. It detects commands and their flags, displays command selection and flag definition in a form, sets flag values based on form input, and executes the command.

Example:

```go
package main

import (
 "github.com/kubex-ecosystem/xtui/cmd/cli"
    "github.com/spf13/cobra"
)

func main() {
    cmd := &cobra.Command{
        Use: "example-form-command",
        RunE: func(cmd *cobra.Command, args []string) error {
            return cli.NavigateAndExecuteFormCommand(cmd, args)
        },
    }

    if err := cmd.Execute(); err != nil {
        panic(err)
    }
}
```

#### NavigateAndExecuteViewCommand

The `NavigateAndExecuteViewCommand` function handles table-based navigation and execution. It detects commands and their flags, displays command selection and flag definition in a table view, sets flag values based on table input, and executes the command.

Example:

```go
package main

import (
 "github.com/kubex-ecosystem/xtui/cmd/cli"
    "github.com/spf13/cobra"
)

func main() {
    cmd := &cobra.Command{
        Use: "example-view-command",
        RunE: func(cmd *cobra.Command, args []string) error {
            return cli.NavigateAndExecuteViewCommand(cmd, args)
        },
    }

    if err := cmd.Execute(); err != nil {
        panic(err)
    }
}
```

## CLI Examples

### Install Applications Command

```sh
go run main.go app-install --application app1 --application app2 --path /usr/local/bin --yes --quiet
```

### Table View Command

```sh
go run main.go table-view
```

### Input Form Command

```sh
go run main.go input-form
```

### Loader Form Command

```sh
go run main.go loader-form
```

## Module Examples

### Log Viewer

```go
package main

import (
 "github.com/kubex-ecosystem/xtui/wrappers"
)

func main() {
    if err := wrappers.LogViewer(); err != nil {
        panic(err)
    }
}
```

### Application Manager

```go
package main

import (
 "github.com/kubex-ecosystem/xtui/wrappers"
)

func main() {
    args := []string{"app1", "app2", "/usr/local/bin", "true", "true"}
    if err := wrappers.InstallDependenciesWithUI(args...); err != nil {
        panic(err)
    }
}
```

## Hotkeys

The following keyboard shortcuts are supported out of the box:

- **q, Ctrl+C:** Exit the application.
- **Enter:** Copy selected row or submit form.
- **Ctrl+R:** Change cursor mode.
- **Tab/Shift+Tab, Up/Down Arrows:** Navigate between form fields or table rows.
- **Ctrl+E:** Export data to CSV.
- **Ctrl+Y:** Export data to YAML.
- **Ctrl+J:** Export data to JSON.
- **Ctrl+X:** Export data to XML.

## Form Handling

**xtui** provides an intuitive API for managing forms with validations:

- **Field Validation:** Enforce required fields, minimum/maximum length, and custom validators.
- **Password Input:** Securely handle password fields with hidden characters.
- **Dynamic Properties:** Automatically adapt form inputs based on external configurations.

### Example

```go
field := types.InputField{
    Ph:  "Email",
    Tp:  "text",
    Req: true,
    Err: "Valid email is required!",
    Vld: func(value string) error {
        if !strings.Contains(value, "@") {
            return fmt.Errorf("Invalid email format")
        }
        return nil
    },
}
```

## Data Export

**xtui** supports exporting table data in multiple formats:

- **CSV:** Saves data as a comma-separated values file.
- **YAML:** Outputs data in a structured YAML format.
- **JSON:** Encodes data into a compact JSON format.
- **XML:** Exports data as XML for interoperability.

### Example (Exporting Data)

To export data to a file, simply use the respective hotkey (e.g., `Ctrl+E` for CSV). Files will be saved with predefined names, such as `exported_data.csv`.

## Testing

To test the new navigation functionalities in the `xtui` module, you can follow these steps:

- Run the unit tests provided in the repository. For example, you can run the tests in `cmd/cli/form_commands.go` and `cmd/cli/views_commands.go` using the `go test` command.
- Use the `NavigateAndExecuteCommand` function in `cmd/cli/app_commands.go` to test command navigation and execution. You can create a new command and call this function with the command and arguments.
- Test the form-based navigation by running the `input-form` command defined in `cmd/cli/form_commands.go`. This command uses the `NavigateAndExecuteFormCommand` function to handle form inputs and execute the command.
- Test the table-based navigation by running the `table-view` command defined in `cmd/cli/views_commands.go`. This command uses the `NavigateAndExecuteViewCommand` function to handle table views and execute the command.
- Test the loader-based navigation by running the `loader-form` command defined in `cmd/cli/form_commands.go`. This command uses the `wrappers.StartLoader` function to display a loader screen and execute the command.

## Contributing

We welcome contributions of all kinds! Whether it’s reporting issues, improving documentation, or submitting new features, your help is appreciated. Please check out our [contributing guidelines](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
