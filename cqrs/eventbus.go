package cqrs


var eventListeners = make(map[EventType][]func(Event))
var commandListeners = make(map[CommandType][]func(Command))

func BroadcastEvent(e *Event){
	for _, f := range append(eventListeners[e.Type], eventListeners[""]...) {
		f(*e)
	}
}

func BroadcastCommand(c *Command){
	for _, f := range append(commandListeners[c.Type], commandListeners[""]...) {
		f(*c)
	}
}

func RegisterEventListener(f func(Event),typ ... EventType){
	for _, t := range typ {
		eventListeners[t] = append(eventListeners[t],f)
	}
}

func RegisterCommandListener(f func(Command),typ ... CommandType){
	for _, t := range typ {
		commandListeners[t] = append(commandListeners[t],f)
	}
}

