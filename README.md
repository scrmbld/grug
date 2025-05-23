# Grug

Grug is a simple command line template renderer for Go templates. Weighing in at around 100 lines of code, Grug enables you to replace your bloated static site generator with your own build system, such as a GNU makefile or simple shell script.

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
`grug -i src/input -o src/output -include dist`

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

`article1.html` and `index.html` have been rendered and written in the output directory, maintaining the directory structure they had in the input directory.

### Flags

`-i` -- input directory \[REQUIRED\]
`-o` -- output directory \[REQUIRED\]
`-include` -- include directory
`-v` -- verbose. Prints out all loaded input and include files and all of the written output files.

### Examples

Grug is best enjoyed within an existing build system, such as `make` or a shell script. This allows other build steps to be run, such as copying over static content or running the tailwind preprocessor.

TODO: Add some examples
