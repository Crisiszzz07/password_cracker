package main

import (
	"bufio"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type InfoReport struct { //estructura para los datos dinámicos de reporte
	Date       string
	Hash       string
	Dictionary string
	Time       string
	Result     string
}

func check(e error) { //función para errores en caso de que hayan
	if e != nil { //si hay un error, no tendrá el nil por defecto
		panic(e)
	}
}

func txtFile(filename string, info InfoReport) error {
	file, err := os.Create(filename) //NO USAR OS.OPEN PORQUE SOLO TIENE PERMISO DE LECTURA

	check(err)

	defer file.Close()

	const plantilla = `==================================================
🛡️ AUDITORY REPORT - PASSWORD CRACKER
==================================================
Event's date'   : %s
Target Hash     : %s
Dictionary used : %s
Attack time     : %s
--------------------------------------------------
[+] RESULT      : %s
==================================================
`

	//Acá usé el blank identifier porque Fprintf retorna dos valores: un int y un error, pero no tengo intención de usar ese int para nada de momento
	_, err = fmt.Fprintf(file, plantilla,
		info.Date,
		info.Hash,
		info.Dictionary,
		info.Time,
		info.Result)

	if err != nil { //en caso de que haya algún error con la escritura sobre el archivo
		return fmt.Errorf("There was a problem wrting the report: %v", err)
	}

	return nil //nil porque esta función retorna un valor de error, si fue exitoso, no necesita retornar un error mas que el nil

}

func jsonFile(filename string) error {
	return nil
}

func getFileExtension(filename string, info InfoReport) error {
	ext := filepath.Ext(filename)
	switch ext { //TODO: probablemente existe una manera más efectiva de soportar y validar las extensiones pero de momento uso lo que he aprendido
	case ".txt": //ARCHIVOS TXT
		err := txtFile(filename, info)
		check(err)

	case ".json": //ARCHIVOS JSON

	default:
		return fmt.Errorf("Unsupported file extension: %s", ext)
	}

	return nil
	//return filepath.Ext(filename)
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
	scanner := bufio.NewScanner(f) //NewScanner para leer las lineas del archivo
	for scanner.Scan() {           //se necesita un for para escanear cada linea del contenido del archivo

		//TODO: Revisar cómo poder contar realmente el tiempo de ejecución
		/*Función defer que intenté para contar la duración total del escaneo pero no sé aun
		start := time.Now()
		defer func(start time.Time) (duration time.Duration) {
			duration = time.Since(start)
			return duration
		}(start)

		*/
		//es mejor usar directamente la función de Bytes del scanner en vez de .Text
		//porque así no se genera una doble asignación porque uno es de tipo string y lo que pide el []byte es
		//tipo bytes
		huellaBinaria := sha256.Sum256(scanner.Bytes())
		//se pasa de binario a texto
		huellaTexto := fmt.Sprintf("%x", huellaBinaria)
		if targetHash == huellaTexto {

			fmt.Println("La contraseña es:", scanner.Text())
			if *filePtr != "" { //en caso de que se haya decidido crear un archivo

				dateNow := time.Now().Format("2006-01-02 15:04:05")
				dataReport := InfoReport{
					Date:       dateNow,
					Hash:       targetHash,
					Dictionary: *dicPtr,
					Time:       "ms",
					Result:     scanner.Text(),
				}
				err = getFileExtension(*filePtr, dataReport)

			}
			return
		}

	}

	fmt.Println("No se encontró ninguna contraseña que coincidiera.")

}
