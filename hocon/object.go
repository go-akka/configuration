package hocon

import (
	"bytes"
	"fmt"
	"strings"
)

type HoconObject struct {
	items map[string]*HoconValue
}

func NewHoconObject() *HoconObject {
	return &HoconObject{
		items: make(map[string]*HoconValue),
	}
}

func (p *HoconObject) GetString() string {
	panic("This element is an object and not a string.")
}

func (p *HoconObject) IsArray() bool {
	return false
}

func (p *HoconObject) GetArray() []*HoconValue {
	panic("This element is an object and not an array.")
}

func (p *HoconObject) GetKey(key string) *HoconValue {
	value, _ := p.items[key]
	return value
}

func (p *HoconObject) GetOrCreateKey(key string) *HoconValue {

	if value, exist := p.items[key]; exist {
		return value
	}

	child := NewHoconValue()
	p.items[key] = child
	return child
}

func (p *HoconObject) IsString() bool {
	return false
}

func (p *HoconObject) String() string {
	return p.ToString(0)
}

func (p *HoconObject) ToString(indent int) string {
	tmp := strings.Repeat(" ", indent*2)
	buf := bytes.NewBuffer(nil)
	for k, v := range p.items {
		key := p.quoteIfNeeded(k)
		buf.WriteString(fmt.Sprintf("%s%s : %s\r\n", tmp, key, v.ToString(indent)))
	}
	return buf.String()
}

func (p *HoconObject) Merge(other *HoconObject) {
	thisItems := p.items
	otherItems := other.items

	for otherkey, otherValue := range otherItems {
		if thisValue, exist := thisItems[otherkey]; exist {
			if thisValue.IsObject() && otherValue.IsObject() {
				thisValue.GetObject().Merge(otherValue.GetObject())
			}
		} else {
			p.items[otherkey] = otherValue
		}
	}
}

func (p *HoconObject) quoteIfNeeded(text string) string {
	if strings.IndexByte(text, ' ') >= 0 ||
		strings.IndexByte(text, '\t') >= 0 {
		return "\"" + text + "\""
	}
	return text
}
