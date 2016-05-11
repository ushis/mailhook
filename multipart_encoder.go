package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
)

type UnsupportedType struct {
	t reflect.Type
}

func (err UnsupportedType) Error() string {
	return fmt.Sprintf("unsupported type: %s", err.t.String())
}

type Marshaler interface {
	MarshalMultipart() string
}

type FileMarshaler interface {
	MarshalMultipart() (string, []byte)
}

type MultipartEncoder struct {
	*multipart.Writer
}

func NewMultipartEncoder(w io.Writer) *MultipartEncoder {
	return &MultipartEncoder{multipart.NewWriter(w)}
}

func (me *MultipartEncoder) Encode(field string, v interface{}) error {
	return me.encode(field, reflect.ValueOf(v))
}

func (me *MultipartEncoder) encode(field string, v reflect.Value) error {
	if v.CanInterface() {
		i := v.Interface()

		if m, ok := i.(Marshaler); ok {
			return me.WriteField(field, m.MarshalMultipart())
		}

		if m, ok := i.(FileMarshaler); ok {
			filename, body := m.MarshalMultipart()
			return me.encodeFile(field, filename, body)
		}
	}

	switch v.Kind() {
	case reflect.Ptr:
		return me.encode(field, v.Elem())
	case reflect.Bool:
		return me.WriteField(field, strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return me.WriteField(field, strconv.FormatInt(v.Int(), 10))
	case reflect.Float32, reflect.Float64:
		return me.WriteField(field, strconv.FormatFloat(v.Float(), 'f', -1, 64))
	case reflect.String:
		return me.WriteField(field, v.String())
	case reflect.Slice, reflect.Array:
		return me.encodeSlice(field, v)
	case reflect.Map:
		return me.encodeMap(field, v)
	case reflect.Struct:
		return me.encodeStruct(field, v)
	case reflect.Interface:
		return me.Encode(field, v.Interface())
	default:
		return UnsupportedType{v.Type()}
	}
}

func (me *MultipartEncoder) encodeSlice(field string, v reflect.Value) error {
	for i := 0; i < v.Len(); i++ {
		if err := me.encode(fmt.Sprintf("%s[]", field), v.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (me *MultipartEncoder) encodeMap(field string, v reflect.Value) error {
	if v.Type().Key().Kind() != reflect.String {
		return UnsupportedType{v.Type()}
	}

	for _, k := range v.MapKeys() {
		if err := me.encode(joinFields(field, k.String()), v.MapIndex(k)); err != nil {
			return err
		}
	}
	return nil
}

func (me *MultipartEncoder) encodeStruct(field string, v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name, _ := parseTag(f.Tag.Get("multipart"))

		if name == "-" {
			continue
		}

		if len(name) == 0 {
			name = f.Name
		}

		if err := me.encode(joinFields(field, name), v.Field(i)); err != nil {
			return err
		}
	}
	return nil
}

func (me *MultipartEncoder) encodeFile(field, filename string, body []byte) error {
	part, err := me.CreateFormFile(field, filename)

	if err != nil {
		return err
	}
	_, err = part.Write(body)
	return err
}

func joinFields(a, b string) string {
	if len(a) == 0 {
		return b
	}
	return fmt.Sprintf("%s[%s]", a, b)
}

func parseTag(tag string) (string, []string) {
	opts := strings.Split(tag, ",")

	for i, opt := range opts {
		opts[i] = strings.TrimSpace(opt)
	}

	if len(opts) == 0 {
		return "", opts
	}
	return opts[0], opts[1:]
}
