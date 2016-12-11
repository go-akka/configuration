package hocon

type IncludeCallback func(filename string) *HoconRoot

type Parser struct {
	reader   *HoconTokenizer
	root     *HoconValue
	callback IncludeCallback

	substitutions []*HoconSubstitution
}

func Parse(text string, callback IncludeCallback) *HoconRoot {
	return new(Parser).parseText(text, callback)
}

func (p *Parser) parseText(text string, callback IncludeCallback) *HoconRoot {
	p.callback = callback
	p.root = NewHoconValue()
	p.reader = NewHoconTokenizer(text)
	p.reader.PullWhitespaceAndComments()
	p.parseObject(p.root, true, "")

	root := NewHoconRoot(p.root)

	c := NewConfigFromRoot(root)
	for _, sub := range p.substitutions {
		res := c.GetValue(sub.Path)
		if res == nil {
			panic("Unresolved substitution:" + sub.Path)
		}
		sub.ResolvedValue = res
	}

	return NewHoconRoot(p.root, p.substitutions...)
}

func (p *Parser) parseObject(owner *HoconValue, root bool, currentPath string) {
	if !owner.IsObject() {
		owner.NewValue(NewHoconObject())
	}

	currentObject := owner.GetObject()

	for !p.reader.EOF() {
		t := p.reader.PullNext()

		switch t.tokenType {
		case TokenTypeInclude:
			included := p.callback(t.value)
			substitutions := included.substitutions
			for _, substitution := range substitutions {
				substitution.Path = currentPath + "." + substitution.Path
			}
			p.substitutions = append(p.substitutions, substitutions...)
			otherObj := included.value.GetObject()
			owner.GetObject().Merge(otherObj)
		case TokenTypeEoF:
		case TokenTypeKey:
			value := currentObject.GetOrCreateKey(t.value)
			nextPath := t.value
			if len(currentPath) > 0 {
				nextPath = currentPath + "." + t.value
			}
			p.parseKeyContent(value, nextPath)
			if !root {
				return
			}

		case TokenTypeObjectEnd:
			return
		}
	}
}

func (p *Parser) parseKeyContent(value *HoconValue, currentPath string) {
	for !p.reader.EOF() {
		t := p.reader.PullNext()
		switch t.tokenType {
		case TokenTypeDot:
			p.parseObject(value, false, currentPath)
		case TokenTypeAssign:
			{
				if !value.IsObject() {
					value.Clear()
				}
			}
			p.ParseValue(value, currentPath)
			return
		case TokenTypeObjectStart:
			p.parseObject(value, true, currentPath)
			return
		}
	}
}

func (p *Parser) ParseValue(owner *HoconValue, currentPath string) {
	if p.reader.EOF() {
		panic("End of file reached while trying to read a value")
	}

	p.reader.PullWhitespaceAndComments()
	for p.reader.isValue() {
		t := p.reader.PullValue()

		switch t.tokenType {
		case TokenTypeEoF:
		case TokenTypeLiteralValue:
			if owner.IsObject() {
				owner.Clear()
			}
			lit := NewHoconLiteral(t.value)
			owner.AppendValue(lit)
		case TokenTypeObjectStart:
			p.parseObject(owner, true, currentPath)
		case TokenTypeArrayStart:
			arr := p.ParseArray(currentPath)
			owner.AppendValue(&arr)
		case TokenTypeSubstitute:
			sub := p.ParseSubstitution(t.value)
			p.substitutions = append(p.substitutions, sub)
			owner.AppendValue(sub)
		}

		if p.reader.IsSpaceOrTab() {
			p.ParseTrailingWhitespace(owner)
		}
	}
	p.ignoreComma()
}

func (p *Parser) ParseTrailingWhitespace(owner *HoconValue) {
	ws := p.reader.PullSpaceOrTab()
	if len(ws.value) > 0 {
		wsList := NewHoconLiteral(ws.value)
		owner.values = append(owner.values, wsList)
	}
}

func (p *Parser) ParseSubstitution(value string) *HoconSubstitution {
	return NewHoconSubstitution(value)
}

func (p *Parser) ParseArray(currentPath string) HoconArray {
	arr := NewHoconArray()
	for !p.reader.EOF() && !p.reader.IsArrayEnd() {
		v := NewHoconValue()
		p.ParseValue(v, currentPath)
		arr.values = append(arr.values, v)
		p.reader.PullWhitespaceAndComments()
	}
	p.reader.PullArrayEnd()
	return *arr
}

func (p *Parser) ignoreComma() {
	if p.reader.IsComma() {
		p.reader.PullComma()
	}
}
