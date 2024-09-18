package Comandos

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// estructura ---------------------------------------
// Master Boot Record (MBR)
type Mbr = struct {
	Mbr_tamano         [100]byte
	Mbr_fecha_creacion [100]byte
	Mbr_dsk_signature  [100]byte
	Dsk_fit            [100]byte
	Mbr_partition      [4]Partition
}

type Partition = struct {
	Part_status      [100]byte
	Part_type        [100]byte
	Part_fit         [100]byte
	Part_start       [100]byte
	Part_size        [100]byte
	Part_name        [100]byte
	Part_correlative [100]byte
	Part_id          [100]byte
}

func NewMbr() Mbr {
	var mbr Mbr
	return mbr
}

/* MKDISK */
func Mkdisk(commandArray []string) {
	//parametros obligatorios
	//size
	//path
	//unit
	//fit

	Salida_comando += "MENSAJE: El comando MKDISK aqui inicia" + "\n"

	// Variables para los valores de los parametros
	val_size := 0
	val_fit := ""
	val_unit := ""
	val_path := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_size := false
	band_fit := false
	band_unit := false
	band_path := false
	band_error := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> SIZE */
		case strings.Contains(data, "size="):
			// Valido si el parametro ya fue ingresado
			if band_size {
				fmt.Print("Invalido: El parametro -size ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_size = true

			// Conversion a entero
			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size

			// ERROR de conversion
			if err != nil {
				MsgError(err)
			}

			// Valido que el tamaño sea positivo
			if val_size < 0 {
				band_error = true
				fmt.Println("Invalido: El valor del parametro -size no es valido...")
				break
			}
		/* PARAMETRO OPCIONAL -> FIT */
		case strings.Contains(data, "fit="):
			// Valido si el parametro ya fue ingresado
			if band_fit {
				fmt.Println("Invalido: El parametro -fit ya fue ingresado...")
				band_error = true
				break
			}

			// Le quito las comillas y lo paso a minusculas
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" { // Best Fit
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" { // First Fit
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" { // Worst Fit
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "w"
			} else {
				fmt.Println("Invalido: El valor del parametro -fit no es valido...")
				band_error = true
				break
			}
		/* Pametro opcional -> UNIT */
		case strings.Contains(data, "unit="):
			// Valido si el parametro ya fue ingresado
			if band_unit {
				fmt.Println("Invalido: El parametro -unit ya fue ingresado...")
				band_error = true
				break
			}

			// Reemplaza comillas y lo paso a minusculas
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)

			// valido que tenga unidades validas
			if val_unit == "k" || val_unit == "m" { // Kilobytes o Megabytes
				// Activo la bandera del parametro
				band_unit = true
			} else {
				// Parametro no valido
				fmt.Println("Invalido: El valor del parametro -unit no es valido...")
				band_error = true
				break
			}
		/* PARAMETRO OBLIGATORIO -> PATH */
		// No se utiliza en el proyecto, es una ruta fija para el disco
		case strings.Contains(data, "path="):
			// Valido si el parametro ya fue ingresado
			if band_path {
				fmt.Println("Invalido: El parametro -path ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
			//verifico que el archivo tenga la extension .mia
			if !strings.Contains(val_path, ".mia") {
				fmt.Println("Invalido: El archivo no tiene la extension correcta...")
				band_error = true
				break
			}
		/* PARAMETRO NO VALIDO */
		default:
			fmt.Println("Invalido: El parametro " + data + " no es valido...")
		}
	}

	// Verifico si no hay errores
	if !band_error {
		// Verifico que el parametro "Path" (Obligatorio) este ingresado
		if band_path {
			// Verifico que el parametro "Size" (Obligatorio) este ingresado
			if band_size {
				total_size := 1024
				master_boot_record := Mbr{}

				// Disco -> Archivo Binario
				crear_disco(val_path)

				// Fecha
				fecha := time.Now()
				str_fecha := fecha.Format("02/01/2006 15:04:05")

				// Copio valor al Struct
				copy(master_boot_record.Mbr_fecha_creacion[:], str_fecha)

				// Numero aleatorio
				rand.Seed(time.Now().UnixNano())
				min := 0
				max := 100
				num_random := rand.Intn(max-min+1) + min

				// Copio valor al Struct
				copy(master_boot_record.Mbr_dsk_signature[:], strconv.Itoa(int(num_random)))

				//verifico si existe el parametro fit
				if band_fit {
					// Copio valor al Struct
					copy(master_boot_record.Dsk_fit[:], val_fit)
				} else {
					// Si no especifica -> First Fit
					copy(master_boot_record.Dsk_fit[:], "f")
				}

				// Verifico si existe el parametro "Unit" (Opcional)
				if band_unit {
					// Megabytes
					if val_unit == "m" {
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
						total_size = val_size * 1024
					} else {
						// Kilobytes
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024)))
						total_size = val_size
					}
				} else {
					// Si no especifica -> Megabytes
					copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
					total_size = val_size * 1024
				}

				// Inicializar Parcticiones
				for i := 0; i < 4; i++ {
					copy(master_boot_record.Mbr_partition[i].Part_status[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_type[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_fit[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_start[:], "-1")
					copy(master_boot_record.Mbr_partition[i].Part_size[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_name[:], "")
					copy(master_boot_record.Mbr_partition[i].Part_correlative[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_id[:], "0")
				}
				// pasar de entero a string
				// Convierto de entero a string
				str_total_size := strconv.Itoa(total_size)

				// Crear el archivo con el tamaño especificado
				err := llenarArchivoCeros(val_path, str_total_size)

				// ERROR
				if err != nil {
					MsgError(err)
				}

				// Se escriben los datos en disco
				// Apertura del archivo
				disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

				// ERROR
				if err != nil {
					MsgError(err)
				}

				// Conversion de struct a bytes
				mbr_byte := struct_a_bytes(master_boot_record)

				// Se posiciona al inicio del archivo para guardar la informacion del disco
				newpos, err := disco.Seek(0, os.SEEK_SET)

				// ERROR
				if err != nil {
					MsgError(err)
				}

				// Escritura de struct en archivo binario
				_, err = disco.WriteAt(mbr_byte, newpos)

				// ERROR
				if err != nil {
					MsgError(err)
				}
				Salida_comando += "Se creo el disco correctamente..." + "\n"
				err = disco.Close()
				if err != nil {
					return
				}
			}
		}
	}
	Salida_comando += "MENSAJE: El comando MKDISK aqui termina" + "\n"
}

func crear_disco(ruta string) {
	aux, err := filepath.Abs(ruta)
	Salida_comando += "Creando disco en la ruta: " + aux + "\n"

	// ERROR
	if err != nil {
		MsgError(err)
	}

	// Crea el directorio de forma recursiva
	err = os.MkdirAll(filepath.Dir(aux), os.ModePerm)
	if err != nil {
		MsgError(err)
	}
	Salida_comando += "Directorio creado correctamente..." + "\n"

	// Cambia los permisos del directorio
	err = os.Chmod(filepath.Dir(aux), 0777)
	if err != nil {
		MsgError(err)
	}

	// Verifica si existe la ruta para el archivo
	if _, err := os.Stat(filepath.Dir(aux)); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			fmt.Println("Invalido: ", err)
		}
	}
}

func llenarArchivoCeros(val_path string, str_total_size string) error {
	var cmd *exec.Cmd

	// Convertir str_total_size a bytes
	totalSizeKB, err := strconv.Atoi(str_total_size)
	if err != nil {
		return err
	}
	totalSizeBytes := totalSizeKB * 1024

	if runtime.GOOS == "windows" {
		// Comando para Windows
		cmd = exec.Command("cmd", "/C", "fsutil file createnew "+val_path+" "+strconv.Itoa(totalSizeBytes))
	} else {
		// Comando para Linux
		cmd = exec.Command("/bin/sh", "-c", "dd if=/dev/zero of=\""+val_path+"\" bs=1024 count="+str_total_size)
		cmd.Dir = "/"
	}

	_, err = cmd.Output()
	return err

}
