package Comandos

import (
	"fmt"
	"strings"
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
	/*
		p := ""
		particion := GetMount("MKFS", id, &p)
		n := math.Floor(float64(particion.Part_size-int64(unsafe.Sizeof(Structs.SuperBloque{}))) / float64(4+unsafe.Sizeof(Structs.Inodos{})+3*unsafe.Sizeof(Structs.BloquesArchivos{})))

		spr := Structs.NewSuperBloque()
		spr.S_magic = 0xEF53
		spr.S_inode_size = int64(unsafe.Sizeof(Structs.Inodos{}))
		spr.S_block_size = int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))
		spr.S_inodes_count = int64(n)
		spr.S_free_inodes_count = int64(n)
		spr.S_blocks_count = int64(3 * n)
		spr.S_free_blocks_count = int64(3 * n)
		fecha := time.Now().String()
		copy(spr.S_mtime[:], fecha)
		spr.S_mnt_count = spr.S_mnt_count + 1
		spr.S_filesystem_type = 2
		ext2(spr, particion, int64(n), p)*/

}
