package mankidown_test

import (
	"testing"

	"github.com/revelaction/mankidown"
)

func TestParseValid(t *testing.T) {

	good := []byte(`# note1 tag1 
## Field1  

First note has two fields, Field1 and Field2

## Field2  

text

# note2 tag2
## 

text

## 

Text`)

	md := mankidown.NewParser()
	_, err := md.Parse(good)

	if err != nil {
		t.Errorf("\ngot Error %s\nwant nil", err)
	}
}

func TestParseH2WordsAfterFirstNote(t *testing.T) {
	withH2Word := []byte(`# note1 tag1 
## Front  

text

## Back  

text

# note2 tag2
## Front 

text 2 

## Back 

text 2`)

	md := mankidown.NewParser()
	_, err := md.Parse([]byte(withH2Word))

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}

func TestParseNoteWithRepeatedField(t *testing.T) {
	withRepeatedField := []byte(`
# note1 tag1 
## Front  

text

## Back  

text

## Front  

text
`)

	md := mankidown.NewParser()
	_, err := md.Parse(withRepeatedField)

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}

func TestParseNoteWithDifferentNumberOfFields(t *testing.T) {

	withDifferentNumberOfFields := []byte(`
# note1 tag1 
## Front  

First note has two fields, Front and Back

## Back  

text

# note2 tag2
## 

second note has 3 fields

## 

second note has 3 fields

## 

second note has 3 fields
`)

	md := mankidown.NewParser()
	_, err := md.Parse(withDifferentNumberOfFields)

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}

// Text after H1 before H2 are allowed
func TestParseContentBetweenH1AndH2(t *testing.T) {

	withTextBetween := []byte(`# note1 tag1 

some text between H1 and H2

## Field1  

First note has two fields, Field1 and Field2

## Field2  

text

# note2 tag2
## 

text

## 

Text`)

	md := mankidown.NewParser()
	_, err := md.Parse(withTextBetween)

	if err != nil {
		t.Errorf("\ngot %s\nwant nil", err)
	}
}

// first note should define all note fields, should no have empty H2
func TestParseNoFieldInFirstNoteH2(t *testing.T) {

	withoutFields := []byte(`# note1 tag1
## 

First note has two fields, Field1 and Field2

## 

text

# note2 tag2
## 

text

## 

Text`)

	md := mankidown.NewParser()
	_, err := md.Parse(withoutFields)

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}

func TestParseNoFieldsInNote(t *testing.T) {

	withoutFields := []byte(`# tag1 tag2
`)

	md := mankidown.NewParser()
	_, err := md.Parse(withoutFields)

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}

func TestParseNoNoteDeclaration(t *testing.T) {

	withoutNote := []byte(`##
`)

	md := mankidown.NewParser()
	_, err := md.Parse(withoutNote)

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}

func TestParseNoNoteBeforeFields(t *testing.T) {

	withoutNote := []byte(`## field1

some text

## field2

text

##
`)

	md := mankidown.NewParser()
	_, err := md.Parse(withoutNote)

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}

func TestParseNoNoteAndFieldDeclarations(t *testing.T) {

	withoutNotes := []byte(`###
`)

	md := mankidown.NewParser()
	_, err := md.Parse(withoutNotes)

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}

func TestParseNoFieldsAfterNote(t *testing.T) {

	withoutFields := []byte(`#
###
`)

	md := mankidown.NewParser()
	_, err := md.Parse(withoutFields)

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}

func TestParseNotesWithAndWithoutGuids(t *testing.T) {

	WithAndWithout := []byte(`# note1 tag1 guid:note1
## field1 

First note has two fields, Field1 and Field2

## field2 

text

# note2 tag2 noguid
## 

text

## 

Text`)

	md := mankidown.NewParser()
	_, err := md.Parse(WithAndWithout)

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}

func TestParseNotesWithDuplicateGuid(t *testing.T) {

	WithAndWithout := []byte(`# note1 tag1 guid:note1
## field1 

First note has two fields, Field1 and Field2

## field2 

text

# note2 tag2 guid:note1
## 

text

## 

Text`)

	md := mankidown.NewParser()
	_, err := md.Parse(WithAndWithout)

	t.Logf("Error is %q", err)
	if err == nil {
		t.Errorf("\ngot %s\nwant Error", err)
	}
}
