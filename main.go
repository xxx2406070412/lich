// lich project main.go
package main

import (
	"fmt"
)

func init() {
	return
}

func main() {

	fmt.Println("lich start")
	pstMux.Listen(":5100")
	fmt.Println("http server start")
}
