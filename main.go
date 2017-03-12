// Echo4 выводит аргументы командной строки
package main

import (
	"flag"
	"fmt"
	// "strings"
)

var n = flag.Bool("n", false, "пропуск символа новой строки")

func main() {
	flag.Parse()
	fmt.Println(n)
	// fmt.Print(strings.Join(flag.Args(), *sep))
	if *n {
		fmt.Println("flag is not empty")
	}
}
