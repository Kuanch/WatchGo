package utils

import (
	"fmt"

	"github.com/Kuanch/mjpeg"
	"gocv.io/x/gocv"
)

var (
	webcamErr error
	webcam    *gocv.VideoCapture
)

func VideoFeed(deviceID int, stream *mjpeg.Stream) {
	webcam, webcamErr = gocv.OpenVideoCapture(deviceID)
	if webcamErr != nil {
		return
	}
	defer webcam.Close()
	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}
