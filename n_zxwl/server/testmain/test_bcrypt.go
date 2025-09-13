package main

import (
	"fmt"
	"mymod/utils"
)

func main() {
	var psw string
	psw, _ = utils.HashPassword("rootroot")
	fmt.Println(psw)
}
