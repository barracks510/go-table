package table

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const (
	Right = iota
	Left
	Center
)

var (
	ErrBadType           = errors.New("Type not supported")
	ErrColumnsNotMatched = errors.New("columns not matched")
)

func Table(header io.Reader, body io.Reader) io.Reader {
	out := new(bytes.Buffer)
	out.WriteString("<table>\n")
	if header != nil {
		out.WriteString("<thead>\n")
		io.Copy(out, header)
		out.WriteString("</thead>\n")
	}
	out.WriteString("\n<tbody>\n")
	io.Copy(out, body)
	out.WriteString("</tbody>\n</table>\n")
	return out
}

func TableRow(text io.Reader) io.Reader {
	out := new(bytes.Buffer)
	out.WriteString("<tr>\n")
	io.Copy(out, text)
	out.WriteString("\n</tr>\n")
	return out
}

func TableHeaderCell(text io.Reader, align int) io.Reader {
	out := new(bytes.Buffer)
	switch align {
	case Left:
		out.WriteString("<th align=\"left\">")
	case Right:
		out.WriteString("<th align=\"right\">")
	case Center:
		out.WriteString("<th align=\"center\">")
	default:
		out.WriteString("<th>")
	}

	io.Copy(out, text)
	out.WriteString("</th>")
	return out
}

func TableCell(text io.Reader, align int) io.Reader {
	out := new(bytes.Buffer)
	switch align {
	case Left:
		out.WriteString("<td align=\"left\">")
	case Right:
		out.WriteString("<td align=\"right\">")
	case Center:
		out.WriteString("<td align=\"center\">")
	default:
		out.WriteString("<td>")
	}

	io.Copy(out, text)
	out.WriteString("</td>")
	return out
}

func convert(src map[interface{}]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range src {
		m[fmt.Sprintf("%v", k)] = v
	}
	return m
}

func makeCell(value interface{}) io.Reader {
	switch value := value.(type) {
	case map[interface{}]interface{}:
		return makeTable(convert(value))
	case map[string]interface{}:
		return makeTable(value)
	case []interface{}:
		return makeArray(value)
	case string:
		return bytes.NewBufferString(value)
	default:
		// fmt.Printf("Value: %T %v\n", value, value)
		return bytes.NewBufferString(fmt.Sprintf("%v", value))
	}
}

func makeArray(data []interface{}) io.Reader {
	var row []io.Reader
	for _, value := range data {
		row = append(row, TableCell(makeCell(value), Center))
	}
	return Table(nil, TableRow(io.MultiReader(row...)))
}

func makeTable(data map[string]interface{}) io.Reader {
	var thead []io.Reader
	var tbody []io.Reader
	for key, value := range data {
		thead = append(thead, TableHeaderCell(makeCell(key), Center))
		tbody = append(tbody, TableCell(makeCell(value), Center))
	}
	return Table(TableRow(io.MultiReader(thead...)), TableRow(io.MultiReader(tbody...)))
}

func MakeTable(raw interface{}) (io.Reader, error) {
	if data, ok := raw.(map[string]interface{}); ok {
		return makeTable(data), nil
	}
	return nil, ErrBadType
}

func MakeMultiRowTable(head []interface{}, body [][]interface{}) (io.Reader, error) {
	var thead []io.Reader
	var tbody []io.Reader
	var trow []io.Reader
	columnCount := len(head)
	for _, value := range head {
		thead = append(thead, TableHeaderCell(makeCell(value), Center))
	}
	for _, row := range body {
		trow = nil
		if columnCount != len(row) {
			return nil, ErrColumnsNotMatched
		}
		for _, value := range row {
			trow = append(trow, TableCell(makeCell(value), Center))
		}
		tbody = append(tbody, TableRow(io.MultiReader(trow...)))
	}

	return Table(TableRow(io.MultiReader(thead...)), io.MultiReader(tbody...)), nil
}
