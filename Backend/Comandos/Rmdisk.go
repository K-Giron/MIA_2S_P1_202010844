package Comandos

import (
	"fmt"
	"os"
	"strings"
)

/*RMDISK*/
func Rmdisk(commandArray []string) {
	Salida_comando += "MENSAJE: El comando RMDISK aqui inicia" + "\n"

	// Variables para los valores de los parametros
	val_path := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_path := false
	band_error := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			// Valido si el parametro ya fue ingresado
			if band_path {
				fmt.Print("Invalido: El parametro -path ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			fmt.Println("Invalido: El parametro " + data + " no es valido")
		}
	}

	// Verifico si no hay errores
	if !band_error {
		// Verifico que el parametro "Path" (Obligatorio) este ingresado
		if band_path {
			// Verifico si existe el archivo
			_, e := os.Stat(val_path)

			if e != nil {
				if os.IsNotExist(e) {
					fmt.Print("Invalido: El archivo no existe")
					band_path = false
				}
			} else {
				// Elimino el archivo
				err := os.Remove(val_path)

				// ERROR
				if err != nil {
					Salida_comando += "ERROR Al eliminar el archivo" + "\n"
					MsgError(err)
				} else {
					Salida_comando += "MENSAJE: El archivo fue eliminado" + "\n"
				}
				band_path = false
			}
		} else {
			fmt.Print("Invalido: El parametro -path es obligatorio")
		}
	}
	Salida_comando += "MENSAJE: El comando RMDISK termina aqui" + "\n"
}
