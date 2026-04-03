package internal

import "unicode"

const (
	TOKEN_EOF        = "EOF"
	TOKEN_IDENT      = "IDENT"
	TOKEN_TEXTO      = "TEXTO"
	TOKEN_PAGINA     = "PAGINA"
	TOKEN_FIN        = "FIN"
	TOKEN_PRIVADO    = "PRIVADO"
	TOKEN_CONTENEDOR = "CONTENEDOR"
	TOKEN_BOTON      = "BOTON"
	TOKEN_CAMPO      = "CAMPO"
	TOKEN_AS         = "AS"
	TOKEN_TABLA_WEB  = "TABLA_WEB"
	TOKEN_ENLACE     = "ENLACE"
	TOKEN_A          = "A"
)

type Token struct {
	Type    string
	Literal string
}

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

func NewLexer(input string) *Lexer { l := &Lexer{input: input}; l.readChar(); return l }

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	var tok Token
	switch l.ch {
	case '"':
		tok = Token{TOKEN_TEXTO, l.readString()}
	case 0:
		tok = Token{TOKEN_EOF, ""}
	default:
		if unicode.IsLetter(rune(l.ch)) {
			lit := l.readIdentifier()
			return Token{lookupIdent(lit), lit}
		}
		tok = Token{TOKEN_IDENT, string(l.ch)}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	p := l.pos + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[p:l.pos]
}

func (l *Lexer) readIdentifier() string {
	p := l.pos
	for unicode.IsLetter(rune(l.ch)) || l.ch == '_' {
		l.readChar()
	}
	return l.input[p:l.pos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func lookupIdent(ident string) string {
	m := map[string]string{
		"pagina": TOKEN_PAGINA, "fin": TOKEN_FIN, "privado": TOKEN_PRIVADO,
		"contenedor": TOKEN_CONTENEDOR, "boton": TOKEN_BOTON, "campo": TOKEN_CAMPO,
		"as": TOKEN_AS, "tabla_web": TOKEN_TABLA_WEB,
		"enlace": TOKEN_ENLACE, "a": TOKEN_A,
	}
	if t, ok := m[ident]; ok {
		return t
	}
	return TOKEN_IDENT
}
