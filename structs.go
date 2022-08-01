package main

type Track struct {
  TrackName      string `json:"track_name"`
  ArtistName     string `json:"artist_name"`
  AlbumName      string `json:"album_name"`
  TrackLength    int    `json:"track_length"`
  Instrumental   int    `json:"instrumental"`
  HasLyrics      int    `json:"has_lyrics"`
  HasSubtitles   int    `json:"has_subtitles"`
}

type Lyrics struct {
  Restricted int    `json:"restricted,omitempty"`
  LyricsBody string `json:"lyrics_body,omitempty"`  
}

type Synced struct {
  Lines []Lines
}

type Lines struct {
  Text string `json:"text,omitempty"`
  Time Time `json:"time,omitempty"`
}

type Time struct {
  Total float64 `json:"total,omitempty"`
  Minutes int `json:"minutes,omitempty"`
  Seconds int `json:"seconds,omitempty"`
  Hundredths int `json:"hundredths,omitempty"`
}
