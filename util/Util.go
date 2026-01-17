package util

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func RunAtFramerate(fps int, totalFrames int, frameFunc func(frameNumber int)) {
	frameDuration := time.Second / time.Duration(fps)

	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()

	frameCount := 0

	for frameCount < totalFrames {
		<-ticker.C
		frameFunc(frameCount)
		frameCount++
	}
}
