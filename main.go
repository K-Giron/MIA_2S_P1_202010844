package main

import (
	"P1_MIA/Comandos"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	fmt.Print("Iniciando el repositorio de Go")
	analizar()

}

// Obtiene y lee el comando
func analizar() {
	finalizar := false
	fmt.Println("Tarea 2: MKDISK, EXECUTE, REP y RMDISK")
	reader := bufio.NewReader(os.Stdin)

	//  Pide constantemente un comando
	for !finalizar {
		fmt.Print("Ingrese un comando: ")
		// Lee hasta que presione ENTER
		comando, _ := reader.ReadString('\n')
		if strings.Contains(comando, "exit") {
			/* SALIR */
			finalizar = true
		} else if strings.Contains(comando, "EXIT") {
			/* SALIR */
			finalizar = true
		} else {
			// Si no es vacio o el comando EXIT
			if comando != "" && comando != "exit\n" && comando != "EXIT\n" {
				// Obtener comando y parametros
				SplitComando(comando)
			}
		}
	}
}

// Separa los diferentes comando con sus parametros si tienen
func SplitComando(comando string) {
	var commandArray []string
	// Elimina los saltos de linea y retornos de carro
	comando = strings.Replace(comando, "\n", "", 1)
	comando = strings.Replace(comando, "\r", "", 1)

	// Banderas para verficar comentarios
	bandComentario := false

	if strings.Contains(comando, "pause") {
		// Comando sin Parametros
		commandArray = append(commandArray, comando)
	} else if strings.Contains(comando, "#") {
		// Comentario
		bandComentario = true
		fmt.Println(comando)
	} else {
		// Comando con Parametros
		commandArray = strings.Split(comando, " -")
	}

	// Ejecuta el comando leido si no es un comentario
	if !bandComentario {
		ejecutarComando(commandArray)
	}
}

// Identifica y ejecuta el comando encontrado
func ejecutarComando(commandArray []string) {
	// Convierte el comando a minusculas
	data := strings.ToLower(commandArray[0])

	// Identifica el comando a ejecutar
	switch data {
	case "mkdisk":
		fmt.Println("MKDISK")
		Comandos.Mkdisk(commandArray)
	case "rmdisk":
		fmt.Println("RMDISK")
		Comandos.Rmdisk(commandArray)
	case "fdisk":
		fmt.Println("FDISK")
		Comandos.Fdisk(commandArray)
	case "mount":
		fmt.Println("MOUNT")
		Comandos.Mount(commandArray)
	case "umount":
		fmt.Println("UMOUNT")
		//Comandos.Mount(commandArray)
	case "rep":
		fmt.Println("REP")
		Comandos.Rep(commandArray)
	case "execute":
		fmt.Println("EXEC")
		Execute(commandArray)
	case "pause":
		fmt.Println("PAUSE")
	default:
		fmt.Println("Comando no reconocido")

	}
}

/* EXECUTE */
func Execute(commandArray []string) {
	fmt.Println("MENSAJE: El comando EXECUTE aqui inicia")
	//variable para el path
	var path = ""

	//verificar que el comando traiga parametros sino mandar error
	if len(commandArray) == 1 {
		fmt.Println("ERROR: El comando no trae parametros")
		return
	}

	//recorrer los parametros
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			// Reemplaza comillas
			path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			fmt.Println("ERROR: Parametro no valido...")
		}
	}
	//verificar que el path tenga la extension correcta
	if !strings.Contains(path, ".mia") {
		fmt.Println("ERROR: El path no tiene la extension correcta")
		return
	}
	fmt.Println("El path es: ", path)
	//abrir el archivo
	disco, err := os.OpenFile(path, os.O_RDWR, 0660)

	//lee el archivo linea por linea y ejecuta los comandos
	if err != nil {
		fmt.Println("ERROR: No se pudo abrir el archivo")
		Comandos.MsgError(err)
	}
	scanner := bufio.NewScanner(disco)
	for scanner.Scan() {
		//si la linea esta vacia no se ejecuta
		if scanner.Text() == "" {
			continue
		}
		linea := scanner.Text()
		SplitComando(linea)
	}
	disco.Close()
	fmt.Println("MENSAJE: El comando EXECUTE aqui finaliza")

}
