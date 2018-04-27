package command

import (
	"os"
	"os/exec"
	"time"
	log "github.com/Sirupsen/logrus"
)

const videoPath = "/tmp/appvideo.mov"
const videoLimitTime = time.Minute * 10

var cmd *exec.Cmd

type CmdStatus int
const (
	CmdStatusInProgress CmdStatus = 1
	CmdStatusDone       CmdStatus = 0
)

var currentCmdStatus CmdStatus

// ATTENTION!
// EXPERIMENTAL CODE!

func StartVideo(deviceID string) error {
	err := os.Remove(videoPath)
	if err != nil {
		log.Println(err)
	}
	cmd = exec.Command("xcrun", "simctl", "io", deviceID, "recordVideo", videoPath)
	err = cmd.Start()
	if err != nil {
		return err
	}
	currentCmdStatus = CmdStatusInProgress

	go cmd.Wait()
	go func() {
		select {
		case <-time.After(videoLimitTime):
			FinishVideo()
		}
	}()

	return nil
}

func FinishVideo() (*os.File, error) {
	if currentCmdStatus != CmdStatusDone {
		log.Println("Finishing video...")
		err := cmd.Process.Signal(os.Interrupt)
		if err != nil {
			return nil, err
		}
		time.Sleep(100 * time.Millisecond)
		currentCmdStatus = CmdStatusDone
		return os.Open(videoPath)
	}
	return nil, nil
}
