package utils

import (
	"fmt"

	"watch_go/stream"

	"gocv.io/x/gocv"
)

var (
	webcamErr error
	webcam    *gocv.VideoCapture
)

func VideoFeed(deviceID int, streamer *stream.Stream) {
	webcam, webcamErr = gocv.OpenVideoCapture(deviceID)
	if webcamErr != nil {
		fmt.Println(webcamErr)
		return
	}
	defer webcam.Close()
	img := gocv.NewMat()
	defer img.Close()

	for {
		if <-streamer.Disconnect {
			fmt.Println("Streamer disconnect")
			break
		}
		if ok := webcam.Read(&img); !ok {
			fmt.Println("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		streamer.UpdateJPEG(buf)
	}
}
