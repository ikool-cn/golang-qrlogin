package main

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	cp := NewChannelPool()
	ch := cp.GetChannelById("channelid")
	fmt.Println(ch)
	ch = cp.GetChannelById("channelid")
	fmt.Println(ch)

	cp.DestroyChannelById("channelid")

	ch = cp.GetChannelById("channelid")
	fmt.Println(ch)

}
