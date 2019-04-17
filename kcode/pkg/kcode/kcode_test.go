package kcode

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// In Golang, tests are based on the testing package:
// https://blog.alexellis.io/golang-writing-unit-tests/
// Characteristics
// 1. The first and only parameter needs to be t *testing.T
// 2. It begins with the word Test followed by a word or phrase
// starting with a capital letter
// 3. (usually the method under test i.e. TestValidateClient)
// 4. Calls t.Error or t.Fail to indicate a failure
// 5. Must be saved in a file named packagename_test.go
// 6. Launched using:
// $ go test kcode
// For more verbose output:
// $ go test -v kcode
// To run benchmarks:
// $ go test -bench=. kcode
// Note: If you have code and tests in the same folder
// then you cannot execute your program with go run *.go.

func validateBlock(filename string, t *testing.T, expecting []string, verbose bool) {
	log.SetOutput(ioutil.Discard)
	xml := GetXML(filename)
	bcount := BlockCount(xml)
	filedata := ReadFile(filename)
	flags := KCodeFlags{Spells: false, Blocks: true, Parts: false, Scene: false}
	_, blocks, _, _ := ProcessKcodeFileString(filedata, flags, verbose)
	fmt.Println(blocks)
	for _, block := range expecting {
		assert.Contains(t, blocks, block)
	}
	assert.Equal(t, len(blocks), bcount)
}

func validateSpell(filename string, t *testing.T, expecting []string, verbose bool) {
	log.SetOutput(ioutil.Discard)
	fmt.Printf("--- %s ---\n", filename)
	xml := GetXML(filename)
	scount := SpellCount(xml)
	filedata := ReadFile(filename)
	flags := KCodeFlags{Spells: true, Blocks: false, Parts: false, Scene: false}
	spells, _, _, _ := ProcessKcodeFileString(filedata, flags, verbose)
	fmt.Println(spells)
	for _, spell := range expecting {
		assert.Contains(t, spells, spell)
	}
	assert.Equal(t, len(spells), scount)
}

func validateParts(filename string, t *testing.T, expecting []string, verbose bool) {
	log.SetOutput(ioutil.Discard)
	fmt.Printf("--- %s ---\n", filename)
	filedata := ReadFile(filename)
	scount := PartCount(filedata)
	fmt.Printf("PartCount: %d\n", scount)
	flags := KCodeFlags{Spells: false, Blocks: false, Parts: true, Scene: false}
	_, _, parts2, _ := ProcessKcodeFileString(filedata, flags, verbose)
	fmt.Printf("Parts: %s\n", parts2)
	for _, part := range expecting {
		assert.Contains(t, parts2, part)
	}
	assert.Equal(t, len(parts2), scount)
}

func processDirectory(dir string, flags KCodeFlags, verbose bool) {
	files := ListFilesInDirectory(dir)
	for _, f := range files {
		filename := dir + "/" + f.Name()
		flags := KCodeFlags{Spells: true, Blocks: true, Parts: false, Scene: false}
		_, _, _, _ = ProcessKcodeFile(filename, flags, verbose)
	}
}

func validateDirectory(dir string, t *testing.T, flags KCodeFlags, verbose bool) {
	files := ListFilesInDirectory(dir)
	for _, f := range files {
		filename := dir + "/" + f.Name()
		if filepath.Ext(filename) == ".kcode" {
			expecting := []string{}
			if flags.Spells {
				validateSpell(filename, t, expecting, verbose)
			}
			if flags.Blocks {
				validateBlock(filename, t, expecting, verbose)
			}
		}
	}
}

func TestIndividualBlocks(t *testing.T) {
	// test block handling
	verbose := false
	// colorvaria
	expecting := []string{"events_onFlick", "objects_setColor"}
	validateBlock("challenges/001_colovaria.kcode", t, expecting, verbose)
	// bus challenge
	expecting = []string{"events_onGesture", "position_set"}
	validateBlock("challenges/057_bus.kcode", t, expecting, verbose)
	// pumpkins challenge
	expecting = []string{"events_onGesture", "objects_scale"}
	validateBlock("challenges/1022_pumpkins.kcode", t, expecting, verbose)
	// big beans challenge
	expecting = []string{"events_whileFlick", "objects_scale", "wand_vibrate"}
	//speaker#speaker_play speaker#speaker_sample
	validateBlock("challenges/020_big_beans.kcode", t, expecting, verbose)
}

func TestIndividualSpells(t *testing.T) {
	// test spell handling
	verbose := false
	// colorvaria
	expecting := []string{}
	validateSpell("challenges/001_colovaria.kcode", t, expecting, verbose)
	// bus challenge
	expecting = []string{"wingardiumLeviosa", "reparo", "accio"}
	validateSpell("challenges/057_bus.kcode", t, expecting, verbose)
	// pumpkins challenge
	expecting = []string{"engorgio", "reducio"}
	validateSpell("challenges/1022_pumpkins.kcode", t, expecting, verbose)
	// big beans challenge
	expecting = []string{}
	validateSpell("challenges/020_big_beans.kcode", t, expecting, verbose)
}

func TestIndividualParts(t *testing.T) {
	// test part handling
	verbose := false
	// colorvaria
	expecting := []string{}
	validateParts("challenges/001_colovaria.kcode", t, expecting, verbose)
	// big beans challenge
	expecting = []string{"speaker"}
	validateParts("challenges/020_big_beans.kcode", t, expecting, verbose)
}

func TestValidateAllChallenges(t *testing.T) {
	verbose := false
	flags := KCodeFlags{Spells: true, Blocks: true, Parts: true, Scene: true}
	validateDirectory("challenges", t, flags, verbose)
}

func BenchmarkAllChallengeSpells(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	verbose := false
	flags := KCodeFlags{Spells: true, Blocks: false, Parts: false, Scene: false}
	for i := 0; i < b.N; i++ {
		processDirectory("challenges", flags, verbose)
	}
}

func BenchmarkAllChallengeBlocks(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	verbose := false
	flags := KCodeFlags{Spells: false, Blocks: true, Parts: false, Scene: false}
	for i := 0; i < b.N; i++ {
		processDirectory("challenges", flags, verbose)
	}
}
