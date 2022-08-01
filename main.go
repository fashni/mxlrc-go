package main

import (
  "log"
  "os"
  "strings"
  "github.com/alexflint/go-arg"
)

func main() {
  var args struct {
    Song string `arg:"-s,--song" help:"song information"`
    Outdir string `arg:"-o,--outdir" help:"output directory" default:"lyrics"`
    Token string `arg:"--token" help:"musixmatch token" default:""`
  }
  arg.MustParse(&args)
  s := strings.Split(args.Song, ",")
  err := os.MkdirAll(args.Outdir, os.ModePerm)
  if err != nil {
    log.Fatal(err)
  }

  var token string
  if token = args.Token; args.Token == ""{
    token = "2203269256ff7abcb649269df00e14c833dbf4ddfb5b36a1aae8b0"
  }

  tr := Track{
    ArtistName: s[0],
    TrackName: s[1],
  }

  mx := Musixmatch{
    Token: token,
  }

  log.Printf("Searching song: %s - %s", tr.ArtistName, tr.TrackName)
  song, err := mx.findLyrics(tr)
  if err != nil {
    log.Fatal(err)
    return
  }

  log.Println("Formatting Lyrics")
  writeLRC(song, "", args.Outdir)
  log.Printf("Lyrics saved: %s/%s - %s.lrc", args.Outdir, song.Track.ArtistName, song.Track.TrackName)
}
