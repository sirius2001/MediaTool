package main

import (
	"fmt"
	"os"

	"github.com/sirius2001/MediaTool/pkg"
	"github.com/sirius2001/MediaTool/trans"
)

func main() {
	ts, err := pkg.ParaseM3u8("/home/sirius/Desktop/MediaTool/example/mp4/playlist.m3u8")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ts)
	
	mp4File, err := os.OpenFile("test.mp4", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	trans.LocalTsM3u8TransMp4("/home/sirius/Desktop/MediaTool/example/mp4/playlist.m3u8", mp4File)

}
