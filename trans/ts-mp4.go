package trans

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log/slog"
	"os"

	"github.com/yapingcat/gomedia/go-mp4"
	"github.com/yapingcat/gomedia/go-mpeg2"
)

func TsTransMp4(mp4file *os.File, tsFiles []*os.File) {
	hasAudio := false
	hasVideo := false
	var atid uint32 = 0
	var vtid uint32 = 0

	muxer, err := mp4.CreateMp4Muxer(mp4file)
	if err != nil {
		fmt.Println(err)
		return
	}

	demuxer := mpeg2.NewTSDemuxer()

	demuxer.OnFrame = func(cid mpeg2.TS_STREAM_TYPE, frame []byte, pts uint64, dts uint64) {
		//add video trank
		if cid == mpeg2.TS_STREAM_H264 {
			if !hasVideo {
				vtid = muxer.AddVideoTrack(mp4.MP4_CODEC_H264)
				hasVideo = true
			}
			err := muxer.Write(vtid, frame, uint64(pts), uint64(dts))
			if err != nil {
				fmt.Println(err)
			}
		} else if cid == mpeg2.TS_STREAM_H265 {
			if !hasVideo {
				vtid = muxer.AddVideoTrack(mp4.MP4_CODEC_H264)
				hasVideo = true
			}
			err := muxer.Write(vtid, frame, uint64(pts), uint64(dts))
			if err != nil {
				fmt.Println(err)
			}
		}

		//add audio trank
		if cid == mpeg2.TS_STREAM_AAC {
			if !hasAudio {
				atid = muxer.AddAudioTrack(mp4.MP4_CODEC_AAC)
				hasAudio = true
			}
			err := muxer.Write(atid, frame, uint64(pts), uint64(dts))
			if err != nil {
				fmt.Println(err)
			}
		} else if cid == mpeg2.TS_STREAM_AUDIO_MPEG1 || cid == mpeg2.TS_STREAM_AUDIO_MPEG2 {
			if !hasAudio {
				atid = muxer.AddAudioTrack(mp4.MP4_CODEC_MP3)
				hasAudio = true
			}
			err := muxer.Write(atid, frame, uint64(pts), uint64(dts))
			if err != nil {
				fmt.Println(err)
			}
		}

	}

	for _, tsFile := range tsFiles {
		buf, err := ioutil.ReadAll(tsFile)
		if err != nil {
			slog.Error("TsTransMp4", "read tsFile err", err)
			continue
		}
		fmt.Println(demuxer.Input(bytes.NewReader(buf)))
		muxer.WriteTrailer()
	}

}
