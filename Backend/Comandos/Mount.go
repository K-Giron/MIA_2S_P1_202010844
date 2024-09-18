package Comandos

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Arreglo para guardar las particiones montadas
var ParticionesMontadas []ParticionMontada

// Estructura para guardar las particiones montadas
type ParticionMontada struct {
	id        string
	Direccion string
	Nombre    string
	Letra     string
	Num       int
	Estado    int // 0 = desmontada, 1 = montada
}

func Mount(commandArray []string) {
	Salida_comando += "MENSAJE: El comando MOUNT aqui inicia" + "\n"
	// Variables para los valores de los parametros
	val_path := ""
	val_name := ""
	// Banderas para verificar los parametros y ver si se repiten
	band_path := false
	band_name := false
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
			if band_path {
				fmt.Println("Invalido: El parametro -path ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO OBLIGATORIO -> NAME */
		case strings.Contains(data, "name="):
			// Valido si el parametro ya fue ingresado
			if band_name {
				fmt.Println("Invalido: El parametro -name ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_name = true

			// Reemplaza comillas
			val_name = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			fmt.Println("Invalido: El parametro " + data + " no es valido")
		}
	}

	// Verifico si no hay errores
	if !band_error {
		// Verifico que el parametro "Path" (Obligatorio) este ingresado
		if band_path && band_name {
			index_p := buscar_particion_p_e(val_path, val_name)
			if index_p != -1 {
				f, err := os.OpenFile(val_path, os.O_RDWR, 0660)

				if err == nil {
					if buscarParticionM(val_path, val_name, ParticionesMontadas[:]) {
						Salida_comando += "Particion ya montada" + "\n"
					} else {
						mbr_empty := Mbr{}

						// Calculo del tamaño de struct en bytes
						mbr2 := struct_a_bytes(mbr_empty)
						sstruct := len(mbr2)

						// Lectrura del archivo binario desde el inicio
						lectura := make([]byte, sstruct)
						f.Seek(0, io.SeekStart)
						f.Read(lectura)

						// Conversion de bytes a struct
						master_boot_record := bytes_a_struct_mbr(lectura)

						// Colocamos la particion ocupada
						copy(master_boot_record.Mbr_partition[index_p].Part_status[:], "2")

						//numero de particion
						num := buscarNumero(val_path, ParticionesMontadas[:])
						//letra de disco
						letra := buscarLetra(val_path, ParticionesMontadas[:])
						//terminacion de carnet
						id := "44" + strconv.Itoa(num) + letra

						// guardar el id en el mbr de la particion
						copy(master_boot_record.Mbr_partition[index_p].Part_id[:], id)
						//guardar el correlative en el mbr de la particion
						copy(master_boot_record.Mbr_partition[index_p].Part_correlative[:], strconv.Itoa(1))

						// Conversion de struct a bytes
						mbr_byte := struct_a_bytes(master_boot_record)

						//Se posiciona al inicio del archivo
						f.Seek(0, io.SeekStart)
						f.Write(mbr_byte)
						f.Close()

						//guardar particion montada
						ParticionesMontadas = append(ParticionesMontadas, ParticionMontada{id, val_path, val_name, letra, num, 1})
						Salida_comando += "Particion ya montada " + "\n"
						//concatenar el recorrido de las particiones montadas
						for i := 0; i < len(ParticionesMontadas); i++ {
							Salida_comando += "ID: " + ParticionesMontadas[i].id + " Direccion: " + ParticionesMontadas[i].Direccion + " Nombre: " + ParticionesMontadas[i].Nombre + " Letra: " + ParticionesMontadas[i].Letra + " Numero: " + strconv.Itoa(ParticionesMontadas[i].Num) + " Estado: " + strconv.Itoa(ParticionesMontadas[i].Estado) + "\n"
						}
					}

				} else {
					fmt.Println("Invalido: No se encuentra el disco")
				}
			} else {
				// Buscar en las particiones logicas
				index_p := buscar_particion_l(val_path, val_name)
				if index_p != -1 {
					f, err := os.OpenFile(val_path, os.O_RDWR, 0660)

					if err == nil {
						ebr_empty := Ebr{}

						// Calculo del tamaño de struct en bytes
						ebr2 := struct_a_bytes(ebr_empty)
						sstruct := len(ebr2)

						// Lectrura del archivo binario desde el inicio
						lectura := make([]byte, sstruct)
						f.Seek(int64(index_p), io.SeekStart)
						f.Read(lectura)

						// Conversion de bytes a struct
						extended_boot_record := bytes_a_struct_ebr(lectura)

						// Colocamos la particion ocupada
						copy(extended_boot_record.Part_status[:], "2")

						// Conversion de struct a bytes
						mbr_byte := struct_a_bytes(extended_boot_record)

						// Se posiciona al inicio del archivo para guardar la informacion del disco
						f.Seek(int64(index_p), io.SeekStart)
						f.Write(mbr_byte)
						f.Close()

						if buscarParticionM(val_path, val_name, ParticionesMontadas[:]) {
							Salida_comando += "Particion ya montada" + "\n"
						} else {
							//generar id
							num := buscarNumero(val_path, ParticionesMontadas[:])
							letra := buscarLetra(val_path, ParticionesMontadas[:])
							id := "44" + strconv.Itoa(num) + letra

							//guardar particion montada
							ParticionesMontadas[num] = ParticionMontada{id, val_path, val_name, letra, num, 1}
							Salida_comando += "Particion logica montada con exito" + "\n"
							//concatenar el recorrido de las particiones montadas
							for i := 0; i < len(ParticionesMontadas); i++ {
								Salida_comando += "ID: " + ParticionesMontadas[i].id + " Direccion: " + ParticionesMontadas[i].Direccion + " Nombre: " + ParticionesMontadas[i].Nombre + " Letra: " + ParticionesMontadas[i].Letra + " Numero: " + strconv.Itoa(ParticionesMontadas[i].Num) + " Estado: " + strconv.Itoa(ParticionesMontadas[i].Estado) + "\n"
							}
						}

					} else {
						fmt.Println("Invalido: No se encuentra el disco")
					}
				} else {
					fmt.Println("Invalido: No se encuentra la particion")
				}
			}

		} else {
			fmt.Println("Invalido: El comando MOUNT requiere los parametros obligatorios [-path] y [-name]")
		}
	}
	Salida_comando += "MENSAJE: El comando MOUNT aqui termina" + "\n"
}

func buscarParticionM(direccion string, nombre string, lista []ParticionMontada) bool {
	for i := 0; i < len(lista); i++ {
		if lista[i].Direccion == direccion && lista[i].Nombre == nombre {
			return true
		}
	}
	return false
}

func buscarNumero(direccion string, lista []ParticionMontada) int {
	retorno := 1

	for i := 0; i < len(lista); i++ {
		if direccion == lista[i].Direccion {
			return retorno
		}
		retorno++
	}

	return retorno
}

func buscarLetra(direccion string, lista []ParticionMontada) string {
	num_ascii := 65

	for i := 0; i < len(lista); i++ {
		if lista[i].Direccion == direccion && string(rune(num_ascii)) == lista[i].Letra {
			num_ascii++
		}
	}
	return string(rune(num_ascii))
}

// Busca particiones Primarias o Extendidas
func buscar_particion_p_e(direccion string, nombre string) int {
	// Apertura del archivo
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		mbr_empty := Mbr{}

		// Calculo del tamaño de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)

		s_part_status := ""
		s_part_name := ""

		// Recorro las 4 particiones
		for i := 0; i < 4; i++ {
			// Antes de comparar limpio la cadena
			s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
			s_part_status = strings.Trim(s_part_status, "\x00")

			if s_part_status != "1" {
				// Antes de comparar limpio la cadena
				s_part_name = string(master_boot_record.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")
				if s_part_name == nombre {
					return i
				}

			}

		}
	}

	return -1
}

func buscar_particion_l(direccion string, nombre string) int {
	// Apertura del archivo
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		extendida := -1
		mbr_empty := Mbr{}

		// Calculo del tamaño de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)

		s_part_type := ""

		// Recorro las 4 particiones
		for i := 0; i < 4; i++ {
			// Antes de comparar limpio la cadena
			s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
			s_part_type = strings.Trim(s_part_type, "\x00")

			if s_part_type != "e" {
				extendida = i
				break
			}
		}

		// Si hay extendida
		if extendida != -1 {
			ebr_empty := Ebr{}

			ebr2 := struct_a_bytes(ebr_empty)
			sstruct := len(ebr2)

			// Lectrura del archivo binario desde el inicio
			lectura := make([]byte, sstruct)

			s_part_start := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
			s_part_start = strings.Trim(s_part_start, "\x00")
			i_part_start, _ := strconv.Atoi(s_part_start)

			f.Seek(int64(i_part_start), io.SeekStart)

			n_leidos, _ := f.Read(lectura)
			pos_actual, _ := f.Seek(0, io.SeekCurrent)

			s_part_size := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
			s_part_size = strings.Trim(s_part_size, "\x00")
			i_part_size, _ := strconv.Atoi(s_part_size)

			for (n_leidos != 0) && (pos_actual < int64(i_part_start+i_part_size)) {
				n_leidos, _ = f.Read(lectura)
				pos_actual, _ = f.Seek(0, io.SeekCurrent)

				// Conversion de bytes a struct
				extended_boot_record := bytes_a_struct_ebr(lectura)

				s_part_name_ext := string(extended_boot_record.Part_name[:])
				s_part_name_ext = strings.Trim(s_part_name_ext, "\x00")

				ebr_empty_byte := struct_a_bytes(ebr_empty)

				if s_part_name_ext == nombre {
					return int(pos_actual) - len(ebr_empty_byte)
				}
			}
		}
		f.Close()
	}

	return -1
}

/*
func GetMount(op string, id string, p *string) Partition {
	if !(id[0] == '4' && id[1] == '4') {
		Error(op, "El primer identificador no es válido.")
		return Partition{}
	}
	//letra := id[len(id)-1]
	id = strings.ReplaceAll(id, "44", "")
	i, _ := strconv.Atoi(string(id[0] - 1))
	if i < 0 {
		Error(op, "El primer identificador no es válido.")
		return Partition{}
	}
	for j := 10; j < 26; j++ {
		if ParticionesMontadas[i].Estado == 1 {
			//comparar letra
			if ParticionesMontadas[i].Letra == string(rune(j)) {

				path := ""
				for k := 0; k < len(ParticionesMontadas[i].Direccion); k++ {
					if ParticionesMontadas[i].Direccion[k] != 0 {
						path += string(ParticionesMontadas[i].Direccion[k])
					}
				}

				file, error := os.Open(strings.ReplaceAll(path, "\"", ""))
				if error != nil {
					Error(op, "Error al abrir el archivo")
					return Partition{}
				}
				disk := NewMbr()
				file.Seek(0, 0)

				data := make([]byte, int(unsafe.Sizeof(disk)))
				buffer := bytes.NewBuffer(data)
				err_ := binary.Read(buffer, binary.BigEndian, &disk)

				if err_ != nil {
					Error("FDSIK", "Error al leer el archivo")
					return Partition{}
				}
				file.Close()

				nombreParticion := ""
				for k := 0; k < len(ParticionesMontadas[i].Nombre); k++ {
					if ParticionesMontadas[i].Nombre[k] != 0 {
						nombreParticion += string(ParticionesMontadas[i].Nombre[k])
					}
				}
				*p = path
				return *BuscarParticiones(disk, nombreParticion, path)

			}
		}
	}
	return Partition{}
}
*/
