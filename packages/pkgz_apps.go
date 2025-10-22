// Package packages provides utilities for managing applications and dependencies.
package packages

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	gl "github.com/kubex-ecosystem/logz/logger"
	cmp "github.com/kubex-ecosystem/xtui/components"
	t "github.com/kubex-ecosystem/xtui/types"
)

// AppInfo show information about an installed application.
type AppInfo struct {
	Name        string // Name
	Version     string // Version
	Method      string // Method
	Status      string // Status
	Description string // Description
}

// AppsTableHandler lida com a tabela de aplicativos.
type AppsTableHandler struct {
	t.TableDataHandler
	apps []AppInfo // Lista de aplicativos
}

// GetHeaders retorna os cabeçalhos da tabela de aplicativos.
// Retorna um slice de strings com os cabeçalhos.
func (h *AppsTableHandler) GetHeaders() []string {
	return []string{"Name", "Version", "Method", "Status", "Description"}
}

// GetRows retorna as linhas da tabela de aplicativos.
// Retorna um slice de slices de strings com as linhas da tabela.
func (h *AppsTableHandler) GetRows() [][]string {
	var rows [][]string //nolint:prealloc
	for _, app := range h.apps {
		rows = append(rows, []string{app.Name, app.Version, app.Method, app.Status, app.Description})
	}
	return rows
}

// CheckDeps verifica se as dependências estão instaladas.
// Recebe o caminho do arquivo de validação e a versão.
// Retorna um booleano indicando se as dependências estão instaladas.
func CheckDeps(validationFilePath string, version string) bool {
	validationFilePath = fmt.Sprintf("%s_%s", validationFilePath, version)
	_, err := os.Stat(validationFilePath)
	return err == nil
}

// GenDepsScript gera um script para verificar e instalar dependências.
// Recebe uma lista de dependências, o caminho do script, o caminho do arquivo de validação e a versão.
// Retorna um erro, se houver.
func GenDepsScript(depsList []string, scriptPath string, validationFilePath string, version string) error {
	validationFilePath = fmt.Sprintf("%s_%s", validationFilePath, version)
	scriptContent := `#!/bin/bash

	# Função para verificar se um comando está disponível
	command_exists() {
	    dpkg -l "$1" &> /dev/null
	}

	# Lista de dependências
	dependencies=(
	`
	if len(depsList) == 0 {
		gl.Log("error", "Dependencies list is empty")
		return fmt.Errorf("lista de dependências vazia")
	}
	for _, dep := range depsList {
		scriptContent += fmt.Sprintf("    \"%s\"\n", dep)
	}
	scriptContent += `)

	# Verifica e instala dependências
	for dep in "${dependencies[@]}"; do
	    if ! command_exists $dep; then
	        echo "$dep não está instalado. Deseja instalar? (s/n)"
	        read -r response
	        if [[ "$response" == "s" ]]; then
	            sudo apt-get update
	            sudo apt-get install -y $dep
	        else
	            echo "Instalação de $dep cancelada."
	        fi
	    else
	        echo "$dep já está instalado."
	    fi
	done

	# Cria arquivo de validação
	mkdir -p $(dirname ` + validationFilePath + `)
	touch ` + validationFilePath + `

	echo "Todas as dependências foram verificadas e instaladas."

	sleep 3

	printf "\033[H\033[2J"

	`
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755) //nolint:gosec
	if err != nil {
		gl.Log("error", "Error writing deps script: "+err.Error())
		return fmt.Errorf("erro ao escrever o script de dependências")
	}
	return nil
}

// InstallAppsShell executa um script de instalação de aplicativos.
// Recebe o caminho do script.
// Retorna um erro, se houver.
func InstallAppsShell(scriptPath string) error {
	if scriptPath == "" {
		return fmt.Errorf("caminho do script de instalação não informado")
	}
	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin // Permite interação do usuário
	err := cmd.Run()
	if err != nil {
		gl.Log("error", "Error running deps script: "+err.Error())
		return err
	}
	return nil
}

// InstallApps instala aplicativos usando a interface do usuário.
// Recebe uma lista de argumentos.
// Retorna um erro, se houver.
func InstallApps(args ...string) error {
	return InstallDepsWithUI(args...)
}

// getInstalledAppsHandler obtém os aplicativos instalados filtrados por nome, status ou método de instalação.
// Recebe o nome, status e método de instalação.
// Retorna um ponteiro para AppsTableHandler e um erro, se houver.
func getInstalledAppsHandler(name string, status string, method string) (t.TableDataHandler, error) {
	// Filtra os aplicativos instalados por nome, status ou method de instalação (auto/manual) usando dpkg-query e se preciso grep
	//nameFilter := ""
	//if name != "" {
	//	nameFilter = fmt.Sprintf("| grep -i %s", name)
	//}
	//statusFilter := ""
	//if status != "" {
	//	statusFilter = fmt.Sprintf("| grep -i %s", status)
	//}
	//methodFilter := ""
	//if method != "" {
	//	methodFilter = fmt.Sprintf("| grep -i %s", method)
	//}
	cmd := exec.Command("bash", "-c", "dpkg-query -W -f='${Package}\\t${Version}\\t${Status}\\t${Description}\\n'") //nolint:gosec
	output, err := cmd.Output()
	if err != nil {
		gl.Log("error", "Error getting installed apps: "+err.Error())
		return nil, fmt.Errorf("erro ao obter aplicativos instalados")
	}
	var apps []AppInfo //nolint:prealloc
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, "\t", 4)
		if len(fields) < 4 {
			continue
		}

		status := "installed"
		if strings.Contains(fields[2], "deinstall") {
			status = "residual"
		}

		apps = append(apps, AppInfo{
			Name:        fields[0],
			Version:     fields[1],
			Method:      "auto", // Assumindo que todos são automáticos, pode ser ajustado conforme necessário
			Status:      status,
			Description: fields[3],
		})
	}

	return &AppsTableHandler{apps: apps}, nil
}

// ShowInstalledAppsTable exibe a tabela de aplicativos instalados.
// Recebe uma lista de argumentos.
// Retorna um erro, se houver.
func ShowInstalledAppsTable(args ...string) error {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}
	status := ""
	if len(args) > 1 {
		status = args[1]
	}
	method := ""
	if len(args) > 2 {
		method = args[2]
	}
	handler, err := getInstalledAppsHandler(name, status, method)
	if err != nil {
		return err
	}
	customStyles := map[string]lipgloss.Color{
		"header": lipgloss.Color("#01BE85"),
		"row":    lipgloss.Color("#252"),
	}

	return cmp.StartTableScreen(handler, customStyles)
}

// installGoogleAuthenticator instala o Google Authenticator.
// Retorna um erro, se houver.
func installGoogleAuthenticator() error {
	cmd := exec.Command("sudo", "apt-get", "install", "-y", "libpam-google-authenticator")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
