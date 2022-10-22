package main

import (
	"fmt"
	"log"
	"os"

	conv "kukulkan.converter/pkg/toMP3"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("ERROR: need some args like -f 'filepath'")
	}
	var (
		cn = conv.New(args[len(args)-1]) // get last element as path
		ok = cn.Check()
	)

	if ok {
		log.Println(cn.Run())
		ok = <-cn.Done
	}
	fmt.Println("- - -")
	for _, v := range cn.Info {
		fmt.Println(v)
	}
	if !ok {
		log.Fatalf("ERROR: cannot convert, %v", cn.Err)
	}

	fmt.Println("DONE.")
}
