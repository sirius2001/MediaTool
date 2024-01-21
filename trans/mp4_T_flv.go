package trans

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/yapingcat/gomedia/go-flv"
	"github.com/yapingcat/gomedia/go-mp4"
)

func Mp4TransFlv(mp4Path string, flvPath string) {

	mp4File, err := os.Open(mp4Path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer mp4File.Close()

	demuxer := mp4.CreateMp4Demuxer(mp4File)
	if infos, err := demuxer.ReadHead(); err != nil && err != io.EOF {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", infos)
	}

	flvFile, err := os.OpenFile(flvPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		slog.Error("Mp4TransFlv", "OpenFile", err)
		return
	}
	defer flvFile.Close()

	fw := flv.CreateFlvWriter(flvFile)
	fw.WriteFlvHeader()

	for {
		pkg, err := demuxer.ReadPacket()
		if err != nil {
			break
		}
		fmt.Printf("track:%d,cid:%+v,pts:%d dts:%d\n", pkg.TrackId, pkg.Cid, pkg.Pts, pkg.Dts)
		if pkg.Cid == mp4.MP4_CODEC_H264 {
			fw.WriteH264(pkg.Data, uint32(pkg.Pts), uint32(pkg.Dts))
		} else if pkg.Cid == mp4.MP4_CODEC_AAC {
			fw.WriteAAC(pkg.Data, uint32(pkg.Pts), uint32(pkg.Dts))
		} else if pkg.Cid == mp4.MP4_CODEC_H265 {
			fw.WriteH265(pkg.Data, uint32(pkg.Pts), uint32(pkg.Dts))
		}
	}

}
