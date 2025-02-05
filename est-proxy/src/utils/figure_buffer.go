package utils

import (
	"est-proxy/src/config"
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

func (fb *FigureBuffer) Remove(figureId string) bool {
	_, flag := fb.data.Get(figureId)
	if flag == false {
		return false
	}
	log.Printf("Deleted figure %s from buffer", figureId) //for debug TODO: remove this
	fb.data.Del(figureId)
	return true
}

func (fb *FigureBuffer) ServeFlush(callback FlushFunc) {
	for {
		time.Sleep(config.BUFFERED_FIGURE_LIVE_TIME)

		if fb.bufferedFigures.Len() == 0 {
			continue
		}

		inWorkFigures := lockfree.NewQueue()

		for fb.bufferedFigures.Len() > 0 {
			figureId := fb.bufferedFigures.Deque().(string)
			figureData, figureExists := fb.data.Get(figureId)
			if !figureExists {
				log.Printf("Figure %s not found in buffer", figureId) //for debug TODO: remove this
				continue
			}
			if time.Now().Sub(figureData.(figureUpdateData).time) > config.BUFFERED_FIGURE_LIVE_TIME {
				fb.safeFlushCall(callback, figureId, figureData.(figureUpdateData).data)
				fb.data.Del(figureId)
			} else {
				inWorkFigures.Enque(figureId)
			}
		}

		for inWorkFigures.Len() > 0 {
			fb.bufferedFigures.Enque(inWorkFigures.Deque().(string))
		}
	}
}

func (fb *FigureBuffer) safeFlushCall(callback FlushFunc, figureId string, message []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered after figure flush callback: %v", r)
		}
	}()
	callback(figureId, message)
}

type FlushFunc func(figureId string, message []byte)

type figureUpdateData struct {
	data []byte
	time time.Time
}
