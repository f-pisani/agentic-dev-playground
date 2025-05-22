package feedbinapi

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// structToURLValues converts a struct to url.Values using `url` tags.
// It supports basic types (string, int, bool), pointers to them, slices of them, and time.Time.
func structToURLValues(s interface{}) (url.Values, error) {
	values := url.Values{}
	if s == nil {
		return values, nil
	}

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or a pointer to a struct, got %T", s)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("url")

		if tag == "" || tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		paramName := parts[0]
		omitempty := false
		commaSeparated := false

		for _, part := range parts[1:] {
			if part == "omitempty" {
				omitempty = true
			}
			if part == "comma" { // For list of IDs
				commaSeparated = true
			}
		}

		// Handle pointer fields: if nil and omitempty, skip. Otherwise, dereference.
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				if omitempty {
					continue
				}
				// If not omitempty, add empty value or handle as error?
				// For now, let's add empty if not omitempty and nil
				values.Add(paramName, "")
				continue
			}
			field = field.Elem() // Dereference pointer
		}

		// Skip if omitempty and zero value
		if omitempty && field.IsZero() && field.Kind() != reflect.Bool { // Bool zero value (false) can be significant
			continue
		}

		// Handle time.Time specifically for ISO 8601 formatting
		if field.Type() == reflect.TypeOf(time.Time{}) {
			if t, ok := field.Interface().(time.Time); ok {
				if omitempty && t.IsZero() {
					continue
				}
				values.Add(paramName, t.Format(time.RFC3339Nano)) // ISO 8601 like format
				continue
			}
		}

		switch field.Kind() {
		case reflect.String:
			values.Add(paramName, field.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			values.Add(paramName, strconv.FormatInt(field.Int(), 10))
		case reflect.Bool:
			values.Add(paramName, strconv.FormatBool(field.Bool()))
		case reflect.Slice:
			if commaSeparated {
				var strVals []string
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					switch elem.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						strVals = append(strVals, strconv.FormatInt(elem.Int(), 10))
					case reflect.String:
						strVals = append(strVals, elem.String())
					default:
						return nil, fmt.Errorf("unsupported slice element type for comma separation: %s", elem.Kind())
					}
				}
				if len(strVals) > 0 {
					values.Add(paramName, strings.Join(strVals, ","))
				} else if !omitempty {
					values.Add(paramName, "") // Add empty if not omitempty and slice is empty
				}
			} else {
				// For slices not comma separated, add multiple values with the same paramName
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					switch elem.Kind() {
					// Add cases as needed, e.g., string, int
					case reflect.String:
						values.Add(paramName, elem.String())
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						values.Add(paramName, strconv.FormatInt(elem.Int(), 10))
					default:
						return nil, fmt.Errorf("unsupported slice element type for param %s: %s", paramName, elem.Kind())
					}
				}
			}
		default:
			// Not supporting other types like struct, map, etc. for query params for now
			if !omitempty { // Only error if not omitempty and type is not supported
				return nil, fmt.Errorf("unsupported field type for param %s: %s", paramName, field.Kind())
			}
		}
	}
	return values, nil
}
