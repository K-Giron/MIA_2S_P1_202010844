package Comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func ValidarDatosMkfs(commandArray []string) {
	Salida_comando += "MENSAJE:  MKFS" + "\n"
	id := ""
	tipo := "full"
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]
		switch {
		case strings.Contains(data, "id="):
			id = val_data
		case strings.Contains(data, "type="):
			// pasar a minusculas
			val_data = strings.ToLower(val_data)
			// validar que sea full
			if val_data == "full" {
				tipo = val_data
			} else {
				Error("MKFS", "Tipo de formateo no valido")
			}
		default:
			Error("MKFS", "Parametro no valido")
		}

	}
	if id == "" {
		Error("MKFS", "EL comando requiere el parámetro id obligatoriamente")
		return
	}
	mkfs(id, tipo)
}

func mkfs(id string, tipo string) {
	fmt.Println("MENSAJE: Formateando partición con id: " + id + " y tipo: " + tipo)

	p := ""

	particion := GetMount("MKFS", id, &p)

	// Convertir el array de bytes a int64
	partSizeStr := strings.Trim(string(particion.Part_size[:]), "\x00")
	partSizeInt, err := strconv.ParseInt(partSizeStr, 10, 64)
	if err != nil {
		fmt.Println("Error al convertir Part_size a int64")
		return
	}

	// Calcular el número de inodos
	n := math.Floor(float64(partSizeInt-int64(unsafe.Sizeof(SuperBloque{}))) / float64(4+unsafe.Sizeof(Inodos{})+3*unsafe.Sizeof(BloquesArchivos{})))

	spr := NewSuperBloque()
	spr.S_magic = 0xEF53
	spr.S_inode_size = int64(unsafe.Sizeof(Inodos{}))
	spr.S_block_size = int64(unsafe.Sizeof(BloquesCarpetas{}))
	spr.S_inodes_count = int64(n)
	spr.S_free_inodes_count = int64(n)
	spr.S_blocks_count = int64(3 * n)
	spr.S_free_blocks_count = int64(3 * n)
	fecha := time.Now().String()
	copy(spr.S_mtime[:], fecha)
	spr.S_mnt_count = spr.S_mnt_count + 1
	spr.S_filesystem_type = 2
	ext2(spr, particion, int64(n), p)

}

func ext2(spr SuperBloque, p Partition, n int64, path string) {
	// Convertir el array de bytes a int64
	partStartStr := strings.Trim(string(p.Part_start[:]), "\x00")
	partStartInt, err := strconv.ParseInt(partStartStr, 10, 64)
	if err != nil {
		fmt.Println("Error al convertir Part_start a int64")
		return
	}

	spr.S_bm_inode_start = partStartInt + int64(unsafe.Sizeof(SuperBloque{}))
	spr.S_bm_block_start = spr.S_bm_inode_start + n
	spr.S_inode_start = spr.S_bm_block_start + (3 * n)
	spr.S_block_start = spr.S_bm_inode_start + (n * int64(unsafe.Sizeof(Inodos{})))

	file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco.")
		return
	}

	partStartStr = strings.Trim(string(p.Part_start[:]), "\x00")
	partStartInt, err = strconv.ParseInt(partStartStr, 10, 64)
	if err != nil {
		fmt.Println("Error al convertir Part_start a int64")
		return
	}
	file.Seek(partStartInt, 0)
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, spr)
	EscribirBytes(file, binario2.Bytes())

	zero := '0'
	file.Seek(spr.S_bm_inode_start, 0)
	for i := 0; i < int(n); i++ {
		var binarioZero bytes.Buffer
		binary.Write(&binarioZero, binary.BigEndian, zero)
		EscribirBytes(file, binarioZero.Bytes())
	}

	file.Seek(spr.S_bm_block_start, 0)
	for i := 0; i < 3*int(n); i++ {
		var binarioZero bytes.Buffer
		binary.Write(&binarioZero, binary.BigEndian, zero)
		EscribirBytes(file, binarioZero.Bytes())
	}

	inode := NewInodos()
	//INICIALIZANDO EL INODO
	inode.I_uid = -1
	inode.I_gid = -1
	inode.I_size = -1
	for i := 0; i < len(inode.I_block); i++ {
		inode.I_block[i] = -1
	}
	inode.I_type = -1
	inode.I_perm = -1

	file.Seek(spr.S_inode_start, 0)
	for i := 0; i < int(n); i++ {
		var binarioInodos bytes.Buffer
		binary.Write(&binarioInodos, binary.BigEndian, inode)
		EscribirBytes(file, binarioInodos.Bytes())
	}

	folder := NewBloquesCarpetas()

	for i := 0; i < len(folder.B_content); i++ {
		folder.B_content[i].B_inodo = -1
	}

	file.Seek(spr.S_block_start, 0)
	for i := 0; i < int(n); i++ {
		var binarioFolder bytes.Buffer
		binary.Write(&binarioFolder, binary.BigEndian, folder)
		EscribirBytes(file, binarioFolder.Bytes())
	}
	file.Close()

	recuperado := NewSuperBloque()
	//ABRIR ARCHIVO
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)

	file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco.")
		return
	}

	partStartStr = strings.Trim(string(p.Part_start[:]), "\x00")
	partStartInt, err = strconv.ParseInt(partStartStr, 10, 64)
	if err != nil {
		fmt.Println("Error al convertir Part_start a int64")
		return
	}
	file.Seek(partStartInt, 0)
	data := LeerBytes(file, int(unsafe.Sizeof(SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &recuperado)
	if err_ != nil {
		Error("FDSIK", "Error al leer el archivo")
		return
	}
	file.Close()

	inode.I_uid = 1
	inode.I_gid = 1
	inode.I_size = 0
	fecha := time.Now().String()
	copy(inode.I_atime[:], fecha)
	copy(inode.I_ctime[:], fecha)
	copy(inode.I_mtime[:], fecha)
	inode.I_type = 0
	inode.I_perm = 664
	inode.I_block[0] = 0

	fb := NewBloquesCarpetas()
	copy(fb.B_content[0].B_name[:], ".")
	fb.B_content[0].B_inodo = 0
	copy(fb.B_content[1].B_name[:], "..")
	fb.B_content[1].B_inodo = 0
	copy(fb.B_content[2].B_name[:], "users.txt")
	fb.B_content[2].B_inodo = 1

	dataArchivo := "1,G,root\n1,U,root,root,123\n"
	inodetmp := NewInodos()
	inodetmp.I_uid = 1
	inodetmp.I_gid = 1
	inodetmp.I_size = int64(unsafe.Sizeof(dataArchivo) + unsafe.Sizeof(BloquesCarpetas{}))

	copy(inodetmp.I_atime[:], fecha)
	copy(inodetmp.I_ctime[:], fecha)
	copy(inodetmp.I_mtime[:], fecha)
	inodetmp.I_type = 1
	inodetmp.I_perm = 664
	inodetmp.I_block[0] = 1

	inode.I_size = inodetmp.I_size + int64(unsafe.Sizeof(BloquesCarpetas{})) + int64(unsafe.Sizeof(Inodos{}))

	var fileb BloquesArchivos
	copy(fileb.B_content[:], dataArchivo)

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco.")
		return
	}
	file.Seek(spr.S_bm_inode_start, 0)
	caracter := '1'

	var bin1 bytes.Buffer
	binary.Write(&bin1, binary.BigEndian, caracter)
	EscribirBytes(file, bin1.Bytes())
	EscribirBytes(file, bin1.Bytes())

	file.Seek(spr.S_bm_block_start, 0)
	var bin2 bytes.Buffer
	binary.Write(&bin2, binary.BigEndian, caracter)
	EscribirBytes(file, bin2.Bytes())
	EscribirBytes(file, bin1.Bytes())

	file.Seek(spr.S_inode_start, 0)

	var bin3 bytes.Buffer
	binary.Write(&bin3, binary.BigEndian, inode)
	EscribirBytes(file, bin3.Bytes())

	file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Inodos{})), 0)
	var bin4 bytes.Buffer
	binary.Write(&bin4, binary.BigEndian, inodetmp)
	EscribirBytes(file, bin4.Bytes())

	file.Seek(spr.S_block_start, 0)

	var bin5 bytes.Buffer
	binary.Write(&bin5, binary.BigEndian, fb)
	EscribirBytes(file, bin5.Bytes())

	//fmt.Println(spr.S_block_start + int64(unsafe.Sizeof(Structs.BloquesCarpetas{})))

	file.Seek(spr.S_block_start+int64(unsafe.Sizeof(BloquesCarpetas{})), 0)
	var bin6 bytes.Buffer
	binary.Write(&bin6, binary.BigEndian, fileb)
	EscribirBytes(file, bin6.Bytes())

	file.Close()

	nombreParticion := ""
	for i := 0; i < len(p.Part_name); i++ {
		if p.Part_name[i] != 0 {
			nombreParticion += string(p.Part_name[i])
		}
	}
	Salida_comando += "Se formateo la particion: " + nombreParticion + "\n"
}
