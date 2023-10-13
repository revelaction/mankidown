# 
## Front

Which mankidown flag below (XXX) should you use to indicate the note type (in this case Cloze):

```console
$ mankidown --deck mydeck --XXX Basic mynotes.md

```

## Back

```console
$ mankidown --deck mydeck --note-type Basic mynotes.md

```

The short version of the flag is `-n`

# 
## 

Which mankidown flag below should you use to indicate the anki deck:

```console
$ mankidown --XXX mydeck -n Basic mynotes.md

```

## 

```console
$ mankidown --deck mydeck -n Basic mynotes.md

```
The short version of the flag is `-d`

# 
## 

Which mankidown flag below (`XXX`) should you use to add the tag `"anki"` to all notes in `mynotes.md`:

```console
$ mankidown --deck mydeck -n Basic --XXX anki mynotes.md

```

## 

```console
$ mankidown --deck mydeck -n Basic --tag anki mynotes.md

```

The short version of the flag is `-t`
