package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	cwd = strings.ReplaceAll(cwd, "\\", "/")

	var args struct {
		Title                string `arg:"positional,required"`
		Template             string
		Rel_Dir              string
		Output_Dir           string
		Default_Variable_Val string
	}
	args.Rel_Dir = cwd
	args.Template = cwd + "/template.html"
	args.Output_Dir = cwd + "/output"

	arg.MustParse(&args)

	fmt.Println("Got args:")
	fmt.Println("Title: ", args.Title)
	fmt.Println("Template: ", args.Template)
	fmt.Println("Rel_Dir: ", args.Rel_Dir)
	fmt.Println("OutputDir: ", args.Output_Dir)
	fmt.Println("DefaultVariableVal: ", args.Default_Variable_Val)
}
