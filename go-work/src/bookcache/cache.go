// cache
package bookcache

import (
	"container/list"
	"fmt"
	"time"
)

type Lru struct {
	mp  map[string]*list.Element //book_id:pElement
	lst *list.List               //type entry
	cap int                      //capacity
}
type entry struct { //list item
	key   string
	value string
}

var cache *Lru = newCache(1000)

func newCache(capacity int) *Lru {
	return &Lru{make(map[string]*list.Element), list.New(), capacity}
}
func (lru *Lru) get(key string) (val string, ok bool) {
	//only touch,no move
	ele, ok := lru.mp[key]
	if ok {
		ent := ele.Value.(entry)
		val = ent.value
	}
	return
}
func (lru *Lru) set(key, val string) {
	ele, ok := lru.mp[key]
	if ok {
		ent := ele.Value.(entry)
		ent.value = val
	}
}
func (lru *Lru) addToFront(key, val string) {
	//check if exist
	_, ok := lru.get(key)
	if ok {
		lru.moveToFront(key)
	} else {
		count := lru.lst.Len()
		if count >= lru.cap {
			//full,remove last
			lru.removeBack()
		}
		//add new to front
		ele := lru.lst.PushFront(entry{key, val})
		lru.mp[key] = ele
	}

}
func (lru *Lru) moveToFront(key string) {
	ele, ok := lru.mp[key]
	if ok {
		lru.lst.MoveToFront(ele)
	}
}
func (lru *Lru) removeBack() {
	ent := lru.lst.Remove(lru.lst.Back())
	delete(lru.mp, ent.(entry).key)
}
func (lru *Lru) remove(key string) {
	ele, ok := lru.mp[key]
	if ok {
		lru.lst.Remove(ele)
		delete(lru.mp, key)
	}
}

func (lru *Lru) Get(key string) (contents string, err error) {
	val, ok := lru.get(key)
	if !ok {
		//read from disk
		val, err = ReadFile(key)
		if err != nil { //not exist
			return
		}
		//cache it
		lru.addToFront(key, val)
	} else {
		//hit,move to front
		lru.moveToFront(key)
	}
	return val, nil
}
func (lru *Lru) Put(key, val string) (err error) {
	//write to disk any way
	err = SaveFile(key, val)
	if err != nil {
		return
	}

	//catch it
	lru.addToFront(key, val)
	//move to top of list
	return
}

//////////////////////////////////////////////////////////////////////////
//for concurrent sync
const (
	READ int = iota
	WRITE
)

type Message struct {
	Type     int //READ or WRITE
	Value    string
	Key      string
	Error    error
	RespChan chan Message //response channel
}

var requestChannel = make(chan Message)

func NewMessage(typ int) Message {
	return Message{Type: typ, RespChan: make(chan Message), Error: nil}
}
func Init(capacity int) {
	cache = newCache(capacity)
}
func StartLoop() {
	go doLoop()
}
func doLoop() { //waiting at RequestChannel
	for {
		msg := <-requestChannel
		switch msg.Type {
		case READ:
			msg.Value, msg.Error = cache.Get(msg.Key)
			fmt.Printf("msg READ key=%v map=%v, list=%v time=%v\n",
				msg.Key, len(cache.mp), cache.lst.Len(), time.Now().Unix())
		case WRITE:
			msg.Error = cache.Put(msg.Key, msg.Value)
			fmt.Printf("msg WRITE key=%v map=%v, list=%v time=%v\n",
				msg.Key, len(cache.mp), cache.lst.Len(), time.Now().Unix())
		}
		go sendToChannel(msg.RespChan, msg) //return result
	}
}
func sendToChannel(ch chan Message, msg Message) {
	ch <- msg
}
func SendMessage(msg Message) Message {
	go sendToChannel(requestChannel, msg)
	return <-msg.RespChan //waiting result
}

//
//////////////////////////////////////////////////////////////////////////
