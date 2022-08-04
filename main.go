package main

import (
  "fmt"
  "log"
  "os"
  "time"

  "github.com/alexflint/go-arg"
)

func main() {
  var args struct {
    Song     []string `arg:"positional,required" help:"song information in [ artist,title ] format (required)"`
    Outdir   string   `arg:"-o,--outdir" help:"output directory, default: lyrics" default:"lyrics"`
    Cooldown int      `arg:"-c,--cooldown" help:"cooldown time in seconds, default: 30" default:"30"`
    Token    string   `arg:"-t,--token" help:"musixmatch token" default:""`
  }
  arg.MustParse(&args)

  inputs, mode := parseInput(args.Song, args.Outdir)
  if mode == "dir" {
    args.Outdir = ""
  } else {
    if err := os.MkdirAll(args.Outdir, os.ModePerm); err != nil {
      log.Fatal(err)
    }
  }

  var token string
  if token = args.Token; args.Token == "" {
    token = "2203269256ff7abcb649269df00e14c833dbf4ddfb5b36a1aae8b0"
  }
  mx := Musixmatch{
    Token: token,
  }

  for idx, input := range inputs {
    log.Printf("searching song: %s - %s", input.Track.ArtistName, input.Track.TrackName)
    song, err := mx.findLyrics(input.Track)
    if err != nil {
      log.Println(err)
      continue
    }

    log.Println("formatting Lyrics")
    writeLRC(song, input.Filename, input.Outdir)

    if idx+1 < len(inputs) {
      for i := args.Cooldown; i >= 0; i-- {
        fmt.Printf("    Please wait... %ds    \r", i)
        time.Sleep(time.Second)
      }
      fmt.Println("")
    }
  }
}
