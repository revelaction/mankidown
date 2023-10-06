// Copyright (c) 2023 revelaction

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/revelaction/mankidown"
)

const outFileExt = ".txt"

const usage = `Usage:
    mankidown [-d DECK] [-n NOTE-TYPE] [-n GUID-PREFIX] [-t TAG] [-o OUTPUT] INPUT

Options:
    -d, --deck                  The anki deck to import the notes
    -n, --note-type             The anki note type to import the notes
    -o, --output OUTPUT         The OUTPUT file to be imported by Anki 
    -p, --guid-prefix           A prefix to build the guid of all notes note of INPUT
    -t, --tag                   A Tag to all notes of INPUT. Can be repeated 

INPUT should be a valid markdown file with the following structure:

- Each H1 Heading element is imported as an Anki card
- H1 Heading elements can contain (in the header line) anki tags separated by space.  
- H2 Heading elements are imported as a "note type" Anki field. 
- The H2 Headings in the first note should contain a word indicating the Anki field to be mapped.
- The H2 Headings in the rest of the notes can not contain words. Its Anki
  field is derived from th first note and its position.

OUTPUT defaults to the basename of INPUT file plus the
extension ".txt". Anki import UI will allow just a list of compatible
extensions, including ".txt". 

If OUTPUT exists, it will be overwritten.

DECK is the anki deck to import the notes. If not set, manual adjustment is
required in the Anki import UI.

NOTE-TYPE is the note type to import the notes. If not set, manual adjustment
is required in the Anki import UI.

TAG is a tag to be applied to all notes of INPUT. The option can be repeated.

GUID-PREFIX is a prefix of the "guid" field for each note in the OUTPUT file.
The "guid" field is unique for each note of INPUT. mankidown will contruct a
guid for each note with the  concatenation of the GUID-PREFIX and a sequence
integer. Anki uses the "guid" to find duplicates and update those notes
preserving schedule times. GUID-PREFIX defaults to OUTPUT.

Examples:
    $ mankidown mynotes.md # will write mynotes.txt, which can be imported in Anki
    $ mankidown --deck mydeck -n mytype -p go -t go -t anki mynotes.md`

type multiFlag []string

func (f *multiFlag) String() string {
    return fmt.Sprint(*f) 
}

func (f *multiFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {

	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }

	conf := new(mankidown.Config)
	var outFlag, guidPrefix string
	var tagFlags multiFlag

	flag.StringVar(&outFlag, "o", "", "output to `FILE` (default stdout)")
	flag.StringVar(&outFlag, "output", "", "output to `FILE` (default stdout)")
	flag.StringVar(&conf.Deck, "d", "", "Export to the Anki Deck")
	flag.StringVar(&conf.Deck, "deck", "", "Export to the Anki Deck")
	flag.StringVar(&conf.NoteType, "n", "", "Export to the Anki Note type")
	flag.StringVar(&conf.NoteType, "note-type", "", "Export to the Anki Note type")
	flag.StringVar(&guidPrefix, "p", "", "Export with the note prefix")
	flag.StringVar(&guidPrefix, "guid-prefix", "", "Export with the note prefix")
	flag.Var(&tagFlags, "t", "Add tag")
	flag.Var(&tagFlags, "tags", "Add tag")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() > 1 {
		errorf("too many INPUT arguments %q", strings.Join(flag.Args(), " "))
	}

	if len(tagFlags) > 0 {
		conf.Tags = tagFlags
	}

	var in io.Reader
	var inFile string = flag.Arg(0)
	if inFile != "" && inFile != "-" {
		f, err := os.Open(inFile)
		if err != nil {
			errorf("failed to open input file %q: %v", inFile, err)
		}
		defer f.Close()
		in = f
	}

	var outFile string

	if outFlag != "" {
		outFile = outFlag
	} else {
		baseName := filepath.Base(inFile)
		outFile = strings.TrimSuffix(baseName, filepath.Ext(baseName)) + outFileExt
	}

	if guidPrefix != "" {
		conf.GuidPrefix = guidPrefix
	} else {
		conf.GuidPrefix = outFile
	}

	var out io.Writer
	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		errorf("failed to open output file %q: %v", outFlag, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			errorf("failed to close output file %q: %v", outFlag, err)
		}
	}()
	out = f

	markdown, err := io.ReadAll(in)
	if err != nil {
		errorf("failed to read input file %q: %v", inFile, err)
	}

	md := mankidown.NewParser()
	notes, err := md.Parse(markdown)
	if err != nil {
		errorf("failed to parse input file %q: %v", inFile, err)
	}

	ex := mankidown.NewExporter(out, conf)
	ex.Export(notes)
}

// l is a logger with no prefixes.
var l = log.New(os.Stderr, "", 0)

func errorf(format string, v ...interface{}) {
	l.Printf("mankidown: error: "+format, v...)
	os.Exit(1)
}
