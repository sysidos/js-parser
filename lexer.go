package js-parser

import (
	"strings"
	"unicode/utf8"
)

// Lexer states.
const (
	StateInitial         = 0
	StateDiv             = 1
	StateTemplate        = 2
	StateTemplateDiv     = 3
	StateTemplateExpr    = 4
	StateTemplateExprDiv = 5
	StateJsxTypeArgs     = 6
	StateJsxTag          = 7
	StateJsxClosingTag   = 8
	StateJsxText         = 9
)

type Dialect int

const (
	Javascript Dialect = iota
	Typescript
	TypescriptJsx
)

// Lexer uses a generated DFA to scan through a utf-8 encoded input string. If
// the string starts with a BOM character, it gets skipped.
type Lexer struct {
	source string

	ch          rune // current character, -1 means EOI
	offset      int  // character offset
	tokenOffset int  // last token offset
	line        int  // current line number (1-based)
	tokenLine   int  // last token line
	scanOffset  int  // scanning offset
	value       interface{}

	State   int // lexer state, modifiable
	Dialect Dialect
	token   Token // last token
	Stack   []int // stack of JSX states, non-empty for StateJsx*
}

var bomSeq = "\xef\xbb\xbf"

// Init prepares the lexer l to tokenize source by performing the full reset
// of the internal state.
func (l *Lexer) Init(source string) {
	l.source = source

	l.ch = 0
	l.offset = 0
	l.tokenOffset = 0
	l.line = 1
	l.tokenLine = 1
	l.State = 0
	l.Dialect = Javascript
	l.token = UNAVAILABLE
	l.Stack = nil

	if strings.HasPrefix(source, bomSeq) {
		l.offset += len(bomSeq)
	}

	l.rewind(l.offset)
}

// Next finds and returns the next token in l.source. The source end is
// indicated by Token.EOI.
//
// The token text can be retrieved later by calling the Text() method.
func (l *Lexer) Next() Token {
	prevLine := l.tokenLine
restart:
	l.tokenLine = l.line
	l.tokenOffset = l.offset

	state := tmStateMap[l.State]
	hash := uint32(0)
	for state >= 0 {
		var ch int
		if uint(l.ch) < tmRuneClassLen {
			ch = int(tmRuneClass[l.ch])
		} else if l.ch < 0 {
			state = int(tmLexerAction[state*tmNumClasses])
			continue
		} else {
			ch = mapRune(l.ch)
		}
		state = int(tmLexerAction[state*tmNumClasses+ch])
		if state > tmFirstRule {
			hash = hash*uint32(31) + uint32(l.ch)

			if l.ch == '\n' {
				l.line++
			}

			// Scan the next character.
			// Note: the following code is inlined to avoid performance implications.
			l.offset = l.scanOffset
			if l.offset < len(l.source) {
				r, w := rune(l.source[l.offset]), 1
				if r >= 0x80 {
					// not ASCII
					r, w = utf8.DecodeRuneInString(l.source[l.offset:])
				}
				l.scanOffset += w
				l.ch = r
			} else {
				l.ch = -1 // EOI
			}
		}
	}

	rule := tmFirstRule - state
	switch rule {
	case 5:
		hh := hash & 127
		switch hh {
		case 1:
			if hash == 0x2f9501 && "enum" == l.source[l.tokenOffset:l.offset] {
				rule = 40
				break
			}
		case 3:
			if hash == 0xcd244983 && "finally" == l.source[l.tokenOffset:l.offset] {
				rule = 20
				break
			}
			if hash == 0xed412583 && "private" == l.source[l.tokenOffset:l.offset] {
				rule = 55
				break
			}
		case 7:
			if hash == 0x33c587 && "null" == l.source[l.tokenOffset:l.offset] {
				rule = 41
				break
			}
		case 11:
			if hash == 0xc8b && "do" == l.source[l.tokenOffset:l.offset] {
				rule = 16
				break
			}
		case 13:
			if hash == 0x6da5f8d && "yield" == l.source[l.tokenOffset:l.offset] {
				rule = 39
				break
			}
		case 14:
			if hash == 0x36758e && "true" == l.source[l.tokenOffset:l.offset] {
				rule = 42
				break
			}
		case 17:
			if hash == 0xcad56011 && "string" == l.source[l.tokenOffset:l.offset] {
				rule = 62
				break
			}
			if hash == 0xcb7e7191 && "target" == l.source[l.tokenOffset:l.offset] {
				rule = 52
				break
			}
			if hash == 0xcccfb691 && "typeof" == l.source[l.tokenOffset:l.offset] {
				rule = 34
				break
			}
		case 20:
			if hash == 0x375194 && "void" == l.source[l.tokenOffset:l.offset] {
				rule = 36
				break
			}
		case 24:
			if hash == 0xcb197598 && "symbol" == l.source[l.tokenOffset:l.offset] {
				rule = 63
				break
			}
		case 25:
			if hash == 0xb22d2499 && "extends" == l.source[l.tokenOffset:l.offset] {
				rule = 19
				break
			}
		case 27:
			if hash == 0x1a21b && "let" == l.source[l.tokenOffset:l.offset] {
				rule = 48
				break
			}
		case 29:
			if hash == 0xd1d && "if" == l.source[l.tokenOffset:l.offset] {
				rule = 23
				break
			}
		case 30:
			if hash == 0x364e9e && "this" == l.source[l.tokenOffset:l.offset] {
				rule = 31
				break
			}
		case 32:
			if hash == 0x1a9a0 && "new" == l.source[l.tokenOffset:l.offset] {
				rule = 27
				break
			}
		case 33:
			if hash == 0x20a6f421 && "debugger" == l.source[l.tokenOffset:l.offset] {
				rule = 13
				break
			}
		case 34:
			if hash == 0x6749f022 && "abstract" == l.source[l.tokenOffset:l.offset] {
				rule = 64
				break
			}
		case 35:
			if hash == 0x5cb1923 && "false" == l.source[l.tokenOffset:l.offset] {
				rule = 43
				break
			}
		case 37:
			if hash == 0xb96173a5 && "import" == l.source[l.tokenOffset:l.offset] {
				rule = 24
				break
			}
			if hash == 0xd25 && "in" == l.source[l.tokenOffset:l.offset] {
				rule = 25
				break
			}
		case 39:
			if hash == 0xde312ca7 && "continue" == l.source[l.tokenOffset:l.offset] {
				rule = 12
				break
			}
			if hash == 0x1c727 && "var" == l.source[l.tokenOffset:l.offset] {
				rule = 35
				break
			}
		case 40:
			if hash == 0x3db6c28 && "boolean" == l.source[l.tokenOffset:l.offset] {
				rule = 60
				break
			}
		case 42:
			if hash == 0x3017aa && "from" == l.source[l.tokenOffset:l.offset] {
				rule = 46
				break
			}
			if hash == 0xd2a && "is" == l.source[l.tokenOffset:l.offset] {
				rule = 67
				break
			}
		case 43:
			if hash == 0xb06685ab && "delete" == l.source[l.tokenOffset:l.offset] {
				rule = 15
				break
			}
		case 44:
			if hash == 0x35c3d12c && "instanceof" == l.source[l.tokenOffset:l.offset] {
				rule = 26
				break
			}
		case 46:
			if hash == 0xdbba6bae && "protected" == l.source[l.tokenOffset:l.offset] {
				rule = 56
				break
			}
		case 48:
			if hash == 0x2e7b30 && "case" == l.source[l.tokenOffset:l.offset] {
				rule = 8
				break
			}
			if hash == 0xc97057b0 && "implements" == l.source[l.tokenOffset:l.offset] {
				rule = 53
				break
			}
			if hash == 0xc84e3d30 && "return" == l.source[l.tokenOffset:l.offset] {
				rule = 28
				break
			}
		case 49:
			if hash == 0x6bdcb31 && "while" == l.source[l.tokenOffset:l.offset] {
				rule = 37
				break
			}
		case 50:
			if hash == 0xc32 && "as" == l.source[l.tokenOffset:l.offset] {
				rule = 44
				break
			}
		case 52:
			if hash == 0xb32913b4 && "export" == l.source[l.tokenOffset:l.offset] {
				rule = 18
				break
			}
			if hash == 0xcafbb734 && "switch" == l.source[l.tokenOffset:l.offset] {
				rule = 30
				break
			}
		case 57:
			if hash == 0x2f8d39 && "else" == l.source[l.tokenOffset:l.offset] {
				rule = 17
				break
			}
			if hash == 0x1df56d39 && "interface" == l.source[l.tokenOffset:l.offset] {
				rule = 54
				break
			}
		case 58:
			if hash == 0x368f3a && "type" == l.source[l.tokenOffset:l.offset] {
				rule = 71
				break
			}
		case 59:
			if hash == 0x5a0eebb && "catch" == l.source[l.tokenOffset:l.offset] {
				rule = 9
				break
			}
			if hash == 0x1c1bb && "try" == l.source[l.tokenOffset:l.offset] {
				rule = 33
				break
			}
		case 65:
			if hash == 0x5c13d641 && "default" == l.source[l.tokenOffset:l.offset] {
				rule = 14
				break
			}
		case 66:
			if hash == 0xcc56be42 && "readonly" == l.source[l.tokenOffset:l.offset] {
				rule = 72
				break
			}
		case 70:
			if hash == 0x37b0c6 && "with" == l.source[l.tokenOffset:l.offset] {
				rule = 38
				break
			}
		case 73:
			if hash == 0x18cc9 && "for" == l.source[l.tokenOffset:l.offset] {
				rule = 21
				break
			}
			if hash == 0xc258db49 && "number" == l.source[l.tokenOffset:l.offset] {
				rule = 61
				break
			}
		case 74:
			if hash == 0xef05ac4a && "unknown" == l.source[l.tokenOffset:l.offset] {
				rule = 59
				break
			}
		case 78:
			if hash == 0x5fb304e && "infer" == l.source[l.tokenOffset:l.offset] {
				rule = 75
				break
			}
		case 81:
			if hash == 0xcde68bd1 && "unique" == l.source[l.tokenOffset:l.offset] {
				rule = 74
				break
			}
		case 86:
			if hash == 0x58e7956 && "await" == l.source[l.tokenOffset:l.offset] {
				rule = 6
				break
			}
			if hash == 0x18f56 && "get" == l.source[l.tokenOffset:l.offset] {
				rule = 47
				break
			}
		case 87:
			if hash == 0xdd7 && "of" == l.source[l.tokenOffset:l.offset] {
				rule = 49
				break
			}
		case 88:
			if hash == 0x524f73d8 && "function" == l.source[l.tokenOffset:l.offset] {
				rule = 22
				break
			}
		case 91:
			if hash == 0x4aa3555b && "namespace" == l.source[l.tokenOffset:l.offset] {
				rule = 69
				break
			}
		case 98:
			if hash == 0x1bc62 && "set" == l.source[l.tokenOffset:l.offset] {
				rule = 50
				break
			}
		case 99:
			if hash == 0x5a73763 && "const" == l.source[l.tokenOffset:l.offset] {
				rule = 11
				break
			}
		case 101:
			if hash == 0x414f0165 && "require" == l.source[l.tokenOffset:l.offset] {
				rule = 70
				break
			}
		case 102:
			if hash == 0x693a6e6 && "throw" == l.source[l.tokenOffset:l.offset] {
				rule = 32
				break
			}
		case 105:
			if hash == 0xc5bdb269 && "public" == l.source[l.tokenOffset:l.offset] {
				rule = 57
				break
			}
		case 106:
			if hash == 0x5bee456a && "declare" == l.source[l.tokenOffset:l.offset] {
				rule = 66
				break
			}
		case 108:
			if hash == 0x179ec && "any" == l.source[l.tokenOffset:l.offset] {
				rule = 58
				break
			}
			if hash == 0xc04ba66c && "module" == l.source[l.tokenOffset:l.offset] {
				rule = 68
				break
			}
		case 110:
			if hash == 0xcacdce6e && "static" == l.source[l.tokenOffset:l.offset] {
				rule = 51
				break
			}
		case 118:
			if hash == 0x6139076 && "keyof" == l.source[l.tokenOffset:l.offset] {
				rule = 73
				break
			}
		case 120:
			if hash == 0x5a5a978 && "class" == l.source[l.tokenOffset:l.offset] {
				rule = 10
				break
			}
		case 122:
			if hash == 0xa152d7fa && "constructor" == l.source[l.tokenOffset:l.offset] {
				rule = 65
				break
			}
		case 123:
			if hash == 0x68b6f7b && "super" == l.source[l.tokenOffset:l.offset] {
				rule = 29
				break
			}
		case 124:
			if hash == 0x58d027c && "async" == l.source[l.tokenOffset:l.offset] {
				rule = 45
				break
			}
		case 127:
			if hash == 0x59a58ff && "break" == l.source[l.tokenOffset:l.offset] {
				rule = 7
				break
			}
		}
	}

	token := tmToken[rule]
	space := false
	switch rule {
	case 0:
		if l.offset == l.tokenOffset {
			l.rewind(l.scanOffset)
		}
	case 2: // WhiteSpace: /[\t\x0b\x0c\x20\xa0\ufeff\p{Zs}]/, WhiteSpace: /[\n\r\u2028\u2029]|\r\n/
		space = true
	case 112: // invalid_token: /\?\.[0-9]/
		{
			l.rewind(l.tokenOffset + 1)
			token = QUEST
		}
	}
	if space {
		goto restart
	}

	// There is an ambiguity in the language that a slash can either represent
	// a division operator, or start a regular expression literal. This gets
	// disambiguated at the grammar level - division always follows an
	// expression, while regex literals are expressions themselves. Here we use
	// some knowledge about the grammar to decide whether the next token can be
	// a regular expression literal.
	//
	// See the following thread for more details:
	// http://stackoverflow.com/questions/5519596/when-parsing-javascript-what

	if l.State <= StateTemplateExprDiv {
		// The lowest bit of "l.State" determines how to interpret a forward
		// slash if it happens to be the next character.
		//   unset: start of a regular expression literal
		//   set:   start of a division operator (/ or /=)
		switch token {
		case NEW, DELETE, VOID, TYPEOF, INSTANCEOF, IN, DO, RETURN, CASE, THROW, ELSE:
			l.State &^= 1
		case TEMPLATEHEAD:
			l.State |= 1
			l.pushState(StateTemplate)
		case TEMPLATEMIDDLE:
			l.State = StateTemplate
		case TEMPLATETAIL:
			l.popState()
		case RPAREN, RBRACK:
			// TODO support if (...) /aaaa/;
			l.State |= 1
		case PLUSPLUS, MINUSMINUS:
			if prevLine != l.tokenLine {
				// This is a pre-increment/decrement, so we expect a regular expression.
				l.State &^= 1
			}
			// Otherwise: if we were expecting a regular expression literal before this
			// token, this is a pre-increment/decrement, otherwise, this is a post. We
			// can just propagate the previous value of the lowest bit of the state.
		case LT:
			if l.State&1 == 0 {
				// Start a new JSX tag.
				if l.Dialect != Typescript {
					l.State |= 1
					l.pushState(StateJsxTag)
				}
			} else {
				l.State &^= 1
			}
		case LBRACE:
			l.State &^= 1
			if l.State >= StateTemplate {
				l.pushState(StateTemplateExpr)
			}
		case RBRACE:
			l.State &^= 1
			if l.State >= StateTemplate {
				l.popState()
			}
		case SINGLELINECOMMENT, MULTILINECOMMENT:
			break
		default:
			if token >= punctuationStart && token < punctuationEnd {
				l.State &^= 1
			} else {
				l.State |= 1
			}
		}
	} else {
		// Handling JSX states.
		switch token {
		case DIV:
			if l.State == StateJsxTag && l.token == LT {
				l.State = StateJsxClosingTag
				if len(l.Stack) > 0 {
					l.Stack = l.Stack[:len(l.Stack)-1]
				}
			}
		case GT:
			if l.State == StateJsxTypeArgs || l.State == StateJsxClosingTag || l.token == DIV {
				l.popState()
			} else {
				l.State = StateJsxText
			}
		case LBRACE:
			if l.State != StateJsxTypeArgs {
				l.pushState(StateTemplateExpr)
			}
		case LT:
			if l.Dialect == TypescriptJsx && l.State != StateJsxText && l.token != ASSIGN {
				// Type arguments.
				l.pushState(StateJsxTypeArgs)
			} else {
				// Start a new JSX tag.
				l.pushState(StateJsxTag)
			}
		}
	}
	l.token = token
	return token
}

// Pos returns the start and end positions of the last token returned by Next().
func (l *Lexer) Pos() (start, end int) {
	start = l.tokenOffset
	end = l.offset
	return
}

// Line returns the line number of the last token returned by Next().
func (l *Lexer) Line() int {
	return l.tokenLine
}

// Text returns the substring of the input corresponding to the last token.
func (l *Lexer) Text() string {
	return l.source[l.tokenOffset:l.offset]
}

// Value returns the value associated with the last returned token.
func (l *Lexer) Value() interface{} {
	return l.value
}

// rewind can be used in lexer actions to accept a portion of a scanned token, or to include
// more text into it.
func (l *Lexer) rewind(offset int) {
	if offset < l.offset {
		l.line -= strings.Count(l.source[offset:l.offset], "\n")
	} else {
		if offset > len(l.source) {
			offset = len(l.source)
		}
		l.line += strings.Count(l.source[l.offset:offset], "\n")
	}

	// Scan the next character.
	l.scanOffset = offset
	l.offset = offset
	if l.offset < len(l.source) {
		r, w := rune(l.source[l.offset]), 1
		if r >= 0x80 {
			// not ASCII
			r, w = utf8.DecodeRuneInString(l.source[l.offset:])
		}
		l.scanOffset += w
		l.ch = r
	} else {
		l.ch = -1 // EOI
	}
}

func (l *Lexer) pushState(newState int) {
	l.Stack = append(l.Stack, l.State)
	l.State = newState
}

func (l *Lexer) popState() {
	if ln := len(l.Stack); ln > 0 {
		l.State = l.Stack[ln-1]
		l.Stack = l.Stack[:ln-1]
	} else {
		l.State = StateDiv
	}
}
