package main

import (
  "bufio"
  "errors"
  "log"
  "os"
  "path/filepath"
  "reflect"
  "regexp"
  "sort"
  "strings"

  "github.com/dhowden/tag"
  "golang.org/x/text/unicode/norm"
)

func supportedFType() [8]string {
  return [8]string{".mp3", ".m4a", ".m4b", ".m4p", ".alac", ".flac", ".ogg", ".dsf"}
}

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

func getSongMulti(song_list []string, save_path string, songs *InputsQueue) {
  for _, song := range song_list {
    track := assertInput(song)
    if track == nil {
      log.Printf("invalid input: %s", song)
      continue
    }
    songs.push(Inputs{*track, save_path, ""})
  }
}

func getSongText(text_fn string, save_path string, songs *InputsQueue) {
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
  getSongMulti(song_list, save_path, songs)
}

func getSongDir(dir string, songs *InputsQueue, update bool, limit int, depth int, bfs bool) {
  log.Printf("scanning directory: %s", dir)
  files, err := os.ReadDir(dir)
  if err != nil {
    log.Fatal(err)
  }

  sort.Slice(files, func(i int, j int) bool {
    id1, id2 := files[i].IsDir(), files[j].IsDir()
    if id1 == id2 {
      return files[i].Name() < files[j].Name()
    }
    return !bfs || !id1 && id2
  })

  for _, file := range files {
    if file.IsDir() {
      if depth < limit {
        getSongDir(filepath.Join(dir, file.Name()), songs, update, limit, depth+1, bfs)
      }
      continue
    }
    if filepath.Ext(file.Name()) == ".lrc" {
      continue
    }
    lrc_file := strings.Replace(file.Name(), filepath.Ext(file.Name()), ".lrc", -1)
    if _, err := os.Stat(filepath.Join(dir, lrc_file)); err == nil && !update {
      log.Printf("skipping %s. lyrics file exist.", file.Name())
      continue
    }

    if !isInArray(supportedFType(), strings.ToLower(filepath.Ext(file.Name()))) {
      log.Printf("skipping %s. unsupported file format.", file.Name())
      continue
    }

    f, err := os.Open(filepath.Join(dir, file.Name()))
    if err != nil {
      log.Println("error reading file: ", err)
      continue
    }
    defer f.Close()

    m, err := tag.ReadFrom(f)
    if err != nil {
      log.Println(file.Name(), err)
      continue
    }

    log.Printf("adding %s", file.Name())
    song := Inputs{
      Track:    Track{ArtistName: m.Artist(), TrackName: m.Title()},
      Outdir:   dir,
      Filename: strings.Replace(file.Name(), filepath.Ext(file.Name()), ".lrc", -1),
    }
    songs.push(song)
  }
}

func parseInput(args Args, in *InputsQueue) string {
  if len(args.Song) == 1 {
    fi, err := os.Stat(args.Song[0])
    if err == nil {
      if !fi.IsDir() {
        getSongText(args.Song[0], args.Outdir, in)
        return "text"
      } else {
        getSongDir(args.Song[0], in, args.Update, args.Depth, 0, args.BFS)
        return "dir"
      }
    } else if !errors.Is(err, os.ErrNotExist) {
      log.Fatal(err)
    }
  }
  getSongMulti(args.Song, args.Outdir, in)
  return "cli"
}

func slugify(s string) string {
  re1 := regexp.MustCompile(`[\\\/:*?"<>|]`) // forbidden chars in filename
  re2 := regexp.MustCompile(`[-]+`)          // multiple dashes
  s = norm.NFKC.String(s)
  s = re1.ReplaceAllString(s, "")
  s = re2.ReplaceAllString(s, "-")
  return strings.Trim(s, "-_") // remove trailing and leading dash or underscore
}

func isInArray(arrType interface{}, item interface{}) bool {
  arr := reflect.ValueOf(arrType)
  if arr.Kind() != reflect.Array {
    log.Fatal("invalid data type")
  }
  for i := 0; i < arr.Len(); i++ {
    if arr.Index(i).Interface() == item {
      return true
    }
  }
  return false
}
