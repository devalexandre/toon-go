package decoder

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Parser struct {
	scanner *bufio.Scanner
	lineNum int
	lines    []string
	linePos  int
	hasLines bool
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
		lineNum: 0,
	}
}

func (p *Parser) Parse() (interface{}, error) {
	var lines []string
	for p.scanner.Scan() {
		lines = append(lines, p.scanner.Text())
	}

	if err := p.scanner.Err(); err != nil {
		return nil, err
	}

	p.lines = lines
	p.linePos = 0
	p.hasLines = true

	minIndent := p.findMinIndent()

	return p.parseObject(minIndent)
}

func (p *Parser) findMinIndent() int {
	minIndent := -1
	for _, line := range p.lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		indent := 0
		for indent < len(line) && line[indent] == ' ' {
			indent++
		}

		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent == -1 {
		minIndent = 0
	}

	return minIndent
}

func (p *Parser) parseObject(currentIndent int) (interface{}, error) {
	obj := make(map[string]interface{})

	for p.linePos < len(p.lines) {
		line := p.lines[p.linePos]
		p.linePos++

		if strings.TrimSpace(line) == "" {
			continue
		}

		indent := 0
		for indent < len(line) && line[indent] == ' ' {
			indent++
		}

		if indent < currentIndent {
			p.linePos-- // Put the line back for parent to process
			break
		}

		if indent > currentIndent {
			p.linePos--
			break
		}

		content := strings.TrimSpace(line[currentIndent:])
		if content == "" {
			continue
		}

		if strings.Contains(content, "[") && strings.Contains(content, "}:") {
			name, array, err := p.parseTabularArray(content, indent)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", p.linePos, err)
			}
			obj[name] = array
			continue
		}

		if strings.Contains(content, "[") && strings.HasSuffix(content, ":") && !strings.Contains(content, "{") {
			name, array, err := p.parseRegularArray(content, indent)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", p.linePos, err)
			}
			obj[name] = array
			continue
		}

		if strings.Contains(content, ":") {
			name, value, err := p.parseKeyValue(content, indent)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", p.linePos, err)
			}
			obj[name] = value
			continue
		}
	}

	return obj, nil
}

func (p *Parser) parseKeyValue(line string, indent int) (string, interface{}, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("invalid key-value format")
	}

	key := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])

	if valueStr == "" && p.linePos < len(p.lines) {
		if p.linePos < len(p.lines) {
			nextLine := p.lines[p.linePos]
			nextIndent := 0
			for nextIndent < len(nextLine) && nextLine[nextIndent] == ' ' {
				nextIndent++
			}

			if nextIndent > indent {
				nestedObj, err := p.parseObject(nextIndent)
				if err != nil {
					return "", nil, err
				}
				return key, nestedObj, nil
			}
		}
	}

	value := p.parsePrimitive(valueStr)
	return key, value, nil
}

func (p *Parser) parseTabularArray(header string, indent int) (string, []interface{}, error) {
	colonIndex := strings.Index(header, "}:")
	if colonIndex == -1 {
		return "", nil, fmt.Errorf("invalid tabular array format")
	}

	headerPart := header[:colonIndex+1]

	nameEnd := strings.Index(headerPart, "[")
	if nameEnd == -1 {
		return "", nil, fmt.Errorf("invalid tabular array format: missing count")
	}

	name := strings.TrimSpace(headerPart[:nameEnd])

	countStart := nameEnd + 1
	countEnd := strings.Index(headerPart, "]")
	if countEnd == -1 {
		return "", nil, fmt.Errorf("invalid tabular array format: missing closing bracket")
	}

	countStr := headerPart[countStart:countEnd]
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return "", nil, fmt.Errorf("invalid count: %v", err)
	}

	fieldsStart := countEnd + 2 // Skip ']{'
	if fieldsStart >= len(headerPart) || headerPart[fieldsStart-1] != '{' {
		return "", nil, fmt.Errorf("invalid tabular array format: missing fields")
	}

	fieldsStr := headerPart[fieldsStart : len(headerPart)-1] // Exclude '}'
	fields := strings.Split(fieldsStr, ",")
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}

	var rows [][]string
	rowCount := 0

	for p.linePos < len(p.lines) && rowCount < count {
		rowLine := p.lines[p.linePos]

		rowIndent := 0
		for rowIndent < len(rowLine) && rowLine[rowIndent] == ' ' {
			rowIndent++
		}

		if rowIndent <= indent {
			break
		}

		dataContent := strings.TrimSpace(rowLine[rowIndent:])
		if dataContent == "" {
			p.linePos++
			continue
		}

		values := strings.Fields(dataContent)

		if len(values) != len(fields) {
			return "", nil, fmt.Errorf("row %d: field count mismatch (expected %d, got %d)",
				rowCount+1, len(fields), len(values))
		}

		rows = append(rows, values)
		rowCount++
		p.linePos++
	}

	if rowCount != count {
		return "", nil, fmt.Errorf("array count mismatch: declared %d, found %d", count, rowCount)
	}

	array := make([]interface{}, len(rows))
	for i, row := range rows {
		obj := make(map[string]interface{})
		for j, field := range fields {
			obj[field] = p.parsePrimitive(row[j])
		}
		array[i] = obj
	}

	return name, array, nil
}

func (p *Parser) parseRegularArray(header string, indent int) (string, []interface{}, error) {
	colonIndex := strings.Index(header, ":")
	if colonIndex == -1 {
		return "", nil, fmt.Errorf("invalid regular array format")
	}

	headerPart := strings.TrimSpace(header[:colonIndex])

	nameEnd := strings.Index(headerPart, "[")
	if nameEnd == -1 {
		return "", nil, fmt.Errorf("invalid regular array format: missing count")
	}

	name := strings.TrimSpace(headerPart[:nameEnd])

	countStart := nameEnd + 1
	countEnd := strings.Index(headerPart, "]")
	if countEnd == -1 {
		return "", nil, fmt.Errorf("invalid regular array format: missing closing bracket")
	}

	countStr := headerPart[countStart:countEnd]
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return "", nil, fmt.Errorf("invalid count: %v", err)
	}

	valuesStr := strings.TrimSpace(header[colonIndex+1:])

	var values []interface{}

	if valuesStr != "" {
		valueStrings := strings.Fields(valuesStr)
		for _, valueStr := range valueStrings {
			trimmed := strings.TrimSpace(valueStr)
			values = append(values, p.parsePrimitive(trimmed))
		}
	} else {
		valueCount := 0
		for p.linePos < len(p.lines) && valueCount < count {
			rowLine := p.lines[p.linePos]

			rowIndent := 0
			for rowIndent < len(rowLine) && rowLine[rowIndent] == ' ' {
				rowIndent++
			}

			if rowIndent <= indent {
				break
			}

			dataContent := strings.TrimSpace(rowLine[rowIndent:])
			if dataContent == "" {
				p.linePos++
				continue
			}

			valueStrings := strings.Fields(dataContent)
			for _, valueStr := range valueStrings {
				if valueCount >= count {
					break
				}
				trimmed := strings.TrimSpace(valueStr)
				values = append(values, p.parsePrimitive(trimmed))
				valueCount++
			}

			p.linePos++
		}
	}

	if count != len(values) {
		return "", nil, fmt.Errorf("array count mismatch: declared %d, found %d", count, len(values))
	}

	return name, values, nil
}

func (p *Parser) parsePrimitive(value string) interface{} {
	trimmed := strings.TrimSpace(value)

	if len(trimmed) >= 2 && trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"' {
		return trimmed[1 : len(trimmed)-1]
	}

	if trimmed == "true" {
		return true
	}
	if trimmed == "false" {
		return false
	}

	if trimmed == "null" {
		return nil
	}

	if num, err := strconv.ParseFloat(trimmed, 64); err == nil {
		if float64(int64(num)) == num {
			return int64(num)
		}
		return num
	}

	return trimmed
}
