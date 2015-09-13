package bmpwriter

import (
	"os"
)

type Rgb struct {
	r uint8
	g uint8
	b uint8
}

const (
	B           = byte(66)
	M           = byte(77)
	RES         = uint32(2835)
	HEADER_SIZE = byte(54)
)

func to_uint8_array(n uint32) []uint8 {
	a := []uint8{
		uint8(n),
		uint8(n >> 8),
		uint8(n >> 16),
		uint8(n >> 24),
	}
	return a
}

func set_four_bytes(data []uint8, n uint32, pos int) {
	for i, element := range to_uint8_array(n) {
		data[i+pos] = element
	}
}

func get_bmp_file_header(width uint32, height uint32) []byte {
	bmpfileheader := make([]uint8, 54)
	bmpfileheader[0] = B
	bmpfileheader[1] = M
	filesize := uint32(HEADER_SIZE) + (uint32(3) * width * height)
	set_four_bytes(bmpfileheader, filesize, 2)
	bmpfileheader[10] = HEADER_SIZE
	bmpfileheader[14] = 40
	set_four_bytes(bmpfileheader, width, 18)
	set_four_bytes(bmpfileheader, height, 22)
	bmpfileheader[26] = 1
	bmpfileheader[28] = 24
	set_four_bytes(bmpfileheader, RES, 38)
	set_four_bytes(bmpfileheader, RES, 42)
	return bmpfileheader
}


func write_reverse_with_padding(file *os.File, image []byte, w uint32, h uint32) {
	row_padding := make([]uint8, 3*w%4)
	for y := h; y > 0; y-- {
		row_start_index := (y - 1) * w * 3
		row_end_index := row_start_index + (w * 3)
		file.Write(image[row_start_index:row_end_index])
		if len(row_padding) > 0 {
			file.Write(row_padding)
		}
	}
}

func Write_bmp(path string, image []byte, w uint32, h uint32) {
	bmpfileheader := get_bmp_file_header(w, h)

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(bmpfileheader)
	write_reverse_with_padding(file, image, w, h)
}