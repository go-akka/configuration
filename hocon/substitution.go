package hocon

type HoconSubstitution struct {
	Path          string
	ResolvedValue *HoconValue
}

func NewHoconSubstitution(path string) *HoconSubstitution {
	return &HoconSubstitution{Path: path}
}

func (p *HoconSubstitution) IsString() bool {
	return p.ResolvedValue.IsString()
}

func (p *HoconSubstitution) GetString() string {
	return p.ResolvedValue.GetString()
}

func (p *HoconSubstitution) IsArray() bool {
	return p.ResolvedValue.IsArray()
}
func (p *HoconSubstitution) GetArray() []*HoconValue {
	return p.ResolvedValue.GetArray()
}

func (p *HoconSubstitution) IsObject() bool {
	return p.ResolvedValue.IsObject()
}

func (p *HoconSubstitution) GetObject() *HoconObject {
	return p.ResolvedValue.GetObject()
}
