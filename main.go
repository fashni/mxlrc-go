package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"

  "github.com/alexflint/go-arg"
)

var inputs InputsQueue
var failed InputsQueue

func main() {
  var args Args
  arg.MustParse(&args)

  mode := parseInput(args, &inputs)
  cnt := inputs.len()
  fmt.Printf("\n%d lyrics to fetch\n\n", cnt)

  if mode == "dir" {
    args.Outdir = ""
  } else {
    if err := os.MkdirAll(args.Outdir, os.ModePerm); err != nil {
      log.Fatal(err)
    }
  }

  closeHandler(mode, cnt)
  var token string
  if token = args.Token; args.Token == "" {
    token = "2203269256ff7abcb649269df00e14c833dbf4ddfb5b36a1aae8b0"
  }
  mx := Musixmatch{
    Token: token,
  }

  for !inputs.empty() {
    cur := inputs.next()
    log.Printf("searching song: %s - %s", cur.Track.ArtistName, cur.Track.TrackName)
    song, err := mx.findLyrics(cur.Track)
    if err == nil {
      log.Println("formatting Lyrics")
      success := writeLRC(song, cur.Filename, cur.Outdir)
      cur = inputs.pop()
      if !success {
        log.Println("failed to save lyrics")
        failed.push(cur)
      }
    } else {
      log.Println(err)
      failed.push(inputs.pop())
    }
    timer(args.Cooldown, inputs.len())
  }
  if !failed.empty() {
    failedHandler(mode, cnt)
  }
}

func timer(maxSec int, n int) {
  if n <= 0 {
    return
  }
  for i := maxSec; i >= 0; i-- {
    fmt.Printf("    Please wait... %ds    \r", i)
    time.Sleep(time.Second)
  }
  fmt.Printf("\n\n")
}

func failedHandler(mode string, cnt int) {
  fmt.Printf("\n")
  if !inputs.empty() {
    failed.Queue = append(failed.Queue, inputs.Queue...)
  }
  log.Printf("Succesfully fetch %d out of %d lyrics.", cnt-failed.len(), cnt)
  if failed.empty() {
    return
  }
  log.Printf("Failed to fetch %d lyrics.", failed.len())

  if mode == "dir" {
    log.Println("You can try again with the same command")
  } else {
    t := time.Now().Format("20060102_150405")
    fn := t + "_failed.txt"
    log.Printf("Saving list of failed items in %s. You can try again using this file as the input", fn)

    f, err := os.Create(fn)
    if err != nil {
      log.Fatal(err)
    }
    defer f.Close()

    buffer := bufio.NewWriter(f)
    for !failed.empty() {
      cur := failed.pop()
      _, err := buffer.WriteString(cur.Track.ArtistName + "," + cur.Track.TrackName + "\n")
      if err != nil {
        log.Fatal(err)
      }
    }
    if err := buffer.Flush(); err != nil {
      log.Fatal(err)
    }
  }
}

func closeHandler(mode string, cnt int) {
  c := make(chan os.Signal)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-c
    fmt.Printf("\n")
    failedHandler(mode, cnt)
    os.Exit(0)
  }()
}
