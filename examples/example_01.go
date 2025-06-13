package examples

import (
	tea "github.com/charmbracelet/bubbletea"
	c "github.com/rafa-mori/xtui/components"
)

type AppModel struct {
	tables       []*c.TableRenderer
	currentTable int
}

func (m *AppModel) Init() tea.Cmd {
	return nil
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "ctrl+tab":
			m.currentTable = (m.currentTable + 1) % len(m.tables)
		case "ctrl+shift+tab":
			m.currentTable = (m.currentTable - 1 + len(m.tables)) % len(m.tables)
		case "enter":
			/*config := t*/

			// Abre um formulário baseado na tabela atual e linha selecionada
			//if err := c.ShowForm(/*m.tables[m.currentTable].GetSelectedRow()*/); err != nil {
			// 	return m, err
			//}
		}
	}
	return m, nil
}

func (m *AppModel) View() string {
	return m.tables[m.currentTable].View()
}
