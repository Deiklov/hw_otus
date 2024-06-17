package main

import (
	"log"
	"os"
)

func main() {

	if len(os.Args) < 5 {
		log.Fatal("should pass all arguments ")
	}

	envs, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("bad args %s", err)
	}
	retCode := RunCmd(os.Args[2:], envs)
	os.Exit(retCode)
}
