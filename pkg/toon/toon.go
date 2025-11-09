package toon

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/devalexandre/toon-go/pkg/decoder"
	"github.com/devalexandre/toon-go/pkg/encoder"
)

type Marshaler interface {
	MarshalTOON() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalTOON(data []byte) error
}

func Marshal(v interface{}) ([]byte, error) {
	var buf []byte
	writer := &byteWriter{buf: &buf}
	enc := encoder.NewEncoder(writer, nil)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return *writer.buf, nil
}

func Unmarshal(data []byte, v interface{}) error {
	reader := strings.NewReader(string(data))
	dec := decoder.NewParser(reader)

	result, err := dec.Parse()
	if err != nil {
		return err
	}

	return convertToValue(result, v)
}

func MarshalIndent(v interface{}, indent string) ([]byte, error) {
	opts := &encoder.Options{
		Indent: indent,
	}
	var buf []byte
	writer := &byteWriter{buf: &buf}
	enc := encoder.NewEncoder(writer, opts)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return *writer.buf, nil
}

type byteWriter struct {
	buf *[]byte
}

// Encode serializa um objeto Go para o formato TOON como string
func Encode(v interface{}) (string, error) {
	// Detecta se é slice de map[string]interface{}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		// Objeto simples
		keys := make([]string, 0, rv.Len())
		for _, k := range rv.MapKeys() {
			keys = append(keys, fmt.Sprint(k.Interface()))
		}
		sort.Strings(keys)
		var buf bytes.Buffer
		for _, k := range keys {
			fmt.Fprintf(&buf, "%s: %v\n", k, rv.MapIndex(reflect.ValueOf(k)).Interface())
		}
		return buf.String(), nil
	case reflect.Struct:
		// Struct simples
		t := rv.Type()
		var buf bytes.Buffer
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := rv.Field(i).Interface()
			fmt.Fprintf(&buf, "%s: %v\n", field.Name, value)
		}
		return buf.String(), nil
	case reflect.Slice:
		if rv.Len() == 0 {
			return "[]\n", nil
		}
		first := rv.Index(0)
		// Tabular para slice de struct
		if first.Kind() == reflect.Struct {
			t := first.Type()
			fields := make([]string, t.NumField())
			for i := 0; i < t.NumField(); i++ {
				fields[i] = t.Field(i).Name
			}
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "[%d,]{%s}:\n", rv.Len(), joinComma(fields))
			for i := 0; i < rv.Len(); i++ {
				row := rv.Index(i)
				vals := make([]string, len(fields))
				for j := range fields {
					vals[j] = fmt.Sprint(row.Field(j).Interface())
				}
				fmt.Fprintf(&buf, "  %s\n", joinComma(vals))
			}
			return buf.String(), nil
		}
		// Tabular para slice de map[string]interface{}
		if first.Kind() == reflect.Map {
			fields := make([]string, 0, first.Len())
			for _, k := range first.MapKeys() {
				fields = append(fields, fmt.Sprint(k.Interface()))
			}
			sort.Strings(fields)
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "[%d,]{%s}:\n", rv.Len(), joinComma(fields))
			for i := 0; i < rv.Len(); i++ {
				row := rv.Index(i)
				vals := make([]string, len(fields))
				for j, f := range fields {
					v := row.MapIndex(reflect.ValueOf(f))
					vals[j] = fmt.Sprint(v.Interface())
				}
				fmt.Fprintf(&buf, "  %s\n", joinComma(vals))
			}
			return buf.String(), nil
		}
		return "", fmt.Errorf("Tipo de slice não suportado para Encode: %s", first.Kind())
	default:
		return "", fmt.Errorf("Tipo não suportado para Encode: %s", rv.Kind())
	}
}

func joinComma(arr []string) string {
	return strings.Join(arr, ",")
}

// Decode desserializa uma string TOON para um objeto Go
func Decode(data string, v interface{}) error {
	return Unmarshal([]byte(data), v)
}

func (w *byteWriter) Write(p []byte) (n int, err error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}

type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}

	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func convertToValue(src interface{}, dst interface{}) error {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return &InvalidUnmarshalError{reflect.TypeOf(dst)}
	}
	if dstValue.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(dst)}
	}

	dstValue = dstValue.Elem()
	srcValue := reflect.ValueOf(src)

	return setFieldValue(dstValue, srcValue)
}

func setFieldValue(dst, src reflect.Value) error {
	if !dst.CanSet() {
		return nil
	}

	switch dst.Kind() {
	case reflect.Interface:
		dst.Set(src)

	case reflect.Ptr:
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))
			return nil
		}
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		return setFieldValue(dst.Elem(), src.Elem())

	case reflect.Struct:
		if src.Kind() == reflect.Map {
			return setStructFromMap(dst, src)
		} else if src.Kind() == reflect.Interface {
			if srcMap, ok := src.Interface().(map[string]interface{}); ok {
				srcMapValue := reflect.ValueOf(srcMap)
				return setStructFromMap(dst, srcMapValue)
			}
		}

	case reflect.Map:
		if src.Kind() == reflect.Map {
			return setMapFromMap(dst, src)
		}

	case reflect.Slice:
		if src.Kind() == reflect.Slice {
			return setSliceFromSlice(dst, src)
		}

	case reflect.String:
		if src.Kind() == reflect.String {
			dst.SetString(src.String())
		} else {
			dst.SetString(formatValue(src.Interface()))
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if src.CanInt() {
			dst.SetInt(src.Int())
		} else if src.Kind() == reflect.Float64 {
			floatVal := src.Float()
			if float64(int64(floatVal)) == floatVal {
				dst.SetInt(int64(floatVal))
			}
		} else if src.Kind() == reflect.Interface {
			if val, ok := src.Interface().(int64); ok {
				dst.SetInt(val)
			} else if val, ok := src.Interface().(int); ok {
				dst.SetInt(int64(val))
			} else if val, ok := src.Interface().(float64); ok {
				dst.SetInt(int64(val))
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if src.CanUint() {
			dst.SetUint(src.Uint())
		} else if src.Kind() == reflect.Float64 {
			dst.SetUint(uint64(src.Float()))
		} else if src.Kind() == reflect.Interface {
			if val, ok := src.Interface().(uint64); ok {
				dst.SetUint(val)
			} else if val, ok := src.Interface().(uint); ok {
				dst.SetUint(uint64(val))
			} else if val, ok := src.Interface().(float64); ok {
				dst.SetUint(uint64(val))
			}
		}

	case reflect.Float32, reflect.Float64:
		if src.CanFloat() {
			dst.SetFloat(src.Float())
		} else if src.CanInt() {
			dst.SetFloat(float64(src.Int()))
		} else if src.Kind() == reflect.Interface {
			if val, ok := src.Interface().(float64); ok {
				dst.SetFloat(val)
			} else if val, ok := src.Interface().(float32); ok {
				dst.SetFloat(float64(val))
			} else if val, ok := src.Interface().(int64); ok {
				dst.SetFloat(float64(val))
			}
		}

	case reflect.Bool:
		if src.Kind() == reflect.Bool {
			dst.SetBool(src.Bool())
		} else if src.Kind() == reflect.Interface {
			if val, ok := src.Interface().(bool); ok {
				dst.SetBool(val)
			}
		}
	}

	return nil
}

func setStructFromMap(dst, src reflect.Value) error {
	dstType := dst.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		dstField := dst.Field(i)

		if !dstField.CanSet() {
			continue
		}

		fieldName := field.Tag.Get("toon")
		if fieldName == "" {
			fieldName = field.Tag.Get("json")
		}
		if fieldName == "" {
			fieldName = field.Name // Use actual field name, not lowercase
		}

		srcValue := src.MapIndex(reflect.ValueOf(fieldName))
		if !srcValue.IsValid() {
			continue
		}

		if err := setFieldValue(dstField, srcValue); err != nil {
			return err
		}
	}

	return nil
}

func setMapFromMap(dst, src reflect.Value) error {
	dstKeyType := dst.Type().Key()
	dstValueType := dst.Type().Elem()

	dst.Set(reflect.MakeMap(dst.Type()))

	for _, key := range src.MapKeys() {
		srcValue := src.MapIndex(key)

		dstKey := reflect.New(dstKeyType).Elem()
		if err := setFieldValue(dstKey, key); err != nil {
			return err
		}

		dstValue := reflect.New(dstValueType).Elem()
		if err := setFieldValue(dstValue, srcValue); err != nil {
			return err
		}

		dst.SetMapIndex(dstKey, dstValue)
	}

	return nil
}

func setSliceFromSlice(dst, src reflect.Value) error {
	dstSlice := reflect.MakeSlice(dst.Type(), src.Len(), src.Len())

	for i := 0; i < src.Len(); i++ {
		srcElem := src.Index(i)
		dstElem := dstSlice.Index(i)

		if err := setFieldValue(dstElem, srcElem); err != nil {
			return err
		}
	}

	dst.Set(dstSlice)
	return nil
}

func formatValue(v interface{}) string {
	if v == nil {
		return "null"
	}

	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(val).String()
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(val).String()
	case float32, float64:
		return reflect.ValueOf(val).String()
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "toon: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Ptr {
		return "toon: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "toon: Unmarshal(nil " + e.Type.String() + ")"
}
