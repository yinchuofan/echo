package core

import (
	"container/list"
	"log"
	"sync"
)

// Subscriber  subscribe some channels, save channels subscribed
type Subscriber interface {
	Notify(message string)
}

//key-channel, value-subscriber list
var (
	mapSubscriber = make(map[string]*list.List)
	mapSubMutex   = new(sync.Mutex)
)

// Subscribe subrcibe channels
func Subscribe(channel string, subcriber Subscriber) {
	log.Printf("subscribe channel: %v, subcriber: %v", channel, subcriber)
	mapSubMutex.Lock()
	defer mapSubMutex.Unlock()
	if elem, ok := mapSubscriber[channel]; ok {
		elem.PushBack(subcriber)
	} else {
		subList := list.New()
		subList.PushBack(subcriber)
		mapSubscriber[channel] = subList
	}
}

// Unsubscribe unsubscribe channel
func Unsubscribe(channel string, subcriber Subscriber) {
	log.Printf("unsubscribe channel: %v, subcriber: %v", channel, subcriber)
	mapSubMutex.Lock()
	defer mapSubMutex.Unlock()
	if subList, ok := mapSubscriber[channel]; ok {
		for e := subList.Front(); e != nil; e = e.Next() {
			if e.Value.(Subscriber) == subcriber {
				subList.Remove(e)
				return
			}
		}
	}
}

// Publish publish message to channels
func Publish(channel string, message string) {
	//log.Printf("publish channel: %v, message:\r\n %v", channel, message)
	log.Printf("publish channel: %v \n", channel)
	mapSubMutex.Lock()
	defer mapSubMutex.Unlock()
	if subList, ok := mapSubscriber[channel]; ok {
		for e := subList.Front(); e != nil; e = e.Next() {
			sub := e.Value.(Subscriber)
			if sub != nil {
				go sub.Notify(message)
			}
		}
	} else {
		log.Printf("no one subscribe channel: %v", channel)
	}
}
