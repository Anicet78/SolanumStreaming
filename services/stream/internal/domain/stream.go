package domain

import "errors"

type TorrentLinkQuery struct {
	TorrentLink string `query:"torrent_link"`
}

type StreamResponse struct {
}

var ErrTorrentLoadingTimeout = errors.New("Torrent loading timeout")
