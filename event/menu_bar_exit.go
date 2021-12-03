package event

type MenuBarExit struct {
	handlers []MenuBarExitHandler
}

type MenuBarExitHandler func(MenuBarExitPayload)

type MenuBarExitPayload struct {
}

func (evt *MenuBarExit) Register(handler MenuBarExitHandler) {
	evt.handlers = append(evt.handlers, handler)
}

func (evt *MenuBarExit) Trigger(payload MenuBarExitPayload) {
	for _, handler := range evt.handlers {
		go handler(payload)
	}
}
