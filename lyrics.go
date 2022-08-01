package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "strings"
  "path/filepath"
)

func writeLRC(song Song, filename string, outdir string) {
  var fn string
  if fn = filename; filename == "" {
    fn = fmt.Sprintf("%s - %s.lrc", song.Track.ArtistName, song.Track.TrackName)
  }
  fp := filepath.Join(outdir, fn)

  tags := []string{
    "[by:fashni]",
    fmt.Sprintf("[ar:%s]", song.Track.ArtistName),
    fmt.Sprintf("[ti:%s]", song.Track.TrackName),
  }
  if song.Track.AlbumName != "" {
    tags = append(tags, fmt.Sprintf("[al:%s]", song.Track.AlbumName))
  }
  if song.Track.TrackLength != 0 {
    tags = append(tags, fmt.Sprintf("[length:%02d:%02d]", song.Track.TrackLength/60, song.Track.TrackLength%60))
  }

  if len(song.Subtitle.Lines) > 0 {
    log.Println("Saving synced lyrics")
    writeSyncedLRC(song, fp, tags)
    return
  }
  if song.Lyrics.LyricsBody != "" {
    log.Println("Saving unsynced lyrics")
    writeUnsyncedLRC(song, fp, tags)
    return
  }
}

func writeUnsyncedLRC(song Song, fpath string, tags []string) {
  f, err := os.Create(fpath)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()
  
  buffer := bufio.NewWriter(f)
  for _, tag := range tags {
    _, err := buffer.WriteString(tag + "\n")
    if err != nil {
      log.Fatal(err)
    }
  }

  lines := strings.Split(song.Lyrics.LyricsBody, "\n")
  var text string
  for _, line := range lines {
    if text = line; line == "" {
      text = "♪"
    }
    _, err := buffer.WriteString("[00:00.00]" + text + "\n")
    if err != nil {
      log.Fatal(err)
    }
  }

  if err := buffer.Flush(); err != nil {
    log.Fatal(err)
  }
}

func writeSyncedLRC(song Song, fpath string, tags []string) {
  f, err := os.Create(fpath)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  buffer := bufio.NewWriter(f)
  for _, tag := range tags {
    _, err := buffer.WriteString(tag + "\n")
    if err != nil {
      log.Fatal(err)
    }
  }

  var text string
  var fLine string
  for _, line := range song.Subtitle.Lines {
    if text = line.Text; line.Text == "" {
      text = "♪"
    }
    fLine = fmt.Sprintf("[%02d:%02d.%02d]%s", line.Time.Minutes, line.Time.Seconds, line.Time.Hundredths, text)
    _, err := buffer.WriteString(fLine + "\n")
    if err != nil {
      log.Fatal(err)
    }
  }

  if err := buffer.Flush(); err != nil {
    log.Fatal(err)
  }
}
