package Comandos

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

/* Ejemplo 7 */
// Estructura para el API
type Cmd_API struct {
	Cmd string `json:"cmd"`
}

// ------------variables globales----------------
var ListaMontajes []ParticionMontada
var Salida_comando string = ""
var GraphDot string = ""

// Obtiene y lee el comando
func Analizar() {
	fmt.Println("Bienvenido al proyecto 1")

	mux := http.NewServeMux()

	mux.HandleFunc("/analizar", func(w http.ResponseWriter, r *http.Request) {
		// Configuración de la cabecera
		w.Header().Set("Content-Type", "application/json")

		var Content Cmd_API
		body, _ := io.ReadAll(r.Body)

		// Deserializar JSON a struct
		json.Unmarshal(body, &Content)

		// Ejecuta el comando y captura la salida
		split_cmd(Content.Cmd)

		// Respuesta del servidor con la salida del comando
		response := `{"result": "` + Salida_comando + `" }`
		w.Write([]byte(response))
	})

	fmt.Println("Servidor en el puerto 5000")
	// Configuración de CORS
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":5000", handler))
}

/* Ejemplo 7 */
// Ejecuta comando linea por linea
func split_cmd(cmd string) {
	arr_com := strings.Split(cmd, "\n")

	for i := 0; i < len(arr_com); i++ {
		if arr_com[i] != "" {
			SplitComando(arr_com[i])
			Salida_comando += "\n"
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
		Salida_comando += comando + "\n"
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
		Salida_comando += "---------------------MKDISK--------------------" + "\n"
		Mkdisk(commandArray)
	case "rmdisk":
		Salida_comando += "---------------------RMDISK--------------------" + "\n"
		Rmdisk(commandArray)
	case "fdisk":
		Salida_comando += "---------------------FDISK--------------------" + "\n"
		Fdisk(commandArray)
	case "mount":
		Salida_comando += "---------------------MOUNT--------------------" + "\n"
		Mount(commandArray)
	case "rep":
		Salida_comando += "---------------------REP--------------------" + "\n"
		Rep(commandArray)
	case "mkfs":
		Salida_comando += "---------------------MKFS--------------------" + "\n"
		ValidarDatosMkfs(commandArray)
	case "execute":
		Salida_comando += "---------------------EXECUTE--------------------" + "\n"
		Execute(commandArray)
	case "pause":
		Salida_comando += "---------------------PAUSE--------------------" + "\n"
	default:
		fmt.Println("Comandooo no valido")
	}
}

/* EXECUTE */
func Execute(commandArray []string) {
	Salida_comando += "MENSAJE: El comando EXECUTE aqui inicia" + "\n"
	//variable para el path
	var path = ""

	//verificar que el comando traiga parametros sino mandar error
	if len(commandArray) == 1 {
		fmt.Println("Invalido: El comando EXECUTE no tiene parametros")
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
			fmt.Println("Invalido: El parametro " + data + " no es valido")
		}
	}
	//verificar que el path tenga la extension correcta
	if !strings.Contains(path, ".mia") {
		fmt.Println("Invalido: El archivo no tiene la extension correcta")
		return
	}
	//abrir el archivo
	disco, err := os.OpenFile(path, os.O_RDWR, 0660)

	//lee el archivo linea por linea y ejecuta los comandos
	if err != nil {
		fmt.Println("Invalido: No se pudo abrir el archivo")
		MsgError(err)
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
	Salida_comando += "MENSAJE: El comando EXECUTE aqui termina" + "\n"
}
