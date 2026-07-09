package main

import (
	"bufio"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
)

func check(e error) { //función para errores en caso de que hayan
	if e != nil { //si hay un error, no tendrá el nil por defecto
		panic(e)
	}
}

// FUNCIÓN PARA CREAR ARCHIVOS
func createFile(filename string) (*os.File, error) {
	f, err := os.Create(filename)
	check(err)
	//defer f.Close()

	return f, err
}

func main() {

	//CONFIGURACIÓN DE FLAGS: El programa configura mas no lee aún el contenido de la terminal
	//los flag.type (ej: flag.String) reserva el espacio de la memoria para el cotenido y entrega el puntero

	//flag para targetHash

	hashPtr := flag.String("hash", "", "OBLIGATORY: hash password that wants to be cracked")

	//flag para diccionario
	dicPtr := flag.String("d", "", "OBLIGATORY: dictionary to use")

	//flag para almacenar un archivo resultante con los resultados que coincidan
	filePtr := flag.String("f", "", "OPTIONAL: file to use")
	//TODO
	//flag para distintas opciones de algoritmos
	flag.Parse() //Go lee el contenido en la terminal al hacer parsing y almacena teniendo en
	//cuenta el espacio de memoria al que hace referencia el puntero

	//para evitar el panic por intentar leer un archivo vacío, se hace una validación de
	//argumentos que son obligatorios
	if *hashPtr == "" || *dicPtr == "" {
		fmt.Println("Error: there are some missing arguments. Use -h or -help for details.")
		return
	}

	//se hace una dereferencia para obtener el contenido al que el puntero hace referencia
	targetHash := *hashPtr
	filePath := *dicPtr

	f, err := os.Open(filePath) //luego con el paquete de os y la función de Open se lee el contenido (que se guarda en f) y el resultado de error (se guard en err, un resultado de nil es igual a 0 error)
	check(err)

	defer f.Close()
	validacion := false            //solo lo uso porque de momento no me ocurré otra forma de imprimir que no hubo coincidencias
	scanner := bufio.NewScanner(f) //NewScanner para leer las lineas del archivo
	for scanner.Scan() {           //se necesita un for para escanear cada linea del contenido del archivo

		//es mejor usar directamente la función de Bytes del scanner en vez de .Text
		//porque así no se genera una doble asignación porque uno es de tipo string y lo que pide el []byte es
		//tipo bytes
		huellaBinaria := sha256.Sum256(scanner.Bytes())
		//se pasa de binario a texto
		huellaTexto := fmt.Sprintf("%x", huellaBinaria)
		if targetHash == huellaTexto {

			fmt.Println("La contraseña es:", scanner.Text())
			validacion = true
			if *filePtr != "" { //en caso de que se haya decidido crear un archivo
				//TODO
				//Mejorar formatos para los tipos de archivos distintos a .txt (ej: json)
				//Agregar validaciones de extensiones
				file, err := createFile(*filePtr)
				if err != nil {
					fmt.Println(err)
				}

				defer file.Close() //TODO: revisar porque dice que puede haber un resource leak por llamarlo en for
				file.WriteString(fmt.Sprintf("%s\n", scanner.Text()))

			}
			break
		}

	}
	if !validacion {
		fmt.Println("No se encontró ninguna contraseña que coincidiera.")
	}
}
