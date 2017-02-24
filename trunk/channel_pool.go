package main

import (
	"log"
	"sync"
)

type ChannelPool struct {
	Channels map[string]chan string
	mutex    sync.Mutex
}

func NewChannelPool() *ChannelPool {
	return &ChannelPool{
		Channels: make(map[string]chan string),
	}
}

func (cp *ChannelPool) GetChannelById(chanId string) chan string {
	cp.mutex.Lock()
	if channel, ok := cp.Channels[chanId]; !ok {
		channel = make(chan string)
		cp.Channels[chanId] = channel
		log.Println("Create a channel, channel_id", chanId)
	}
	cp.mutex.Unlock()
	return cp.Channels[chanId]
}

func (cp *ChannelPool) DestroyChannelById(chanId string) {
	cp.mutex.Lock()
	if _, ok := cp.Channels[chanId]; ok {
		delete(cp.Channels, chanId)
		log.Println("Delete a channel, channel_id", chanId)
	}
	cp.mutex.Unlock()
}
