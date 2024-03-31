package jsonmodifier

import (
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

func (me *jsonModifier) processValue(v gjson.Result, fields map[string]*Field) any {
	if v.IsObject() {
		reply := make(map[string]any)
		for k, v1 := range v.Map() {
			if me.typ == ModifyTypeOnly {
				if f, ok := fields[k]; ok {
					if f.Children != nil {
						reply[k] = me.processValue(v1, f.Children)
					} else {
						reply[k] = json.RawMessage(v1.Raw)
					}
				}
			} else {
				if f, ok := fields[k]; ok {
					if f.Children != nil {
						reply[k] = me.processValue(v1, f.Children)
					}
				} else {
					reply[k] = json.RawMessage(v1.Raw)
				}
			}

		}
		return reply
	} else if v.IsArray() {
		reply := make([]any, 0)
		for _, v1 := range v.Array() {
			reply = append(reply, me.processValue(v1, fields))
		}
		return reply
	} else {
		return json.RawMessage(v.Raw)
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
	v := gjson.ParseBytes(b)
	reply := me.processValue(v, field.Children)
	b, err = json.Marshal(reply)
	return b, err
}
