package main

import (
  "bufio"
  "errors"
  "log"
  "os"
  "regexp"
  "strings"

  "golang.org/x/text/unicode/norm"
)

func assertInput(song string) *Track {
  s := strings.Split(song, ",")
  if len(s) != 2 {
    return nil
  }
  tr := &Track{
    ArtistName: strings.Trim(s[0], " "),
    TrackName:  strings.Trim(s[1], " "),
  }
  return tr
}

func getSongMulti(song_list []string, save_path string) []Inputs {
  var songs []Inputs
  for _, song := range song_list {
    track := assertInput(song)
    if track == nil {
      log.Printf("invalid input: %s", song)
      continue
    }
    songs = append(songs, Inputs{*track, save_path, ""})
  }
  return songs
}

func getSongText(text_fn string, save_path string) []Inputs {
  f, err := os.Open(text_fn)
  if err != nil {
    log.Fatal(err)
  }
  scanner := bufio.NewScanner(f)
  scanner.Split(bufio.ScanLines)
  var song_list []string
  for scanner.Scan() {
    song_list = append(song_list, scanner.Text())
  }
  f.Close()
  return getSongMulti(song_list, save_path)
}

func parseInput(argsong []string, outdir string) ([]Inputs, string) {
  if len(argsong) == 1 {
    fi, err := os.Stat(argsong[0])
    if err == nil {
      if !fi.IsDir() {
        return getSongText(argsong[0], outdir), "text"
      } // else {
        // return getSongDir(argsong[0]), "dir"
      // }
    } else if !errors.Is(err, os.ErrNotExist) {
      log.Fatal(err)
    }
  }
  return getSongMulti(argsong, outdir), "cli"
}

func slugify(s string) string {
  re1 := regexp.MustCompile(`[\\\/:*?"<>|]`) // forbidden chars in filename
  re2 := regexp.MustCompile(`[-]+`)          // multiple dashes
  s = norm.NFKC.String(s)
  s = re1.ReplaceAllString(s, "")
  s = re2.ReplaceAllString(s, "-")
  return strings.Trim(s, "-_") // remove trailing and leading dash or underscore
}
