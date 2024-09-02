package Comandos

import (
	"fmt"
	"os"
	"strings"
)

/*RMDISK*/
func Rmdisk(commandArray []string) {
	fmt.Println("[MENSAJE] El comando RMDISK aqui inicia")

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
				fmt.Println("[ERROR] El parametro -path ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			fmt.Println("[ERROR] Parametro no valido...")
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
					fmt.Println("[ERROR] El archivo no existe...")
					band_path = false
				}
			} else {
				// si existe el archivo
				fmt.Println("[MENSAJE] ¿Desea eliminar el disco [S/N]?: ")

				// Obtengo la opcion ingresada por el usuario
				var opcion string
				fmt.Scanln(&opcion)

				// verifico la opcion ingresada
				if opcion == "S" || opcion == "s" {

					// Elimino el archivo
					err := os.Remove(val_path)

					// ERROR
					if err != nil {
						MsgError(err)
					} else {
						fmt.Println("[SUCCES] El archivo fue eliminado!")
					}

					band_path = false
				} else if opcion == "N" || opcion == "n" {
					fmt.Println("[Mensaje] El archivo no fue eliminado!")
					band_path = false
				} else {
					fmt.Println("[ERROR] Opcion no valida...")
				}
			}
		}
	}
	fmt.Println("[MENSAJE] El comando RMDISK aqui finaliza")
}
