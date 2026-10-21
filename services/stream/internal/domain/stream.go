package domain

import "errors"

type TorrentLinkQuery struct {
	TorrentLink string `query:"torrent_link"`
}

type StreamResponse struct {
}

var ErrTorrentLoadingTimeout = errors.New("Torrent loading timeout")
var ErrTorrentNoFile = errors.New("Torrent has no files")
var ErrNoVideo = errors.New("No video")
var ErrUnsupportedContainer = errors.New("Unsupported format")
var ErrBadEBML = errors.New("mkv: invalid EBML data")
var ErrWrongFile = errors.New("mkv: not a matroska file")
var ErrSegmentNotFound = errors.New("mkv: segment not found")
var ErrTrackNotFound = errors.New("mkv: video track not found")
var ErrNoCues = errors.New("mkv: no cues (no seek index)")
var ErrCuesBadPosition = errors.New("mkv: cues not at expected position")
var ErrCuesEmpty = errors.New("mkv: cues empty for video track")
