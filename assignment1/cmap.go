package main

type reduceMessage struct {
	functor ReduceFunc
	key     string
	val     int
}

type pairStrInt struct {
	s string
	i int
}

// channelMap struct to control data store
type channelMap struct {
	askChanReq    chan string
	askChanRes    chan int
	addChanReq    chan string
	addChanRes    chan bool
	reduceChanReq chan reduceMessage
	reduceChanRes chan pairStrInt
	store         map[string]int
	stopChan      chan bool
}

// NewChannelMap Return a new channelMap
func NewChannelMap() EmergingMap {
	return &channelMap{
		askChanReq:    make(chan string, ASK_BUFFER_SIZE),
		askChanRes:    make(chan int, ASK_BUFFER_SIZE),
		addChanReq:    make(chan string, ADD_BUFFER_SIZE),
		addChanRes:    make(chan bool, ADD_BUFFER_SIZE),
		reduceChanReq: make(chan reduceMessage),
		reduceChanRes: make(chan pairStrInt),
		store:         make(map[string]int),
		stopChan:      make(chan bool),
	}
}

func (cmap *channelMap) Listen() {
	for {
		select {
		case ask := <-cmap.askChanReq:
			cmap.askChanRes <- cmap.store[ask]
		case add := <-cmap.addChanReq:
			cmap.store[add]++
			cmap.addChanRes <- true
		case msg := <-cmap.reduceChanReq:
			accumStr := msg.key
			accumInt := msg.val
			// Go has iteration on maps, so iterate over each item and accum
			// I don't think go maps are ordered though
			for key, val := range cmap.store {
				accumStr, accumInt = msg.functor(accumStr, accumInt, key, val)
			}
			cmap.reduceChanRes <- pairStrInt{s: accumStr, i: accumInt}
		case <-cmap.stopChan:
			return
		}
	}
}
func (cmap *channelMap) Stop() {
	cmap.stopChan <- true
}

func (cmap *channelMap) Reduce(functor ReduceFunc, accum_str string, accum_int int) (string, int) {
	msg := reduceMessage{
		functor: functor,
		key:     accum_str,
		val:     accum_int,
	}
	cmap.reduceChanReq <- msg
	res := <-cmap.reduceChanRes
	return res.s, res.i
}

func (cmap *channelMap) AddWord(word string) {
	cmap.addChanReq <- word
	<-cmap.addChanRes
}
func (cmap *channelMap) GetCount(word string) int {
	cmap.askChanReq <- word
	return <-cmap.askChanRes
}
