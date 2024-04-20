package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	pagegen "pagegen/core"

	"github.com/alexflint/go-arg"
)

func main() {
	cwd, err := os.Getwd()
	pagegen.CheckErr(err)

	cwd = strings.ReplaceAll(cwd, "\\", "/")

	var args struct {
		Title                string `arg:"positional,required"`
		Template             string
		Rel_Dir              string
		Output_Dir           string
		Default_Variable_Val string
		Content_File         string
		Verbose              bool `arg:"-v" help:"Enable verbose logging"`
		DoubleVerbose        bool `arg:"--vv" help:"Enable verbose and debug logging"`
	}
	args.Rel_Dir = cwd
	args.Template = cwd + "/template.html"
	args.Output_Dir = cwd + "/output"
	args.Verbose = false
	args.DoubleVerbose = false

	arg.MustParse(&args)

	if args.Content_File == "" {
		args.Content_File = cwd + "/" + args.Title + "-content.yml"
	}

	if args.DoubleVerbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else if args.Verbose {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	} else {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}

	slog.Info("Got args:")
	slog.Info("Title: ", "title", args.Title)
	slog.Info("Template: ", "template", args.Template)
	slog.Info("Rel_Dir: ", "reldir", args.Rel_Dir)
	slog.Info("OutputDir: ", "outdir", args.Output_Dir)
	slog.Info("ContentFile: ", "contents", args.Content_File)
	slog.Info("DefaultVariableVal: ", "default var val", args.Default_Variable_Val)

	fmt.Println("Processing...")
	slog.Info("Reading contents of ", "contentfile", args.Content_File)
	contents := pagegen.ContentReader(args.Content_File)
	slog.Info("Done reading content!")

	slog.Info("Parsing variables...")
	contents = pagegen.VariablesParser(contents)
	slog.Info("Done parsing variables!")

	slog.Info("Parsing template...")
	generatedPage := pagegen.TemplateParser(contents, args.Template)
	slog.Info("Done parsing template!")
	fmt.Println(generatedPage)

	fmt.Println("Writing to output file...")

	//CreateDir(args.Output_Dir)
	filename := args.Output_Dir + "/" + args.Title + ".html"
	slog.Info("Writing content to", "outfile", filename)
	//WriteStringToFile(filename, generatedPage)

	fmt.Println("Output written to ", filename)
}
