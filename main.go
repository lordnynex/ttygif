package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var version = "0.0.3"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	input := flag.String("in", "ttyrecord", "input ttyrec file")
	output := flag.String("out", "tty.gif", "output gif file")
	speed := flag.Float64("s", 1.0, "play speed")
	noloop := flag.Bool("noloop", false, "play only once")
	help := flag.Bool("help", false, "usage")
	v := flag.Bool("v", false, "version")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *v {
		fmt.Println(version)
		os.Exit(0)
	}

	generator := NewGifGenerator()
	generator.Speed = *speed
	generator.NoLoop = *noloop
	err := generator.Generate(*input, *output)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	absPath, err := filepath.Abs(*output)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s created!\n", absPath)
}
