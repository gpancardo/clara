package internal

type Nodo interface{}

type Pagina struct {
	Ruta      string
	EsPrivada bool
	Hijos     []Nodo
}

type Elemento struct {
	Type  string
	Valor string // El texto visible (Label, Placeholder o texto del enlace)
	Ident string // El ID técnico (nombre del campo en JSON o ruta del enlace)
	Hijos []Nodo
}

type Parser struct {
	l    *Lexer
	cur  Token
	peek Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) Parsear() []Nodo {
	var prog []Nodo
	for p.cur.Type != TOKEN_EOF {
		if p.cur.Type == TOKEN_PAGINA {
			prog = append(prog, p.parsePagina())
		}
		p.nextToken()
	}
	return prog
}

func (p *Parser) parsePagina() *Pagina {
	pag := &Pagina{Ruta: "/"}

	// Buscamos la ruta y si es privada antes de entrar a los hijos
	for p.cur.Type != TOKEN_TEXTO && p.cur.Type != TOKEN_FIN && p.cur.Type != TOKEN_EOF {
		if p.cur.Type == TOKEN_PRIVADO {
			pag.EsPrivada = true
		}
		p.nextToken()
	}

	if p.cur.Type == TOKEN_TEXTO {
		pag.Ruta = p.cur.Literal
	}

	p.nextToken()

	// Parsear el cuerpo de la página
	for p.cur.Type != TOKEN_FIN && p.cur.Type != TOKEN_EOF {
		el := p.parseElemento()
		if el != nil {
			pag.Hijos = append(pag.Hijos, el)
		}
		p.nextToken()
	}
	return pag
}

func (p *Parser) parseElemento() Nodo {
	// Evitar procesar tokens de cierre como elementos
	if p.cur.Type == TOKEN_FIN || p.cur.Type == TOKEN_EOF {
		return nil
	}

	el := &Elemento{Type: p.cur.Type}

	// Lógica según el tipo de componente
	switch p.cur.Type {

	case TOKEN_CAMPO:
		if p.peek.Type == TOKEN_TEXTO {
			p.nextToken()
			el.Valor = p.cur.Literal
		}
		// Soporte para: campo "Nombre" as "user"
		if p.peek.Type == TOKEN_AS {
			p.nextToken() // salta "as"
			if p.peek.Type == TOKEN_TEXTO {
				p.nextToken()
				el.Ident = p.cur.Literal
			}
		}

	case TOKEN_ENLACE:
		if p.peek.Type == TOKEN_TEXTO {
			p.nextToken()
			el.Valor = p.cur.Literal // Texto del link
		}
		// Soporte para: enlace "Ir" a "/ruta"
		if p.peek.Type == TOKEN_A {
			p.nextToken() // salta "a"
			if p.peek.Type == TOKEN_TEXTO {
				p.nextToken()
				el.Ident = p.cur.Literal // Destino
			}
		}

	case TOKEN_BOTON, TOKEN_TABLA_WEB:
		if p.peek.Type == TOKEN_TEXTO {
			p.nextToken()
			el.Valor = p.cur.Literal
		}

	case TOKEN_CONTENEDOR:
		p.nextToken() // entramos al contenedor
		for p.cur.Type != TOKEN_FIN && p.cur.Type != TOKEN_EOF {
			hijo := p.parseElemento()
			if hijo != nil {
				el.Hijos = append(el.Hijos, hijo)
			}
			p.nextToken()
		}
	}

	return el
}
