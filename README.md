<p align="center"><img alt="mankidown" src="data/logo.png"/></p>

[![Go Reference](https://pkg.go.dev/badge/github.com/revelaction/mankidown)](https://pkg.go.dev/github.com/revelaction/mankidown)
[![Test](https://github.com/revelaction/mankidown/actions/workflows/test.yml/badge.svg)](https://github.com/revelaction/mankidown/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/revelaction/mankidown)](https://goreportcard.com/report/github.com/revelaction/mankidown)
[![GitHub Release](https://img.shields.io/github/v/release/revelaction/mankidown?style=flat)]() 

`mankidown` is a simple command-line tool to convert [markdown](https://commonmark.org/help/) documents to [Anki](https://apps.ankiweb.net/) flashcards 

# Features 

- mankidown output is a plain text file that meet the conditions of the [Anki import process](https://docs.ankiweb.net/importing/text-files.html). To import the output file, click the File menu and then "Import".
- mankidown supports the [Anki text file headers](https://docs.ankiweb.net/importing/text-files.html#file-headers) to simplify the import process.
- mankidown renders the markdown contents as html.
- mankidown supports per file and per note [Anki tags](https://docs.ankiweb.net/importing/text-files.html#adding-tags).
- mankidown supports [custom Anki notes types](https://docs.ankiweb.net/editing.html#adding-a-note-type). Each `note type field` is mapped to a markdown H2 Heading element.

# Installation

On Linux, macOS, and FreeBSD and Windows you can use the [pre-built binaries](https://github.com/revelaction/mankidown/releases/) 

If your system has a supported version of Go, you can build from source

```console
go install github.com/revelaction/mankidown/cmd/mankidown@latest
```

# Run mankidown

```console
# convert to cards for mydeck and Basic note type
mankidown --deck mydeck -n Basic mynotes.md

# Also add tags to all notes
mankidown --deck mydeck -n Basic -t tag1 -t tag2 mynotes.md
```

# Usage

## 1) Write notes in markdown

Create a markdown file called `mynotes.md` with your favorite editor. The structure of the markdown file is the same regardless of the anki note type. 
In this example, we write two notes for the anki `Basic` Note Type. This note type is the anki default and has two note fields: `Front` and `Back`.

````markdown
# anki mandidown
## Front  

What mankidown flag should you use to indicate the `note type`? 

## Back  

Use the `-n, --note-type` flag:

```
mankidown -n Basic mynotes.md 
```
# anki mandidown
## 

What mankidown flag shoudl you use indicate the `Deck`? 

## 

Use the `-d, --deck` flag:

```
mankidown -d myDeck mynotes.md 
```
````

The markdown file above will create two anki cards for the `Basic` note type.

The structure of the markdown file is simple:

- Each H1 Heading element is imported as an Anki note
- H1 Heading elements can contain (in the header line) Anki tags separated by space. Anki will apply those tags to the note. 
- H2 Heading elements are imported as Anki note type fields. 
- The H2 Headings of the first note SHOULD contain the Anki field to be mapped.
- The H2 Headings of the rest of the notes should not contain words. Its Anki
  field is derived from the first note.

See the [markidown format definition](data/examples/mankidown-format.md) in the `data/examples` directory.


## 2) Run mankidown

Run mankidown indicating the Anki note type, the deck, and optional tags:

```console
mankidown --deck mydeck -n Basic mynotes.md
```

This will produce a `mynotes.txt` file that can be imported in Anki.

## 3) Import the output file in anki

In the Anki app, click the File menu and then "Import". For the desktop app:

![anki import](data/anki-import.png)

# Examples

## Cloze

There is an [example markdown Cloze file](data/examples/cloze.md) in the `data/examples` directory.
The anki `Cloze` type has two fields: [`Text` and `Extra`](https://docs.ankiweb.net/editing.html#cloze-deletion)

To write a `Cloze` note just use the anki convention `{{c1::<hidden response>}}` after the `Text` H2 Heading. For example:

````markdown
# 
## Text
Which mankidown flag below should you use to indicate the note type (in this case Cloze):

```console
mankidown --deck mydeck --{{c1::note-type}} Cloze mynotes.md
```
````

You can run the Cloze example file like this:

```console
mankidown -d myDeck -n Cloze data/examples/cloze.md 
```

and import the .txt file in anki.

## Custom

# GUID, updating cards, duplicates

mankidown makes use of the [GUID Column](https://docs.ankiweb.net/importing/text-files.html#guid-column) of the Anki import process.

Anki notes imported from mankidown have the Anki GUID (Globally Unique Identifier) set by mankidown.

There are two ways mankidown uses to build the GUID:

#### Per note GUID

Each note can declare its Anki GUID in the H1 header by using the prefix `guid:` followed by the guid:

```markdown
# tag1 tag2 guid:563ho8
## Front

## Back
```

The markdown file above will produce one anki note with the guid `563ho8`.

If per note GUIDs are used, all notes of the markdown file are required to have one.

Per note GUIDs should be preferred.

#### Position in document

if there are no `guid:`s per note, mankidown will make a guid by appending an
integer to the name of the markdown file.  That integer is just the note number
(position) in the markdown file:

```console
mymarkdownfile0, mymarkdownfile1, mymarkdownfile2
```

It's possible to use a Guid prefix instead of the file name in the command line.

```console
$ mankidown -d myDeck -n Cloze -p myguidprefix data/examples/cloze.md 
```

This method is only useful if there are no or little deletion of notes in the
document, as they will change the guid of the following notes.

# Tags

[Anki Tags](https://docs.ankiweb.net/editing.html?highlight=tags#using-tags) are a useful way to keep your collection organized. 
mankidown can generate `document tags`  and `note tags`

#### Document Tags

By using the `-t, --tag` flag in the command line, mankidown adds the tag to
all notes of the document. The flag can be repeated

```console
$ mankidown --deck mydeck -n Basic -t tag1 -t tag2 mynotes.md
```

Additionally mankidown splits the words of the markdown document name and adds
them as tags for all notes. If the markdown file is called
`english-vocabulary.md`, mankidown will add `english` and `vocabulary` to all
notes of the file.

#### Note Tags

The title of the H1 tag can be used to add tags to the note separated by spaces:

```markdown
# tag1 tag2
## Front
This is a Basic note with two tags 
## Back

the tags are `tag1` ans `tag2`
```

# Notes with media

## Images

Put the image files in the [Anki collection.media folder](https://docs.ankiweb.net/files.html#file-locations).
In the markdown file, use the standard image format without paths:

    ![myimage](myimage.jpg)

## Audio

Put the audio files in the [Anki collection.media folder](https://docs.ankiweb.net/files.html#file-locations).
In the markdown file, use the standard audio format without paths:

    [sound:mysound.mp3]

# Command line options
```console
mankidown [-d DECK] [-n NOTE-TYPE] [-n GUID-PREFIX] [-t TAG] [-o OUTPUT] INPUT

Options:
    -d, --deck                  The anki deck to import the notes
    -n, --note-type             The anki note type to import the notes
    -o, --output OUTPUT         The OUTPUT file to be imported by Anki 
    -p, --guid-prefix           A prefix to build the guid of all notes note of INPUT
    -t, --tag                   A Tag to all notes of INPUT. Can be repeated 
        --version               The version 

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
```
