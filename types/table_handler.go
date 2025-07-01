package types

type TableDataHandler interface {
	GetHeaders() []string
	GetRows() [][]string
	GetArrayMap() map[string][]string
	GetHashMap() map[string]string
	GetObjectMap() []map[string]string
	GetByteMap() map[string][]byte
}

func NewTableHandler(headers []string, rows [][]string) TableDataHandler {
	return &TableHandler{
		Headers: headers,
		Rows:    rows,
	}
}

func NewTableHandlerFromRows(headers []string, rows [][]string) TableDataHandler {
	if len(headers) == 0 {
		headers = make([]string, len(rows[0]))
		for i := range headers {
			headers[i] = "Column " + string(i+1)
		}
	}
	return &TableHandler{
		Headers: headers,
		Rows:    rows,
	}
}

func NewTableHandlerFromMap(data map[string][]string) TableDataHandler {
	headers := make([]string, 0, len(data))
	rows := make([][]string, 0, len(data))
	for key, values := range data {
		headers = append(headers, key)
		rows = append(rows, append([]string{key}, values...))
	}
	return &TableHandler{
		Headers: headers,
		Rows:    rows,
	}
}

type TableHandler struct {
	TableDataHandler
	Headers []string
	Rows    [][]string
}
type TableHandlerWithContext struct {
	TableHandler
	Context string
}

func (h *TableHandler) GetHeaders() []string { return h.Headers }
func (h *TableHandler) GetRows() [][]string  { return h.Rows }
func (h *TableHandler) GetArrayMap() map[string][]string {
	m := make(map[string][]string)
	for _, row := range h.Rows {
		m[row[0]] = row[1:]
	}
	return m
}
func (h *TableHandler) GetHashMap() map[string]string {
	m := make(map[string]string)
	for _, row := range h.Rows {
		m[row[0]] = row[1]
	}
	return m
}
func (h *TableHandler) GetObjectMap() []map[string]string {
	var m []map[string]string
	for _, row := range h.Rows {
		m = append(m, map[string]string{row[0]: row[1]})
	}
	return m
}
func (h *TableHandler) GetByteMap() map[string][]byte {
	m := make(map[string][]byte)
	for _, row := range h.Rows {
		m[row[0]] = []byte(row[1])
	}
	return m
}
