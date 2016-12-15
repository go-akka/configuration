package hocon

type TokenType int

const (
	TokenTypeNone TokenType = iota
	TokenTypeComment
	TokenTypeKey
	TokenTypeLiteralValue
	TokenTypeAssign
	TokenTypeObjectStart
	TokenTypeObjectEnd
	TokenTypeDot
	TokenTypeEoF
	TokenTypeArrayStart
	TokenTypeArrayEnd
	TokenTypeComma
	TokenTypeSubstitute
	TokenTypeInclude
)

var (
	DefaultToken = Token{}
)

type Token struct {
	tokenType TokenType
	value     string
}

func NewToken(v interface{}) *Token {

	switch value := v.(type) {
	case string:
		{
			return &Token{tokenType: TokenTypeLiteralValue, value: value}
		}
	case TokenType:
		{
			return &Token{tokenType: value}
		}
	}

	return nil
}

func (p *Token) Key(key string) *Token {
	return &Token{tokenType: TokenTypeKey, value: key}
}

func (p *Token) Substitution(path string) *Token {
	return &Token{tokenType: TokenTypeSubstitute, value: path}
}

func (p *Token) LiteralValue(value string) *Token {
	return &Token{tokenType: TokenTypeLiteralValue, value: value}
}

func (p *Token) Include(path string) *Token {
	return &Token{tokenType: TokenTypeInclude, value: path}
}
