package Comandos

import "unsafe"

type SuperBloque struct {
	S_filesystem_type   int64
	S_inodes_count      int64
	S_blocks_count      int64
	S_free_blocks_count int64
	S_free_inodes_count int64
	S_mtime             [16]byte
	S_mnt_count         int64
	S_magic             int64
	S_inode_size        int64
	S_block_size        int64
	S_firts_ino         int64
	S_first_blo         int64
	S_bm_inode_start    int64
	S_bm_block_start    int64
	S_inode_start       int64
	S_block_start       int64
}

func NewSuperBloque() SuperBloque {
	var spr SuperBloque
	spr.S_magic = 0xEF53
	spr.S_inode_size = int64(unsafe.Sizeof(Inodos{}))
	spr.S_block_size = int64(unsafe.Sizeof(BloquesCarpetas{}))
	spr.S_firts_ino = 0
	spr.S_first_blo = 0
	return spr
}

type Inodos struct {
	I_uid   int64
	I_gid   int64
	I_size  int64
	I_atime [16]byte
	I_ctime [16]byte
	I_mtime [16]byte
	I_block [16]int64
	I_type  int64
	I_perm  int64
}

func NewInodos() Inodos {
	var inode Inodos
	inode.I_uid = -1
	inode.I_gid = -1
	inode.I_size = -1
	for i := 0; i < 16; i++ {
		inode.I_block[i] = -1
	}
	inode.I_type = -1
	inode.I_perm = -1
	return inode
}

type BloquesArchivos struct {
	B_content [64]byte
}

type BloquesCarpetas struct {
	B_content [4]Content
}

func NewBloquesCarpetas() BloquesCarpetas {
	var bl BloquesCarpetas
	for i := 0; i < len(bl.B_content); i++ {
		bl.B_content[i] = NewContent()
	}
	return bl
}

type Content struct {
	B_name  [12]byte
	B_inodo int64
}

func NewContent() Content {
	var cont Content
	cont.B_inodo = -1
	return cont
}
