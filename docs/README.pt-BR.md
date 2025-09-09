[//]: # (![XTui Banner](./assets/banner.png))

# XTui

---

**Uma biblioteca de interface de usuário de terminal (TUI) para Go, de alto desempenho e fácil de usar, permitindo que
desenvolvedores criem aplicações de terminal interativas e visualmente atraentes com o mínimo de esforço.**

---

![Versão Go](https://img.shields.io/github/go-mod/go-version/faelmori/xtui)
![Licença](https://img.shields.io/github/license/faelmori/xtui)
![Status do Build](https://img.shields.io/github/actions/workflow/status/faelmori/xtui/build.yml)

## Índice

- [Introdução](#introducao)
- [Recursos](#recursos)
- [Instalação](#instalacao)
- [Uso](#uso)
- [Exemplos CLI](#exemplos-cli)
- [Exemplos de Módulo](#exemplos-de-modulo)
- [Atalhos de Teclado](#atalhos-de-teclado)
- [Manipulação de Formulários](#manipulacao-de-formularios)
- [Exportação de Dados](#exportacao-de-dados)
- [Funcionalidades de Navegação de Comando](#funcionalidades-de-navegacao-de-comando)
- [Testes](#testes)
- [Contribuindo](#contribuindo)
- [Licença](#licenca)

## Introducao

**xtui** é uma biblioteca TUI para Go, de alto desempenho e fácil de usar. Permite que desenvolvedores criem aplicações
de terminal interativas e visualmente atraentes com o mínimo de esforço, mantendo flexibilidade e performance.

## Recursos

- **API Intuitiva** – Simplifica a criação de interfaces ricas para terminal.
- **Estilos Personalizáveis** – Personalize componentes de UI conforme suas necessidades com estilos e configurações
  customizadas.
- **Manipulação Interativa de Formulários** – Gerencie entradas de formulário com validação, proteção de senha e
  navegação.
- **Filtragem, Ordenação e Navegação de Dados** – Suporte embutido para operações em tabelas.
- **Atalhos de Teclado** – Proporciona uma experiência eficiente com teclas de atalho pré-definidas.
- **Visualizações Paginadas** – Permite navegação fluida em grandes volumes de dados.
- **Exportação Multi-formato** – Exporte dados para CSV, YAML, JSON e XML.
- **Registro de Erros** – Integrado com a biblioteca **logz** para rastreamento e depuração de erros.

## Instalacao

Para instalar o **xtui**, execute o comando:

```sh
go get github.com/kubex-ecosystem/xtui
```

## Uso

Aqui está um exemplo rápido demonstrando como usar o **xtui** para exibir tabelas:

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
            {Name: "ID", Placeholder: "ID Único"},
            {Name: "Nome", Placeholder: "Nome do Usuário"},
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

Para interações baseadas em formulário:

```go
package main

import (
    "github.com/kubex-ecosystem/xtui"
    "github.com/kubex-ecosystem/xtui/types"
)

func main() {
    config := types.Config{
        Title: "Cadastro de Usuário",
        Fields: types.FormFields{
            Inputs: []types.InputField{
                {Ph: "Nome", Tp: "text", Req: true, Err: "Nome é obrigatório!"},
                {Ph: "Senha", Tp: "password", Req: true, Err: "Senha é obrigatória!"},
            },
        },
    }
    
    result, err := xtui.ShowForm(config)
    if err != nil {
        panic(err)
    }
    println("Formulário enviado:", result)
}
```

### Funcionalidades de Navegação de Comando

O módulo `xtui` fornece várias funcionalidades de navegação de comando para aprimorar a experiência do usuário. Estas
incluem `NavigateAndExecuteCommand`, `NavigateAndExecuteFormCommand` e `NavigateAndExecuteViewCommand`.

#### NavigateAndExecuteCommand

A função `NavigateAndExecuteCommand` gerencia a navegação e execução de comandos. Ela detecta comandos e seus flags,
exibe seleção de comando e definição de flags em um formulário, define valores dos flags com base na entrada do
formulário e executa o comando.

Exemplo:

```go
package main

import (
    "github.com/kubex-ecosystem/xtui/cmd/cli"
    "github.com/spf13/cobra"
)

func main() {
    cmd := &cobra.Command{
        Use: "exemplo-comando",
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

A função `NavigateAndExecuteFormCommand` gerencia navegação e execução baseadas em formulário. Ela detecta comandos e
seus flags, exibe seleção de comando e definição de flags em um formulário, define valores dos flags com base na entrada
do formulário e executa o comando.

Exemplo:

```go
package main

import (
    "github.com/kubex-ecosystem/xtui/cmd/cli"
    "github.com/spf13/cobra"
)

func main() {
    cmd := &cobra.Command{
        Use: "exemplo-form-comando",
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

A função `NavigateAndExecuteViewCommand` gerencia navegação e execução baseadas em tabela. Ela detecta comandos e seus
flags, exibe seleção de comando e definição de flags em uma visualização de tabela, define valores dos flags com base na
entrada da tabela e executa o comando.

Exemplo:

```go
package main

import (
    "github.com/kubex-ecosystem/xtui/cmd/cli"
    "github.com/spf13/cobra"
)

func main() {
    cmd := &cobra.Command{
        Use: "exemplo-view-comando",
        RunE: func(cmd *cobra.Command, args []string) error {
            return cli.NavigateAndExecuteViewCommand(cmd, args)
        },
    }

    if err := cmd.Execute(); err != nil {
        panic(err)
    }
}
```

## Exemplos CLI

### Comando para Instalar Aplicações

```sh
go run main.go app-install --application app1 --application app2 --path /usr/local/bin --yes --quiet
```

### Comando de Visualização de Tabela

```sh
go run main.go table-view
```

### Comando de Formulário de Entrada

```sh
go run main.go input-form
```

### Comando de Formulário com Loader

```sh
go run main.go loader-form
```

## Exemplos de Módulo

### Visualizador de Logs

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

### Gerenciador de Aplicações

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

## Atalhos de Teclado

Os seguintes atalhos de teclado são suportados nativamente:

- **q, Ctrl+C:** Sair da aplicação.
- **Enter:** Copiar linha selecionada ou enviar formulário.
- **Ctrl+R:** Alterar modo do cursor.
- **Tab/Shift+Tab, Setas Cima/Baixo:** Navegar entre campos do formulário ou linhas da tabela.
- **Ctrl+E:** Exportar dados para CSV.
- **Ctrl+Y:** Exportar dados para YAML.
- **Ctrl+J:** Exportar dados para JSON.
- **Ctrl+X:** Exportar dados para XML.

## Manipulacao de Formularios

O **xtui** fornece uma API intuitiva para gerenciar formulários com validações:

- **Validação de Campos:** Exige campos obrigatórios, tamanho mínimo/máximo e validadores customizados.
- **Campo de Senha:** Manipula campos de senha de forma segura, ocultando caracteres.
- **Propriedades Dinâmicas:** Adapta automaticamente os campos do formulário com base em configurações externas.

### Exemplo

```go
field := types.InputField{
Ph:  "Email",
Tp:  "text",
Req: true,
Err: "E-mail válido é obrigatório!",
Vld: func (value string) error {
if !strings.Contains(value, "@") {
return fmt.Errorf("Formato de e-mail inválido")
}
return nil
},
}
```

## Exportacao de Dados

O **xtui** suporta exportação de dados de tabela em múltiplos formatos:

- **CSV:** Salva os dados em arquivo separado por vírgulas.
- **YAML:** Exporta os dados em formato YAML estruturado.
- **JSON:** Codifica os dados em formato JSON compacto.
- **XML:** Exporta os dados como XML para interoperabilidade.

### Exemplo

Para exportar dados para um arquivo, basta usar o atalho correspondente (ex: `Ctrl+E` para CSV). Os arquivos serão
salvos com nomes predefinidos, como `exported_data.csv`.

## Testes

Para testar as novas funcionalidades de navegação do módulo `xtui`, siga os passos:

* Execute os testes unitários fornecidos no repositório. Por exemplo, rode os testes em `cmd/cli/form_commands.go` e
  `cmd/cli/views_commands.go` usando o comando `go test`.
* Use a função `NavigateAndExecuteCommand` em `cmd/cli/app_commands.go` para testar navegação e execução de comandos.
  Você pode criar um novo comando e chamar essa função com o comando e argumentos.
* Teste a navegação baseada em formulário rodando o comando `input-form` definido em `cmd/cli/form_commands.go`. Esse
  comando usa a função `NavigateAndExecuteFormCommand` para manipular entradas de formulário e executar o comando.
* Teste a navegação baseada em tabela rodando o comando `table-view` definido em `cmd/cli/views_commands.go`. Esse
  comando usa a função `NavigateAndExecuteViewCommand` para manipular visualizações de tabela e executar o comando.
* Teste a navegação baseada em loader rodando o comando `loader-form` definido em `cmd/cli/form_commands.go`. Esse
  comando usa a função `wrappers.StartLoader` para exibir uma tela de carregamento e executar o comando.

## Contribuindo

Contribuições de todos os tipos são bem-vindas! Seja reportando problemas, melhorando a documentação ou enviando novas
funcionalidades, sua ajuda é apreciada. Confira nosso [guia de contribuição](CONTRIBUTING.md) para mais detalhes.

## Licenca

Este projeto está licenciado sob a Licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
