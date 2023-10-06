package mankidown

import (
	"fmt"
	"io"
	"strings"
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

	guidColumnName = "Mankidown-Guid"
	tagsColumnName = "Tags" 
)

type Config struct {
	GuidPrefix   string 
	InFileName string
	Deck       string
	NoteType   string
	Tags       []string
}

type Exporter struct {
	out    io.Writer
	config *Config
}

func NewExporter(out io.Writer, config *Config) *Exporter {
	return &Exporter{out: out, config: config}
}

func (ex *Exporter) AppendHeaders() {
	// separator
	fmt.Fprintln(ex.out, HeaderSeparator)

	// html
	fmt.Fprintln(ex.out, HeaderHtml)

	// guid column
	fmt.Fprintln(ex.out, HeaderGuidColumn)

	// tag column
	fmt.Fprintln(ex.out, HeaderTagColumn)

	// deck
	if "" != ex.config.Deck {
		fmt.Fprintf(ex.out, HeaderDeck, ex.config.Deck)
	}

	// notetype
	if "" != ex.config.NoteType {
		fmt.Fprintf(ex.out, HeaderNoteType, ex.config.NoteType)
	}

	// Tags
	if len(ex.config.Tags) > 0 {
		tagsStr := strings.Join(ex.config.Tags, " ")
		fmt.Fprintf(ex.out, HeaderTags, tagsStr)
	}
}

func (ex *Exporter) AppendHeaderColumns(columns []string) {

	cols := []string{}
	// prepend id as first (anki duplication)
	cols = append(cols, guidColumnName) //1)
	cols = append(cols, tagsColumnName) //2)

	for _, fieldName := range columns {
		cols = append(cols, fieldName)
	}

	c := strings.Join(cols, separatorChar)
	fmt.Fprintf(ex.out, "#columns:%s\n", c)
}

func buildIdField(n *Note, config *Config) string {

	if "" != config.GuidPrefix {
		return config.GuidPrefix + n.Id()
	}
	return n.Id()
}

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

func (ex *Exporter) Export(notes *Notes) {
	ex.AppendHeaders()
	ex.AppendHeaderColumns(notes.FieldNames())

	notelines := []string{}
	for _, note := range notes.Notes {

		// 1 field) id
		fields := []string{}
		fields = append(fields, buildIdField(note, ex.config))

		// 2 field) tags
		fields = append(fields, buildTagsField(note))

		for _, field := range note.Fields() {
			fields = append(fields, buildFieldAsLine(field.Html))
		}

		noteline := strings.Join(fields, separatorChar)
		fmt.Fprintf(ex.out, "%s\n", noteline)

		notelines = append(notelines, noteline)
	}
}
