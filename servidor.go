package main
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
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		
		if r.Method == "POST" {
			r.ParseForm()
			nuevo := make(map[string]string)
			for k, v := range r.PostForm { 
				if k == "pass" { 
					h := sha256.Sum256([]byte(v[0]))
					nuevo[k] = fmt.Sprintf("%x", h)
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
		
		tmpl, err := template.ParseFiles("vista_indice.html")
		if err != nil {
			http.Error(w, "Error al cargar la vista", 500)
			return
		}
		tmpl.Execute(w, regs)
	})
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		
		if r.Method == "POST" {
			r.ParseForm()
			nuevo := make(map[string]string)
			for k, v := range r.PostForm { 
				if k == "pass" { 
					h := sha256.Sum256([]byte(v[0]))
					nuevo[k] = fmt.Sprintf("%x", h)
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
		
		tmpl, err := template.ParseFiles("vista_rootsignup.html")
		if err != nil {
			http.Error(w, "Error al cargar la vista", 500)
			return
		}
		tmpl.Execute(w, regs)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		
		if r.Method == "POST" {
			r.ParseForm()
			nuevo := make(map[string]string)
			for k, v := range r.PostForm { 
				if k == "pass" { 
					h := sha256.Sum256([]byte(v[0]))
					nuevo[k] = fmt.Sprintf("%x", h)
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
		
		tmpl, err := template.ParseFiles("vista_rootlogin.html")
		if err != nil {
			http.Error(w, "Error al cargar la vista", 500)
			return
		}
		tmpl.Execute(w, regs)
	})
	fmt.Println("🚀 Clara activa en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}