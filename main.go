package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
)

func check(e error) { //función para errores en caso de que hayan
	if e != nil { //si hay un error, no tendrá el nil por defecto
		panic(e)
	}
}

func main() {
	targetHash := "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"
	filePath := "diccionario.txt"
	f, err := os.Open(filePath) //luego con el paquete de os y la función de Open se lee el contenido (que se guarda en f) y el resultado de error (se guard en err, un resultado de nil es igual a 0 error)
	check(err)

	defer f.Close()
	scanner := bufio.NewScanner(f) //NewScanner para leer las lineas del archivo
	for scanner.Scan() {           //se necesita un for para escanear cada linea del contenido del archivo

		//es mejor usar directamente la función de Bytes del scanner en vez de .Text
		//porque así no se genera una doble asignación porque uno es de tipo string y lo que pide el []byte es
		//tipo bytes
		huellaBinaria := sha256.Sum256([]byte(scanner.Bytes()))
		//se pasa de binario a texto
		huellaTexto := fmt.Sprintf("%x", huellaBinaria)
		if targetHash == huellaTexto {

			fmt.Println("La contraseña es:", scanner.Text())
			break
		}
	}
}
