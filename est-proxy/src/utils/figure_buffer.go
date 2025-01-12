package utils

import (
	"github.com/dustinxie/lockfree"
	"log"
	"time"
)

type FigureBuffer struct {
	data            lockfree.HashMap
	bufferedFigures lockfree.Queue
}

func NewFigureBuffer() *FigureBuffer {
	return &FigureBuffer{
		data:            lockfree.NewHashMap(),
		bufferedFigures: lockfree.NewQueue(),
	}
}

func (fb *FigureBuffer) Add(figureId string, newData []byte) {
	it, flag := fb.data.Get(figureId)
	if flag == false {
		fb.data.Set(figureId, figureUpdateData{
			data: newData,
			time: time.Now(),
		})
		fb.bufferedFigures.Enque(figureId)
		return
	}
	fb.data.Set(figureId, figureUpdateData{
		data: append(it.(figureUpdateData).data, newData...),
		time: time.Now(),
	})
}

func (fb *FigureBuffer) ServeFlush(callback FlushFunc) {
	for {
		time.Sleep(3 * time.Second)

		if fb.bufferedFigures.Len() == 0 {
			continue
		}

		inWorkFigures := lockfree.NewQueue()

		for fb.bufferedFigures.Len() > 0 {
			figureId := fb.bufferedFigures.Deque().(string)
			figureData, figureExists := fb.data.Get(figureId)
			if !figureExists {
				log.Printf("Figure %s not found in buffer", figureId)
				continue
			}
			if time.Now().Sub(figureData.(figureUpdateData).time) > 3*time.Second {
				callback(figureId, figureData.(figureUpdateData).data)
				fb.data.Del(figureId)
				log.Printf("%s has been removed from the buffer", figureId)
			} else {
				inWorkFigures.Enque(figureId)
			}
		}

		for inWorkFigures.Len() > 0 {
			fb.bufferedFigures.Enque(inWorkFigures.Deque().(string))
		}
	}
}

type FlushFunc func(figureId string, message []byte)

type figureUpdateData struct {
	data []byte
	time time.Time
}
