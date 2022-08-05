# MxLRC
[![build](https://github.com/fashni/mxlrc-go/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/fashni/mxlrc-go/actions/workflows/build.yml)

Command line tool to fetch synced lyrics from [Musixmatch](https://www.musixmatch.com/) and save it as *.lrc file.

---

## Python version
[Check it here](https://github.com/fashni/MxLRC)

---

## Download
### Standalone binary
**TBA**

### Build from source
Required Go 1.17+
```
go install github.com/fashni/mxlrc-go@latest
```

---

## Usage
```
Usage: mxlrc-go [--outdir OUTDIR] [--cooldown COOLDOWN] [--token TOKEN] SONG [SONG ...]

Positional arguments:
  SONG                        song information in [ artist,title ] format (required)

Options:
  --outdir OUTDIR, -o OUTDIR  output directory, default: lyrics [default: lyrics]
  --cooldown COOLDOWN, -c COOLDOWN
                              cooldown time in seconds, default: 30 [default: 30]
  --token TOKEN, -t TOKEN     musixmatch token
  --help, -h                  display this help and exit
```

## Example:
### One song
```
mxlrc-go adele,hello
```
### Multiple song and custom output directory
```
mxlrc-go adele,hello "the killers,mr. brightside" -o some_directory
```
### With a text file and custom cooldown time
```
mxlrc-go example_input.txt -c 20
```
### Directory Mode (recursive) **(TBA)**
```
mxlrc-go "Dream Theater"
```
> **_This option overrides the `-o/--outdir` argument which means the lyrics will be saved in the same directory as the given input._**

> **_The `-d/--depth` argument limit the depth of subdirectory to scan. Use `-d 0` or `--depth 0` to only scan the specified directory._**

---

## How to get the Musixmatch Token
Follow steps 1 to 5 from the guide [here](https://spicetify.app/docs/faq#sometimes-popup-lyrics-andor-lyrics-plus-seem-to-not-work) to get a new Musixmatch token.

## Credits
* [Spicetify Lyrics Plus](https://github.com/spicetify/spicetify-cli/tree/master/CustomApps/lyrics-plus)
