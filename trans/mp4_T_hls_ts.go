package trans

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/grafov/m3u8"
	"github.com/yapingcat/gomedia/go-mp4"
	"github.com/yapingcat/gomedia/go-mpeg2"
)

type Ts struct {
	Uri        string
	DurationMs uint
	TsFile     *os.File
}

func Mp4TransTs(mp4Path string, tsPrefix string, durationMs uint64, name string) {
	table := make(map[int]uint)
	var i int
	var start uint64
	var end uint64
	var tsFile *os.File
	var vid uint
	var aid uint
	var lastPacket *mp4.AVPacket
	var tsPath string

	os.Mkdir(tsPrefix, os.ModePerm)

	ts := make([]*Ts, 0)
	mp4File, err := os.Open(mp4Path)
	if err != nil {
		slog.Error("Open", "err", err)
		return
	}

	demuxer := mp4.CreateMp4Demuxer(mp4File)
	tracks, err := demuxer.ReadHead()
	if err != nil {
		slog.Error("Mp4TransTs", "err", err)
		return
	}

	tsMuxer := mpeg2.NewTSMuxer()
	for _, track := range tracks {
		if track.Cid == mp4.MP4_CODEC_H264 {
			vid = uint(tsMuxer.AddStream(mpeg2.TS_STREAM_H264))
			table[track.TrackId] = vid
		}

		if track.Cid == mp4.MP4_CODEC_H265 {
			vid = uint(tsMuxer.AddStream(mpeg2.TS_STREAM_H265))
			table[track.TrackId] = vid
		}

		if track.Cid == mp4.MP4_CODEC_AAC {
			aid = uint(tsMuxer.AddStream(mpeg2.TS_STREAM_AAC))
			table[track.TrackId] = aid
		}

	}

	tsPath = tsPrefix + fmt.Sprintf("%s_%d.ts", name, i)

	tsFile, err = os.OpenFile(tsPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		slog.Error("OpenFile", "err", err)
		return
	}

	for {

		packet, err := demuxer.ReadPacket()
		if err != nil {
			tsPath, _ = filepath.Abs(tsPath)
			ts = append(ts, &Ts{
				Uri:        tsPath,
				TsFile:     tsFile,
				DurationMs: uint(lastPacket.Dts - start),
			})
			tsFile.Close()
			break
		}
		lastPacket = packet

		if start == 0 {
			tsMuxer.OnPacket = func(pkg []byte) {
				tsFile.Write(pkg)
			}
			start = packet.Dts
			end = packet.Dts + uint64(durationMs)
		}

		slog.Info("ReadPacket", "Dts", packet.Dts, "Pts", packet.Pts, "start", start, "end", end)

		if packet.Dts >= end {
			i++
			tsFile.Close()
			tsPath, _ = filepath.Abs(tsPath)
			ts = append(ts, &Ts{
				Uri:        tsPath,
				TsFile:     tsFile,
				DurationMs: uint(packet.Dts - start),
			})

			start = packet.Dts
			end = packet.Dts + durationMs
			tsPath = tsPrefix + fmt.Sprintf("%s_%d.ts", name, i)

			tsFile, err = os.OpenFile(tsPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
			if err != nil {
				slog.Error("OpenFile", "err", err)
				return
			}
			pid := table[packet.TrackId]
			tsMuxer.Write(uint16(pid), packet.Data, packet.Pts, packet.Dts)

			continue
		}

		pid := table[packet.TrackId]
		tsMuxer.Write(uint16(pid), packet.Data, packet.Pts, packet.Dts)
	}

	gennerlM3u8(ts)
}

func gennerlM3u8(ts []*Ts) {

	file, err := os.OpenFile("playlist.m3u8", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}

	playList, _ := m3u8.NewMediaPlaylist(uint(len(ts)), uint(len(ts)))
	for _, v := range ts {
		slog.Info("ts", "url", v.Uri, "duration", v.DurationMs)
		// 将分段添加到播放列表
		playList.Append(v.Uri, float64(v.DurationMs)/1000.0, "")
	}

	playList.Close()

	// 生成 M3U8 文件
	m3u8PlayList := playList.Encode().Bytes()
	file.Write(m3u8PlayList)
}
