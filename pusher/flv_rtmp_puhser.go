package pusher

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/yapingcat/gomedia/go-codec"
	"github.com/yapingcat/gomedia/go-flv"
	"github.com/yapingcat/gomedia/go-rtmp"
)

type FlvPuhser struct {
}

// Will push the last file under mp4sPath to the specified rtmp server
func (p FlvPuhser) FlvStart(flvPath string, rtmpUrl string, loop bool) {
	url, _ := url.Parse(rtmpUrl)
	c, err := net.Dial("tcp4", url.Host) // like 127.0.0.1:1935
	if err != nil {
		fmt.Println(err)
	}

	cli := rtmp.NewRtmpClient(rtmp.WithComplexHandshake(),
		rtmp.WithComplexHandshakeSchema(rtmp.HANDSHAKE_COMPLEX_SCHEMA0),
		rtmp.WithEnablePublish())

	cli.OnError(func(code, describe string) {
		fmt.Printf("rtmp code:%s ,describe:%s\n", code, describe)
	})

	isReady := make(chan struct{})
	cli.OnStatus(func(code, level, describe string) {
		fmt.Printf("rtmp onstatus code:%s ,level %s ,describe:%s\n", code, describe)
	})
	cli.OnStateChange(func(newState rtmp.RtmpState) {
		if newState == rtmp.STATE_RTMP_PUBLISH_START {
			fmt.Println("ready for publish")
			close(isReady)
		}
	})
	cli.SetOutput(func(bytes []byte) error {
		_, err := c.Write(bytes)
		return err
	})
	go func() {
		<-isReady
		fmt.Println("start to read file")
		FlvPushRtmp(flvPath, cli, loop)

	}()

	cli.Start(rtmpUrl)
	buf := make([]byte, 4096)
	n := 0
	for err == nil {
		n, err = c.Read(buf)
		if err != nil {
			continue
		}
		fmt.Println("read byte", n)
		cli.Input(buf[:n])
	}
	fmt.Println(err)
}

func FlvPushRtmp(fileName string, cli *rtmp.RtmpClient, loop bool) {
	flvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer flvFile.Close()
	reader := flv.CreateFlvReader()

	reader.OnFrame = func(cid codec.CodecID, frame []byte, pts, dts uint32) {
		if cid == codec.CODECID_VIDEO_H264 {
			time.Sleep(20 * time.Millisecond)
			pts := video_pts_adjust.adjust(int64(pts))
			dts := video_dts_adjust.adjust(int64(dts))
			cli.WriteVideo(codec.CODECID_VIDEO_H264, frame, uint32(pts), uint32(dts))
		} else if cid == codec.CODECID_AUDIO_AAC {
			pts := audio_ts_adjust.adjust(int64(pts))
			cli.WriteAudio(codec.CODECID_AUDIO_AAC, frame, uint32(pts), uint32(pts))
		} else if cid == codec.CODECID_VIDEO_H265 {
			time.Sleep(20 * time.Millisecond)
			pts := video_pts_adjust.adjust(int64(pts))
			dts := video_dts_adjust.adjust(int64(dts))
			cli.WriteAudio(codec.CODECID_AUDIO_MP3, frame, uint32(pts), uint32(dts))
		}
	}

	for {
		buffer := make([]byte, 1024)
		n, err := flvFile.Read(buffer)
		if err != nil {
			break
		}
		reader.Input(buffer[:n])
	}

	select {}

}
