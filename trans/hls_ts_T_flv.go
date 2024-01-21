package trans

import (
	"log/slog"
	"os"

	"github.com/yapingcat/gomedia/go-flv"
	"github.com/yapingcat/gomedia/go-mpeg2"
)

func HlsTsTransFlv(tsPathArry []string, flvPath string) error {
	flvFile, err := os.OpenFile(flvPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer flvFile.Close()
	if err != nil {
		slog.Info("HlsTsTransFlv", "OpenFile", err)
		return err
	}

	flvWriter := flv.CreateFlvWriter(flvFile)
	flvWriter.WriteFlvHeader()

	demuxer := mpeg2.NewTSDemuxer()
	demuxer.OnFrame = func(cid mpeg2.TS_STREAM_TYPE, frame []byte, pts, dts uint64) {
		if cid == mpeg2.TS_STREAM_AAC {
			flvWriter.WriteAAC(frame, uint32(pts), uint32(dts))
		}

		if cid == mpeg2.TS_STREAM_H264 {
			flvWriter.WriteH264(frame, uint32(pts), uint32(dts))
		}

		if cid == mpeg2.TS_STREAM_H265 {
			flvWriter.WriteH265(frame, uint32(pts), uint32(dts))
		}
	}

	for _, tsPath := range tsPathArry {
		tsFile, err := os.Open(tsPath)
		if err != nil {
			slog.Warn("tsFile", tsPath, "not exit")
			continue
		}
		demuxer.Input(tsFile)
	}

	return nil
}
