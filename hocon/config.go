package hocon

import (
	"strings"
)

type Config struct {
	root          *HoconValue
	substitutions []*HoconSubstitution
	fallback      *Config
}

func NewConfigFromRoot(root *HoconRoot) Config {
	if root.Value() == nil {
		panic("The root value cannot be null.")
	}

	return Config{
		root:          root.Value(),
		substitutions: root.Substitutions(),
	}
}

func NewConfigFromConfig(source, fallback *Config) Config {
	if source == nil {
		panic("The source configuration cannot be null.")
	}

	return Config{
		root:     source.root,
		fallback: fallback,
	}
}

func (p *Config) IsEmpty() bool {
	return p.root == nil || p.root.IsEmpty()
}

func (p *Config) Root() HoconValue {
	return *p.root
}

func (p *Config) Copy() Config {

	var fallback Config
	if p.fallback != nil {
		fallback = p.fallback.Copy()
	}
	return Config{
		fallback:      &fallback,
		root:          p.root,
		substitutions: p.substitutions,
	}
}

func (p *Config) GetValue(path string) *HoconValue {
	return p.GetNode(path)
}

func (p *Config) GetNode(path string) *HoconValue {
	elements := splitDottedPathHonouringQuotes(path)
	currentNode := p.root

	if currentNode == nil {
		panic("Current node should not be null")
	}

	for _, key := range elements {
		currentNode = currentNode.GetChildObject(key)
		if currentNode == nil {
			if p.fallback != nil {
				return p.fallback.GetNode(path)
			}
			return nil
		}
	}
	return currentNode
}

func (p *Config) GetString(path string) string {
	if obj := p.GetNode(path); obj != nil {
		return obj.GetString()
	}
	return ""
}

// TODO:
func splitDottedPathHonouringQuotes(path string) []string {
	return strings.Split(path, ".")
}
