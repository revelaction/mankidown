package mankidown

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"strconv"
	"strings"
)

type Parser struct {
	mdp goldmark.Markdown
}

// A Field contains the parsed html content inside a H2 Heading element
type Field struct {
	Html string
}

// A Note contains parsed mardown html that will be mapped to a anki note type
// field in the output file
type Note struct {
	id     string
	tags   []string
	fields []Field
}

func newNote() *Note {
	return &Note{
		tags:   []string{},
		fields: []Field{},
	}
}

func (n *Note) Id() string {
	return n.id
}

func (n *Note) Fields() []Field {
	return n.fields
}

func (n *Note) Tags() []string {
	return n.tags
}

func (n *Note) addField(f Field) {
	n.fields = append(n.fields, f)
}

func (n *Note) numFieds() int {
	return len(n.fields)
}

func (n *Note) addTags(tags string) {

	words := strings.Fields(tags)

    n.tags = append(n.tags, words...)
}

func (n *Note) addId(id string) {
	n.id = id
}

// A Notes contains the parsed notes elements
type Notes struct {
	Notes      []*Note
	fieldNames []string
}

func newNotes() *Notes {
	return &Notes{
		Notes: []*Note{},
	}
}

func (n *Notes) FieldNames() []string {
	return n.fieldNames
}

func (n *Notes) addNote(nt *Note) error {
	err := n.validateNote(nt)
	if err != nil {
		return err
	}

	n.Notes = append(n.Notes, nt)
	return nil
}

func (n *Notes) validateNote(nt *Note) error {

	if nt.numFieds() == 0 {
		return fmt.Errorf("no fields in note %d", n.numNotes()+1)
	}

	if n.numFieds() != nt.numFieds() {
		return fmt.Errorf("invalid number of fields in note %d (want %d, have %d)", n.numNotes()+1, n.numFieds(), nt.numFieds())
	}

	if nt.Id() == "" {
		return fmt.Errorf("no notes declared at the start by parsing note %d", n.numNotes()+1)
	}

	return nil
}

func (n *Notes) numFieds() int {
	return len(n.fieldNames)
}

func (n *Notes) numNotes() int {
	return len(n.Notes)
}

func (ns *Notes) addFieldName(fieldName string) error {

	// after first note:
	if ns.numNotes() > 0 {
		if fieldName != "" {
			return fmt.Errorf("invalid presence of fields (%q) in note %d", fieldName, ns.numNotes()+1)
		} else {
			return nil
		}
	}

	// first note:
	if fieldName == "" {
		return fmt.Errorf("missing fields in note %d", ns.numNotes()+1)
	}

	for _, fn := range ns.fieldNames {
		if fn == fieldName {
			return fmt.Errorf("Note %d: Field %q already exist", ns.numNotes()+1, fieldName)
		}
	}

	ns.fieldNames = append(ns.fieldNames, fieldName)
	return nil
}

func (p *Parser) Parse(markdown []byte) (*Notes, error) {

	root := p.mdp.Parser().Parse(text.NewReader(markdown))

	var fieldBuf bytes.Buffer
	var insideNoteField bool = false

	nt := newNote()
	notes := newNotes()

	err := ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		var err error

		if !entering {

			if isFieldEnd(n, insideNoteField) {
				nt.addField(Field{Html: fieldBuf.String()})
				insideNoteField = false
			}

			if isNoteEnd(n) {
				err = notes.addNote(nt)
				if err != nil {
					return ast.WalkStop, err
				}

			}

			return ast.WalkSkipChildren, nil
		}

		// entering
		// Ignore Document node entering
		if n.Kind() == ast.KindDocument {
			return ast.WalkContinue, nil
		}

		if isNoteStart(n) {
			nt = newNote()

			// Tags are defined in the H1 header
			if tags := string(n.Text(markdown)); tags != "" {
				nt.addTags(tags)
			}

			// note guid suffix
			noteNum := len(notes.Notes) + 1
			nt.addId(strconv.Itoa(noteNum))

			return ast.WalkSkipChildren, nil
		}

		if isFieldStart(n) {
			fieldText := string(n.Text(markdown))
			err = notes.addFieldName(fieldText)
			if err != nil {
				return ast.WalkStop, err
			}

			fieldBuf = bytes.Buffer{}

			insideNoteField = true

			return ast.WalkSkipChildren, nil
		}

		// Render the node contents
		if err = p.mdp.Renderer().Render(&fieldBuf, markdown, n); err != nil {
			return ast.WalkStop, fmt.Errorf("render error %v", err)
		}

		return ast.WalkSkipChildren, nil
	})

	if err != nil {
		return nil, err
	}

	return notes, nil
}

func NewParser() *Parser {

	md := goldmark.New(
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	return &Parser{mdp: md}
}

func isNoteStart(n ast.Node) bool {

	h, ok := n.(*ast.Heading)
	if !ok {
		return false
	}

	if h.Level == 1{
		return true
	}

	return false
}

// the Note end is determined by either:
// - a following (sibling) h1, but no coming fron preceding Document start
// - the following is Document final (not entering)
func isNoteEnd(n ast.Node) bool {

	_, ok := n.(*ast.Document)
	if ok {
		return true
	}

	if nil == n.NextSibling() {
		return false
	}

	switch v := n.NextSibling().(type) {
	case *ast.Document:
		return true
	case *ast.Heading:
		if v.Level == 1 {
			return true
		}
	}

	return false
}

func isFieldStart(n ast.Node) bool {
	h, ok := n.(*ast.Heading)
	if !ok {
		return false
	}

	if h.Level == 2{
		return true
	}

	return false
}

// the field end is determined by either:
// - a following (sibling) h2, but no coming fron preceding h1
// - the following is Document final (not entering)
func isFieldEnd(n ast.Node, inside bool) bool {

	if !inside {
		return false
	}

	_, ok := n.(*ast.Document)
	if ok {
		return true
	}

	if nil == n.NextSibling() {
		return false
	}

	switch v := n.NextSibling().(type) {
	case *ast.Document:
		return true
	case *ast.Heading:
		if v.Level == 2 {
			return true
		}
		if v.Level == 1{
			return true
		}
	}

	return false
}
