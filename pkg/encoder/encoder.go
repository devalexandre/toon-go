package encoder

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Options struct {
	Indent         string
	ForceTabular   bool
	MaxArraySize   int
	TokenOptimized bool
}

type Encoder struct {
	writer io.Writer
	opts   *Options
}

func NewEncoder(w io.Writer, opts *Options) *Encoder {
	if opts == nil {
		opts = &Options{
			Indent:         "  ",
			MaxArraySize:   1000,
			TokenOptimized: true,
		}
	}
	return &Encoder{
		writer: w,
		opts:   opts,
	}
}

func (e *Encoder) Encode(v interface{}) error {
	return e.encodeValue(v, 0, "")
}

func (e *Encoder) encodeValue(v interface{}, depth int, fieldName string) error {
	if v == nil {
		_, err := e.writer.Write([]byte("null"))
		return err
	}

	rv := reflect.ValueOf(v)
	kind := rv.Kind()

	if kind == reflect.Ptr {
		if rv.IsNil() {
			_, err := e.writer.Write([]byte("null"))
			return err
		}
		return e.encodeValue(rv.Elem().Interface(), depth, fieldName)
	}

	switch kind {
	case reflect.Bool:
		var err error
		if rv.Bool() {
			_, err = e.writer.Write([]byte("true"))
		} else {
			_, err = e.writer.Write([]byte("false"))
		}
		return err

	case reflect.String:
		str := rv.String()
		if e.needsQuotes(str) {
			_, err := e.writer.Write([]byte(fmt.Sprintf("\"%s\"", str)))
			return err
		}
		_, err := e.writer.Write([]byte(str))
		return err

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err := e.writer.Write([]byte(strconv.FormatInt(rv.Int(), 10)))
		return err

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		_, err := e.writer.Write([]byte(strconv.FormatUint(rv.Uint(), 10)))
		return err

	case reflect.Float32, reflect.Float64:
		str := strconv.FormatFloat(rv.Float(), 'g', -1, 64)
		_, err := e.writer.Write([]byte(str))
		return err

	case reflect.Slice, reflect.Array:
		return e.encodeArray(rv, depth, fieldName)

	case reflect.Map:
		return e.encodeMap(rv, depth, fieldName)

	case reflect.Struct:
		return e.encodeStruct(rv, depth, fieldName)

	default:
		str := fmt.Sprintf("%v", v)
		if e.needsQuotes(str) {
			_, err := e.writer.Write([]byte(fmt.Sprintf("\"%s\"", str)))
			return err
		}
		_, err := e.writer.Write([]byte(str))
		return err
	}
}

func (e *Encoder) needsQuotes(s string) bool {
	if s == "" {
		return true
	}
	return true
}

func (e *Encoder) encodeArray(rv reflect.Value, depth int, fieldName string) error {
	length := rv.Len()
	indent := strings.Repeat(e.opts.Indent, depth)

	if length == 0 {
		if fieldName != "" {
			_, err := fmt.Fprintf(e.writer, "%s%s[0]:\n", indent, fieldName)
			return err
		}
		_, err := fmt.Fprintf(e.writer, "%s[]\n", indent)
		return err
	}

	slice := make([]interface{}, length)
	for i := 0; i < length; i++ {
		slice[i] = rv.Index(i).Interface()
	}

	if e.opts.TokenOptimized && e.shouldUseTabularFormat(slice) {
		return e.encodeTabularArray(slice, depth, fieldName)
	}

	if fieldName != "" {
		if _, err := fmt.Fprintf(e.writer, "%s%s[%d]:\n", indent, fieldName, length); err != nil {
			return err
		}
	}

	for _, item := range slice {
		if err := e.encodeValue(item, depth+1, ""); err != nil {
			return err
		}
		if _, err := e.writer.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) shouldUseTabularFormat(slice []interface{}) bool {
	if len(slice) < 2 {
		return false
	}

	var firstFields []string

	for i, item := range slice {
		if obj, ok := item.(map[string]interface{}); ok {
			for _, v := range obj {
				if !e.isPrimitiveType(v) {
					return false
				}
			}

			fields := make([]string, 0, len(obj))
			for key := range obj {
				fields = append(fields, key)
			}
			sort.Strings(fields)

			if i == 0 {
				firstFields = fields
			} else {
				if len(fields) != len(firstFields) {
					return false
				}
				for j, field := range fields {
					if field != firstFields[j] {
						return false
					}
				}
			}
		} else if obj, ok := item.(map[string]string); ok {
			fields := make([]string, 0, len(obj))
			for key := range obj {
				fields = append(fields, key)
			}
			sort.Strings(fields)

			if i == 0 {
				firstFields = fields
			} else {
				if len(fields) != len(firstFields) {
					return false
				}
				for j, field := range fields {
					if field != firstFields[j] {
						return false
					}
				}
			}
		} else {
			return false
		}
	}

	return len(firstFields) > 0
}

func (e *Encoder) isPrimitiveType(v interface{}) bool {
	switch v.(type) {
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool, nil:
		return true
	default:
		return false
	}
}

func (e *Encoder) encodeTabularArray(slice []interface{}, depth int, fieldName string) error {
	if len(slice) == 0 {
		return nil
	}

	var fields []string
	if obj, ok := slice[0].(map[string]interface{}); ok {
		fields = make([]string, 0, len(obj))
		for key := range obj {
			fields = append(fields, key)
		}
		sort.Strings(fields)
	} else if obj, ok := slice[0].(map[string]string); ok {
		fields = make([]string, 0, len(obj))
		for key := range obj {
			fields = append(fields, key)
		}
		sort.Strings(fields)
	}

	if len(fields) == 0 {
		return e.encodeArray(reflect.ValueOf(slice), depth, fieldName)
	}

	indent := strings.Repeat(e.opts.Indent, depth)

	if _, err := fmt.Fprintf(e.writer, "%s%s[%d]{%s}:\n", indent, fieldName, len(slice), strings.Join(fields, ",")); err != nil {
		return err
	}

	for _, item := range slice {
		if _, err := fmt.Fprintf(e.writer, "%s%s", indent, e.opts.Indent); err != nil {
			return err
		}

		if obj, ok := item.(map[string]interface{}); ok {
			values := make([]string, len(fields))
			for i, field := range fields {
				value := obj[field]
				values[i] = e.formatValueForTabular(value)
			}
			if _, err := e.writer.Write([]byte(strings.Join(values, " "))); err != nil {
				return err
			}
		} else if obj, ok := item.(map[string]string); ok {
			values := make([]string, len(fields))
			for i, field := range fields {
				value := obj[field]
				values[i] = e.formatValueForTabular(value)
			}
			if _, err := e.writer.Write([]byte(strings.Join(values, " "))); err != nil {
				return err
			}
		}
		if _, err := e.writer.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) formatValueForTabular(v interface{}) string {
	if v == nil {
		return "null"
	}

	switch val := v.(type) {
	case string:
		if e.needsQuotes(val) {
			return fmt.Sprintf("\"%s\"", val)
		}
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", val)
	case float32, float64:
		return fmt.Sprintf("%v", val)
	default:
		str := fmt.Sprintf("%v", val)
		if e.needsQuotes(str) {
			return fmt.Sprintf("\"%s\"", str)
		}
		return str
	}
}

func (e *Encoder) encodeMap(rv reflect.Value, depth int, fieldName string) error {
	indent := strings.Repeat(e.opts.Indent, depth)

	keys := rv.MapKeys()
	keyStrings := make([]string, len(keys))
	for i, key := range keys {
		keyStrings[i] = key.String()
	}
	sort.Strings(keyStrings)

	if fieldName != "" {
		if _, err := fmt.Fprintf(e.writer, "%s%s:\n", indent, fieldName); err != nil {
			return err
		}
	}

	for _, keyStr := range keyStrings {
		var value reflect.Value
		for _, key := range keys {
			if key.String() == keyStr {
				value = rv.MapIndex(key)
				break
			}
		}

		if !value.IsValid() {
			continue
		}

		valueInterface := value.Interface()

		if e.isNestedType(valueInterface) {
			if err := e.encodeValue(valueInterface, depth+1, keyStr); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(e.writer, "%s%s: ", indent+e.opts.Indent, keyStr); err != nil {
				return err
			}
			if err := e.encodeValue(valueInterface, 0, ""); err != nil {
				return err
			}
			if _, err := e.writer.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Encoder) encodeStruct(rv reflect.Value, depth int, structFieldName string) error {
	indent := strings.Repeat(e.opts.Indent, depth)

	rt := rv.Type()
	numField := rv.NumField()

	if structFieldName != "" {
		if _, err := fmt.Fprintf(e.writer, "%s%s:\n", indent, structFieldName); err != nil {
			return err
		}
	}

	type fieldInfo struct {
		name  string
		value reflect.Value
		field reflect.StructField
	}

	var validFields []fieldInfo

	for i := 0; i < numField; i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		if !value.CanInterface() {
			continue
		}

		fieldName := field.Name
		if toonTag := field.Tag.Get("toon"); toonTag != "" && toonTag != "-" {
			if commaIndex := strings.Index(toonTag, ","); commaIndex != -1 {
				fieldName = toonTag[:commaIndex]
			} else {
				fieldName = toonTag
			}
		} else if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
				fieldName = jsonTag[:commaIndex]
			} else {
				fieldName = jsonTag
			}
		}

		validFields = append(validFields, fieldInfo{
			name:  fieldName,
			value: value,
			field: field,
		})
	}

	for _, fieldInfo := range validFields {
		fieldValue := fieldInfo.value.Interface()
		fieldName := fieldInfo.name

		if e.isNestedType(fieldValue) {
			if err := e.encodeValue(fieldValue, depth+1, fieldName); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(e.writer, "%s%s: ", indent+e.opts.Indent, fieldName); err != nil {
				return err
			}
			if err := e.encodeValue(fieldValue, 0, ""); err != nil {
				return err
			}
			if _, err := e.writer.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Encoder) isNestedType(v interface{}) bool {
	if v == nil {
		return false
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
		return true
	case reflect.Ptr:
		return !rv.IsNil() && e.isNestedType(rv.Elem().Interface())
	default:
		return false
	}
}
