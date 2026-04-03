package main

import (
	"clara/internal"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: clara <archivo.cla>")
		return
	}

	abs, _ := filepath.Abs(os.Args[1])
	os.Chdir(filepath.Dir(abs))

	if _, err := os.Stat("db.json"); os.IsNotExist(err) {
		os.WriteFile("db.json", []byte("[]"), 0644)
	}

	input, _ := os.ReadFile(abs)
	l := internal.NewLexer(string(input))
	p := internal.NewParser(l)
	internal.Generar(p.Parsear())

	fmt.Println("🛠️  Compilación exitosa.")
	cmd := exec.Command("go", "run", "servidor.go")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	cmd.Run()
}
