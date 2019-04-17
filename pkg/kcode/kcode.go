package kcode

// kcode.go
// --------
// Description:
// Module for parsing Kano .kcode files and retrieving spell and block details from them.
//
// API:
// BlockCount(xml []byte) int
// SpellCount(xml []byte) int
// PartCount(kcode []string) int
// SceneCount(kcode []byte) int
// ValidateFile(filename string, verbose bool) (int,int,int,int,bool)
// ValidateString(kcode []byte, verbose bool) (int,int,int,int,bool)
// Extract(pblocks *[]string, jstr []byte, flags KCodeFlags, verbose bool) ([]string)
// ExtractBlocks(jstr []byte, verbose bool) []string
// ExtractSpells(jstr []byte, verbose bool) []string
// ExtractXML(jsdata []byte) ([]byte, error)
// ExtractParts(jsdata []byte) ([]string, error)
// ExtractScene(jsdata []byte) (string, error)
// DumpXML(kcode []byte, prettyPrint bool, verbose bool)
// IsDirectory(filename string) bool
// ExistsFile(filename string) bool
// ReadFile(filename string) (data []byte)
// ListFilesInDirectory(dirname string) []os.FileInfo
// GetParts(filename string) (parts []string)
// GetXML(filename string) (kcode []byte)
// ProcessKcodeFile(filename string, flags KCodeFlags, verbose bool) (spells []string, blocks []string, parts []string)
// ProcessKcodeFileString(data []byte, flags KCodeFlags, verbose bool) (spells []string, blocks []string, parts []string, scene string)
// InitLogging(verbose bool)
//

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"

	xml2json "github.com/basgys/goxml2json"
	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"
	"github.com/yosssi/gohtml"
)

var (
	errInvalidJSON = errors.New("invalid JSON")
)

type KCodeFlags struct {
	Blocks bool `json:"blocks"`
	Spells bool `json:"spells"`
	Scene  bool `json:"scene"`
	Parts  bool `json:"parts"`
}

// Struct for Kano Code .kcode files
// A .kcode file contains: i) XML, ii) Scenes, iii) Parts
// XML
// ---
// XML can be located differently depending on
// Hence there are two kinds of KCode structures:
// 1. v1 for PK and MSK creations where XML kcode is here:
//	xml = kcode.get('code').get('snapshot').get('blocks')
// 2. v2 for Wand creations where XML kcode is here:
// Scenes
// ------
// String
// "scene":"honeydukesbeans"}

type KCodeLegacy struct {
	Code  []interface{} `json:"source"`
	Parts []interface{} `json:"parts"`
	Scene string        `json:"scene"`
}

// KCode struct ...
type KCode struct {
	Source string      `json:"source"`
	Parts  []KCodePart `json:"parts"`
	Scene  string      `json:"scene"`
}

// Parts
// -----
// Parts are a list of dictionaries of parts as follows:
// "parts":[{"id":"speaker","name":"Speaker","type":"speaker","tagName":"kano-part-speaker",
// "userStyle":{},"userProperties":{},"nonvolatileProperties":[],"position":{"x":54.444437662760414,"y":40.833333333333336},
// "partType":"hardware","supportedHardware":[]}]
type KCodePart struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Tag      string `json:"tagName"`
	PartType string `json:"partType"`
}

// ---------- Utils  ----------

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func check(err string, e error) {
	if e != nil {
		msg := fmt.Sprintf("%s ('%s')\n", err, e)
		fmt.Printf(msg)
		log.Panic(err)
		panic(e)
	}
}

// ---------- Block and Spell handling  ----------

// BlockCount counts the number of `<block` found in input kcode XML
func BlockCount(xml []byte) int {
	r := regexp.MustCompile(`(<block)`)
	res := r.FindAllStringSubmatch(string(xml), -1)
	return len(res)
}

// SpellCount counts the number of `events_onGesture` found in input kcode XML
func SpellCount(xml []byte) int {
	r := regexp.MustCompile(`(events_onGesture)`)
	res := r.FindAllStringSubmatch(string(xml), -1)
	return len(res)
}

// PartCount counts the number of `partType` found in input
func PartCount(kcode []byte) int {
	r := regexp.MustCompile(`(partType)`)
	res := r.FindAllStringSubmatch(string(kcode), -1)
	return len(res)
}

// SceneCount counts value of `scene` found in input
func SceneCount(kcode []byte) int {
	// just return the length of the array slice
	r := regexp.MustCompile(`(scene: )`)
	res := r.FindAllStringSubmatch(string(kcode), -1)
	return len(res)
}

// ValidateFile validates that the number of blocks and spells found matches expected counts
func ValidateFile(filename string, verbose bool) (int, int, int, int, int, int, int, int, bool) {
	log.Info(fmt.Sprintf("------- Reading file '%s' ------\n", filename))
	filedata := ReadFile(filename)
	xml := GetXML(filename)
	return ValidateString(filedata, xml, verbose)
}

// ValidateString validates that the number of blocks and spells found matches expected counts
func ValidateString(filedata []byte, kcode []byte, verbose bool) (int, int, int, int, int, int, int, int, bool) {
	nexpectedblocks := BlockCount(kcode)
	nexpectedspells := SpellCount(kcode)
	nexpectedparts := PartCount(filedata)
	nexpectedscene := SceneCount(filedata)
	//fmt.Printf("%d %d %d %d\n", nexpectedblocks, nexpectedspells, nexpectedparts, nexpectedscene)
	flags := KCodeFlags{Spells: true, Blocks: true, Parts: true, Scene: true}
	spells, blocks, parts, scene := ProcessKcodeFileString(filedata, flags, verbose)
	nfoundspells, nfoundblocks, nfoundparts, nfoundscene := len(spells), len(blocks), len(parts), len(scene)
	isvalid := (nfoundblocks == nexpectedblocks) && (nfoundspells == nexpectedspells) && (nfoundparts == nexpectedparts)
	return nexpectedspells, nfoundspells, nexpectedblocks, nfoundblocks, nexpectedparts, nfoundparts, nexpectedscene, nfoundscene, isvalid
}

func dumpString(str string, verbose bool) {
	log.Info(str)
	if verbose {
		fmt.Print(str)
	}
}

func processBlock(pblocks *[]string, pvalue *[]byte, flags KCodeFlags, datatype jsonparser.ValueType, verbose bool) {
	//fmt.Println(string(*pvalue))
	t, _, _, err := jsonparser.Get(*pvalue, "-type")
	check("extractType", err)
	id, _, _, err := jsonparser.Get(*pvalue, "-id")
	check("extractId", err)

	parseSpell := func(pblocks *[]string, value []byte) {
		spell, _, _, err := jsonparser.Get(value, "field", "#content")
		check("events_onGesture", err)
		// Print out information about the spell we just found
		dumpString(fmt.Sprintf("SPELL: type=%s, spell=%s, id=%s\n", string(t), spell, string(id)), verbose)
		// Update *pblocks with spell
		*pblocks = append(*pblocks, string(spell))
	}

	parseBlock := func(t []byte, pblocks *[]string, value []byte) {
		log.Info(fmt.Sprintf("Found block '%s': ", string(t)))
		x, _, _, _ := jsonparser.Get(value, "-x")
		y, _, _, _ := jsonparser.Get(value, "-y")
		statement, stdatatype, _, _ := jsonparser.Get(value, "statement")
		next, ndatatype, _, _ := jsonparser.Get(value, "next")
		v, vdatatype, _, _ := jsonparser.Get(value, "value")
		if flags.Blocks {
			// Print out information about the block we just found
			dumpString(fmt.Sprintf("BLOCK: type=%s, id=%s, x=%s, y=%s, statement=%d, next=%d, val=%d, bdtype=%s\n",
				string(t), string(id), string(x), string(y), len(statement), len(next), len(v), datatype), verbose)
			// Update *pblocks
			*pblocks = append(*pblocks, string(t))
		}
		// We now need to recurse on any "statement", "next" or "value" keys found in this block
		if len(statement) > 0 {
			//log.Info(Sprintf("Found statement block:\n%s\n", string(statement)))
			processValue(pblocks, &statement, flags, stdatatype, verbose)
		}
		if len(next) > 0 {
			//log.Info(fmt.Sprintf("Found next block:\n%s\n", string(next)))
			processValue(pblocks, &next, flags, ndatatype, verbose)
		}
		if len(v) > 0 {
			//log.Info(fmt.Sprintf("Found value block:\n%s\n", string(v)))
			processValue(pblocks, &v, flags, vdatatype, verbose)
		}
	}

	switch string(t) {
	case "events_onGesture":
		// We found a spell!
		// Not going to do anything with it for now, just fall through
		// https://stackoverflow.com/questions/45268681/golangs-fallthrough-seems-unexpected
		if flags.Spells {
			parseSpell(pblocks, *pvalue)
		}
		// Need to fallthrough to process any "statement" and "next" found within spell block
		fallthrough
	default:
		// Not every block has all these attributes so can't enforce strict checking.
		// Commenting out requires dealing with err per:
		// https://stackoverflow.com/questions/21743841/how-to-avoid-annoying-error-declared-and-not-used
		parseBlock(t, pblocks, *pvalue)
	}
}

func processValue(pblocks *[]string, pvalue *[]byte, flags KCodeFlags, datatype jsonparser.ValueType, verbose bool) {
	// See here for how to do a function within a function in Go:
	// https://stackoverflow.com/questions/21961615/why-doesnt-go-allow-nested-function-declarations-functions-inside-functions
	parseVal := func(value []byte) {
		//fmt.Printf("Parsing object value %s\n", string(*pvalue))
		t, _, _, _ := jsonparser.Get(value, "-type")
		id, _, _, _ := jsonparser.Get(value, "-id")
		name, _, _, _ := jsonparser.Get(value, "-name")
		block, bdatatype, _, _ := jsonparser.Get(value, "block")
		statement, stdatatype, _, _ := jsonparser.Get(value, "statement")
		next, ndatatype, _, _ := jsonparser.Get(value, "next")
		v, vdatatype, _, _ := jsonparser.Get(value, "value")
		shadow, shdatatype, _, _ := jsonparser.Get(value, "shadow")
		dumpString(fmt.Sprintf("VALUE: name=%s,type=%s,id=%s,block=%d,val=%d,shadow=%d,statement=%d,next=%d,dtype=%s\n",
			string(name), string(t), string(id), len(block), len(v), len(shadow), len(statement), len(next), datatype), verbose)
		// We now need to recurse on any "value" keys found in this block
		if len(block) > 0 {
			//log.Info(fmt.Sprintf("Found BLOCK in VALUE: %s\n",string(block)))
			processBlock(pblocks, &block, flags, bdatatype, verbose)
		}
		if len(shadow) > 0 {
			//log.Info(fmt.Sprintf("Found value block:\n%s\n", string(valueBlock)))
			processValue(pblocks, &shadow, flags, shdatatype, verbose)
		}
		if len(v) > 0 {
			//log.Info(fmt.Sprintf("Found value block:\n%s\n", string(valueBlock)))
			processValue(pblocks, &v, flags, vdatatype, verbose)
		}
		if len(statement) > 0 {
			//log.Info(fmt.Sprintf("Found statement block:\n%s\n", string(statementBlock)))
			processValue(pblocks, &statement, flags, stdatatype, verbose)
		}
		if len(next) > 0 {
			//log.Info(fmt.Sprintf("Found NEXT in VALUE:\n%s\n", string(next)))
			processValue(pblocks, &next, flags, ndatatype, verbose)
		}
	}

	switch datatype {
	case jsonparser.Object:
		parseVal(*pvalue)
	case jsonparser.Array:
		jsonparser.ArrayEach(*pvalue, func(value []byte, datatype jsonparser.ValueType, offset int, err error) {
			parseVal(value)
		})
	default:
		log.Panic(fmt.Sprintf("processValue - unknown datatype=%s", datatype))
	}
}

// Extract is the main logic function to process input code
func Extract(pblocks *[]string, jstr []byte, flags KCodeFlags, verbose bool) []string {
	//log.Info(fmt.Sprintf("%s,%s\n",string(jstr),typeof(jstr)))
	// The top level input for kcode could contain:
	// a. Single Object block
	// b. Array of Object blocks
	// You determine which one by checking datatype return value in call to Get
	block, datatype, _, err := jsonparser.Get(jstr, "xml", "block")
	check("block", err)
	switch datatype {
	case jsonparser.Object:
		log.Info("---- extract: Top level single block ----")
		processBlock(pblocks, &block, flags, datatype, verbose)
	case jsonparser.Array:
		log.Info("---- extract: Top level array of blocks ----")
		// You can use `ArrayEach` helper to iterate items in block [item1, item2 .... itemN]
		jsonparser.ArrayEach(jstr, func(block []byte, dataType jsonparser.ValueType, offset int, err error) {
			processBlock(pblocks, &block, flags, datatype, verbose)
		}, "xml", "block")
	default:
		log.Panic(fmt.Sprintf("extract - unknown datatype=%s", datatype))
	}
	return *pblocks
}

// ExtractBlocks extracts the blocks from input
func ExtractBlocks(jstr []byte, verbose bool) []string {
	blocks := make([]string, 0)
	flags := KCodeFlags{Spells: false, Blocks: true, Parts: false, Scene: false}
	return Extract(&blocks, jstr, flags, verbose)
}

// ExtractSpells extracts the spells from input
func ExtractSpells(jstr []byte, verbose bool) []string {
	spells := make([]string, 0)
	flags := KCodeFlags{Spells: true, Blocks: false, Parts: false, Scene: false}
	return Extract(&spells, jstr, flags, verbose)
}

// ExtractXML extracts the kcode from input
func ExtractXML(jsdata []byte) ([]byte, error) {
	var kc KCode
	if err := json.Unmarshal(jsdata, &kc); err != nil {
		return nil, errInvalidJSON
	}
	// Extract and return []byte "source" from .kcode file which is well formed JSON per KCode struct
	kcode := []byte(kc.Source)
	return kcode, nil
}

// ExtractParts extracts the parts from input
func ExtractParts(jsdata []byte) ([]string, error) {
	var kc KCode
	if err := json.Unmarshal(jsdata, &kc); err != nil {
		return nil, errInvalidJSON
	}
	// Extract and return KCodeParts from .kcode file which is well formed JSON per KCode struct
	parts := kc.Parts
	//fmt.Printf("%s\n", parts)
	lparts := make([]string, 0)
	for _, part := range parts {
		//fmt.Printf("%s,%s,%s", part.Id, part.Name, part.Tag)
		lparts = append(lparts, string(part.Id))
	}
	return lparts, nil
}

// ExtractScene extracts the scene from input
func ExtractScene(jsdata []byte) (string, error) {
	var kc KCode
	if err := json.Unmarshal(jsdata, &kc); err != nil {
		return "", errInvalidJSON
	}
	// Extract and return sring "source" from .kcode file which is well formed JSON per KCode struct
	scene := kc.Scene
	return scene, nil
}

// ---------- XML handling  ----------

// DumpXML prettyprint .kcode XML
func DumpXML(xml []byte, prettyPrint bool, verbose bool) {
	if prettyPrint {
		log.Info(gohtml.Format(string(xml)))
		if verbose {
			fmt.Println(gohtml.Format(string(xml)))
		}
	} else {
		log.Info(string(xml))
		if verbose {
			fmt.Println(string(xml))
		}
	}
}

// ---------- File handling  ----------

// IsDirectory check whether the file is a directory
func IsDirectory(filename string) bool {
	fi, err := os.Stat(filename)
	check("isDirectory", err)
	directory := false
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// do directory stuff
		log.Info("directory")
		directory = true
	case mode.IsRegular():
		// do file stuff
		log.Info("file")
	default:
		log.Error("Unknown")
	}
	return directory
}

// ExistsFile check whether the file exists
func ExistsFile(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// ReadFile open and read file and return contents as a []byte
func ReadFile(filename string) (data []byte) {
	data, err := ioutil.ReadFile(filename)
	check("readFile", err)
	return data
}

// ListFilesInDirectory list files in directory
func ListFilesInDirectory(dirname string) []os.FileInfo {
	files, err := ioutil.ReadDir(dirname)
	check("listdir", err)
	return files
}

// GetParts extracts parts as array slice from file
func GetParts(filename string) (parts []string) {
	data := ReadFile(filename)
	// Convert it to JSON and extract kcode XML
	parts, err := ExtractParts(data)
	check("getparts", err)
	return
}

// GetXML extracts kcode XML as []byte
func GetXML(filename string) (kcode []byte) {
	data := ReadFile(filename)
	// Convert it to JSON and extract kcode XML
	kcode, err := ExtractXML(data)
	check("getkcode", err)
	return
}

// ProcessKcodeFile process .kcode in file and return []string of spells and/or []string of blocks
func ProcessKcodeFile(filename string, flags KCodeFlags, verbose bool) (spells []string, blocks []string, parts []string, scene string) {
	log.Info(fmt.Sprintf("------- Reading file '%s' ------\n", filename))
	data := ReadFile(filename)
	// Convert it to JSON and extract kcode XML
	log.Info("------- Extracting kcode XML ------")
	spells, blocks, parts, scene = ProcessKcodeFileString(data, flags, verbose)
	return
}

// ProcessKcodeFileString process .kcode in string and return []string of spells and/or blocks, []string of parts and string scene.
func ProcessKcodeFileString(data []byte, flags KCodeFlags, verbose bool) (spells []string, blocks []string, parts []string, scene string) {
	xml, _ := ExtractXML(data)
	DumpXML(xml, true, verbose)
	// Parse the XML and convert to JSON
	log.Info("------- Parsing kcode XML to JSON ------")
	// https://stackoverflow.com/questions/44065935/cannot-use-type-byte-as-type-io-reader
	jsn, err := xml2json.Convert(bytes.NewReader(xml))
	check("new reader", err)
	//log.Info(jsn.String())
	log.WithFields(log.Fields{
		"json": jsn.String(),
	}).Info("XML to JSON")

	if flags.Blocks {
		// Extract blocks from JSON
		//fmt.Printf("------- Extracting blocks from .kcode file '%s' ------\n", filename)
		log.Info("------- Extracting blocks ------")
		blocks = ExtractBlocks(jsn.Bytes(), verbose)
		log.Info(fmt.Sprintf("--- Found %d blocks ---\n", len(blocks)))
		for i, block := range blocks {
			log.Info(fmt.Sprintf("%d] %s\n", i+1, block))
		}
	}
	if flags.Spells {
		// Extract spells from JSON
		//fmt.Printf("------- Extracting spells from .kcode file '%s' ------\n", filename)
		log.Info("------- Extracting spells ------")
		spells = ExtractSpells(jsn.Bytes(), verbose)
		log.Info(fmt.Sprintf("--- Found %d spells ---\n", len(spells)))
		for i, spell := range spells {
			log.Info(fmt.Sprintf("%d] %s\n", i+1, spell))
		}
	}
	if flags.Parts {
		parts, _ = ExtractParts(data)
	}
	if flags.Scene {
		scene, _ = ExtractScene(data)
	}

	return
}

// ----------- logging ---------------

// InitLogging initialise logging
func InitLogging(verbose bool) {
	// Logging setup has been changed as follows:
	// 1. verbose=true: write all log messages to logfile.log. Also add extra logging to stdout.
	// 2. verbose=false: disable all log messages and keep stdout minimal.
	// Why?  Logging to stdout isn't very readable using logrus.  More meant for endpoint consumption.
	//if verbose {// log to stdout
	//	fmt.Println("Logging to stdout")
	//	log.SetFormatter(&log.JSONFormatter{})
	//	log.SetOutput(os.Stdout)
	//} else {
	if verbose {
		filename := "logfile.log"
		//fmt.Printf("Logging to file '%s'\n", filename)
		// Delete the logfile if it already exists then recreate
		if ExistsFile(filename) {
			err := os.Remove(filename)
			check("deleteFile", err)
		}
		//f, err := os.OpenFile(filename, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
		f, err := os.OpenFile(filename, os.O_CREATE, 0644)
		//log.SetFormatter(&log.JSONFormatter{PrettyPrint:true})
		log.SetFormatter(&log.JSONFormatter{})
		if err != nil {
			// Cannot open log file. Logging to stderr
			fmt.Println(err)
		} else {
			log.SetOutput(f)
		}
	} else { // https://stackoverflow.com/questions/10571182/go-disable-a-log-logger
		//log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
}
