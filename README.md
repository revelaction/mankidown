`mankidown` is a simple command-line tool to convert markdown documents to Anki flashcards 

# Installation

If your system has a supported version of Go, you can build from source.

    go install github.com/revelaction/mankidown/cmd@latest

# Features 

- `mankidown` output is a plain text file that meet the conditions of the [Anki import process](https://docs.ankiweb.net/importing/text-files.html). To import the output file, click the File menu and then "Import".
- `mankidown` supports the [Anki text file headers](https://docs.ankiweb.net/importing/text-files.html#file-headers) to simplify the import process.
- `mankidown` renders the markdown contents as html.
- `mankidown` supports per file and per note [Anki tags](https://docs.ankiweb.net/importing/text-files.html#adding-tags).
- `mankidown` supports [custom Anki notes types](https://docs.ankiweb.net/editing.html#adding-a-note-type). Each `note type field` is mapped to a markdown H2 Heading element.

# Usage

## Write cards in markdown

Write a markdown file f.ex. `mynotes.md` with notes for the Anki `Basic` Note Type (with `Front` and `Back` fields):

    # anki mandidown
    ## Front  

    What mankidown flag should you use to indicate the `note type`? 

    ## Back  

    Use the `-n, --note-type` flag:
    
    ```
     $ mankidown -n Basic mynotes.md 
    ```
    # anki mandidown
    ## 

    What mankidown flag shoudl you use indicate the `Deck`? 

    ## 

    Use the `-d, --deck` flag:
    
    ```
     $ mankidown -d myDeck mynotes.md 
    ```

The markdown file above will create two anki cards for the `Basic` note type.

The structure of the file is simple:

- Each H1 Heading element is imported as an Anki card
- H1 Heading elements can contain (in the header line) anki tags separated by space.  
- H2 Heading elements are imported as a note type `field`. 
- The H2 Headings in the first note should contain a word indicating the Anki field to be mapped.
- The H2 Headings in the rest of the notes can not contain words. Its Anki
  field is derived from the first note.

## Run mankidown

```
mankidown --deck mydeck -n mytype  mynotes.md
```

This will produce a `mynotes.txt` file that can be imported in Anki.

## Import the output file

In the Anki app, click the File menu and then "Import". For the desktop:

# Duplicated cards, update

# Media

- All content of the markdown is parsed as html. You can use images and sound

# command line options
```
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
```
