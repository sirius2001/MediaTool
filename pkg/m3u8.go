package pkg

import (
	"fmt"
	"log"
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

	tsPath := make([]string, 0)

	// 解析 M3U8 文件
	pl, _, err := m3u8.DecodeFrom(m3u8File, true)
	if err != nil {
		log.Fatal(err)
	}

	switch pl := pl.(type) {
	case *m3u8.MediaPlaylist:
		for _, segment := range pl.Segments {
			if segment == nil {
				break
			}
			tsPath = append(tsPath, segment.URI)
		}
	default:
		return nil, fmt.Errorf("unexpected playlist type")
	}

	return tsPath, nil
}
