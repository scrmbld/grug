# Grug

Grug is a simple command line template renderer for Go templates. Weighing in at less than 200 lines of code, Grug enables you to replace your bloated static site generator with your own build system, such as a GNU makefile or simple shell script.

## Installation

TODO  
Right now there isn't a real system to install grug. Just clone the repository, run `go build .`, and copy the resulting `grug` executable to your PATH.

## Usage

Grug takes an input directory and an output directory. It gets all of the files in the input directory, loads them as templates, and renders each of the file-level templates. The rendered templates are written to the same relative path in the output directory. Optionally, grug can be provided an include directory for templates that are not their own pages.

Take the following file tree:

```
├── dist
└── src
    ├── include
    │   └── base.html
    └── pages
        ├── articles
        │   └── article1.html
        └── index.html
```

After running

`grug -i src/pages -o dist -include src/include`

We get the resulting file tree:

```
├── dist
│   ├── articles
│   │   └── article1.html
│   └── index.html
└── src
    ├── include
    │   └── base.html
    └── pages
        ├── articles
        │   └── article1.html
        └── index.html

```

`article1.html` and `index.html` have been rendered and written in `dist`, maintaining the directory structure they had in the input directory (`src/pages`).

### Flags

`-i [input dir]` -- input directory \[REQUIRED\]  
`-o [output dir]` -- output directory \[REQUIRED\]  
`-include [include dir]` -- include directory. Default `[input dir]/_include`.  
`-v` -- verbose. Prints out all loaded input and include files and all of the written output files.  

### Examples

Grug is best enjoyed within an existing build system, such as `make` or a shell script. This allows other build steps to be run, such as copying over static content or running the tailwind preprocessor.

TODO: Add some examples

## Builtins

Grug makes additional functions available from templates on top of what Go already builds in: `mkSlice` and `mkMap`. This provides everything you need to avoid writing the same HTML multiple times.

#### mkSlice

Returns a slice (the Go equivalent of a Python list or JS array) containing the given arguments. All of the arguments must be of the same type.

```
<!-- this sets $slice equal to ["hello", "world"] -->
{{$slice := mkSlice "hello" "world"}}
```

#### mkMap

Takes a string containing a JSON object and turns it into a map.

```
<!-- this sets $map equal to map[field1:val1]-->
{{$map := mkMap `{"field1": "val1"}`}}
```
