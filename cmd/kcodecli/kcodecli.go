package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	kcode "github.com/malminhas/kcode/kcode/pkg/kcode"
	//kcode "../../pkg/kcode"
	docopt "github.com/docopt/docopt-go"
)

// Build instructions:
// $ go build -ldflags="-s -w" ..\..\cmd\kcodecli.go -o kcodecli.exe

func dumpBlocks(blocks []string) {
	for i, block := range blocks {
		fmt.Printf("block %d: %s\n", i+1, block)
	}
}

func dumpSpells(spells []string) {
	for i, spell := range spells {
		fmt.Printf("spell %d: %s\n", i+1, spell)
	}
}

func dumpParts(parts []string) {
	for i, part := range parts {
		fmt.Printf("part %d: %s\n", i+1, part)
	}
}

func processDirectory(dir string, flags kcode.KCodeFlags, verbose bool) {
	files := kcode.ListFilesInDirectory(dir)
	for _, f := range files {
		fname := dir + "/" + f.Name()
		spells, blocks, parts, scene := kcode.ProcessKcodeFile(fname, flags, verbose)
		if flags.Spells {
			dumpSpells(spells)
		}
		if flags.Blocks {
			dumpBlocks(blocks)
		}
		if flags.Parts {
			fmt.Printf("%s", parts)
		}
		if flags.Scene {
			fmt.Printf("%s", scene)
		}
	}
}

func validateDirectory(dir string, verbose bool) {
	files := kcode.ListFilesInDirectory(dir)
	for _, f := range files {
		fname := dir + "/" + f.Name()
		if filepath.Ext(fname) == ".kcode" {
			expectedSpells, foundSpells, expectedBlocks, foundBlocks, expectedParts, foundParts, expectedScene, foundScene, valid := kcode.ValidateFile(fname, verbose)
			if valid {
				fmt.Printf("SUCCEEDED in validating '%s'.\nExpected and found %d spells and %d blocks\n", fname, foundSpells, foundBlocks)
			} else {
				fmt.Printf("FAILED to validate '%s'.\nExpected %d spells and found %d.\nExpected %d blocks and found %d\n",
					fname, expectedSpells, foundSpells, expectedBlocks, foundBlocks)
				fmt.Printf("Expected %d parts and found %d.\nExpected %d len scene and found %d\n",
					expectedParts, foundParts, expectedScene, foundScene)
			}
		}
	}
}

// ---------- opts handling  ----------

func procOpts(opts *docopt.Opts) {
	//opts, _ := docopt.ParseDoc(usage)
	//fmt.Println(typeof(opts))
	//fmt.Println(opts)
	var conf struct {
		Blocks   bool   `docopt:"blocks"`
		Spells   bool   `docopt:"spells"`
		Parts    bool   `docopt:"parts"`
		Scene    bool   `docopt:"scene"`
		Validate bool   `docopt:"validate"`
		File     string `docopt:"<file>"`
		Verbose  bool   `docopt:"--verbose"`
	}
	opts.Bind(&conf)

	fname := conf.File
	verbose := conf.Verbose

	if len(fname) > 0 {
		kcode.InitLogging(verbose)
		start := time.Now()
		if conf.Blocks {
			// Note there is no ternary operator in Go:
			// https://stackoverflow.com/questions/19979178/what-is-the-idiomatic-go-equivalent-of-cs-ternary-operator
			flags := kcode.KCodeFlags{Spells: false, Blocks: true, Parts: false, Scene: false}
			if kcode.IsDirectory(fname) { // The file passed in is a directory
				fmt.Println(fmt.Sprintf("Extracting 'blocks' from all .kcode files in target directory '%s'...", fname))
				processDirectory(fname, flags, verbose)
			} else {
				fmt.Println(fmt.Sprintf("Extracting 'blocks' in target .kcode file '%s'...", fname))
				_, blocks, _, _ := kcode.ProcessKcodeFile(fname, flags, verbose)
				dumpBlocks(blocks)
			}
		} else if conf.Spells {
			flags := kcode.KCodeFlags{Spells: true, Blocks: false, Parts: false, Scene: false}
			if kcode.IsDirectory(fname) { // The file passed in is a directory
				fmt.Println(fmt.Sprintf("Seeking 'spells' in target directory '%s'...", fname))
				processDirectory(fname, flags, verbose)
			} else {
				fmt.Println(fmt.Sprintf("Seeking 'spells' in .kcode file '%s'...", fname))
				spells, _, _, _ := kcode.ProcessKcodeFile(fname, flags, verbose)
				dumpSpells(spells)
			}
		} else if conf.Parts {
			flags := kcode.KCodeFlags{Spells: false, Blocks: false, Parts: true, Scene: false}
			if kcode.IsDirectory(fname) { // The file passed in is a directory
				fmt.Println(fmt.Sprintf("Seeking 'parts' in target directory '%s'...", fname))
				processDirectory(fname, flags, verbose)
			} else {
				fmt.Println(fmt.Sprintf("Seeking 'parts' in .kcode file '%s'...", fname))
				_, _, parts, _ := kcode.ProcessKcodeFile(fname, flags, verbose)
				dumpParts(parts)
			}
		} else if conf.Scene {
			flags := kcode.KCodeFlags{Spells: false, Blocks: false, Parts: false, Scene: true}
			if kcode.IsDirectory(fname) { // The file passed in is a directory
				fmt.Println(fmt.Sprintf("Seeking 'scene' in target directory '%s'...", fname))
				processDirectory(fname, flags, verbose)
			} else {
				fmt.Println(fmt.Sprintf("Seeking 'scene' in .kcode file '%s'...", fname))
				_, _, _, scene := kcode.ProcessKcodeFile(fname, flags, verbose)
				fmt.Printf("%s\n", scene)
			}
		} else if conf.Validate {
			if kcode.IsDirectory(fname) { // The file passed in is a directory
				fmt.Println(fmt.Sprintf("Validating .kcode files in target directory '%s'...", fname))
				validateDirectory(fname, verbose)
			} else {
				fmt.Println(fmt.Sprintf("Validating .kcode file '%s'...", fname))
				expectedSpells, foundSpells, expectedBlocks, foundBlocks,
					expectedParts, foundParts, expectedScene, foundScene, valid := kcode.ValidateFile(fname, verbose)
				if valid {
					fmt.Printf("SUCCEEDED in validating '%s'.\nExpected %d spells and found %d\n", fname, expectedSpells, foundSpells)
					fmt.Printf("Expected %d blocks and found %d\n", expectedBlocks, foundBlocks)
					fmt.Printf("Expected %d parts and found %d\nExpected scence len %d and found len %d\n",
						expectedParts, foundParts, expectedScene, foundScene)
				} else {
					fmt.Printf("FAILED to validate '%s'.\nExpected %d spells and found %d.\nExpected %d blocks and found %d\n",
						fname, expectedSpells, foundSpells, expectedBlocks, foundBlocks)
					fmt.Printf("Expected %d parts and found %d.\nExpected scence len %d and found len %d\n",
						expectedParts, foundParts, expectedScene, foundScene)
				}
			}
		} else {
			fmt.Println(opts)
		}
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Printf("========== FINISHED ===========\nElapsed time = %s", elapsed)
	} else {
		fmt.Println("No file passed in")
	}
}

// ---------- main  ----------

func main() {
	usage := `KCode parser
------------
Usage:
  kcodecli blocks <file> [--verbose]
  kcodecli spells <file> [--verbose]
  kcodecli parts <file> [--verbose] 
  kcodecli scene <file> [--verbose] 
  kcodecli validate <file> [--verbose]
  kcodecli --help | --version

Options:
  --help    	Show this screen.
  --version     Show version.

Examples:
  1. Find spells in 'mycreation.kcode':
  kcodecli spells mycreation.kcode
  2. Find spells in 'spelldir' directory:
  kcodecli spells -d spelldir
`
	// Process error handling
	version := "1.0"
	opts, _ := docopt.ParseArgs(usage, os.Args[1:], version)
	procOpts(&opts)
}
