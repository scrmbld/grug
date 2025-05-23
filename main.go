package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type model struct {
	Header string
}

func getFiles(basePath string) ([]string, error) {
	var result []string
	walkFunc := func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			result = append(result, s)
		}
		return nil
	}
	err := filepath.WalkDir(basePath, walkFunc)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	// command line flags
	var (
		inputDir   string
		outputDir  string
		includeDir string
		verbose    bool
	)
	flag.StringVar(&inputDir, "i", "", "-i [input dir]")
	flag.StringVar(&outputDir, "o", "", "-o [output dir]")
	flag.StringVar(&includeDir, "include", "", "-include [include_dir]")
	flag.BoolVar(&verbose, "v", false, "-v")
	flag.Parse()

	if len(inputDir) == 0 {
		fmt.Println("Must specify -i\nusage: -i [input dir]")
		os.Exit(1)
	}
	if len(outputDir) == 0 {
		fmt.Println("Must specify -o\nuseage: -o [ouput dir]")
		os.Exit(1)
	}

	inputDir = filepath.Clean(inputDir)
	outputDir = filepath.Clean(outputDir)
	includeDir = filepath.Clean(includeDir)

	// get all of the files we need to load into our templater
	inputfiles, err := getFiles(inputDir)
	if err != nil {
		fmt.Println(err)
		fmt.Println("quitting...")
		os.Exit(1)
	}
	if verbose {
		fmt.Println("input files:", inputfiles)
	}
	var includefiles []string
	if len(includeDir) != 0 {
		includefiles, err = getFiles(includeDir)
		if err != nil {
			fmt.Println(err)
			fmt.Println("quitting...")
			os.Exit(1)
		}
	}
	if verbose {
		fmt.Println("include files:", includefiles)
	}

	// load all of the templates
	ts, err := template.ParseFiles(append(inputfiles, includefiles...)...)
	if err != nil {
		fmt.Println(err)
		fmt.Println("quitting...")
		os.Exit(1)
	}

	// render the templates and write them to the outputDir
	for _, file := range inputfiles {
		templateName := filepath.Base(file)
		outputPath := filepath.Join(outputDir, strings.Replace(file, inputDir, "", 1))
		// make sure the output file's directory exists
		os.MkdirAll(filepath.Dir(outputPath), 0o777)
		// open the file
		outFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			fmt.Println(err)
			fmt.Println("quitting...")
			os.Exit(1)
		}
		err = ts.ExecuteTemplate(outFile, templateName, nil)
		if err != nil {
			fmt.Println(err)
			fmt.Println("quitting...")
			os.Exit(1)
		}
		if verbose {
			fmt.Printf("%s written successfully\n", outputPath)
		}
	}
}
