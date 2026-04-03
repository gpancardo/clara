package internal

import (
	"fmt"
	"os"
	"strings"
)

func Generar(nodos []Nodo) {
	var handlers strings.Builder

	// Recorremos todos los nodos para generar los archivos HTML y preparar los handlers
	for _, n := range nodos {
		p, ok := n.(*Pagina)
		if !ok {
			continue
		}

		id := strings.ReplaceAll(p.Ruta, "/", "root")
		if id == "root" {
			id = "indice"
		} // Evitamos nombres ambiguos

		// Generamos el archivo HTML para esta página específica
		htmlBody := generarHTML(p.Hijos)
		os.WriteFile("vista_"+id+".html", []byte(wrap(htmlBody)), 0644)

		checkPrivado := ""
		if p.EsPrivada {
			checkPrivado = `
		c, err := r.Cookie("sesion")
		if err != nil || c.Value == "" { 
			http.Redirect(w, r, "/", 303)
			return 
		}`
		}

		// Añadimos la lógica de esta página al constructor de handlers
		handlers.WriteString(fmt.Sprintf(`
	http.HandleFunc("%s", func(w http.ResponseWriter, r *http.Request) {
		%s
		if r.Method == "POST" {
			r.ParseForm()
			nuevo := make(map[string]string)
			for k, v := range r.PostForm { 
				if k == "pass" { 
					h := sha256.Sum256([]byte(v[0]))
					nuevo[k] = fmt.Sprintf("%%x", h)
				} else { 
					nuevo[k] = v[0] 
				}
			}
			// Si el usuario envió "admin" en un campo, le damos sesión
			for _, v := range nuevo {
				if v == "admin" {
					http.SetCookie(w, &http.Cookie{Name: "sesion", Value: "admin", Path: "/"})
				}
			}

			if len(nuevo) > 0 { guardar(nuevo) }
			http.Redirect(w, r, r.URL.Path, 303)
			return
		}
		
		datos, _ := os.ReadFile("db.json")
		var regs []map[string]string
		json.Unmarshal(datos, &regs)
		
		tmpl, err := template.ParseFiles("vista_%s.html")
		if err != nil {
			http.Error(w, "Error al cargar la vista", 500)
			return
		}
		tmpl.Execute(w, regs)
	})`, p.Ruta, checkPrivado, id))
	}

	// IMPORTANTE: El código final del servidor se construye FUERA del bucle for
	serverCode := fmt.Sprintf(`package main
import (
	"net/http"
	"os"
	"encoding/json"
	"html/template"
	"fmt"
	"crypto/sha256"
)

type Reg map[string]string

func guardar(d map[string]string) {
	var db []map[string]string
	f, _ := os.ReadFile("db.json")
	json.Unmarshal(f, &db)
	db = append(db, d)
	b, _ := json.Marshal(db)
	os.WriteFile("db.json", b, 0644)
}

func main() {
	%s
	fmt.Println("🚀 Clara activa en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}`, handlers.String())

	os.WriteFile("servidor.go", []byte(serverCode), 0644)
}

func generarHTML(nodos []Nodo) string {
	res := ""
	for _, n := range nodos {
		el, ok := n.(*Elemento)
		if !ok {
			continue
		}

		switch el.Type {
		case TOKEN_CONTENEDOR:
			res += "<main>" + generarHTML(el.Hijos) + "</main>"
		case TOKEN_ENLACE:
			res += fmt.Sprintf("<a href='%s' style='display:block; margin: 10px 0;'>%s</a>", el.Ident, el.Valor)

		case TOKEN_BOTON:
			// Mejoramos el botón para que se vea como un llamado a la acción
			res += fmt.Sprintf("<button type='submit' style='width:100%%'>%s</button>", el.Valor)
		case TOKEN_CAMPO:
			tipo := "text"
			nombre := el.Ident
			if nombre == "" {
				nombre = "dato"
			} // Fallback si no hay "as"
			if nombre == "pass" {
				tipo = "password"
			}

			res += fmt.Sprintf("<label>%s</label><input name='%s' type='%s' required>", el.Valor, nombre, tipo)
		case TOKEN_TABLA_WEB:
			res += `
			<div style="overflow-x:auto;">
				<table>
					<thead><tr><th>Datos Registrados</th></tr></thead>
					<tbody>
						{{range .}}
						<tr>
							<td>
								{{range $k, $v := .}}
									<strong>{{$k}}:</strong> {{$v}} | 
								{{end}}
							</td>
						</tr>
						{{end}}
					</tbody>
				</table>
			</div>`

		}
	}
	return "<form method='POST'>" + res + "</form>"
}

func wrap(b string) string {
	return `<!DOCTYPE html>
<html lang="es">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/water.css@2/out/water.css">
	<title>App Clara</title>
</head>
<body>
	` + b + `
</body>
</html>`
}
