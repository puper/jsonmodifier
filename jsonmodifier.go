package jsonmodifier

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/tidwall/gjson"
)

func JsonModify(v any) *jsonModifier {
	return &jsonModifier{any: v}
}

const (
	ModifyTypeOnly = iota
	ModifyTypeExcept
)

type jsonModifier struct {
	any
	typ    int
	fields []string
}

func (me *jsonModifier) Only(fields ...string) *jsonModifier {
	me.typ = ModifyTypeOnly
	me.fields = fields
	return me
}

func (me *jsonModifier) Except(fields ...string) *jsonModifier {
	me.typ = ModifyTypeExcept
	me.fields = fields
	return me
}

type Field struct {
	Name     string
	Children map[string]*Field
}

func (me *jsonModifier) processValue(v gjson.Result, fields map[string]*Field, buf *bytes.Buffer) {
	if v.IsObject() {
		buf.WriteByte('{')
		isFirst := true
		for k, v1 := range v.Map() {
			if me.typ == ModifyTypeOnly {
				if f, ok := fields[k]; ok {
					if isFirst {
						isFirst = false
					} else {
						buf.WriteByte(',')
					}
					buf.WriteString(`"` + EncodeKey(k) + `":`)
					if f.Children != nil {
						me.processValue(v1, f.Children, buf)
					} else {
						buf.WriteString(v1.Raw)
					}
				}
			} else {
				if f, ok := fields[k]; ok {
					if f.Children != nil {
						if isFirst {
							isFirst = false
						} else {
							buf.WriteByte(',')
						}
						buf.WriteString(`"` + EncodeKey(k) + `":`)
						me.processValue(v1, f.Children, buf)
					}
				} else {
					if isFirst {
						isFirst = false
					} else {
						buf.WriteByte(',')
					}
					buf.WriteString(`"` + EncodeKey(k) + `":` + v1.Raw)
				}
			}

		}
		buf.WriteByte('}')
	} else if v.IsArray() {
		buf.WriteByte('[')
		for _, v1 := range v.Array() {
			me.processValue(v1, fields, buf)
		}
		buf.WriteByte(']')
	} else {
		buf.WriteString(v.Raw)
	}
}

func (me *jsonModifier) MarshalJSON() ([]byte, error) {
	field := &Field{
		Children: map[string]*Field{},
	}
	for _, f := range me.fields {
		tmp := field
		for _, c := range strings.Split(f, ".") {
			if tmp.Children == nil {
				tmp.Children = map[string]*Field{}
			}
			if tmp.Children[c] == nil {
				tmp.Children[c] = &Field{
					Name: c,
				}
			}
			tmp = tmp.Children[c]
		}
	}
	b, err := json.Marshal(me.any)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	me.processValue(gjson.ParseBytes(b), field.Children, buf)
	return buf.Bytes(), nil
}

func EncodeKey(key string) string {
	//key = strings.ReplaceAll(key, `\`, `\\`)
	//key = strings.ReplaceAll(key, `"`, `\\"`)
	return key

}
