package types

import (
	"bufio"
	"encoding/asn1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"reflect"
	"strings"

	gl "github.com/kubex-ecosystem/logz"
	"github.com/pelletier/go-toml/v2"
	"github.com/subosito/gotenv"
	"gopkg.in/yaml.v3"
)

// IMapper is a generic interface for serializing and deserializing objects of type T.
type IMapper[T any] interface {
	// SerializeToFile serializes an object of type T to a file in the specified format.
	SerializeToFile(format string)
	// DeserializeFromFile deserializes an object of type T from a file in the specified format.
	DeserializeFromFile(format string) (T, error)
	// Serialize converts an object of type T to a byte array in the specified format.
	Serialize(format string) ([]byte, error)
	// Deserialize converts a byte array in the specified format to an object of type T.
	Deserialize(data []byte, format string) (T, error)
}

// Mapper is a generic struct that implements the IMapper interface for serializing and deserializing objects.
type Mapper[T any] struct {
	filePath string
	object   T
}

// NewMapperTypeWithObject creates a new instance of Mapper.
func NewMapperTypeWithObject[T any](object *T, filePath string) *Mapper[T] {
	return &Mapper[T]{filePath: filePath, object: *object}
}

// NewMapperType creates a new instance of Mapper.
func NewMapperType[T any](object *T, filePath string) *Mapper[T] {
	return &Mapper[T]{filePath: filePath, object: *object}
}

// NewMapperPtr creates a new instance of Mapper.
func NewMapperPtr[T any](object *T, filePath string) *Mapper[*T] {
	return &Mapper[*T]{filePath: filePath, object: object}
}

// NewMapper creates a new instance of Mapper.
func NewMapper[T any](object *T, filePath string) IMapper[T] {
	return NewMapperType(object, filePath)
}

// Serialize converts an object of type T to a byte array in the specified format.
func (m *Mapper[T]) Serialize(format string) ([]byte, error) {
	switch format {
	case "json":
		return json.Marshal(m.object)
	case "yaml":
		return yaml.Marshal(m.object)
	case "xml":
		return xml.Marshal(m.object)
	case "toml":
		return toml.Marshal(m.object)
	case "asn":
		return asn1.Marshal(m.object)
	case "env":
		if env, ok := reflect.ValueOf(m.object).Interface().(map[string]string); ok {
			if strM, strMErr := gotenv.Marshal(env); strMErr != nil {
				gl.Log("error", fmt.Sprintf("Error serializing to env: %v", strMErr))
				return nil, fmt.Errorf("erro ao serializar para env: %v", strMErr)
			} else {
				return []byte(strM), nil
			}
		} else {
			return nil, fmt.Errorf("tipo não suportado para env: %T", m.object)
		}
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

// Deserialize converts a byte array in the specified format to an object of type T.
func (m *Mapper[T]) Deserialize(data []byte, format string) (T, error) {
	if len(data) == 0 {
		return *new(T), fmt.Errorf("os dados estão vazios")
	}
	if !reflect.ValueOf(m.object).IsValid() || reflect.ValueOf(m.object).IsNil() {
		m.object = *new(T)
		//val := reflect.ValueOf(m.object)
		//if val.Kind() == reflect.Ptr {
		//	if val.IsNil() {
		//		val.Set(reflect.New(val.Type()))
		//	}
		//	m.object = val.Interface().(T)
		//} else if val.Kind() == reflect.Map {
		//	if val.IsNil() {
		//		val.Set(reflect.MakeMap(val.Type()))
		//	}
		//	m.object = val.Interface().(T)
		//} else if val.Kind() == reflect.Slice {
		//	if val.IsNil() {
		//		sliceType := reflect.SliceOf(reflect.TypeOf(m.object).Elem())
		//		val = reflect.New(sliceType).Elem()
		//		// Initialize the slice with zero length and capacity
		//		val.Set(reflect.MakeSlice(sliceType, 0, 0))
		//	} else {
		//		// Initialize the slice with zero length and capacity
		//		val.Set(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(m.object).Elem()), 0, 0))
		//	}
		//	m.object = val.Interface().(T)
		//} else if val.Kind() == reflect.Array {
		//	if val.IsNil() {
		//		arrayType := reflect.ArrayOf(val.Len(), val.Type().Elem())
		//		val = reflect.New(arrayType).Elem()
		//		// Initialize the array with zero length and capacity
		//		val.Set(reflect.MakeSlice(arrayType, 0, 0))
		//	} else {
		//		// Initialize the array with zero length and capacity
		//		val.Set(reflect.MakeSlice(reflect.ArrayOf(val.Len(), val.Type().Elem()), 0, 0))
		//	}
		//	m.object = val.Interface().(T)
		//} else {
		//	m.object = reflect.New(reflect.TypeOf(m.object)).Interface().(T)
		//}
	}
	var err error
	switch format {
	case "json", "js":
		err = json.Unmarshal(data, m.object)
	case "yaml", "yml":
		err = yaml.Unmarshal(data, m.object)
	case "xml", "html":
		err = xml.Unmarshal(data, m.object)
	case "toml", "tml":
		err = toml.Unmarshal(data, m.object)
	case "asn", "asn1":
		_, err = asn1.Unmarshal(data, m.object)
	case "env", "envs", ".env", "environment":
		value := reflect.ValueOf(m.object)
		envT, ok := value.Interface().(*map[string]string)
		if !ok {
			envTB, ok := value.Interface().(map[string]string)
			if !ok {
				gl.Log("error", fmt.Sprintf("Type not valid for env: %T", m.object))
				return *new(T), fmt.Errorf("envT não é válido")
			} else {
				envT = &envTB
			}
		}
		if envT != nil {
			if !reflect.ValueOf(envT).IsValid() {
				gl.Log("error", fmt.Sprintf("Type not valid for env: %T", m.object))
				return *new(T), fmt.Errorf("envT não é válido")
			}
			env := *envT
			if strM, strMErr := gotenv.Unmarshal(string(data)); strMErr != nil {
				return *new(T), fmt.Errorf("erro ao desserializar de env: %v", strMErr)
			} else {
				for k, v := range strM {
					env[k] = v
				}
			}
		} else {
			// We not set err to nil, also we not set err to another value here.
			// This is a special case where we want to return nil, to allow the caller to
			// know that the object is really nil.
			gl.Log("error", fmt.Sprintf("Nil type for env: %T", m.object))
		}
	default:
		err = fmt.Errorf("formato não suportado: %s", format)
	}
	if err != nil {
		return *new(T), fmt.Errorf("erro ao desserializar os dados: %v", err)
	}
	return m.object, nil
}

// SerializeToFile serializes an object of type T to a file in the specified format.
func (m *Mapper[T]) SerializeToFile(format string) {
	if dataSer, dataSerErr := m.Serialize(format); dataSerErr != nil {
		gl.Log("error", fmt.Sprintf("Error serializing object: %v", dataSerErr.Error()))
		return
	} else {
		gl.Log("debug", fmt.Sprintf("Serialized object: %s", string(dataSer)))
		orf, orfErr := os.OpenFile(m.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if orfErr != nil {
			gl.Log("error", fmt.Sprintf("Error opening file: %v", orfErr.Error()))
			return
		}
		defer func() {
			if err := orf.Close(); err != nil {
				gl.Log("error", fmt.Sprintf("Error closing file: %v", err.Error()))
				return
			}
		}()
		if _, err := orf.WriteString(fmt.Sprintf("%s\n", string(dataSer))); err != nil {
			gl.Log("error", fmt.Sprintf("Error writing to file: %v", err.Error()))
			return
		}
	}
}

// DeserializeFromFile deserializes an object of type T from a file in the specified format.
func (m *Mapper[T]) DeserializeFromFile(format string) (T, error) {
	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		gl.Log("error", fmt.Sprintf("File does not exist: %v", err.Error()))
		return *new(T), err
	}

	inputFile, err := os.Open(m.filePath)
	if err != nil {
		gl.Log("error", "Error opening file: %v", err.Error())
		return *new(T), err
	}

	defer func(inputFile *os.File) {
		gl.Log("debug", "Closing input file")
		if closeErr := inputFile.Close(); closeErr != nil {
			gl.Log("error", fmt.Sprintf("Error closing file: %v", closeErr.Error()))
			return
		}
	}(inputFile)

	reader := bufio.NewReader(inputFile)
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanLines)
	objSlice := make([]T, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if desObj, desObjErr := m.Deserialize([]byte(line), format); desObjErr != nil {
			gl.Log("error", fmt.Sprintf("Error deserializing line: %v", desObjErr.Error()))
			return *new(T), err
		} else {
			gl.Log("debug", fmt.Sprintf("Deserialized line: %s", line))
			gl.Log("debug", fmt.Sprintf("Deserialized object: %v", desObj))
			if err = scanner.Err(); err != nil {
				gl.Log("error", fmt.Sprintf("Error reading file: %v", err.Error()))
				return *new(T), err
			}
			objSlice = append(objSlice, desObj)
		}
	}
	if err := scanner.Err(); err != nil {
		gl.Log("error", fmt.Sprintf("Error reading file: %v", err.Error()))
		return *new(T), err
	}
	gl.Log("debug", "File closed successfully")

	var isSliceOrMap int
	value := reflect.ValueOf(m.object)

	switch reflect.TypeFor[T]().Kind() {
	case reflect.Slice, reflect.SliceOf(reflect.TypeFor[map[string]string]()).Kind():
		if value.Len() == 0 {
			// If the slice is empty, assign the deserialized object to the slice
			m.object = *reflect.ValueOf(objSlice).Interface().(*T)
			return m.object, nil
		}
	case reflect.Map:
		if value.Len() == 0 {
			// If the map is empty, assign the deserialized object to the map
			m.object = *reflect.ValueOf(objSlice).Interface().(*T)
			return m.object, nil
		}
		isSliceOrMap = 1
	default:
		// If the type is neither a slice nor a map, assign the first object to m.object
		if len(objSlice) == 0 {
			gl.Log("error", "No objects found in the file")
			return *new(T), fmt.Errorf("nenhum objeto encontrado no arquivo")
		}
		if len(objSlice) > 1 {
			gl.Log("error", "Multiple objects found in the file")
			return *new(T), fmt.Errorf("múltiplos objetos encontrados no arquivo")
		}
		m.object = objSlice[0]
	}

	for _, obj := range objSlice {
		if isSliceOrMap == 0 {
			if reflect.TypeOf(m.object).Kind() == reflect.Slice {
				m.object = *reflect.AppendSlice(reflect.ValueOf(m.object), reflect.ValueOf(obj)).Interface().(*T)
			} else {
				// Check if is a pointer
				if reflect.TypeOf(m.object).Kind() != reflect.Ptr {
					m.object = *reflect.Append(reflect.ValueOf(m.object), reflect.ValueOf(obj)).Interface().(*T)
				} else {
					// Check if is a map
					if reflect.TypeOf(m.object).Kind() == reflect.Map {
						m.object = *reflect.Append(reflect.ValueOf(m.object), reflect.ValueOf(obj)).Interface().(*T)
					} else {
						m.object = obj
					}
				}
			}
		} else {
			newMap := reflect.ValueOf(obj)
			iter := newMap.MapRange()
			for iter.Next() {
				value.SetMapIndex(iter.Key(), iter.Value())
			}
		}
	}

	m.object = value.Interface().(T)
	gl.Log("debug", fmt.Sprintf("File %s deserialized successfully", m.filePath))
	return m.object, nil
}

func SanitizeQuotesAndSpaces(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "'", "\"")
	input = strings.Trim(input, "\"")
	return input
}

func IsEqual(a, b string) bool {
	a, b = SanitizeQuotesAndSpaces(a), SanitizeQuotesAndSpaces(b)
	ptsEqual := levenshtein(a, b)
	maxLen := maxL(len(a), len(b))
	threshold := maxLen / 4
	return ptsEqual <= threshold
}

func maxL(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func levenshtein(s, t string) int {
	m, n := len(s), len(t)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}
	prevRow := make([]int, n+1)
	for j := 0; j <= n; j++ {
		prevRow[j] = j
	}
	for i := 1; i <= m; i++ {
		currRow := make([]int, n+1)
		currRow[0] = i
		for j := 1; j <= n; j++ {
			cost := 1
			if s[i-1] == t[j-1] {
				cost = 0
			}
			currRow[j] = min(prevRow[j]+1, currRow[j-1]+1, prevRow[j-1]+cost)
		}
		prevRow = currRow
	}
	return prevRow[n]
}
