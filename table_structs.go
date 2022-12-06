package tablewriter

import (
	"errors"
	"fmt"
	"reflect"
)

// SetStructs sets header and rows from slice of struct.
//
// If something that is not a slice is passed, an error will be returned.
//
// The tag specified by "tablewriter" for the struct becomes the header.
// If not specified or empty, the field name will be used.
//
// The field of the first element of the slice is used as the header.
// If the element implements fmt.Stringer, the result will be used.
// And the slice contains nil, it will be skipped without rendering.
func (t *Table) SetStructs(v interface{}) error {
	if v == nil {
		return errors.New("nil value")
	}

	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)

	switch vt.Kind() {
	case reflect.Slice, reflect.Array:
		return t.setStructSlice(vv)
	default:
		return fmt.Errorf("invalid type %T", v)
	}
}

func (t *Table) setStructSlice(vv reflect.Value) error {
	if vv.Len() < 1 {
		return errors.New("empty value")
	}

	// check first element to set header
	first := vv.Index(0)
	e, err := getElementType(first)
	if err != nil {
		return err
	}

	n := e.NumField()
	headers := make([]string, n)
	for i := 0; i < n; i++ {
		f := e.Field(i)
		header := f.Tag.Get("tablewriter") // TODO: option to set user-defined tag
		if header == "" {
			header = f.Name
		}
		headers[i] = header
	}

	t.header = headers

	for i := 0; i < vv.Len(); i++ {
		item := reflect.Indirect(vv.Index(i))
		itemType := reflect.TypeOf(item)
		switch itemType.Kind() {
		case reflect.Struct:
			// OK
		default:
			return fmt.Errorf("invalid item type %v", itemType.Kind())
		}
		if !item.IsValid() {
			// skip rendering
			continue
		}
		nf := item.NumField()
		if n != nf {
			return errors.New("invalid num of field")
		}
		rows := make([]string, nf)
		for j := 0; j < nf; j++ {
			f := reflect.Indirect(item.Field(j))
			if f.Kind() == reflect.Ptr {
				f = f.Elem()
			}
			if f.IsValid() {
				if s, ok := f.Interface().(fmt.Stringer); ok {
					rows[j] = s.String()
					continue
				}
				rows[j] = fmt.Sprint(f)
			} else {
				rows[j] = "nil"
			}
		}

		t.Append(rows)
	}

	return nil
}

func getElementType(first reflect.Value) (reflect.Type, error) {
	e := first.Type()
	switch e.Kind() {
	case reflect.Struct:
		// OK
	case reflect.Ptr:
		if first.IsNil() {
			return e, errors.New("the first element is nil")
		}
		e = first.Elem().Type()
		if e.Kind() != reflect.Struct {
			return e, fmt.Errorf("invalid kind %s", e.Kind())
		}
	default:
		return e, fmt.Errorf("invalid kind %s", e.Kind())
	}

	return e, nil
}
