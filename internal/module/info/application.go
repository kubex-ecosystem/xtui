// Package info provides functionality to read and parse the application manifest.
package info

import (
	_ "embed"
	"encoding/json"
	"os"

	gl "github.com/kubex-ecosystem/logz"
	// "github.com/kubex-ecosystem/xtui/internal/module/info"
)

//go:embed manifest.json
var manifestJSONData []byte

// var application Manifest

type Reference struct {
	Name            string `json:"name"`
	ApplicationName string `json:"application"`
	Bin             string `json:"bin"`
	Version         string `json:"version"`
}

type mmanifest struct {
	Manifest
	Name            string   `json:"name"`
	ApplicationName string   `json:"application"`
	Bin             string   `json:"bin"`
	Version         string   `json:"version"`
	Repository      string   `json:"repository"`
	Aliases         []string `json:"aliases,omitempty"`
	Homepage        string   `json:"homepage,omitempty"`
	Description     string   `json:"description,omitempty"`
	Main            string   `json:"main,omitempty"`
	Author          string   `json:"author,omitempty"`
	License         string   `json:"license,omitempty"`
	Keywords        []string `json:"keywords,omitempty"`
	Platforms       []string `json:"platforms,omitempty"`
	LogLevel        string   `json:"log_level,omitempty"`
	Debug           bool     `json:"debug,omitempty"`
	ShowTrace       bool     `json:"show_trace,omitempty"`
	Private         bool     `json:"private,omitempty"`
}
type Manifest interface {
	GetName() string
	GetVersion() string
	GetAliases() []string
	GetRepository() string
	GetHomepage() string
	GetDescription() string
	GetMain() string
	GetBin() string
	GetAuthor() string
	GetLicense() string
	GetKeywords() []string
	GetPlatforms() []string
	IsPrivate() bool
}

func (m *mmanifest) GetName() string        { return m.Name }
func (m *mmanifest) GetVersion() string     { return m.Version }
func (m *mmanifest) GetAliases() []string   { return m.Aliases }
func (m *mmanifest) GetRepository() string  { return m.Repository }
func (m *mmanifest) GetHomepage() string    { return m.Homepage }
func (m *mmanifest) GetDescription() string { return m.Description }
func (m *mmanifest) GetMain() string        { return m.Main }
func (m *mmanifest) GetBin() string {
	if m.Bin == "" {
		m.Bin, _ = os.Executable()
	}
	return m.Bin
}
func (m *mmanifest) GetAuthor() string      { return m.Author }
func (m *mmanifest) GetLicense() string     { return m.License }
func (m *mmanifest) GetKeywords() []string  { return m.Keywords }
func (m *mmanifest) GetPlatforms() []string { return m.Platforms }
func (m *mmanifest) IsPrivate() bool        { return m.Private }

// lazy cache
var (
	cachedManifest Manifest
	cachedControl  *Control
)

// GetManifest lazy, sem init() com side-effects
func GetManifest() (Manifest, error) {
	if cachedManifest != nil {
		return cachedManifest, nil
	}

	if len(manifestJSONData) == 0 {
		return nil, gl.Errorf("manifest.json: embed is empty")
	}

	var m mmanifest
	if err := json.Unmarshal(manifestJSONData, &m); err != nil {
		return nil, gl.Errorf("manifest.json: %v", err)
	}
	cachedManifest = &m
	return &m, nil
}

// FS secOrder quiser permitir override por FS externo:
type FS interface {
	ReadFile(name string) ([]byte, error)
}

// LoadFromFS lazy, sem init() com side-effects
// TODO: o que acontece se não achar? deve dar erro? ou criar um default?
// De acordo com o pacote control, o control deve ser criado em cada módulo.
// Apenas o xtui (app) não cria um control.json
func LoadFromFS(fs FS) (
	Manifest, // Manifest is the default. But can be other implementations
	*Control, // ControlIface is the interface that the Control implements. We return a pointer to Control because it has side effects when it is loaded
	error,
) {
	var m Manifest
	var c Control

	if b, err := fs.ReadFile("manifest.json"); err == nil {
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, nil, gl.Errorf("manifest.json: %v", err)
		}
	} else {
		return nil, nil, gl.Errorf("manifest.json: %v", err)
	}
	if b, err := fs.ReadFile("control.json"); err == nil {
		if err := json.Unmarshal(b, &c); err != nil {
			return m, &c, gl.Errorf("control.json: %v", err)
		}
	} else {
		return m, nil, gl.Errorf("control.json: %v", err)
	}

	return m, &c, nil
}
