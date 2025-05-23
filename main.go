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

func subdirMatch(base string, match string) bool {
	basePath := strings.Split(base, "/")
	matchPath := strings.Split(match, "/")
	if len(matchPath) > len(basePath) {
		return false
	}
	for i, _ := range matchPath {
		if basePath[i] != matchPath[i] {
			return false
		}
	}
	return true
}

func anySubdirMatch(base string, matches []string) bool {
	for _, v := range matches {
		if subdirMatch(base, v) {
			return true
		}
	}
	return false
}

func getFiles(basePath string, ignoreDirs []string) ([]string, error) {
	var result []string
	walkFunc := func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && !anySubdirMatch(s, ignoreDirs) {
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

	if inputDir == "" {
		fmt.Println("Must specify -i\nusage: -i [input dir]")
		os.Exit(1)
	}
	if outputDir == "" {
		fmt.Println("Must specify -o\nuseage: -o [ouput dir]")
		os.Exit(1)
	}

	// get all of the files we need to load into our templater
	inputDir = filepath.Clean(inputDir)
	inputfiles, err := getFiles(inputDir, []string{includeDir})
	if err != nil {
		fmt.Println(err)
		fmt.Println("quitting...")
		os.Exit(1)
	}
	if verbose {
		fmt.Println("input files:", inputfiles)
	}

	var includefiles []string
	if includeDir == "" {
		includeDir = filepath.Join(inputDir, "_include")
	}
	includeDir = filepath.Clean(includeDir)
	includefiles, err = getFiles(includeDir, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println("quitting...")
		os.Exit(1)
	}
	if verbose {
		fmt.Println("include files:", includefiles)
	}

	// also run clean on the output path
	outputDir = filepath.Clean(outputDir)

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
