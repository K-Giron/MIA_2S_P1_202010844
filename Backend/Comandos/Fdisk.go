package Comandos

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// estructura ---------------------------------------
type Ebr = struct {
	Part_status [100]byte
	Part_fit    [100]byte
	Part_start  [100]byte
	Part_size   [100]byte
	Part_next   [100]byte
	Part_name   [100]byte
}

/* FDISK */
func Fdisk(commandArray []string) {
	//acceder a la salida_comandos en menu.go
	Salida_comando += "Comando FDISK inicia" + "\n"

	// Variables para los valores de los parametros
	val_size := 0
	val_unit := ""
	val_path := ""
	val_type := ""
	val_fit := ""
	val_name := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_size := false
	band_unit := false
	band_path := false
	band_type := false
	band_fit := false
	band_name := false
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
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_size = true

			// Conversion a entero
			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size
			Salida_comando += "Size: " + strconv.Itoa(val_size) + "\n"
			// ERROR de conversion
			if err != nil {
				MsgError(err)
				band_error = true
			}

			// Valido que el tamaño sea positivo
			if val_size < 0 {
				band_error = true
				break
			}
		/* PARAMETRO OPCIONAL -> UNIT */
		case strings.Contains(data, "unit="):
			// Valido si el parametro ya fue ingresado
			if band_unit {
				band_error = true
				break
			}

			// Reemplaza comillas y lo paso a minusculas
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)
			Salida_comando += "Unit: " + val_unit + "\n"
			if val_unit == "b" || val_unit == "k" || val_unit == "m" {
				// Activo la bandera del parametro
				band_unit = true
			} else {
				// Parametro no valido
				band_error = true
				break
			}
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
			Salida_comando += "Path: " + val_path + "\n"
		/* PARAMETRO OPCIONAL -> TYPE */
		case strings.Contains(data, "type="):
			if band_type {
				band_error = true
				break
			}

			// Reemplaza comillas y lo paso a minusculas
			val_type = strings.Replace(val_data, "\"", "", 2)
			val_type = strings.ToLower(val_type)
			Salida_comando += "Type: " + val_type + "\n"
			if val_type == "p" || val_type == "e" || val_type == "l" {
				// Activo la bandera del parametro
				band_type = true
			} else {
				// Parametro no valido
				band_error = true
				break
			}
		/* PARAMETRO OPCIONAL -> FIT */
		case strings.Contains(data, "fit="):
			// Valido si el parametro ya fue ingresado
			if band_fit {
				band_error = true
				break
			}

			// Le quito las comillas y lo paso a minusculas
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "w"
			} else {
				band_error = true
				break
			}
			Salida_comando += "Fit: " + val_fit + "\n"
		/* PARAMETRO OBLIGATORIO -> NAME */
		case strings.Contains(data, "name="):
			// Valido si el parametro ya fue ingresado
			if band_name {
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_name = true

			// Reemplaza comillas
			val_name = strings.Replace(val_data, "\"", "", 2)
			Salida_comando += "Name: " + val_name + "\n"
		/* PARAMETRO NO VALIDO */
		default:
			fmt.Print("Invalido: Parametro no valido" + "\n")
		}
	}

	// Verifico si no hay errores
	if !band_error {
		if band_size {
			if band_path {
				if band_name {
					if band_type {
						if val_type == "p" {
							// Primaria
							crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
						} else if val_type == "e" {
							// Extendida
							crear_particion_extendia(val_path, val_name, val_size, val_fit, val_unit)
						} else {
							// Logica
							crear_particion_logica(val_path, val_name, val_size, val_fit, val_unit)
						}
					} else {
						// Si no lo indica se tomara como Primaria
						crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
					}
				} else {
					fmt.Println("Invalido: El parametro -name no fue ingresado")
				}
			} else {
				fmt.Println("Invalido: El parametro -path no fue ingresado")
			}
		} else {
			fmt.Println("Invalido: El parametro -size no fue ingresado")
		}
	}
	Salida_comando += "Comando FDISK finalizado" + "\n"
}

// Crea la Particion Primaria
func crear_particion_primaria(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := Mbr{}
	var empty [100]byte

	// Verifico si tiene Ajuste
	if fit != "" {
		aux_fit = fit
	} else {
		// Por default es Peor ajuste
		aux_fit = "w"
	}

	// Verifico si tiene Unidad
	if unit != "" {
		aux_unit = unit

		// *Bytes
		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			// *Kilobytes
			size_bytes = size * 1024
		} else {
			// *Megabytes
			size_bytes = size * 1024 * 1024
		}
	} else {
		// Por default Kilobytes
		size_bytes = size * 1024
	}

	// Abro el archivo para lectura con opcion a modificar
	// * O_RDWR = Lectura y Escritura
	// * 0660 = Permisos de lectura y escritura
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		fmt.Println("Invalido: No se pudo abrir el archivo")
	} else {
		// Bandera para ver si hay una particion disponible
		band_particion := false
		// Valor del numero de particion
		num_particion := 0

		// Calculo del tamaño de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		// Se posiciona al inicio del archivo para leer la informacion del disco
		// * 0 = Desde el inicio
		// * SEEK_SET = Posicionamiento
		f.Seek(0, os.SEEK_SET)
		// Lectura de conjunto de bytes en archivo binario
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_start := ""

			// Recorro las 4 particiones (Logica de Particiones)
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
				// Le quito los caracteres null
				s_part_start = strings.Trim(s_part_start, "\x00")

				// Verifico si en las particiones hay espacio
				// * -1 = No hay particion
				if s_part_start == "-1" {
					band_particion = true
					num_particion = i
					break
				}
			}

			// Si hay una particion disponible
			if band_particion {
				espacio_usado := 0
				s_part_size := ""
				i_part_size := 0
				s_part_status := ""

				// Recorro las 4 particiones
				for i := 0; i < 4; i++ {
					// Obtengo el espacio utilizado
					s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
					// Le quito los caracteres null
					s_part_size = strings.Trim(s_part_size, "\x00")
					i_part_size, _ = strconv.Atoi(s_part_size)

					// Obtengo el status de la particion
					s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
					// Le quito los caracteres null
					s_part_status = strings.Trim(s_part_status, "\x00")

					if s_part_status != "1" {
						// Le sumo el valor al espacio
						espacio_usado += i_part_size
					}
				}

				/* Tamaño del disco */

				// Obtengo el tamaño del disco
				s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
				// Le quito los caracteres null
				s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
				i_tamaño_disco, _ := strconv.Atoi(s_tamaño_disco)

				espacio_disponible := i_tamaño_disco - espacio_usado
				Salida_comando += "[ESPACIO DISPONIBLE] " + strconv.Itoa(espacio_disponible) + " Bytes" + "\n"
				Salida_comando += "[ESPACIO NECESARIO] " + strconv.Itoa(size_bytes) + " Bytes" + "\n"

				// Verifico que haya espacio suficiente
				if espacio_disponible >= size_bytes {
					// Verifico si no existe una particion con ese nombre
					if !existe_particion(direccion, nombre) {
						// Antes de comparar limpio la cadena
						s_dsk_fit := string(master_boot_record.Dsk_fit[:])
						s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")

						/*  Primer Ajuste  */
						if s_dsk_fit == "f" {
							copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], aux_fit)

							// Si esta iniciando
							if num_particion == 0 {
								// Guardo el inicio de la particion y dejo un espacio de separacion
								mbr_empty_byte := struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
							} else {
								// Obtengo el inicio de la particion anterior
								s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])
								// Le quito los caracteres null
								s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
								i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

								// Obtengo el tamaño de la particion anterior
								s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])
								// Le quito los caracteres null
								s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
								i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant+1))
							}

							copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], nombre)

							// Se guarda de nuevo el MBR

							// Conversion de struct a bytes
							mbr_byte := struct_a_bytes(master_boot_record)

							// Se posiciona al inicio del archivo para guardar la informacion del disco
							// * 0 = Desde el inicio
							// * SEEK_SET = Posicionamiento
							f.Seek(0, os.SEEK_SET)
							f.Write(mbr_byte)

							// Obtengo el inicio de la particion
							s_part_start = string(master_boot_record.Mbr_partition[num_particion].Part_start[:])
							// Le quito los caracteres null
							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)

							// Se posiciona en el inicio de la particion
							// * int64(i_part_start) = Posicionamiento de donde inicia la particion
							// * SEEK_SET = Posicionamiento
							f.Seek(int64(i_part_start), os.SEEK_SET)

							// Lo llena de unos
							for i := 1; i < size_bytes; i++ {
								f.Write([]byte{1})
							}
							Salida_comando += "[SUCCES] La Particion primaria fue creada con exito!" + "\n"
						} else if s_dsk_fit == "b" {
							/*  Mejor Ajuste  */
							best_index := num_particion

							// Variables para conversiones
							s_part_start_act := ""
							s_part_status_act := ""
							s_part_size_act := ""
							i_part_size_act := 0
							s_part_start_best := ""
							i_part_start_best := 0
							s_part_start_best_ant := ""
							i_part_start_best_ant := 0
							s_part_size_best := ""
							i_part_size_best := 0
							s_part_size_best_ant := ""
							i_part_size_best_ant := 0

							for i := 0; i < 4; i++ {
								// Obtengo el inicio de la particion actual
								s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start_act = strings.Trim(s_part_start_act, "\x00")

								// Obtengo el size de la particion actual
								s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
								// Le quito los caracteres null
								s_part_status_act = strings.Trim(s_part_status_act, "\x00")

								// Obtengo la posicion de la particion actual
								s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_size_act = strings.Trim(s_part_size_act, "\x00")
								i_part_size_act, _ = strconv.Atoi(s_part_size_act)

								if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
									if i != num_particion {
										// Obtengo el tamaño de la particion del mejor indice
										s_part_size_best = string(master_boot_record.Mbr_partition[best_index].Part_size[:])
										// Le quito los caracteres null
										s_part_size_best = strings.Trim(s_part_size_best, "\x00")
										i_part_size_best, _ = strconv.Atoi(s_part_size_best)

										// Obtengo la posicion de la particion actual
										s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
										// Le quito los caracteres null
										s_part_size_act = strings.Trim(s_part_size_act, "\x00")
										i_part_size_act, _ = strconv.Atoi(s_part_size_act)

										if i_part_size_best > i_part_size_act {
											best_index = i
											break
										}
									}
								}
							}

							// Primaria
							copy(master_boot_record.Mbr_partition[best_index].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[best_index].Part_fit[:], aux_fit)

							// Si esta iniciando
							if best_index == 0 {
								// Guardo el inicio de la particion y dejo un espacio de separacion
								mbr_empty_byte := struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
							} else {
								// Obtengo el inicio de la particion actual
								s_part_start_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_start[:])
								// Le quito los caracteres null
								s_part_start_best_ant = strings.Trim(s_part_start_best_ant, "\x00")
								i_part_start_best_ant, _ = strconv.Atoi(s_part_start_best_ant)

								// Obtengo el inicio de la particion actual
								s_part_size_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_size[:])
								// Le quito los caracteres null
								s_part_size_best_ant = strings.Trim(s_part_size_best_ant, "\x00")
								i_part_size_best_ant, _ = strconv.Atoi(s_part_size_best_ant)

								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(i_part_start_best_ant+i_part_size_best_ant))
							}

							copy(master_boot_record.Mbr_partition[best_index].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[best_index].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[best_index].Part_name[:], nombre)

							// Se guarda de nuevo el MBR

							// Conversion de struct a bytes
							mbr_byte := struct_a_bytes(master_boot_record)

							// Se posiciona al inicio del archivo para guardar la informacion del disco
							f.Seek(0, os.SEEK_SET)
							f.Write(mbr_byte)

							// Obtengo el inicio de la particion best
							s_part_start_best = string(master_boot_record.Mbr_partition[best_index].Part_start[:])
							// Le quito los caracteres null
							s_part_start_best = strings.Trim(s_part_start_best, "\x00")
							i_part_start_best, _ = strconv.Atoi(s_part_start_best)

							// Conversion de struct a bytes

							// Se posiciona en el inicio de la particion
							f.Seek(int64(i_part_start_best), os.SEEK_SET)

							// Lo llena de unos
							for i := 1; i < size_bytes; i++ {
								f.Write([]byte{1})
							}
							Salida_comando += "[SUCCES] La Particion primaria fue creada con exito!" + "\n"
						} else {
							/*  Peor ajuste  */
							worst_index := num_particion

							// Variables para conversiones
							s_part_start_act := ""
							s_part_status_act := ""
							s_part_size_act := ""
							i_part_size_act := 0
							s_part_start_worst := ""
							i_part_start_worst := 0
							s_part_start_worst_ant := ""
							i_part_start_worst_ant := 0
							s_part_size_worst := ""
							i_part_size_worst := 0
							s_part_size_worst_ant := ""
							i_part_size_worst_ant := 0

							for i := 0; i < 4; i++ {
								// Obtengo el inicio de la particion actual
								s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start_act = strings.Trim(s_part_start_act, "\x00")

								// Obtengo el size de la particion actual
								s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
								// Le quito los caracteres null
								s_part_status_act = strings.Trim(s_part_status_act, "\x00")

								// Obtengo la posicion de la particion actual
								s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_size_act = strings.Trim(s_part_size_act, "\x00")
								i_part_size_act, _ = strconv.Atoi(s_part_size_act)

								if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
									if i != num_particion {
										// Obtengo el tamaño de la particion del mejor indice
										s_part_size_worst = string(master_boot_record.Mbr_partition[worst_index].Part_size[:])
										// Le quito los caracteres null
										s_part_size_worst = strings.Trim(s_part_size_worst, "\x00")
										i_part_size_worst, _ = strconv.Atoi(s_part_size_worst)

										// Obtengo la posicion de la particion actual
										s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
										// Le quito los caracteres null
										s_part_size_act = strings.Trim(s_part_size_act, "\x00")
										i_part_size_act, _ = strconv.Atoi(s_part_size_act)

										if i_part_size_worst < i_part_size_act {
											worst_index = i
											break
										}
									}
								}
							}

							// Particiones Primarias
							copy(master_boot_record.Mbr_partition[worst_index].Part_type[:], "p")
							copy(master_boot_record.Mbr_partition[worst_index].Part_fit[:], aux_fit)

							// Se esta iniciando
							if worst_index == 0 {
								// Guardo el inicio de la particion y dejo un espacio de separacion
								mbr_empty_byte := struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
							} else {
								// Obtengo el inicio de la particion anterior
								s_part_start_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_start[:])
								// Le quito los caracteres null
								s_part_start_worst_ant = strings.Trim(s_part_start_worst_ant, "\x00")
								i_part_start_worst_ant, _ = strconv.Atoi(s_part_start_worst_ant)

								// Obtengo el tamaño de la particion anterior
								s_part_size_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_size[:])
								// Le quito los caracteres null
								s_part_size_worst_ant = strings.Trim(s_part_size_worst_ant, "\x00")
								i_part_size_worst_ant, _ = strconv.Atoi(s_part_size_worst_ant)

								copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(i_part_start_worst_ant+i_part_size_worst_ant))
							}

							copy(master_boot_record.Mbr_partition[worst_index].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[worst_index].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[worst_index].Part_name[:], nombre)

							// Se guarda de nuevo el MBR

							// Conversion de struct a bytes
							mbr_byte := struct_a_bytes(master_boot_record)

							// Escribe desde el inicio del archivo
							f.Seek(0, os.SEEK_SET)
							f.Write(mbr_byte)

							// Obtengo el inicio de la particion best
							s_part_start_worst = string(master_boot_record.Mbr_partition[worst_index].Part_start[:])
							// Le quito los caracteres null
							s_part_start_worst = strings.Trim(s_part_start_worst, "\x00")
							i_part_start_worst, _ = strconv.Atoi(s_part_start_worst)

							// Se posiciona en el inicio de la particion
							f.Seek(int64(i_part_start_worst), os.SEEK_SET)

							// Lo llena de unos
							for i := 1; i < size_bytes; i++ {
								f.Write([]byte{1})
							}
							Salida_comando += "[SUCCES] La Particion primaria fue creada con exito!" + "\n"
						}
					} else {
						fmt.Println("Invalido: Ya existe una particion con ese nombre")
					}
				} else {
					fmt.Println("Invalido: No hay particiones disponibles")
					fmt.Println("MENSAJE: Se recomienda eliminar alguna partición disponible")
				}
			} else {
				fmt.Println("Invalido: No hay particiones disponibles")
			}
		} else {
			fmt.Println("Invalido: El disco no ha sido creado")
		}

		f.Close()
	}
}

// Crea la Particion Extendida
func crear_particion_extendia(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := Mbr{}
	ebr_empty := Ebr{}
	var empty [100]byte

	// Verifico si tiene Ajuste
	if fit != "" {
		aux_fit = fit
	} else {
		// Por default es Peor ajuste
		aux_fit = "w"
	}

	// Verifico si tiene Unidad
	if unit != "" {
		aux_unit = unit

		// *Bytes
		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			// *Kilobytes
			size_bytes = size * 1024
		} else {
			// *Megabytes
			size_bytes = size * 1024 * 1024
		}
	} else {
		// Por default Kilobytes
		size_bytes = size * 1024
	}

	// Abro el archivo para lectura con opcion a modificar
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		MsgError(err)
	} else {
		// Procede a leer el archivo
		band_particion := false
		band_extendida := false
		num_particion := 0

		// Calculo del tamaño de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, os.SEEK_SET)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_type := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					band_extendida = true
					break
				}
			}

			// Si aun no ha creado la extendida
			if !band_extendida {
				s_part_start := ""
				s_part_status := ""
				s_part_size := ""
				i_part_size := 0

				// Recorro las 4 particiones
				for i := 0; i < 4; i++ {
					// Antes de comparar limpio la cadena
					s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
					s_part_start = strings.Trim(s_part_start, "\x00")

					s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
					s_part_status = strings.Trim(s_part_status, "\x00")

					s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
					s_part_size = strings.Trim(s_part_size, "\x00")
					i_part_size, _ = strconv.Atoi(s_part_size)

					// Verifica si existe una particion disponible
					if s_part_start == "-1" || (s_part_status == "1" && i_part_size >= size_bytes) {
						band_particion = true
						num_particion = i
						break
					}
				}

				// Si hay una particion
				if band_particion {
					espacio_usado := 0

					// Recorro las 4 particiones
					for i := 0; i < 4; i++ {
						s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
						s_part_status = strings.Trim(s_part_status, "\x00")

						if s_part_status != "1" {
							// Obtengo el espacio utilizado
							s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
							// Le quito los caracteres null
							s_part_size = strings.Trim(s_part_size, "\x00")
							i_part_size, _ = strconv.Atoi(s_part_size)

							// Le sumo el valor al espacio
							espacio_usado += i_part_size
						}
					}

					/* Tamaño del disco */

					// Obtengo el tamaño del disco
					s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
					// Le quito los caracteres null
					s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
					i_tamaño_disco, _ := strconv.Atoi(s_tamaño_disco)

					espacio_disponible := i_tamaño_disco - espacio_usado

					Salida_comando += "[ESPACIO DISPONIBLE] " + strconv.Itoa(espacio_disponible) + " Bytes" + "\n"
					Salida_comando += "[ESPACIO NECESARIO] " + strconv.Itoa(size_bytes) + " Bytes" + "\n"

					// Verifico que haya espacio suficiente
					if espacio_disponible >= size_bytes {
						// Verifico si no existe una particion con ese nombre
						if !existe_particion(direccion, nombre) {
							// Antes de comparar limpio la cadena
							s_dsk_fit := string(master_boot_record.Dsk_fit[:])
							s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")

							/*  Primer Ajuste  */
							if s_dsk_fit == "f" {
								copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "e")
								copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], aux_fit)

								// Si esta iniciando
								if num_particion == 0 {
									// Guardo el inicio de la particion y dejo un espacio de separacion
									mbr_empty_byte := struct_a_bytes(mbr_empty)
									copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
								} else {
									// Obtengo el inicio de la particion anterior
									s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])
									// Le quito los caracteres null
									s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
									i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

									// Obtengo el tamaño de la particion anterior
									s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])
									// Le quito los caracteres null
									s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
									i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

									copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant+1))
								}

								copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], strconv.Itoa(size_bytes))
								copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
								copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], nombre)

								// Se guarda de nuevo el MBR

								// Conversion de struct a bytes
								mbr_byte := struct_a_bytes(master_boot_record)

								// Escribe en la posicion inicial del archivo
								f.Seek(0, os.SEEK_SET)
								f.Write(mbr_byte)

								// Obtengo el tamaño de la particion
								s_part_start = string(master_boot_record.Mbr_partition[num_particion].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								// Se posiciona en el inicio de la particion
								f.Seek(int64(i_part_start), os.SEEK_SET)

								extended_boot_record := Ebr{}
								copy(extended_boot_record.Part_fit[:], aux_fit)
								copy(extended_boot_record.Part_status[:], "0")
								copy(extended_boot_record.Part_start[:], s_part_start)
								copy(extended_boot_record.Part_size[:], "0")
								copy(extended_boot_record.Part_next[:], "-1")
								copy(extended_boot_record.Part_name[:], "")
								ebr_byte := struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								ebr_empty_byte := struct_a_bytes(ebr_empty)

								// Lo corro una posicion de donde se encuentra
								pos_actual, _ := f.Seek(0, os.SEEK_CUR)
								f.Seek(int64(pos_actual+1), os.SEEK_SET)

								// Lo llena de unos
								for i := 1; i < (size_bytes - int(binary.Size(ebr_empty_byte))); i++ {
									f.Write([]byte{1})
								}
								Salida_comando += "[SUCCES] La Particion extendida fue creada con exito!" + "\n"
							} else if s_dsk_fit == "b" {
								/*  Mejor Ajuste  */
								best_index := num_particion

								// Variables para conversiones
								s_part_start_act := ""
								s_part_status_act := ""
								s_part_size_act := ""
								i_part_size_act := 0
								s_part_start_best := ""
								i_part_start_best := 0
								s_part_start_best_ant := ""
								i_part_start_best_ant := 0
								s_part_size_best := ""
								i_part_size_best := 0
								s_part_size_best_ant := ""
								i_part_size_best_ant := 0

								for i := 0; i < 4; i++ {
									// Obtengo el inicio de la particion actual
									s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
									// Le quito los caracteres null
									s_part_start_act = strings.Trim(s_part_start_act, "\x00")

									// Obtengo el size de la particion actual
									s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
									// Le quito los caracteres null
									s_part_status_act = strings.Trim(s_part_status_act, "\x00")

									// Obtengo la posicion de la particion actual
									s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
									// Le quito los caracteres null
									s_part_size_act = strings.Trim(s_part_size_act, "\x00")
									i_part_size_act, _ = strconv.Atoi(s_part_size_act)

									if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
										if i != num_particion {
											// Obtengo el tamaño de la particion del mejor indice
											s_part_size_best = string(master_boot_record.Mbr_partition[best_index].Part_size[:])
											// Le quito los caracteres null
											s_part_size_best = strings.Trim(s_part_size_best, "\x00")
											i_part_size_best, _ = strconv.Atoi(s_part_size_best)

											// Obtengo la posicion de la particion actual
											s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
											// Le quito los caracteres null
											s_part_size_act = strings.Trim(s_part_size_act, "\x00")
											i_part_size_act, _ = strconv.Atoi(s_part_size_act)

											if i_part_size_best > i_part_size_act {
												best_index = i
												break
											}
										}
									}
								}

								// Extendida
								copy(master_boot_record.Mbr_partition[best_index].Part_type[:], "e")
								copy(master_boot_record.Mbr_partition[best_index].Part_fit[:], aux_fit)

								// Si esta iniciando
								if best_index == 0 {
									// Guardo el inicio de la particion y dejo un espacio de separacion
									mbr_empty_byte := struct_a_bytes(mbr_empty)
									copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
								} else {
									// Obtengo el inicio de la particion actual
									s_part_start_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_start[:])
									// Le quito los caracteres null
									s_part_start_best_ant = strings.Trim(s_part_start_best_ant, "\x00")
									i_part_start_best_ant, _ = strconv.Atoi(s_part_start_best_ant)

									// Obtengo el inicio de la particion actual
									s_part_size_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_size[:])
									// Le quito los caracteres null
									s_part_size_best_ant = strings.Trim(s_part_size_best_ant, "\x00")
									i_part_size_best_ant, _ = strconv.Atoi(s_part_size_best_ant)

									copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(i_part_start_best_ant+i_part_size_best_ant+1))
								}

								copy(master_boot_record.Mbr_partition[best_index].Part_size[:], strconv.Itoa(size_bytes))
								copy(master_boot_record.Mbr_partition[best_index].Part_status[:], "0")
								copy(master_boot_record.Mbr_partition[best_index].Part_name[:], nombre)

								// Se guarda de nuevo el MBR

								// Conversion de struct a bytes
								mbr_byte := struct_a_bytes(master_boot_record)

								// Se escribe al inicio del archivo
								f.Seek(0, os.SEEK_SET)
								f.Write(mbr_byte)

								// Obtengo el inicio de la particion best
								s_part_start_best = string(master_boot_record.Mbr_partition[best_index].Part_start[:])
								// Le quito los caracteres null
								s_part_start_best = strings.Trim(s_part_start_best, "\x00")
								i_part_start_best, _ = strconv.Atoi(s_part_start_best)

								// Se posiciona en el inicio de la particion
								f.Seek(int64(i_part_start_best), os.SEEK_SET)

								extended_boot_record := Ebr{}
								copy(extended_boot_record.Part_fit[:], aux_fit)
								copy(extended_boot_record.Part_status[:], "0")
								copy(extended_boot_record.Part_start[:], s_part_start_best)
								copy(extended_boot_record.Part_size[:], "0")
								copy(extended_boot_record.Part_next[:], "-1")
								copy(extended_boot_record.Part_name[:], "")
								ebr_byte := struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								// Lo corro una posicion de donde se encuentra
								pos_actual, _ := f.Seek(0, os.SEEK_CUR)
								f.Seek(int64(pos_actual+1), os.SEEK_SET)

								ebr_empty_byte := struct_a_bytes(mbr_empty)

								// Lo llena de unos
								for i := 1; i < (size_bytes - int(binary.Size(ebr_empty_byte))); i++ {
									f.Write([]byte{1})
								}
								Salida_comando += "[SUCCES] La Particion extendida fue creada con exito!" + "\n"
							} else {
								/*  Peor ajuste  */
								worst_index := num_particion

								// Variables para conversiones
								s_part_start_act := ""
								s_part_status_act := ""
								s_part_size_act := ""
								i_part_size_act := 0
								s_part_start_worst := ""
								i_part_start_worst := 0
								s_part_start_worst_ant := ""
								i_part_start_worst_ant := 0
								s_part_size_worst := ""
								i_part_size_worst := 0
								s_part_size_worst_ant := ""
								i_part_size_worst_ant := 0

								for i := 0; i < 4; i++ {
									// Obtengo el inicio de la particion actual
									s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
									// Le quito los caracteres null
									s_part_start_act = strings.Trim(s_part_start_act, "\x00")

									// Obtengo el size de la particion actual
									s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
									// Le quito los caracteres null
									s_part_status_act = strings.Trim(s_part_status_act, "\x00")

									// Obtengo la posicion de la particion actual
									s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
									// Le quito los caracteres null
									s_part_size_act = strings.Trim(s_part_size_act, "\x00")
									i_part_size_act, _ = strconv.Atoi(s_part_size_act)

									if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
										if i != num_particion {
											// Obtengo el tamaño de la particion del mejor indice
											s_part_size_worst = string(master_boot_record.Mbr_partition[worst_index].Part_size[:])
											// Le quito los caracteres null
											s_part_size_worst = strings.Trim(s_part_size_worst, "\x00")
											i_part_size_worst, _ = strconv.Atoi(s_part_size_worst)

											// Obtengo la posicion de la particion actual
											s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
											// Le quito los caracteres null
											s_part_size_act = strings.Trim(s_part_size_act, "\x00")
											i_part_size_act, _ = strconv.Atoi(s_part_size_act)

											if i_part_size_worst < i_part_size_act {
												worst_index = i
												break
											}
										}
									}
								}

								// Particiones Extendidas
								copy(master_boot_record.Mbr_partition[worst_index].Part_type[:], "e")
								copy(master_boot_record.Mbr_partition[worst_index].Part_fit[:], aux_fit)

								// Se esta iniciando
								if worst_index == 0 {
									// Guardo el inicio de la particion y dejo un espacio de separacion
									mbr_empty_byte := struct_a_bytes(mbr_empty)
									copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
								} else {
									// Obtengo el inicio de la particion actual
									s_part_start_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_start[:])
									// Le quito los caracteres null
									s_part_start_worst_ant = strings.Trim(s_part_start_worst_ant, "\x00")
									i_part_start_worst_ant, _ = strconv.Atoi(s_part_start_worst_ant)

									// Obtengo el inicio de la particion actual
									s_part_size_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_size[:])
									// Le quito los caracteres null
									s_part_size_worst_ant = strings.Trim(s_part_size_worst_ant, "\x00")
									i_part_size_worst_ant, _ = strconv.Atoi(s_part_size_worst_ant)

									copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(i_part_start_worst_ant+i_part_size_worst_ant+1))
								}

								copy(master_boot_record.Mbr_partition[worst_index].Part_size[:], strconv.Itoa(size_bytes))
								copy(master_boot_record.Mbr_partition[worst_index].Part_status[:], "0")
								copy(master_boot_record.Mbr_partition[worst_index].Part_name[:], nombre)

								// Se guarda de nuevo el MBR

								// Conversion de struct a bytes
								mbr_byte := struct_a_bytes(master_boot_record)

								// Se escribe desde el inicio del archivo
								f.Seek(0, os.SEEK_SET)
								f.Write(mbr_byte)

								// Obtengo el inicio de la particion best
								s_part_start_worst = string(master_boot_record.Mbr_partition[worst_index].Part_start[:])
								// Le quito los caracteres null
								s_part_start_worst = strings.Trim(s_part_start_worst, "\x00")
								i_part_start_worst, _ = strconv.Atoi(s_part_start_worst)

								// Se posiciona en el inicio de la particion
								f.Seek(int64(i_part_start_worst), os.SEEK_SET)

								extended_boot_record := Ebr{}
								copy(extended_boot_record.Part_fit[:], aux_fit)
								copy(extended_boot_record.Part_status[:], "0")
								copy(extended_boot_record.Part_start[:], s_part_start_worst)
								copy(extended_boot_record.Part_size[:], "0")
								copy(extended_boot_record.Part_next[:], "-1")
								copy(extended_boot_record.Part_name[:], "")
								ebr_byte := struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								// Lo corro una posicion de donde se encuentra
								pos_actual, _ := f.Seek(0, os.SEEK_CUR)
								f.Seek(int64(pos_actual+1), os.SEEK_SET)

								ebr_empty_byte := struct_a_bytes(mbr_empty)

								// Lo llena de unos
								for i := 1; i < (size_bytes - int(binary.Size(ebr_empty_byte))); i++ {
									f.Write([]byte{1})
								}
								Salida_comando += "[SUCCES] La Particion extendida fue creada con exito!" + "\n"
							}
						} else {
							fmt.Println("Invalido: Ya existe una particion con ese nombre")
						}
					} else {
						fmt.Println("Invalido: No hay particiones disponibles")
						fmt.Println("MENSAJE: Se recomienda eliminar alguna partición disponible")
					}
				} else {
					fmt.Println("Invalido: No hay particiones disponibles")
				}
			} else {
				fmt.Println("Invalido: Ya existe una particion extendida")
			}
		} else {
			fmt.Println("Invalido: El disco no ha sido creado")
		}
		f.Close()
	}
}

// Crea la Particion Logica
func crear_particion_logica(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := Mbr{}
	ebr_empty := Ebr{}
	var empty [100]byte

	// Verifico si tiene Ajuste
	if fit != "" {
		aux_fit = fit
	} else {
		// Por default es Peor ajuste
		aux_fit = "w"
	}

	// Verifico si tiene Unidad
	if unit != "" {
		aux_unit = unit

		// *Bytes
		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			// *Kilobytes
			size_bytes = size * 1024
		} else {
			// *Megabytes
			size_bytes = size * 1024 * 1024
		}
	} else {
		// Por default Kilobytes
		size_bytes = size * 1024
	}

	// Abro el archivo para lectura con opcion a modificar
	// * direccion -> Nombre del disco
	// *os.O_RDWR -> Abre el archivo para lectura y escritura
	// *0660 -> Permisos de lectura y escritura
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		fmt.Println("Invalido: No se pudo abrir el archivo")
	} else {
		// Calculo del tamaño de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		// Se posiciona en el inicio del archivo
		// * 0 -> Posiciona al inicio del archivo
		// * os.SEEK_SET -> Posiciona al inicio del archivo
		f.Seek(0, os.SEEK_SET)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_type := ""
			num_extendida := -1

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					num_extendida = i
					break
				}
			}

			if !existe_particion(direccion, nombre) {
				if num_extendida != -1 {
					s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
					s_part_start = strings.Trim(s_part_start, "\x00")
					i_part_start, _ := strconv.Atoi(s_part_start)

					cont := i_part_start

					// Se posiciona en el inicio de la particion
					f.Seek(int64(cont), os.SEEK_SET)

					// Calculo del tamaño de struct en bytes
					ebr2 := struct_a_bytes(ebr_empty)
					sstruct := len(ebr2)

					// Lectrura del archivo binario desde el inicio
					lectura := make([]byte, sstruct)
					f.Read(lectura)

					// Conversion de bytes a struct
					extended_boot_record := bytes_a_struct_ebr(lectura)

					// Obtencion de datos
					s_part_size_ext := string(extended_boot_record.Part_size[:])
					s_part_size_ext = strings.Trim(s_part_size_ext, "\x00")

					if s_part_size_ext == "0" {
						// Obtencion de datos
						s_part_size := string(master_boot_record.Mbr_partition[num_extendida].Part_size[:])
						s_part_size = strings.Trim(s_part_size, "\x00")
						i_part_size, _ := strconv.Atoi(s_part_size)

						Salida_comando += "[ESPACIO DISPONIBLE] " + strconv.Itoa(i_part_size) + " Bytes" + "\n"
						Salida_comando += "[ESPACIO NECESARIO] " + strconv.Itoa(size_bytes) + " Bytes" + "\n"

						// Si excede el tamaño de la extendida
						if i_part_size < size_bytes {
							fmt.Println("Invalido: El tamaño de la particion excede el tamaño de la extendida")
						} else {
							copy(extended_boot_record.Part_status[:], "0")
							copy(extended_boot_record.Part_fit[:], aux_fit)

							// Posicion actual en el archivo
							pos_actual, _ := f.Seek(0, os.SEEK_CUR)

							copy(extended_boot_record.Part_start[:], strconv.Itoa(int(pos_actual)-int(binary.Size(ebr_empty))+1))
							copy(extended_boot_record.Part_size[:], strconv.Itoa(size_bytes))
							copy(extended_boot_record.Part_next[:], "-1")
							copy(extended_boot_record.Part_name[:], nombre)

							// Obtencion de datos
							s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)

							// Se posiciona en el inicio de la particion
							ebr_byte := struct_a_bytes(extended_boot_record)
							f.Seek(int64(i_part_start), os.SEEK_SET)
							f.Write(ebr_byte)
							Salida_comando += "[SUCCES] La Particion logica fue creada con exito!" + "\n"
						}
					} else {
						fmt.Println("Invalido: Ya existe una particion logica")
					}
				} else {
					fmt.Println("Invalido: No existe una particion extendida")
				}
			} else {
				fmt.Println("Invalido: Ya existe una particion con ese nombre")
			}
		} else {
			fmt.Println("Invalido: El disco no ha sido creado")
		}
		f.Close()
	}
}

/* Ejemplo 4: Fdisk 1.0 */
// Verifica si el nombre de la particion esta disponible
func existe_particion(direccion string, nombre string) bool {
	extendida := -1
	mbr_empty := Mbr{}
	ebr_empty := Ebr{}
	var empty [100]byte
	cont := 0
	fin_archivo := false

	// Abro el archivo para lectura con opcion a modificar
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		MsgError(err)
	} else {
		// Procedo a leer el archivo

		// Calculo del tamano de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		// make -> Crea un slice de bytes con el tamaño indicado (sstruct)
		// ReadAt -> Lee el archivo binario desde la posicion indicada (0) y lo guarda en el slice de bytes
		// Slice de byte es un arreglo de bytes que se puede modificar y con ReadAt se llena con los bytes del archivo
		lectura := make([]byte, sstruct)
		_, err = f.ReadAt(lectura, 0)

		// ERROR
		if err != nil && err != io.EOF {
			MsgError(err)
		}

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)
		sstruct = len(lectura)

		// ERROR
		if err != nil {
			MsgError(err)
		}

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_name := ""
			s_part_type := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				// Obtengo el nombre de la particion
				// [:] -> Convierte el arreglo de bytes a cadena
				s_part_name = string(master_boot_record.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")

				/* Pendiente */
				// Verifico si ya existe una particion con ese nombre
				if s_part_name == nombre {

				}

				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				// Verifico si de tipo extendida
				if s_part_type == "E" {
					extendida = i
				}
			}

			// Lo busco en las extendidas
			if extendida != -1 {
				// Obtengo el inicio de la particion
				s_part_start := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
				// Le quito los caracteres null
				s_part_start = strings.Trim(s_part_start, "\x00")
				i_part_start, err := strconv.Atoi(s_part_start)

				// ERROR
				if err != nil {
					MsgError(err)
					fin_archivo = true
				}

				// Obtengo el espacio de la partcion
				s_part_size := string(master_boot_record.Mbr_partition[extendida].Part_size[:])
				// Le quito los caracteres null
				s_part_size = strings.Trim(s_part_size, "\x00")
				i_part_size, err := strconv.Atoi(s_part_size)

				// ERROR
				if err != nil {
					MsgError(err)
					fin_archivo = true
				}

				// Calculo del tamano de struct en bytes
				ebr2 := struct_a_bytes(ebr_empty)
				sstruct := len(ebr2)

				// Lectrura de conjunto de bytes desde el inicio de la particion
				for !fin_archivo {
					// Lectrura de conjunto de bytes en archivo binario
					lectura := make([]byte, sstruct)
					n_leidos, err := f.ReadAt(lectura, int64(sstruct*cont+i_part_start))

					// ERROR
					if err != nil {
						MsgError(err)
						fin_archivo = true
					}

					// Posicion actual en el archivo
					// Seek -> Cambia la posicion del puntero de lectura/escritura
					// Seek(offset int64, whence int) (int64, error)
					// whence -> 0: desde el inicio, 1: desde la posicion actual, 2: desde el final
					// os.SEEK_CUR -> Desde la posicion actual
					pos_actual, err := f.Seek(0, os.SEEK_CUR)

					// ERROR
					if err != nil {
						MsgError(err)
						fin_archivo = true
					}

					// Si no lee nada y ya se paso del tamaño de la particion
					if n_leidos == 0 && pos_actual < int64(i_part_start+i_part_size) {
						fin_archivo = true
						break
					}

					// Conversion de bytes a struct
					extended_boot_record := bytes_a_struct_ebr(lectura)
					sstruct = len(lectura)

					if err != nil {
						MsgError(err)
					}

					if extended_boot_record.Part_size == empty {
						fin_archivo = true
					} else {

						Salida_comando += " Nombre: " + string(extended_boot_record.Part_name[:]) + "\n"

						// Antes de comparar limpio la cadena
						s_part_name = string(extended_boot_record.Part_name[:])
						s_part_name = strings.Trim(s_part_name, "\x00")

						// Verifico si ya existe una particion con ese nombre
						if s_part_name == nombre {
							f.Close()
							return true
						}

						// Obtengo el espacio utilizado
						s_part_next := string(extended_boot_record.Part_next[:])
						// Le quito los caracteres null
						s_part_next = strings.Trim(s_part_next, "\x00")
						i_part_next, err := strconv.Atoi(s_part_next)

						// ERROR
						if err != nil {
							MsgError(err)
						}

						// Si ya termino
						if i_part_next != -1 {
							f.Close()
							return false
						}
					}
					cont++
				}
			}
		}
	}
	f.Close()
	return false
}

/*
func newParticion() *Partition {
	return &Partition{
		Part_status: '0',
		Part_type:   'P',
		Part_fit:    'F',
		Part_start:  -1,
		Part_size:   0,
		Part_name:   [16]byte{},
	}
}

func BuscarParticiones(mbr Mbr, name string, path string) *Partition {
	particiones := mbr.Mbr_partition

	ext := false
	extended := newParticion()
	for i := 0; i < len(particiones); i++ {
		particion := particiones[i]
		if strings.TrimSpace(string(particion.Part_name[:])) != "" {
			nombre := ""
			for j := 0; j < len(particion.Part_name); j++ {
				if particion.Part_name[j] != 0 {
					nombre += string(particion.Part_name[j])
				}
			}
			if Comparar(nombre, name) {
				return &particion
			} else if strings.TrimSpace(string(particion.Part_type[:])) == "E" || strings.TrimSpace(string(particion.Part_type[:])) == "e" {
				ext = true
				extended = &particion
			}

		}
	}

	if ext {
		ebrs := GetLogicas(extended, path)
		for i := 0; i < len(ebrs); i++ {
			ebr := ebrs[i]
			if ebr.Part_status == '1' {
				nombre := ""
				for j := 0; j < len(ebr.Part_name); j++ {
					if ebr.Part_name[j] != 0 {
						nombre += string(ebr.Part_name[j])
					}
				}
				if Comparar(nombre, name) {
					tmp := Structs.NewParticion()
					tmp.Part_status = '1'
					tmp.Part_type = 'L'
					tmp.Part_fit = ebr.Part_fit
					tmp.Part_start = ebr.Part_start
					tmp.Part_size = ebr.Part_size
					copy(tmp.Part_name[:], ebr.Part_name[:])
					return &tmp
				}
			}
		}
	}
	return nil
}
*/
