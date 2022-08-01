package main

import (
  "errors"
  "log"
  "encoding/json"
  "io/ioutil"
  "net/http"
  "net/url"
  "github.com/valyala/fastjson"
)

const URL = "https://apic-desktop.musixmatch.com/ws/1.1/macro.subtitles.get"

type Musixmatch struct {
  Token string
}

func (mx Musixmatch) findLyrics(track Track) (Song, error) {
  song := Song{}
  baseURL, _ := url.Parse(URL)
  params := url.Values{
   "format": {"json"},
   "namespace": {"lyrics_richsynched"},
   "subtitle_format": {"mxm"},
   "app_id": {"web-desktop-app-v1.0"},
   "usertoken": {mx.Token},
   "q_album": {track.AlbumName},
   "q_artist": {track.ArtistName},
   "q_artists": {track.ArtistName},
   "q_track": {track.TrackName},
   "track_spotify_id": {""},
   "q_duration": {""},
   "f_subtitle_length": {""},
  }
  baseURL.RawQuery = params.Encode()

  client := http.Client{}
  req , err := http.NewRequest("GET", baseURL.String(), nil)
  if err != nil {
    return song, err
  }

  req.Header = http.Header{
    "authority": {"apic-desktop.musixmatch.com"},
    "cookie": {"x-mxm-token-guid="},
  }

  res , err := client.Do(req)
  if err != nil {
    return song, err
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return song, err
  }

  var p fastjson.Parser
  v, err := p.Parse(string(body))
  if err != nil {
    return song, err
  }

  mtg := v.Get("message", "body", "macro_calls", "matcher.track.get", "message")
  tlg := v.Get("message", "body", "macro_calls", "track.lyrics.get", "message")
  tsg := v.Get("message", "body", "macro_calls", "track.subtitles.get", "message")

  switch mtg.GetInt("header", "status_code") {
    case 200:
      if err := json.Unmarshal(mtg.Get("body", "track").MarshalTo(nil), &song.Track); err != nil {
        return song, err
      }
    case 401:
      return song, errors.New("Cooldown")
    case 404:
      return song, errors.New("No results found")
    default:
      return song, errors.New("Unknown error")
    }

  if song.Track.HasSubtitles == 1 {
    if err := json.Unmarshal(tsg.GetStringBytes("body", "subtitle_list", "0", "subtitle", "subtitle_body"), &song.Subtitle.Lines); err != nil {
      return song, err
    }
  } else {
    log.Println("No synced lyrics found")
    if song.Track.HasLyrics == 1 {
      if err := json.Unmarshal(tlg.Get("body", "lyrics").MarshalTo(nil), &song.Lyrics); err != nil {
        return song, err
      }
      if song.Lyrics.Restricted == 1 {
        return song, errors.New("Restricted lyrics")
      }
    } else {
      return song, errors.New("No lyrics found")
    }
  }
  return song, nil
}
