# guid:md001 tag1
## Text

```console
$ mankidown --deck mydeck --{{c1::note-type}} Cloze mynotes.md

```
## Question

Which mankidown flag below should you use to indicate the note type (in this case Cloze):

## Extra

The short version of the flag is `-n`

## Source

[mankidown](https://github.com/revelaction/mankidown)

# guid:md002 tag2
## 

```console
$ mankidown --{{c1::deck}} mydeck -n Cloze mynotes.md

```

## 
Which mankidown flag below should you use to indicate the anki deck:

## 

The short version of the flag is `-d`

##

[mankidown](https://github.com/revelaction/mankidown)

# guid:md003 tag3 tag4
## 

```console
$ mankidown --deck mydeck -n Cloze --{{c1::tag}} anki mynotes.md

```

## 
Which mankidown flag below should you use to add the tag `"anki"` to all notes in `mynotes.md`:

## 

The short version of the flag is `-t`

## 

[mankidown](https://github.com/revelaction/mankidown)

# guid:md004 tag5 tag6 tag7
## 

```console
$ mankidown --mydeck -n Cloze {{c1::--tag mankidown --tag anki}} mynotes.md

```

##
Fill de dots below to add the tags `"mankidown"` and `"anki"` to all the notes in `mynotes.md`:

## 

The tag `--tag`, `-t` can be repeated.

##

[mankidown](https://github.com/revelaction/mankidown)
