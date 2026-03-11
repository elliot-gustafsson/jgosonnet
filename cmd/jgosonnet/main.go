package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/elliot-gustafsson/jgosonnet"
)

var (
	help  bool
	jpath string
)

func main() {
	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	err := run()
	if err != nil {
		slog.Error(err.Error())
	}
}

func run() error {
	flag.BoolVar(&help, "help", false, "This message")
	flag.BoolVar(&help, "h", false, "This message")
	flag.StringVar(&jpath, "jpath", "", "Specify an additional library search dir")

	flag.Parse()

	if help {
		flag.PrintDefaults()
		return nil
	}

	args := flag.Args()

	if len(args) == 0 {
		return fmt.Errorf("no arguments passed")
	}

	if len(args) > 1 {
		return fmt.Errorf("mulitple arguments passed")
	}

	interpreter := jgosonnet.NewEvaluator()

	if jpath != "" {
		fi, err := os.Stat(jpath)
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			return fmt.Errorf("provided jpath is not a directory")
		}
		interpreter.JPaths([]string{jpath})
	}

	_, err := interpreter.Evaluate(args[0])
	if err != nil {
		return err
	}

	// err = json.NewEncoder(os.Stdout).Encode(val)
	// if err != nil {
	// 	return err
	// }

	return nil
}
