package pkg

import (
	"bytes"
	"io"
	"log/slog"
	"os"

	"github.com/grafov/m3u8"
)

func ParaseM3u8(path string) ([]string, error) {

	m3u8File, err := os.Open(path)
	if err != nil {
		slog.Error("ParaseM3u8", "err", err)
		return nil, err
	}

	data, err := io.ReadAll(m3u8File)
	if err != nil {
		slog.Error("ParaseM3u8", "ReadAll err", err)
		return nil, err
	}

	p, listType, err := m3u8.Decode(*bytes.NewBuffer(data), true)
	if err != nil {
		slog.Error("ParaseM3u8", "m3u8.Decode", err)
		return nil, err
	}

	tsArray := make([]string, 0)

	switch listType {
	case m3u8.MEDIA:
		mediaPlaylist := p.(*m3u8.MediaPlaylist)
		for _, segment := range mediaPlaylist.Segments {
			tsArray = append(tsArray, segment.URI)
		}
	case m3u8.MASTER:
		masterPlaylist := p.(*m3u8.MasterPlaylist)

		for _, variant := range masterPlaylist.Variants {
			tsArray = append(tsArray, variant.URI)
		}
	}
	return tsArray, nil
}
