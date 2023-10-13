# 

## Text

Which mankidown flag below should you use to indicate the note type (in this case Cloze):

```console
$ mankidown --deck mydeck --{{c1::note-type}} Cloze mynotes.md

```

## Extra

The short version of the flag is `-n`

# 
## 

Which mankidown flag below should you use to indicate the anki deck:

```console
$ mankidown --{{c1::deck}} mydeck -n Cloze mynotes.md

```

## 

The short version of the flag is `-d`

# 
## 

Which mankidown flag below should you use to add the tag `"anki"` to all notes in `mynotes.md`:

```console
$ mankidown --deck mydeck -n Cloze --{{c1::tag}} anki mynotes.md

```

## 

The short version of the flag is `-t`

# 
## 

Fill de dots below to add the tags `"mankidown"` and `"anki"` to all the notes in `mynotes.md`:

```console
$ mankidown --mydeck -n Cloze {{c1::--tag mankidown --tag anki}} mynotes.md

```

## 

The tag `--tag`, `-t` can be repeated.
