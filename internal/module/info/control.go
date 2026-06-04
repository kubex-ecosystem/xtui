// Package info gerencia controle e configuração modular, com suporte a arquivos separados por módulo.
package info

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	kbxTypes "github.com/kubex-ecosystem/kbx/types"
	gl "github.com/kubex-ecosystem/logz"
	control "github.com/kubex-ecosystem/xtui/internal/module/control"
)

type ControlIface interface {
	GetName() string
	GetVersion() string
}

// Control representa a configuração de controle de um módulo.
type Control struct {
	Reference     kbxTypes.GlobalRef `json:"-"` // Usado internamente para nome do arquivo, nunca exportado
	SchemaVersion int                `json:"schema_version"`

	JobFlagControl control.JobFlagC  `json:"-"`
	JobFlag        control.JobFlag   `json:"-"`
	JobState       control.JobState  `json:"-"`
	JobStateSec    control.JobStateS `json:"-"`

	SecFlag      control.SecFlag            `json:"-"`
	FlagReg32    control.FlagReg32[uint32]  `json:"-"`
	FlagReg32Arr control.FlagReg32A[uint32] `json:"-"`

	Seq     int64  `json:"seq"`
	EpochNS int64  `json:"epoch_ns"`
	Version string `json:"version"`
}

func (c *Control) GetName() string    { return c.Reference.Name }
func (c *Control) GetVersion() string { return c.Version }

// LoadControlByModule carrega o controle de um arquivo específico do módulo.
func LoadControlByModule(dir string, moduleName string) (*Control, error) {
	file := filepath.Join(dir, fmt.Sprintf("control_%s.json", moduleName))
	f, err := os.Open(file)
	if err != nil {
		return nil, gl.Errorf("erro ao abrir %s: %v", file, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			_ = gl.Errorf("erro ao fechar %s: %v", file, cerr)
		}
	}()
	var c Control
	dec := json.NewDecoder(f)
	if err := dec.Decode(&c); err != nil {
		return nil, gl.Errorf("erro ao decodificar %s: %v", file, err)
	}
	c.Reference = kbxTypes.GlobalRef{Name: moduleName}
	return &c, nil
}

// SaveControl salva o controle do módulo em arquivo separado.
func (c *Control) SaveControl(dir string) error {
	if c.Reference.Name == "" {
		return gl.Errorf("Reference.Name não pode ser vazio para salvar o controle")
	}
	file := filepath.Join(dir, fmt.Sprintf("control_%s.json", c.Reference.Name))
	f, err := os.Create(file)
	if err != nil {
		return gl.Errorf("erro ao criar %s: %v", file, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			_ = gl.Errorf("erro ao fechar %s: %v", file, cerr)
		}
	}()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	// Reference não é exportado
	return enc.Encode(c)
}
