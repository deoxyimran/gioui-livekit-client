package video

import (
	"image"
	"log"
	"sync"

	"gocv.io/x/gocv"
)

type WebcamSource struct {
	mutex          sync.Mutex
	frame          image.Image
	devName        string
	isVideoRunning bool
	bgJobStopSig   chan bool
}

func NewWebcamSource(devName string) WebcamSource {
	return WebcamSource{devName: devName}
}

func (w *WebcamSource) SetDevice(name string) {
	w.devName = name
}

func (v *WebcamSource) videoCaptureBackgroundTask() error {
	cap, err := gocv.VideoCaptureFile(v.devName)
	if err != nil {
		return err
	}
	mat := gocv.NewMat()
	go func() {
		loop := true
		for loop {
			select {
			case b := <-v.bgJobStopSig:
				if b {
					loop = !loop
				}
			default:
				if b := cap.Read(&mat); b {
					v.mutex.Lock()
					v.frame, err = mat.ToImage()
					if err != nil {
						log.Println(err.Error())
					}
					v.mutex.Unlock()
				}
			}
		}
		mat.Close()
		cap.Close()
	}()
	return nil
}

func (w *WebcamSource) StartVideo() error {
	if !w.isVideoRunning {
		w.isVideoRunning = true
		err := w.videoCaptureBackgroundTask()
		if err != nil {
			w.isVideoRunning = false
			return err
		}
	}
	return nil
}

func (w *WebcamSource) StopVideo() {
	if w.isVideoRunning {
		w.isVideoRunning = false
		w.bgJobStopSig <- true
	}
}

func (w *WebcamSource) IsVideoOn() bool {
	return w.isVideoRunning
}

func (w *WebcamSource) GetVideoOutFrame() image.Image {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.frame
}
