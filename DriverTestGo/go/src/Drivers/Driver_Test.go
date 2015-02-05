package main  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "elev.h"
#include "main.c"
*/
import "C"


func main(){
	C.main();
}
