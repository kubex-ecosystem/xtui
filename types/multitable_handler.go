package types

type MultiTableManager struct {
	Handlers []TableDataHandler
	Current  int
}

func (m *MultiTableManager) Next() {
	m.Current = (m.Current + 1) % len(m.Handlers)
}
func (m *MultiTableManager) Previous() {
	m.Current = (m.Current - 1 + len(m.Handlers)) % len(m.Handlers)
}
func (m *MultiTableManager) GetCurrentHandler() TableDataHandler {
	return m.Handlers[m.Current]
}

type MultiTableHandler struct {
	handlers []TableHandler
	current  int // Índice da tabela atual
}

func (h *MultiTableHandler) GetHeaders() []string {
	return h.handlers[h.current].GetHeaders()
}
func (h *MultiTableHandler) GetRows() [][]string {
	return h.handlers[h.current].GetRows()
}
func (h *MultiTableHandler) NextTable() {
	h.current = (h.current + 1) % len(h.handlers)
}
func (h *MultiTableHandler) PreviousTable() {
	h.current = (h.current - 1 + len(h.handlers)) % len(h.handlers)
}
