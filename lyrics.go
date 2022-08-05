package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "path/filepath"
  "strings"
)

func writeLRC(song Song, filename string, outdir string) {
  var fn string
  if fn = filename; filename == "" {
    fn = slugify(fmt.Sprintf("%s - %s", song.Track.ArtistName, song.Track.TrackName)) + ".lrc"
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

  f, err := os.Create(fp)
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

  if len(song.Subtitles.Lines) > 0 {
    log.Println("saving synced lyrics")
    writeSyncedLRC(song, fp, buffer)
    log.Printf("synced lyrics saved: %s", fp)
    return
  }
  if song.Lyrics.LyricsBody != "" {
    log.Println("saving unsynced lyrics")
    writeUnsyncedLRC(song, fp, buffer)
    log.Printf("unsynced lyrics saved: %s", fp)
    return
  }
  if song.Track.Instrumental == 1 {
    log.Println("saving instrumental")
    writeInstrumentalLRC(song, fp, buffer)
    log.Printf("instrumental lyrics saved: %s", fp)
    return
  }
}

func writeUnsyncedLRC(song Song, fpath string, buff *bufio.Writer) {
  lines := strings.Split(song.Lyrics.LyricsBody, "\n")
  var text string
  for _, line := range lines {
    if text = line; line == "" {
      text = "♪"
    }
    _, err := buff.WriteString("[00:00.00]" + text + "\n")
    if err != nil {
      log.Fatal(err)
    }
  }

  if err := buff.Flush(); err != nil {
    log.Fatal(err)
  }
}

func writeSyncedLRC(song Song, fpath string, buff *bufio.Writer) {
  var text string
  var fLine string
  for _, line := range song.Subtitles.Lines {
    if text = line.Text; line.Text == "" {
      text = "♪"
    }
    fLine = fmt.Sprintf("[%02d:%02d.%02d]%s", line.Time.Minutes, line.Time.Seconds, line.Time.Hundredths, text)
    _, err := buff.WriteString(fLine + "\n")
    if err != nil {
      log.Fatal(err)
    }
  }

  if err := buff.Flush(); err != nil {
    log.Fatal(err)
  }
}

func writeInstrumentalLRC(song Song, fpath string, buff *bufio.Writer) {
  line := "[00:00.00]♪ Instrumental ♪"
  _, err := buff.WriteString(line + "\n")
  if err != nil {
    log.Fatal(err)
  }
}
