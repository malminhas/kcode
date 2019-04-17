# kcode analyser
This repo contains a package to analyse Kano .kcode files.
[Here](https://github.com/malminhas/kcode/blob/master/GolangNotes.md) are some notes on Golang development.

## Overview
The Golang `kcode` package provides a way to extract spells, blocks, scenes and parts from valid `.kcode` files.  This allows developers to efficiently extract information from `.kcode` data.
The `kcode` package comprises three files:
* [kcodecli.go](https://github.com/malminhas/kcode/blob/master/cmd/kcodecli/kcodecli.go): a CLI for working with kcode files.
* [kcode.go](https://github.com/malminhas/kcode/blob/master/pkg/kcode/kcode.go): package used by above
* [kcode_test.go](https://github.com/malminhas/kcode/blob/master/pkg/kcode/kcode_test.go): test code and benchmarking for [kcode.go](https://github.com/malminhas/kcode/blob/master/pkg/kcode/kcode.go) package.

## Installation
To install this package, you must have a local version of Golang installed.  See [here](https://golang.org/doc/install) for details of how to do this on your host PC. You will now need to set up your `GOPATH` environment variable to point to the directory you will be working from.  It must be a fully-formed path and point to the top level `kcode` source code.  Here's how I did it on my Mac with code git cloned to corresponding repo:
```
$ export GOPATH=/Users/malm/Desktop/CODE/kcode
```
You will need to get local copies of all dependencies into `GOPATH`.  This can be done individually as follows after you have setup up `GOPATH`:
```
$ go get -v github.com/buger/jsonparser
$ go get -v github.com/basgys/goxml2json
$ go get -v github.com/sirupsen/logrus
$ go get -v github.com/yosssi/gohtml
$ go get -v github.com/stretchr/testify/assert
$ go get -v github.com/docopt/docopt-go
```
and also `kcode` itself:
```
$ go get -v github.com/malminhas/kcode
```
and build the command line `kcodecli.go` executable as follows from the 
cloned source repo:
```
$ go build -ldflags="-s -w" ./cmd/kcodecli/kcodecli.go 
```
At this point you should be able to test your install has worked by printing the help output from [kcodecli.go](https://github.com/malminhas/kcode/blob/master/cmd/kcodecli/kcodecli.go) as follows:
```
$ kcodecli
KCode parser
------------
Usage:
  kcodecli blocks <file> [--verbose]
  kcodecli spells <file> [--verbose]
  kcodecli validate <file> [--verbose]
  kcodecli --help | --version

Options:
  --help        Show this screen.
  --version     Show version.

Examples:
  1. Find spells in 'mycreation.kcode':
  kcodecli spells mycreation.kcode
  2. Find spells in 'spelldir' directory:
  kcodecli spells -d spelldir
```

## Test
To run the package test code, make sure you are in the `GOPATH` directory then run tests in [kcode_test.go](https://github.com/malminhas/kcode/blob/master/pkg/kcode/kcode_test.go) as follows:
```
$ go test ./kcode/pkg/kcode
ok      kcode   2.085s
```
To run these tests verbosely:
```
$ go test -v ./kcode/pkg/kcode
```
To benchmark spell and block extraction over 10 x 70 challenges:
```
$ go test -bench=. ./kcode/pkg/kcode 
...
      10         176916060 ns/op
PASS
ok      kcode   5.749s
```

## Usage
[kcodecli.go](https://github.com/malminhas/kcode/blob/master/cmd/kcodecli/kcodecli.go) is a command line interface (CLI) for working with kcode files.  Here are a couple of example commands illustrating how to use them to pull blocks out of `.kcode` files either directly by naming the file or indirectly by naming a directory:
```
$ kcodecli blocks src\kcode\challenges\1022_pumpkins.kcode
Seeking 'blocks' in target .kcode file 'src\kcode\challenges\1022_pumpkins.kcode'...
------- Extracting blocks from .kcode file 'src\kcode\challenges\1022_pumpkins.kcode' ------
block 1: events_onGesture
block 2: objects_scale
block 3: objects_scale
block 4: events_onGesture
block 5: objects_scale
========== FINISHED ===========
Elapsed time = 6.0065ms
```
And here's how to pull out spells:
```
$ kcodecli spells src\kcode\challenges\1022_pumpkins.kcode
Seeking 'spells' in .kcode file 'src\kcode\challenges\1022_pumpkins.kcode'...
------- Extracting spells from .kcode file 'src\kcode\challenges\1022_pumpkins.kcode' ------
spell 1: engorgio
spell 2: reducio
========== FINISHED ===========
Elapsed time = 1.9895ms
```
You can set the `--verbose` flag to get more detail on the parsing:
```
$ kcodecli blocks src\kcode\challenges\001_colovaria.kcode --verbose
Seeking 'blocks' in target .kcode file 'src\kcode\challenges\001_colovaria.kcode'...
<xml xmlns="http://www.w3.org/1999/xhtml">
  <variables>
  </variables>
  <block type="events_onFlick" id="8e2Z`.);iVBJ-m3P/4L~" x="342" y="266">
    <field name="TYPE">
      up
    </field>
    <statement name="CALLBACK">
      <block type="objects_setColor" id="z$aXI`j$V}|;^jC5jfQ+">
        <value name="TINT">
          <shadow type="objects_get" id="@nG3U_MJv*?}?g[SyNb5">
            <field name="ID">
              all
            </field>
          </shadow>
        </value>
        <value name="TO COLOR">
          <shadow type="colour_picker" id="zVP6}*(g${T/]Ln:r~c]">
            <field name="COLOUR">
              #FF5723
            </field>
          </shadow>
        </value>
      </block>
    </statement>
  </block>
</xml>
------- Extracting blocks from .kcode file 'src\kcode\challenges\001_colovaria.kcode' ------
type=events_onFlick, id=8e2Z`.);iVBJ-m3P/4L~, x=342, y=266, statement=347, next=0
type=objects_setColor, id=z$aXI`j$V}|;^jC5jfQ+, x=, y=, statement=0, next=0
block 1: events_onFlick
block 2: objects_setColor
========== FINISHED ===========
Elapsed time = 6.8075ms
```

## Validation
Validation checking is now in place to check that the number of blocks and spells found in the `.kcode` XML exactly matches what is pulled out by the XML to JSON parsing logic.  The count of each found is compared with the count of instances of `<block` and `events_onGesture` found in the XML.  Here is how validate a particular `.kcode` file:
```
$ kcodecli validate src\kcode\challenges\001_colovaria.kcode
Validating spells and blocks in .kcode file 'src\kcode\challenges\001_colovaria.kcode'...
------- Extracting spells from .kcode file 'src\kcode\challenges\001_colovaria.kcode' ------
------- Extracting blocks from .kcode file 'src\kcode\challenges\001_colovaria.kcode' ------
SUCCEEDED in validating 'src\kcode\challenges\001_colovaria.kcode'.
Expected and found 0 spells and 2 blocks
========== FINISHED ===========
Elapsed time = 6.0601ms
```
Here's a full verbose examination which allows us to compare blocks in the XML with those that are parsed:
```
$ kcodecli validate d:\CODE\go\wandChallengesKcode\009_accio.kcode --verbose
Validating spells and blocks in .kcode file 'd:\CODE\go\wandChallengesKcode\009_accio.kcode'...
<xml xmlns="http://www.w3.org/1999/xhtml">
  <variables>
  </variables>
  <block type="events_onGesture" id="#9_FyuL,Y8#]q*3i{O;z" x="172" y="289">
    <field name="TYPE">
      accio
    </field>
    <statement name="CALLBACK">
      <block type="objects_add" id="Hug`V%+[#:b@+ZE#B/}f">
        <field name="ID">
          Broomstick 1
        </field>
        <field name="NAME">
          Broomstick 1
        </field>
        <value name="POSITION">
          <shadow type="position_create" id="]7GX^mF/Z6F3O;^Qq~2h">
            <value name="X">
              <shadow type="math_number" id="d(5QWM7/G}|6(E8bZ%jM">
                <field name="NUM">
                  400
                </field>
              </shadow>
            </value>
            <value name="Y">
              <shadow type="math_number" id="qQ$@G(b/K{X:7lrgPEex">
                <field name="NUM">
                  300
                </field>
              </shadow>
            </value>
          </shadow>
        </value>
      </block>
    </statement>
  </block>
</xml>
SPELL: type=events_onGesture, spell=accio, id=#9_FyuL,Y8#]q*3i{O;z
VALUE: name=CALLBACK,type=,id=,block=530,val=0,shadow=0,statement=0,next=0,dtype=object
VALUE: name=POSITION,type=,id=,block=0,val=0,shadow=328,statement=0,next=0,dtype=object
VALUE: name=,type=position_create,id=]7GX^mF/Z6F3O;^Qq~2h,block=0,val=258,shadow=0,statement=0,next=0,dtype=object
VALUE: name=X,type=,id=,block=0,val=0,shadow=101,statement=0,next=0,dtype=array
VALUE: name=,type=math_number,id=d(5QWM7/G}|6(E8bZ%jM,block=0,val=0,shadow=0,statement=0,next=0,dtype=object
VALUE: name=Y,type=,id=,block=0,val=0,shadow=101,statement=0,next=0,dtype=array
VALUE: name=,type=math_number,id=qQ$@G(b/K{X:7lrgPEex,block=0,val=0,shadow=0,statement=0,next=0,dtype=object
BLOCK: type=events_onGesture, id=#9_FyuL,Y8#]q*3i{O;z, x=172, y=289, statement=562, next=0, val=0, bdtype=object
VALUE: name=CALLBACK,type=,id=,block=530,val=0,shadow=0,statement=0,next=0,dtype=object
BLOCK: type=objects_add, id=Hug`V%+[#:b@+ZE#B/}f, x=, y=, statement=0, next=0, val=361, bdtype=object
VALUE: name=POSITION,type=,id=,block=0,val=0,shadow=328,statement=0,next=0,dtype=object
VALUE: name=,type=position_create,id=]7GX^mF/Z6F3O;^Qq~2h,block=0,val=258,shadow=0,statement=0,next=0,dtype=object
VALUE: name=X,type=,id=,block=0,val=0,shadow=101,statement=0,next=0,dtype=array
VALUE: name=,type=math_number,id=d(5QWM7/G}|6(E8bZ%jM,block=0,val=0,shadow=0,statement=0,next=0,dtype=object
VALUE: name=Y,type=,id=,block=0,val=0,shadow=101,statement=0,next=0,dtype=array
VALUE: name=,type=math_number,id=qQ$@G(b/K{X:7lrgPEex,block=0,val=0,shadow=0,statement=0,next=0,dtype=object
SUCCEEDED in validating 'd:\CODE\go\wandChallengesKcode\009_accio.kcode'.
Expected and found 1 spells and 2 blocks
========== FINISHED ===========
Elapsed time = 6.9981ms
```
