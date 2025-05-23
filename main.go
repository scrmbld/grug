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

func getFiles(basePath string, ignorePaths []string) ([]string, error) {
	var result []string
	walkFunc := func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// skip directories that are in the ignorePaths -- this means we don't go over files inside them
		if anySubdirMatch(s, ignorePaths) {
			return fs.SkipDir
		}
		// ignore paths might not be a directory, so we have to check here too
		if !d.IsDir() && !anySubdirMatch(s, ignorePaths) {
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
		inputDir        string
		outputDir       string
		includeDir      string
		noIncludePassed bool
		verbose         bool
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
	if includeDir == "" {
		noIncludePassed = true
	}

	inputDir = filepath.Clean(inputDir)
	outputDir = filepath.Clean(outputDir)
	includeDir = filepath.Clean(includeDir)

	// get all of the files we need to load into our templater
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
	if noIncludePassed {
		includeDir = filepath.Join(inputDir, "_include")
	}
	includefiles, err = getFiles(includeDir, nil)
	if err != nil {
		// check for include dir does not exist error
		// we can't just quit on this error because the include dir is optional
		if err, ok := err.(*fs.PathError); ok && err.Op == "lstat" && err.Path == includeDir {
			fmt.Printf("WARNING: Include directory %s does not exist\n", includeDir)
		} else {
			fmt.Println(err)
			fmt.Println("quitting...")
			os.Exit(1)
		}
	} else if verbose {
		// the "include dir does not exist" implies what this will be (namely, empty)
		fmt.Println("include files:", includefiles)
	}

	// load all of the templates
	ts := template.New("demo").Funcs(template.FuncMap(grugFuncMap))
	template.Must(ts.ParseFiles(append(inputfiles, includefiles...)...))
	if err != nil {
		fmt.Println(err)
		fmt.Println("quitting...")
		os.Exit(1)
	}
	ts = ts.Funcs(grugFuncMap)

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
