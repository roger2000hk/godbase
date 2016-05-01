package godbase

type EvtFn func(...interface{})

type EvtSub interface{}

type evtSubs map[EvtSub]EvtFn

type Evt struct {
	subs evtSubs
}

func NewEvt() *Evt {
	return new(Evt).Init()
}

func (self *Evt) Init() *Evt {
	self.subs = make(evtSubs)
	return self
}

func (e *Evt) Publish(args...interface{}) {
	for _, fn := range e.subs {
		fn(args...)
	} 
}

func (e *Evt) Subscribe(k EvtSub, fn EvtFn) {
	e.subs[k] = fn
}

func (e *Evt) Unsubscribe(k EvtSub) {
	delete(e.subs, k)
}
