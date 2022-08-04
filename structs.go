package main

type Track struct {
  TrackName    string `json:"track_name,omitempty"`
  ArtistName   string `json:"artist_name,omitempty"`
  AlbumName    string `json:"album_name,omitempty"`
  TrackLength  int    `json:"track_length,omitempty"`
  Instrumental int    `json:"instrumental,omitempty"`
  HasLyrics    int    `json:"has_lyrics,omitempty"`
  HasSubtitles int    `json:"has_subtitles,omitempty"`
}

type Lyrics struct {
  LyricsBody string `json:"lyrics_body,omitempty"`
}

type Synced struct {
  Lines []Lines
}

type Lines struct {
  Text string `json:"text,omitempty"`
  Time Time   `json:"time,omitempty"`
}

type Time struct {
  Total      float64 `json:"total,omitempty"`
  Minutes    int     `json:"minutes,omitempty"`
  Seconds    int     `json:"seconds,omitempty"`
  Hundredths int     `json:"hundredths,omitempty"`
}

type Song struct {
  Track     Track
  Lyrics    Lyrics
  Subtitles Synced
}

type Inputs struct {
  Track    Track
  Outdir   string
  Filename string
}
