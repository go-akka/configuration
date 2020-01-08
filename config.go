package configuration

import (
	"math/big"
	"strings"
	"time"

	"github.com/go-akka/configuration/hocon"
)

type Config struct {
	root          *hocon.HoconValue
	substitutions []*hocon.HoconSubstitution
	fallback      *Config
}

func NewConfigFromRoot(root *hocon.HoconRoot) *Config {
	if root.Value() == nil {
		panic("The root value cannot be null.")
	}

	return &Config{
		root:          root.Value(),
		substitutions: root.Substitutions(),
	}
}

func NewConfigFromConfig(source, fallback *Config) *Config {
	if source == nil {
		panic("The source configuration cannot be null.")
	}

	return &Config{
		root:     source.root,
		fallback: fallback,
	}
}

func (p *Config) IsEmpty() bool {
	return p == nil || p.root == nil || p.root.IsEmpty()
}

func (p *Config) Root() *hocon.HoconValue {
	return p.root
}

func (p *Config) Copy(fallback ...*Config) *Config {

	var fb *Config

	if p.fallback != nil {
		fb = p.fallback.Copy()
	} else {
		if len(fallback) > 0 {
			fb = fallback[0]
		}
	}
	return &Config{
		fallback:      fb,
		root:          p.root,
		substitutions: p.substitutions,
	}
}

func (p *Config) GetNode(path string) *hocon.HoconValue {
	if p == nil {
		return nil
	}

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

func (p *Config) GetBooleanSafely(path string, defaultVal ...bool) (bool, error) {
	obj := p.GetNode(path)
	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return false, nil
	}
	return obj.GetBooleanSafely()
}

func (p *Config) GetBoolean(path string, defaultVal ...bool) bool {
	val, err := p.GetBooleanSafely(path, defaultVal...)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetByteSize(path string) *big.Int {
	obj := p.GetNode(path)
	if obj == nil {
		return big.NewInt(-1)
	}
	return obj.GetByteSize()
}

func (p *Config) GetInt32Safely(path string, defaultVal ...int32) (int32, error) {
	obj := p.GetNode(path)
	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return 0, nil
	}
	return obj.GetInt32Safely()
}

func (p *Config) GetInt32(path string, defaultVal ...int32) int32 {
	val, err := p.GetInt32Safely(path, defaultVal...)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetInt64Safely(path string, defaultVal ...int64) (int64, error) {
	obj := p.GetNode(path)
	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return 0, nil
	}
	return obj.GetInt64Safely()
}

func (p *Config) GetInt64(path string, defaultVal ...int64) int64 {
	val, err := p.GetInt64Safely(path, defaultVal...)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetString(path string, defaultVal ...string) string {
	obj := p.GetNode(path)
	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return ""
	}
	return obj.GetString()
}

func (p *Config) GetFloat32Safely(path string, defaultVal ...float32) (float32, error) {
	obj := p.GetNode(path)
	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return 0, nil
	}
	return obj.GetFloat32Safely()
}

func (p *Config) GetFloat32(path string, defaultVal ...float32) float32 {
	val, err := p.GetFloat32Safely(path, defaultVal...)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetFloat64Safely(path string, defaultVal ...float64) (float64, error) {
	obj := p.GetNode(path)
	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return 0, nil
	}
	return obj.GetFloat64Safely()
}

func (p *Config) GetFloat64(path string, defaultVal ...float64) float64 {
	val, err := p.GetFloat64Safely(path, defaultVal...)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetTimeDuration(path string, defaultVal ...time.Duration) time.Duration {
	obj := p.GetNode(path)
	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	return obj.GetTimeDuration(true)
}

func (p *Config) GetTimeDurationInfiniteNotAllowed(path string, defaultVal ...time.Duration) time.Duration {
	obj := p.GetNode(path)
	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	return obj.GetTimeDuration(false)
}

func (p *Config) GetBooleanListSafely(path string) ([]bool, error) {
	obj := p.GetNode(path)
	if obj == nil {
		return nil, nil
	}
	return obj.GetBooleanListSafely()
}

func (p *Config) GetBooleanList(path string) []bool {
	val, err := p.GetBooleanListSafely(path)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetFloat32ListSafely(path string) ([]float32, error) {
	obj := p.GetNode(path)
	if obj == nil {
		return nil, nil
	}
	return obj.GetFloat32ListSafely()
}

func (p *Config) GetFloat32List(path string) []float32 {
	val, err := p.GetFloat32ListSafely(path)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetFloat64ListSafely(path string) ([]float64, error) {
	obj := p.GetNode(path)
	if obj == nil {
		return nil, nil
	}
	return obj.GetFloat64ListSafely()
}

func (p *Config) GetFloat64List(path string) []float64 {
	val, err := p.GetFloat64ListSafely(path)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetInt32ListSafely(path string) ([]int32, error) {
	obj := p.GetNode(path)
	if obj == nil {
		return nil, nil
	}
	return obj.GetInt32ListSafely()
}

func (p *Config) GetInt32List(path string) []int32 {
	val, err := p.GetInt32ListSafely(path)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetInt64ListSafely(path string) ([]int64, error) {
	obj := p.GetNode(path)
	if obj == nil {
		return nil, nil
	}
	return obj.GetInt64ListSafely()
}

func (p *Config) GetInt64List(path string) []int64 {
	val, err := p.GetInt64ListSafely(path)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetByteListSafely(path string) ([]byte, error) {
	obj := p.GetNode(path)
	if obj == nil {
		return nil, nil
	}
	return obj.GetByteListSafely()
}

func (p *Config) GetByteList(path string) []byte {
	val, err := p.GetByteListSafely(path)
	if err != nil {
		panic(err)
	}
	return val
}

func (p *Config) GetStringList(path string) []string {
	obj := p.GetNode(path)
	if obj == nil {
		return nil
	}
	return obj.GetStringList()
}

func (p *Config) GetConfig(path string) *Config {
	if p == nil {
		return nil
	}

	value := p.GetNode(path)
	if p.fallback != nil {
		f := p.fallback.GetConfig(path)
		if value == nil && f == nil {
			return nil
		}
		if value == nil {
			return f
		}
		return NewConfigFromRoot(hocon.NewHoconRoot(value)).WithFallback(f)
	}

	if value == nil {
		return nil
	}
	return NewConfigFromRoot(hocon.NewHoconRoot(value))
}

func (p *Config) GetValue(path string) *hocon.HoconValue {
	return p.GetNode(path)
}

func (p *Config) WithFallback(fallback *Config) *Config {
	if fallback == p {
		panic("Config can not have itself as fallback")
	}

	if fallback == nil {
		return p
	}

	mergedRoot := p.root.GetObject().MergeImmutable(fallback.root.GetObject())
	newRoot := hocon.NewHoconValue()

	newRoot.AppendValue(mergedRoot)

	mergedConfig := p.Copy(fallback)

	mergedConfig.root = newRoot

	return mergedConfig
}

func (p *Config) HasPath(path string) bool {
	return p.GetNode(path) != nil
}

func (p *Config) IsObject(path string) bool {
	node := p.GetNode(path)
	if node == nil {
		return false
	}

	return node.IsObject()
}

func (p *Config) IsArray(path string) bool {
	node := p.GetNode(path)
	if node == nil {
		return false
	}

	return node.IsArray()
}

func (p *Config) AddConfig(textConfig string, fallbackConfig *Config) *Config {
	root := hocon.Parse(textConfig, nil)
	config := NewConfigFromRoot(root)
	return config.WithFallback(fallbackConfig)
}

func (p *Config) AddConfigWithTextFallback(config *Config, textFallback string) *Config {
	fallbackRoot := hocon.Parse(textFallback, nil)
	fallbackConfig := NewConfigFromRoot(fallbackRoot)
	return config.WithFallback(fallbackConfig)
}

func (p Config) String() string {
	return p.root.String()
}

func splitDottedPathHonouringQuotes(path string) []string {
	tmp1 := strings.Split(path, "\"")
	var values []string
	for i := 0; i < len(tmp1); i++ {
		tmp2 := strings.Split(tmp1[i], ".")
		for j := 0; j < len(tmp2); j++ {
			if len(tmp2[j]) > 0 {
				values = append(values, tmp2[j])
			}
		}
	}
	return values
}
