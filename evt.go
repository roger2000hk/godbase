package godbase

type EvtFn func(...interface{}) error

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

func (e *Evt) Publish(args...interface{}) error {
	for _, fn := range e.subs {
		if err := fn(args...); err != nil {
			return err
		}
	} 

	return nil
}

func (e *Evt) Subscribe(k EvtSub, fn EvtFn) {
	e.subs[k] = fn
}

func (e *Evt) Unsubscribe(k EvtSub) {
	delete(e.subs, k)
}
