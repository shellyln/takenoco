package csv

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// toCsvRow writes a CSV row to the strings.Builder.
func toCsvRow(sb *strings.Builder, row []string) {
	for i := 0; i < len(row); i++ {
		if i != 0 {
			sb.WriteString(",")
		}
		if 0 <= strings.IndexRune(row[i], ',') ||
			0 <= strings.IndexRune(row[i], '\u000a') ||
			0 <= strings.IndexRune(row[i], '\u000b') ||
			0 <= strings.IndexRune(row[i], '\u000c') ||
			0 <= strings.IndexRune(row[i], '\u000d') ||
			0 <= strings.IndexRune(row[i], '\u0085') {
			sb.WriteString("\"")
			sb.WriteString(strings.ReplaceAll(row[i], "\"", "\"\""))
			sb.WriteString("\"")
		} else {
			sb.WriteString(row[i])
		}
	}
}

// ToCsv converts a 2d slice of string into a CSV string.
func ToCsv(a [][]string) string {
	var sb strings.Builder

	length := len(a)
	for i := 0; i < length; i++ {
		if i != 0 {
			sb.WriteString("\n")
		}
		toCsvRow(&sb, a[i])
	}
	return sb.String()
}

// Marshal converts a slice of struct into a CSV string.
func Marshal(cols []string, a interface{}) (string, error) {
	slice := reflect.ValueOf(a)
	k := slice.Kind()
	switch k {
	case reflect.Slice:
	default:
		return "", errors.New("Data should be slice of struct")
	}

	switch slice.Type().Elem().Kind() {
	case reflect.Struct:
	default:
		return "", errors.New("Data should be slice of struct")
	}

	var sb strings.Builder

	// make mappings
	// CSV column name -> CSV column index
	columns := make(map[string]int)
	for i := 0; i < len(cols); i++ {
		columns[cols[i]] = i
	}
	// Field name -> CSV column name
	mapping := make(map[string]string)
	// Field name -> Type
	colTypes := make(map[string]reflect.Kind)

	for i := 0; i < slice.Len(); i++ {
		rvi := slice.Index(i)
		rti := rvi.Type()
		for j := 0; j < rti.NumField(); j++ {
			rtf := rti.Field(j)
			csvColName := rtf.Tag.Get("csv")
			if csvColName == "" {
				csvColName = rtf.Name
			}
			mapping[rtf.Name] = csvColName
			colTypes[rtf.Name] = rtf.Type.Kind()
		}
		break
	}

	rowBuf := make([]string, len(cols))

	// set data
	for i := 0; i < slice.Len(); i++ {
		rvi := slice.Index(i)

		for j := 0; j < len(rowBuf); j++ {
			rowBuf[j] = ""
		}
		for k, v := range mapping {
			rvf := rvi.FieldByName(k)
			index := columns[v]
			switch colTypes[k] {
			case reflect.String:
				rowBuf[index] = rvf.String()
			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
				rowBuf[index] = strconv.FormatInt(rvf.Int(), 10)
			case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
				rowBuf[index] = strconv.FormatUint(rvf.Uint(), 10)
			case reflect.Float64, reflect.Float32:
				rowBuf[index] = strconv.FormatFloat(rvf.Float(), 'g', -1, 64)
			}
		}
		if i != 0 {
			sb.WriteString("\n")
		}
		toCsvRow(&sb, rowBuf)
	}

	return sb.String(), nil
}

// Convert from CSV string to slice of struct.
func Unmarshal(pa interface{}, s string) error {
	rv := reflect.ValueOf(pa)
	k := rv.Kind()
	switch k {
	case reflect.Ptr:
	default:
		return errors.New("Out parameter should be pointer to slice of struct")
	}

	// dereference the pointer
	rv2 := rv.Elem()
	k = rv2.Kind()
	switch k {
	case reflect.Slice:
	default:
		return errors.New("Out parameter should be pointer to slice of struct")
	}

	switch rv2.Type().Elem().Kind() {
	case reflect.Struct:
	default:
		return errors.New("Out parameter should be pointer to slice of struct")
	}

	data, err := Parse(s)
	if err != nil {
		return err
	}
	if len(data) < 2 {
		return nil
	}

	sliceLength := len(data) - 1
	slice := reflect.MakeSlice(rv2.Type(), sliceLength, sliceLength)

	// make mappings
	// CSV column name -> CSV column index
	columns := make(map[string]int)
	for i := 0; i < len(data[0]); i++ {
		columns[data[0][i]] = i
	}
	// CSV column name -> Field name
	mapping := make(map[string]string)
	// CSV column name -> Type
	colTypes := make(map[string]reflect.Kind)

	for i := 0; i < slice.Len(); i++ {
		rvi := slice.Index(i)
		rti := rvi.Type()
		for j := 0; j < rti.NumField(); j++ {
			rtf := rti.Field(j)
			tag := rtf.Tag.Get("csv")
			if _, ok := columns[tag]; ok {
				mapping[tag] = rtf.Name
				colTypes[tag] = rtf.Type.Kind()
			} else if _, ok := columns[rtf.Name]; ok {
				mapping[rtf.Name] = rtf.Name
				colTypes[rtf.Name] = rtf.Type.Kind()
			} else {
				continue
			}
		}
		break
	}

	// set data
	for i := 0; i < slice.Len(); i++ {
		rvi := slice.Index(i)
		for k, v := range mapping {
			rvf := rvi.FieldByName(v)
			w := data[i+1][columns[k]] // TODO: check slice length
			switch colTypes[k] {
			case reflect.String:
				rvf.SetString(w)
			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
				z, err := strconv.ParseInt(w, 10, 64)
				if err != nil {
					return err
				}
				rvf.SetInt(z)
			case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
				z, err := strconv.ParseUint(w, 10, 64)
				if err != nil {
					return err
				}
				rvf.SetUint(z)
			case reflect.Float64, reflect.Float32:
				z, err := strconv.ParseFloat(w, 64)
				if err != nil {
					return err
				}
				rvf.SetFloat(z)
			}
		}
	}

	// assign to pointer
	rv.Elem().Set(slice)

	return nil
}
