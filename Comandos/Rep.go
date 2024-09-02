package Comandos

import (
	"fmt"
	"io"
	"os"
)

// Muestra los datos en el disco
func Rep() {
	fin_archivo := false
	var empty [100]byte
	mbr_empty := mbr{}
	cont := 0

	fmt.Println("* Reporte de MKDISK: *")

	// Apertura de archivo
	disco, err := os.OpenFile("C:\\Users\\kevin\\Desktop\\P1_MIA\\Discos\\A.mia", os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		MsgError(err)
	}

	// Calculo del tamano de struct en bytes
	mbr2 := struct_a_bytes(mbr_empty)
	sstruct := len(mbr2)

	// RECORRE CADA STRUCT DEL ARCHIVO
	for !fin_archivo {
		// Lectrura de conjunto de bytes en archivo binario
		lectura := make([]byte, sstruct)
		_, err = disco.ReadAt(lectura, int64(sstruct*cont))

		// ERROR
		if err != nil && err != io.EOF {
			MsgError(err)
		}

		// Conversion de bytes a struct
		mbr := bytes_a_struct_mbr(lectura)
		sstruct = len(lectura)

		// ERROR
		if err != nil {
			MsgError(err)
		}

		if mbr.Mbr_tamano == empty {
			fin_archivo = true
		} else {
			fmt.Print("Tamaño: ")
			fmt.Print(string(mbr.Mbr_tamano[:]))
			fmt.Println(" bytes ")
			fmt.Print("Fecha: ")
			fmt.Println(string(mbr.Mbr_fecha_creacion[:]))
			fmt.Print("Signature: ")
			fmt.Println(string(mbr.Mbr_dsk_signature[:]))
		}

		cont++
	}
	disco.Close()

}
