package mankidown

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// Anki header fields.
// See https://docs.ankiweb.net/importing/text-files.html#file-headers
const (
	HeaderSeparator  = "#separator:Pipe"
	HeaderHtml       = "#html:true"
	HeaderGuidColumn = "#guid column:1"
	HeaderTagColumn  = "#tags column:2"
	HeaderDeck       = "#deck:%s\n"
	HeaderNoteType   = "#notetype:%s\n"
	HeaderTags       = "#tags:%s\n"

	separatorChar = `|`

	guidColumnName = "mankidown-Guid"
	tagsColumnName = "Tags"

	outFileExt = ".txt"
)

type Config struct {
	GuidPrefix string
	InFile     string
	Deck       string
	NoteType   string
	Tags       []string
}

type Exporter struct {
	config *Config
}

func NewExporter(config *Config) *Exporter {
	return &Exporter{config: config}
}

func (ex *Exporter) Export(notes *Notes) error {

	var err error
	inFileBaseName := filepath.Base(ex.config.InFile)
	inFile := strings.TrimSuffix(inFileBaseName, filepath.Ext(inFileBaseName))
	outFile := inFile + outFileExt

	var out io.Writer
	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open output file %q: %v", f, err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			err = fmt.Errorf("failed to close output file %q: %v", f, err)
		}
	}()

	out = f

	ex.appendHeaders(out)
	ex.appendHeaderColumns(out, notes.FieldNames())
	ex.appendHeaderTags(out, inFile)

	for _, note := range notes.Notes {

		// 1 field) id
		fields := []string{}
		fields = append(fields, ex.buildIdField(note, outFile))

		// 2 field) tags
		fields = append(fields, buildTagsField(note))

		for _, field := range note.Fields() {
			fields = append(fields, buildFieldAsLine(field.Html))
		}

		noteline := strings.Join(fields, separatorChar)
		fmt.Fprintf(out, "%s\n", noteline)

	}

	return nil
}

func (ex *Exporter) appendHeaders(out io.Writer) {
	// separator
	fmt.Fprintln(out, HeaderSeparator)

	// html
	fmt.Fprintln(out, HeaderHtml)

	// guid column
	fmt.Fprintln(out, HeaderGuidColumn)

	// tag column
	fmt.Fprintln(out, HeaderTagColumn)

	// deck
	if ex.config.Deck != "" {
		fmt.Fprintf(out, HeaderDeck, ex.config.Deck)
	}

	// notetype
	if ex.config.NoteType != "" {
		fmt.Fprintf(out, HeaderNoteType, ex.config.NoteType)
	}

}

func (ex *Exporter) appendHeaderColumns(out io.Writer, columns []string) {

	cols := []string{}
	// prepend id as first (anki duplication)
	cols = append(cols, guidColumnName) //1)
	cols = append(cols, tagsColumnName) //2)

	cols = append(cols, columns...)

	c := strings.Join(cols, separatorChar)
	fmt.Fprintf(out, "#columns:%s\n", c)
}

func (ex *Exporter) appendHeaderTags(out io.Writer, inFile string) {

	tags := ex.config.Tags

	// Split the inFile words
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	tagsFromInFile := strings.FieldsFunc(inFile, f)
	tags = append(tags, tagsFromInFile...)
	fmt.Fprintf(out, HeaderTags, strings.Join(tags, " "))
}

func (ex *Exporter) buildIdField(n *Note, outFile string) string {

	if ex.config.GuidPrefix != "" {
		return ex.config.GuidPrefix + n.Id()
	}

	return outFile + n.Id()
}

// buildTagsField builds the Tags string for a note
func buildTagsField(n *Note) string {

	return strings.Join(n.Tags(), " ")
}

func buildFieldAsLine(html string) string {

	// 1) replace all but the last \n of the field,
	occurrencesCount := strings.Count(html, "\n")
	fieldLine := strings.Replace(html, "\n", "<br>", occurrencesCount-1)

	// 2) replace the last one with nothing
	fieldLine = strings.Replace(fieldLine, "\n", "", -1)

	// 3) delete <br> between tags
	fieldLine = strings.Replace(fieldLine, "><br><", "><", -1)
	return fieldLine
}
