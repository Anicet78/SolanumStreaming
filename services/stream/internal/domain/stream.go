package domain

import "errors"

type TorrentLinkParam struct {
	TorrentLink string `param:"torrent_link"`
}

type StreamResponse struct {
}

var ErrTorrentLoadingTimeout = errors.New("Torrent loading timeout")
