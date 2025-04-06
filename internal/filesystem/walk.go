package filesystem

import (
	"os"
	"path/filepath"
	"strings"
)

type ConfigFile struct {
	Path     string
	Filename string
}

func walkDirectory(root string, depth int) ([]ConfigFile, error) {
	var configFiles []ConfigFile

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignorar diretórios ou arquivos ocultos
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Verificar extensão do arquivo para "config" (ex.: .conf, .json, .yaml, etc.)
		if filepath.Ext(info.Name()) == ".conf" || filepath.Ext(info.Name()) == ".yaml" || filepath.Ext(info.Name()) == ".json" {
			relPath, err := filepath.Rel(root, path)
			if err == nil {
				configFiles = append(configFiles, ConfigFile{Path: relPath, Filename: info.Name()})
			}
		}
		return nil
	})

	// Limitar profundidade (remover arquivos fora do limite)
	if depth > 0 {
		configFiles = filterByDepth(configFiles, depth)
	}

	return configFiles, err
}

func filterByDepth(files []ConfigFile, maxDepth int) []ConfigFile {
	var filteredFiles []ConfigFile
	for _, file := range files {
		segments := strings.Split(filepath.Dir(file.Path), string(filepath.Separator))
		if len(segments) <= maxDepth {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles
}
