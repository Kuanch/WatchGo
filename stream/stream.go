package stream

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Stream represents a single video feed.
type Stream struct {
	c             chan []byte
	frame         []byte
	FrameInterval time.Duration
	Disconnect    chan bool
}

const boundaryWord = "MJPEGBOUNDARY"
const headerf = "\r\n" +
	"--" + boundaryWord + "\r\n" +
	"Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: 0.000000\r\n" +
	"\r\n"

// ServeHTTP responds to HTTP requests with the MJPEG stream, implementing the http.Handler interface.
func (s *Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Stream:", r.RemoteAddr, "connected")
	defer s.disconnect()
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+boundaryWord)

	for {
		s.Disconnect <- false
		b := <-s.c
		_, err := w.Write(b)
		if err != nil {
			break
		}
	}

	log.Println("Stream:", r.RemoteAddr, "disconnected")
}

// UpdateJPEG pushes a new JPEG frame onto the clients.
func (s *Stream) UpdateJPEG(jpeg []byte) {
	header := fmt.Sprintf(headerf, len(jpeg))
	if len(s.frame) < len(jpeg)+len(header) {
		s.frame = make([]byte, (len(jpeg)+len(header))*2)
	}

	copy(s.frame, header)
	copy(s.frame[len(header):], jpeg)

	s.c <- s.frame
}

func (s *Stream) disconnect() {
	fmt.Println("Stream diconnection is called")
	s.Disconnect <- true
}

// NewStream initializes and returns a new Stream.
func NewStream() *Stream {
	return &Stream{
		c:             make(chan []byte),
		frame:         make([]byte, len(headerf)),
		FrameInterval: 50 * time.Millisecond,
		Disconnect:    make(chan bool),
	}
}
