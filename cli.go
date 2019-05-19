package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var copy bool
var keep bool
var original bool
var pattern string
var recursive bool

func init() {
	flag.BoolVar(&original, "o", false, "keep original titles")

	flag.StringVar(&pattern, "p", "%d - %s", "destination pattern")

	flag.BoolVar(&recursive, "r", false, "also rename files in subdirectories")

	//flag.Usage = func() {
	//	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\tsds [flags] [-p pattern] <path>", os.Args[0])
	//	flag.PrintDefaults()
	//}
}

func input(prompt string) string {
	fmt.Print(prompt)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(line)
}
