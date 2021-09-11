package signals

var Signals = make(map[string][]func([]interface{}))

func Connect(name string, cb func([]interface{})) {
	_, ok := Signals[name]
	if !ok {
		Signals[name] = make([]func([]interface{}), 0)
	}
	Signals[name] = append(Signals[name], cb)
}

func Emit(name string, params []interface{}) {
	cbs, ok := Signals[name]
	if !ok {
		return
	}
	for _, cb := range cbs {
		cb(params)
	}
}
