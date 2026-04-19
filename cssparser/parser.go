package cssparser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/css/scanner"
)

type CssParserState uint32

const (
	CssParserStateRoot CssParserState = iota
	CssParserStateSelector
	CssParserStateBlock
	CssParserStatePropName
	CssParserStatePropValue
)

type Declaration struct {
	Property string
	Value    string
}

type Rule struct {
	Selector     string
	Declarations []*Declaration
}

type CssParser struct {
	state              CssParserState
	scanner            *scanner.Scanner
	currentValue       *strings.Builder
	currentRule        *Rule
	currentDeclaration *Declaration
	parenthesesDepth   int

	Error  error
	Result []*Rule
}

func NewCssParser() *CssParser {
	return &CssParser{
		state:            CssParserStateRoot,
		parenthesesDepth: 0,
		Result:           make([]*Rule, 0),
	}
}

func (p *CssParser) Parse(input string) ([]*Rule, error) {
	p.Result = make([]*Rule, 0)
	p.scanner = scanner.New(input)

	for {
		token := p.scanner.Next()

		if p.Error != nil {
			return nil, p.Error
		}

		if token.Type == scanner.TokenEOF {
			if p.state == CssParserStateRoot {
				return p.Result, nil
			} else {
				return nil, errors.New("parser is in an invalid state")
			}
		}

		if token.Type == scanner.TokenAtKeyword {
			return nil, errors.New("parser does not support At-rules")
		}

		if token.Type == scanner.TokenError {
			return nil, errors.New(token.String())
		}

		switch p.state {
		case CssParserStateRoot:
			p.stateRoot(token)
		case CssParserStateSelector:
			p.stateSelector(token)
		case CssParserStateBlock:
			p.stateBlock(token)
		case CssParserStatePropName:
			p.statePropName(token)
		case CssParserStatePropValue:
			p.statePropValue(token)
		}
	}
}

func (p *CssParser) stateRoot(token *scanner.Token) {
	switch token.Type {
	case scanner.TokenIdent, scanner.TokenHash:
		p.toState(CssParserStateSelector, token.Value)
	case scanner.TokenChar:
		switch token.Value {
		case ".", "#", "*", "[", ":", "::":
			p.toState(CssParserStateSelector, token.Value)
		default:
			p.unexpectedToken(token)
		}
	case scanner.TokenS, scanner.TokenComment, scanner.TokenBOM, scanner.TokenCDO, scanner.TokenCDC:
		// just skip spaces and comments
	case scanner.TokenAtKeyword:
		// skip all the AtRule (unsupported)
	default:
		p.unexpectedToken(token)
	}
}

func (p *CssParser) stateSelector(token *scanner.Token) {
	switch token.Type {
	case scanner.TokenIdent,
		scanner.TokenS,
		scanner.TokenHash,
		scanner.TokenString,
		scanner.TokenNumber,
		scanner.TokenFunction,
		scanner.TokenIncludes,
		scanner.TokenDashMatch,
		scanner.TokenPrefixMatch,
		scanner.TokenSuffixMatch,
		scanner.TokenSubstringMatch:
		p.currentValue.WriteString(token.Value)
	case scanner.TokenChar:
		switch token.Value {
		case ">", "+", "~", "*", ",", ".", "[", "]", ")", ":", "|", "=":
			p.currentValue.WriteString(token.Value)
		case "{":
			p.state = CssParserStateBlock
			selector := strings.TrimSpace(p.currentValue.String())
			p.currentValue = nil
			p.currentRule = &Rule{Selector: selector, Declarations: make([]*Declaration, 0)}
			p.Result = append(p.Result, p.currentRule)
		default:
			p.unexpectedToken(token)
		}
	case scanner.TokenComment:
		// just skip comments
	case scanner.TokenAtKeyword:
		// skip all the AtRule (unsupported)
	default:
		p.unexpectedToken(token)
	}
}

func (p *CssParser) stateBlock(token *scanner.Token) {
	switch token.Type {
	case scanner.TokenIdent:
		p.toState(CssParserStatePropName, token.Value)
	case scanner.TokenChar:
		switch token.Value {
		case "}":
			p.state = CssParserStateRoot
			p.currentDeclaration = nil
			p.currentRule = nil
			p.parenthesesDepth = 0
		case "-":
			next := p.scanner.Next()
			// check if it's a css variable
			if next.Type == scanner.TokenIdent {
				p.toState(CssParserStatePropName, token.Value+next.Value)
			} else {
				p.unexpectedToken(token)
			}
		default:
			p.unexpectedToken(token)
		}
	case scanner.TokenS, scanner.TokenComment:
		// just skip spaces and comments
	default:
		p.unexpectedToken(token)
	}
}

func (p *CssParser) statePropName(token *scanner.Token) {
	switch token.Type {
	case scanner.TokenChar:
		if token.Value == ":" {
			p.state = CssParserStatePropValue
			property := p.currentValue.String()
			p.currentValue = new(strings.Builder)
			p.currentDeclaration = &Declaration{Property: property}
			p.currentRule.Declarations = append(p.currentRule.Declarations, p.currentDeclaration)
		} else {
			p.unexpectedToken(token)
		}
	case scanner.TokenS, scanner.TokenComment:
		// just skip spaces and comments
	default:
		p.unexpectedToken(token)
	}
}

func (p *CssParser) statePropValue(t *scanner.Token) {
	switch t.Type {
	case scanner.TokenIdent,
		scanner.TokenS,
		scanner.TokenString,
		scanner.TokenNumber,
		scanner.TokenPercentage,
		scanner.TokenDimension,
		scanner.TokenHash,
		scanner.TokenURI:
		p.currentValue.WriteString(t.Value)
	case scanner.TokenChar:
		if p.parenthesesDepth == 0 && t.Value == "}" {
			p.state = CssParserStateRoot
			p.currentRule = nil
			value := strings.TrimSpace(p.currentValue.String())
			p.currentValue = nil
			p.currentDeclaration.Value = value
			p.currentDeclaration = nil
		} else if p.parenthesesDepth == 0 && t.Value == ";" {
			p.state = CssParserStateBlock
			value := strings.TrimSpace(p.currentValue.String())
			p.currentValue = nil
			p.currentDeclaration.Value = value
			p.currentDeclaration = nil
		} else if t.Value == "(" {
			p.parenthesesDepth += 1
			p.currentValue.WriteString(t.Value)
		} else if t.Value == ")" {
			p.parenthesesDepth -= 1
			p.currentValue.WriteString(t.Value)
		} else if t.Value == "," || t.Value == "/" || t.Value == "!" || t.Value == "+" || t.Value == "-" || t.Value == "*" {
			p.currentValue.WriteString(t.Value)
		} else {
			p.unexpectedToken(t)
		}
	case scanner.TokenFunction:
		p.parenthesesDepth += 1
		p.currentValue.WriteString(t.Value)
	case scanner.TokenComment:
		// just skip comments
	default:
		p.unexpectedToken(t)
	}
}

func (p *CssParser) toState(state CssParserState, value string) {
	p.state = state
	p.currentValue = new(strings.Builder)
	p.currentValue.WriteString(value)
}

func (p *CssParser) unexpectedToken(token *scanner.Token) {
	p.Error = fmt.Errorf("unexpected token %s", token.String())
}
