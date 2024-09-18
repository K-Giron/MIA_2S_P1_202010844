package Comandos

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Variables globales
var graphDot = ""

// Muestra los datos en el disco
func Rep(commandArray []string) {
	Salida_comando += "MENSAJE: El comando REP aqui inicia" + "\n"
	// Variables para los valores de los parametros
	val_name := ""
	val_path := ""
	val_id := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_name := false
	band_path := false
	band_id := false
	band_ruta := false
	band_error := false

	// Limpio la variable global
	graphDot = ""

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
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
		/* PARAMETRO OBLIGATORIO -> ID */
		case strings.Contains(data, "id="):
			// Valido si el parametro ya fue ingresado
			if band_id {
				fmt.Println("Invalido: El parametro -id ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_id = true

			// Reemplaza comillas
			val_id = val_data
		/* PARAMETRO OBLIGATORIO -> RUTA */
		case strings.Contains(data, "ruta="):
			if band_ruta {
				fmt.Println("Invalido: El parametro -ruta ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_ruta = true
		/* PARAMETRO NO VALIDO */
		default:
			fmt.Println("Invalido: El parametro " + data + " no es valido")
		}
	}

	if !band_error {
		if band_path {
			if band_name {
				if band_id {
					var aux *ParticionMontada

					for i := 0; i < len(ParticionesMontadas); i++ {
						if ParticionesMontadas[i].id == val_id {
							aux = &ParticionesMontadas[i]
							break
						}
					}

					if aux != nil {
						if val_name == "disk" {
							graficar_disk(aux.Direccion, val_path)
						}
					} else {
						fmt.Println("Invalido: No se encontro la particion montada con el id:", val_id)
					}
				} else {
					fmt.Println("Invalido: El parametro -id no fue ingresado...")
				}
			} else {
				fmt.Println("Invalido: El parametro -name no fue ingresado...")
			}
		} else {
			fmt.Println("Invalido: El parametro -path no fue ingresado...")
		}
	}
	Salida_comando += "MENSAJE: El comando REP aqui finaliza" + "\n"
}

// Reporte DISK
func graficar_disk(direccion string, destino string) {
	mbr_empty := Mbr{}
	var empty [100]byte

	// Abro el archivo para lectura con opcion a modificar
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// Calculo del tamaño de struct en bytes
	mbr2 := struct_a_bytes(mbr_empty)
	sstruct := len(mbr2)

	// Lectrura del archivo binario desde el inicio
	lectura := make([]byte, sstruct)
	f.Seek(0, io.SeekStart)
	f.Read(lectura)

	// Conversion de bytes a struct
	master_boot_record := bytes_a_struct_mbr(lectura)

	if master_boot_record.Mbr_tamano != empty {
		if err == nil {
			graphDot += "digraph G{\n\n"
			graphDot += "  tbl [\n    shape=box\n    label=<\n"
			graphDot += "     <table border='0' cellborder='2' width='600' height='150' color='dodgerblue1'>\n"
			graphDot += "     <tr>\n"
			graphDot += "     <td height='150' width='110'> MBR </td>\n"

			// Obtengo el espacio utilizado
			s_mbr_tamano := string(master_boot_record.Mbr_tamano[:])
			// Le quito los caracteres null
			s_mbr_tamano = strings.Trim(s_mbr_tamano, "\x00")
			i_mbr_tamano, _ := strconv.Atoi(s_mbr_tamano)
			total := i_mbr_tamano

			var espacioUsado float64
			espacioUsado = 0

			for i := 0; i < 4; i++ {
				// Obtengo el espacio utilizado
				s_part_s := string(master_boot_record.Mbr_partition[i].Part_size[:])
				// Le quito los caracteres null
				s_part_s = strings.Trim(s_part_s, "\x00")
				i_part_s, _ := strconv.Atoi(s_part_s)

				parcial := i_part_s

				// Obtengo el espacio utilizado
				s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
				// Le quito los caracteres null
				s_part_start = strings.Trim(s_part_start, "\x00")

				if s_part_start != "-1" {
					var porcentaje_real float64
					porcentaje_real = (float64(parcial) * 100) / float64(total)
					var porcentaje_aux float64
					porcentaje_aux = (porcentaje_real * 500) / 100

					espacioUsado += porcentaje_real

					// Obtengo el espacio utilizado
					s_part_status := string(master_boot_record.Mbr_partition[i].Part_status[:])
					// Le quito los caracteres null
					s_part_status = strings.Trim(s_part_status, "\x00")

					if s_part_status != "1" {
						// Obtengo el espacio utilizado
						s_part_type := string(master_boot_record.Mbr_partition[i].Part_type[:])
						// Le quito los caracteres null
						s_part_type = strings.Trim(s_part_type, "\x00")

						if s_part_type == "p" {
							graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Primaria <br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"

							if i != 3 {
								// Obtengo el espacio utilizado
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)

								// Obtengo el espacio utilizado
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								p1 := i_part_start + i_part_s

								// Obtengo el espacio utilizado
								s_part_start = string(master_boot_record.Mbr_partition[i+1].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ = strconv.Atoi(s_part_start)

								p2 := i_part_start

								if s_part_start != "-1" {
									if (p2 - p1) != 0 {
										fragmentacion := p2 - p1
										porcentaje_real = float64(fragmentacion) * 100 / float64(total)
										porcentaje_aux = (porcentaje_real * 500) / 100

										graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
									}
								}
							} else {
								// Obtengo el espacio utilizado
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)

								// Obtengo el espacio utilizado
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								p1 := i_part_start + i_part_s

								mbr_empty_byte := struct_a_bytes(mbr_empty)
								mbr_size := total + len(mbr_empty_byte)

								// Si esta libre
								if (mbr_size - p1) != 0 {
									libre := (float64(mbr_size) - float64(p1)) + float64(len(mbr_empty_byte))
									porcentaje_real = (float64(libre) * 100) / float64(total)
									porcentaje_aux = (porcentaje_real * 500) / 100
									graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
								}
							}
						} else {
							// Si es extendida
							graphDot += "     <td  height='200' width='" + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + "'>\n     <table border='0'  height='200' WIDTH='" + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + "' cellborder='1'>\n"
							graphDot += "     <tr>  <td height='60' colspan='15'>Extendida</td>  </tr>\n     <tr>\n"

							// Obtengo el espacio utilizado
							s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
							// Le quito los caracteres null
							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)

							f.Seek(int64(i_part_start), io.SeekStart)

							ebr_empty := Ebr{}

							// Calculo del tamaño de struct en bytes
							ebr2 := struct_a_bytes(ebr_empty)
							sstruct := len(ebr2)

							// Lectrura del archivo binario desde el inicio
							lectura := make([]byte, sstruct)
							f.Read(lectura)

							// Conversion de bytes a struct
							extended_boot_record := bytes_a_struct_ebr(lectura)

							// Obtengo el espacio utilizado
							s_part_size := string(extended_boot_record.Part_size[:])
							// Le quito los caracteres null
							s_part_size = strings.Trim(s_part_size, "\x00")
							i_part_size, _ := strconv.Atoi(s_part_size)

							if i_part_size != 0 {
								// Obtengo el espacio utilizado
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								f.Seek(int64(i_part_start), io.SeekStart)

								band := true

								// Obtengo el espacio utilizado
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)

								// Obtengo el espacio utilizado
								s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ = strconv.Atoi(s_part_start)

								for band {
									// Calculo del tamaño de struct en bytes
									ebr2 := struct_a_bytes(ebr_empty)
									sstruct := len(ebr2)

									// Lectrura del archivo binario desde el inicio
									lectura := make([]byte, sstruct)
									f.Seek(0, io.SeekStart)
									n, _ := f.Read(lectura)

									// Posicion actual en el archivo
									pos_actual, _ := f.Seek(0, io.SeekCurrent)

									if n != 0 && pos_actual < int64(i_part_start)+int64(i_part_s) {
										band = false
										break
									}

									// Obtengo el espacio utilizado
									s_part_s = string(extended_boot_record.Part_size[:])
									// Le quito los caracteres null
									s_part_s = strings.Trim(s_part_s, "\x00")
									i_part_s, _ = strconv.Atoi(s_part_s)

									parcial = i_part_start
									porcentaje_real = float64(parcial) * 100 / float64(total)

									if porcentaje_real != 0 {
										// Obtengo el espacio utilizado
										s_part_status = string(extended_boot_record.Part_status[:])
										// Le quito los caracteres null
										s_part_status = strings.Trim(s_part_status, "\x00")

										if s_part_status != "1" {
											graphDot += "     <td height='140'>EBR</td>\n"
											graphDot += "     <td height='140'>Logica<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
										} else {
											// Espacio no asignado
											graphDot += "      <td height='150'>Libre 1 <br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
										}

										// Obtengo el espacio utilizado
										s_part_next := string(extended_boot_record.Part_next[:])
										// Le quito los caracteres null
										s_part_next = strings.Trim(s_part_next, "\x00")
										i_part_next, _ := strconv.Atoi(s_part_next)

										if i_part_next == -1 {
											// Obtengo el espacio utilizado
											s_part_start := string(extended_boot_record.Part_start[:])
											// Le quito los caracteres null
											s_part_start = strings.Trim(s_part_start, "\x00")
											i_part_start, _ := strconv.Atoi(s_part_start)

											// Obtengo el espacio utilizado
											s_part_size := string(extended_boot_record.Part_size[:])
											// Le quito los caracteres null
											s_part_size = strings.Trim(s_part_size, "\x00")
											i_part_size, _ := strconv.Atoi(s_part_size)

											// Obtengo el espacio utilizado
											s_part_start_mbr := string(master_boot_record.Mbr_partition[i].Part_start[:])
											// Le quito los caracteres null
											s_part_start_mbr = strings.Trim(s_part_start_mbr, "\x00")
											i_part_start_mbr, _ := strconv.Atoi(s_part_start_mbr)

											// Obtengo el espacio utilizado
											s_part_s_mbr := string(master_boot_record.Mbr_partition[i].Part_size[:])
											// Le quito los caracteres null
											s_part_s_mbr = strings.Trim(s_part_s_mbr, "\x00")
											i_part_s_mbr, _ := strconv.Atoi(s_part_s_mbr)

											parcial = (i_part_start_mbr + i_part_s_mbr) - (i_part_size + i_part_start)
											porcentaje_real = (float64(parcial) * 100) / float64(total)

											if porcentaje_real != 0 {
												graphDot += "     <td height='150'>Libre 2<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
											}
											break

										} else {
											// Obtengo el espacio utilizado
											s_part_next := string(extended_boot_record.Part_next[:])
											// Le quito los caracteres null
											s_part_next = strings.Trim(s_part_next, "\x00")
											i_part_next, _ := strconv.Atoi(s_part_next)

											f.Seek(int64(i_part_next), io.SeekStart)
										}
									}

								}
							} else {
								graphDot += "     <td height='140'> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
							}
							graphDot += "     </tr>\n     </table>\n     </td>\n"

							// Verifica que no haya espacio fragemntado
							if i != 3 {
								// Obtengo el espacio utilizado
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)

								// Obtengo el espacio utilizado
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								p1 := i_part_start + i_part_s

								// Obtengo el espacio utilizado
								s_part_start = string(master_boot_record.Mbr_partition[i+1].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ = strconv.Atoi(s_part_start)

								p2 := i_part_start

								if s_part_start != "-1" {
									if (p2 - p1) != 0 {
										fragmentacion := p2 - p1
										porcentaje_real = float64(fragmentacion) * 100 / float64(total)
										porcentaje_aux = (porcentaje_real * 500) / 100

										graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
									}
								}
							} else {
								// Obtengo el espacio utilizado
								s_part_s = string(master_boot_record.Mbr_partition[i].Part_size[:])
								// Le quito los caracteres null
								s_part_s = strings.Trim(s_part_s, "\x00")
								i_part_s, _ = strconv.Atoi(s_part_s)

								// Obtengo el espacio utilizado
								s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
								// Le quito los caracteres null
								s_part_start = strings.Trim(s_part_start, "\x00")
								i_part_start, _ := strconv.Atoi(s_part_start)

								p1 := i_part_start + i_part_s

								mbr_empty_byte := struct_a_bytes(mbr_empty)
								mbr_size := total + len(mbr_empty_byte)

								// Si esta libre
								if (mbr_size - p1) != 0 {
									libre := (float64(mbr_size) - float64(p1)) + float64(len(mbr_empty_byte))
									porcentaje_real = (float64(libre) * 100) / float64(total)
									porcentaje_aux = porcentaje_real * 500 / 100
									graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
								}
							}
						}
					} else {
						graphDot += "     <td height='200' width='" + strconv.FormatFloat(porcentaje_aux, 'g', 3, 64) + "'>Libre<br/> " + strconv.FormatFloat(porcentaje_real, 'g', 3, 64) + " por ciento del Disco </td>\n"
					}
				}
			}

			graphDot += "     </tr> \n     </table>        \n>];\n\n}"

			// Escribe el contenido en un archivo
			err := ioutil.WriteFile("reporte.dot", []byte(graphDot), 0644)
			Salida_comando += "MENSAJE: Generando reporte..." + "\n"
			fmt.Println("[MENSAJE] Generando reporte...")
			if err != nil {
				fmt.Println("[ERROR] Error al escribir en el archivo:", err)
				return
			}

			// Ejecutar el comando para generar la imagen
			cmd := exec.Command("dot", "-Tpng", "reporte.dot", "-o", destino)
			if err := cmd.Run(); err != nil {
				fmt.Println("[ERROR] Error al generar la imagen:", err)
				return
			}
		} else {
			fmt.Println("[Invalido] Error al abrir el archivo:", err)
		}
	} else {
		fmt.Println("[Invalido] El disco no fue encontrado...")
	}
}
