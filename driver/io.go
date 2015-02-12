package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

func io_init() int{
	return int(C.io_init())
}

func io_set_bit(int channel){
	C.io_set_bit(channel)
}

func io_clear_bit(int channel){
	C.io_clear_bit(channel)
}

func io_write_analog(int channel, int value){
	C.io_write_analog(channel, value)
}

func io_read_bit(int channel) int{
	return int(C.io_read_bit(channel))
}

func io_read_analog(int channel) int{
	return int(C.io_read_analog(channel))
}