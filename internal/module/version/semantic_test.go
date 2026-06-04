// Package version testa funcionalidades de versionamento semântico.
//
// Estes testes validam:
// - Parse de versões semânticas (v1.2.3, 1.2.3-beta, etc.)
// - Comparação de versões
// - Verificação de versão máxima
//
// Tipo: Testes unitários
// Previne: Parse incorreto de versões, comparações erradas
package version

import (
	"reflect"
	"testing"
)

// --- Testes de parseVersion ---
// Previne: Parse incorreto de strings de versão

func TestServiceImpl_parseVersion(t *testing.T) {
	svc := &ServiceImpl{}

	tests := []struct {
		name    string
		version string
		want    []int
	}{
		{
			name:    "simple version",
			version: "1.2.3",
			want:    []int{1, 2, 3},
		},
		{
			name:    "version with v prefix",
			version: "v1.2.3",
			want:    []int{1, 2, 3},
		},
		{
			name:    "version with prerelease suffix",
			version: "1.2.3-beta",
			want:    []int{1, 2, 3},
		},
		{
			name:    "version with v prefix and suffix",
			version: "v2.0.0-alpha.1",
			want:    []int{2, 0, 0},
		},
		{
			name:    "major only",
			version: "5",
			want:    []int{5},
		},
		{
			name:    "major.minor only",
			version: "3.14",
			want:    []int{3, 14},
		},
		{
			name:    "four parts",
			version: "1.2.3.4",
			want:    []int{1, 2, 3, 4},
		},
		{
			name:    "empty string",
			version: "",
			want:    nil,
		},
		{
			name:    "invalid version with letters",
			version: "1.2.abc",
			want:    nil,
		},
		{
			name:    "zero version",
			version: "0.0.0",
			want:    []int{0, 0, 0},
		},
		{
			name:    "large numbers",
			version: "100.200.300",
			want:    []int{100, 200, 300},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.parseVersion(tt.version)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

// --- Testes de vrsCompare ---
// Previne: Comparações incorretas de versões

func TestServiceImpl_vrsCompare(t *testing.T) {
	svc := &ServiceImpl{}

	tests := []struct {
		name string
		v1   []int
		v2   []int
		want int
	}{
		{
			name: "equal versions",
			v1:   []int{1, 2, 3},
			v2:   []int{1, 2, 3},
			want: 0,
		},
		{
			name: "v1 less than v2 - major",
			v1:   []int{1, 0, 0},
			v2:   []int{2, 0, 0},
			want: -1,
		},
		{
			name: "v1 greater than v2 - major",
			v1:   []int{2, 0, 0},
			v2:   []int{1, 0, 0},
			want: 1,
		},
		{
			name: "v1 less than v2 - minor",
			v1:   []int{1, 1, 0},
			v2:   []int{1, 2, 0},
			want: -1,
		},
		{
			name: "v1 greater than v2 - minor",
			v1:   []int{1, 3, 0},
			v2:   []int{1, 2, 0},
			want: 1,
		},
		{
			name: "v1 less than v2 - patch",
			v1:   []int{1, 2, 3},
			v2:   []int{1, 2, 4},
			want: -1,
		},
		{
			name: "v1 greater than v2 - patch",
			v1:   []int{1, 2, 5},
			v2:   []int{1, 2, 4},
			want: 1,
		},
		{
			name: "empty v1",
			v1:   []int{},
			v2:   []int{1, 0, 0},
			want: 0, // No comparison possible
		},
		{
			name: "different lengths - v1 shorter",
			v1:   []int{1, 2},
			v2:   []int{1, 2, 0},
			want: 0, // Equal up to length of v1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.vrsCompare(tt.v1, tt.v2)
			if err != nil {
				t.Fatalf("vrsCompare() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("vrsCompare(%v, %v) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

// --- Testes de versionAtMost ---
// Previne: Verificação incorreta de versão máxima

func TestServiceImpl_versionAtMost(t *testing.T) {
	svc := &ServiceImpl{}

	tests := []struct {
		name    string
		version []int
		max     []int
		want    bool
	}{
		{
			name:    "version equals max",
			version: []int{1, 2, 3},
			max:     []int{1, 2, 3},
			want:    true,
		},
		{
			name:    "version less than max",
			version: []int{1, 2, 0},
			max:     []int{1, 2, 3},
			want:    true,
		},
		{
			name:    "version greater than max",
			version: []int{1, 2, 4},
			max:     []int{1, 2, 3},
			want:    false,
		},
		{
			name:    "major version greater",
			version: []int{2, 0, 0},
			max:     []int{1, 9, 9},
			want:    false,
		},
		{
			name:    "minor version greater",
			version: []int{1, 3, 0},
			max:     []int{1, 2, 9},
			want:    false,
		},
		{
			name:    "zero version",
			version: []int{0, 0, 0},
			max:     []int{1, 0, 0},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.versionAtMost(tt.version, tt.max)
			if err != nil {
				t.Fatalf("versionAtMost() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("versionAtMost(%v, %v) = %v, want %v", tt.version, tt.max, got, tt.want)
			}
		})
	}
}

// --- Testes de integração parseVersion + vrsCompare ---
// Previne: Falhas na cadeia de processamento

func TestVersionParseThenCompare(t *testing.T) {
	svc := &ServiceImpl{}

	tests := []struct {
		name    string
		v1Str   string
		v2Str   string
		wantCmp int
	}{
		{
			name:    "compare v1.0.0 with v2.0.0",
			v1Str:   "v1.0.0",
			v2Str:   "v2.0.0",
			wantCmp: -1,
		},
		{
			name:    "compare 1.2.3-beta with 1.2.3",
			v1Str:   "1.2.3-beta",
			v2Str:   "1.2.3",
			wantCmp: 0, // Prerelease suffix is stripped
		},
		{
			name:    "compare identical versions",
			v1Str:   "v3.14.159",
			v2Str:   "3.14.159",
			wantCmp: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := svc.parseVersion(tt.v1Str)
			v2 := svc.parseVersion(tt.v2Str)

			if v1 == nil || v2 == nil {
				t.Fatalf("parseVersion failed for %q or %q", tt.v1Str, tt.v2Str)
			}

			cmp, err := svc.vrsCompare(v1, v2)
			if err != nil {
				t.Fatalf("vrsCompare() error = %v", err)
			}
			if cmp != tt.wantCmp {
				t.Errorf("compare(%q, %q) = %d, want %d", tt.v1Str, tt.v2Str, cmp, tt.wantCmp)
			}
		})
	}
}

// --- Benchmark ---

func BenchmarkParseVersion(b *testing.B) {
	svc := &ServiceImpl{}
	for i := 0; i < b.N; i++ {
		_ = svc.parseVersion("v1.2.3-beta.4")
	}
}

func BenchmarkVersionCompare(b *testing.B) {
	svc := &ServiceImpl{}
	v1 := []int{1, 2, 3}
	v2 := []int{1, 2, 4}
	for i := 0; i < b.N; i++ {
		_, _ = svc.vrsCompare(v1, v2)
	}
}
