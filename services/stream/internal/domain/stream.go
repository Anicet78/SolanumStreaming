package domain

import "errors"

type TorrentLinkQuery struct {
	TorrentLink string `query:"torrent_link"`
}

type StreamResponse struct {
}

var ErrTorrentLoadingTimeout = errors.New("Torrent loading timeout")
var ErrTorrentNoFile = errors.New("Torrent has no files")
var ErrNoSegments = errors.New("No segments")
var ErrNoVideo = errors.New("No video")
var ErrUnsupportedContainer = errors.New("Unsupported format")
var ErrBadEBML = errors.New("mkv: invalid EBML data")
var ErrWrongFile = errors.New("mkv: not a matroska file")
var ErrSegmentNotFound = errors.New("mkv: segment not found")
var ErrTrackNotFound = errors.New("video track not found")
var ErrNoCues = errors.New("mkv: no cues (no seek index)")
var ErrCuesBadPosition = errors.New("mkv: cues not at expected position")
var ErrCuesEmpty = errors.New("mkv: cues empty for video track")
var ErrInvalidData = errors.New("mp4: invalid data")
var ErrMoovNotFoundData = errors.New("mp4: moov not found")
var ErrBoxWithUnknownSize = errors.New("mp4: box with unknown size")
var ErrFmp4NotSupported = errors.New("mp4: fragmented mp4 not supported")
var ErrIncompleteVideoTrack = errors.New("mp4: incomplete video track")
var ErrNoStssTable = errors.New("mp4: no stts table")
var ErrNoKeyframes = errors.New("mp4: no keyframes")
