package event

import "code.rocketnine.space/tslocum/cview"

type UIRedraw struct {
	handlers []UIRedrawHandler
}

type UIRedrawHandler func(UIRedrawPayload)

type UIRedrawPayload struct {
	Components []cview.Primitive
}

func (evt *UIRedraw) Register(handler UIRedrawHandler) {
	evt.handlers = append(evt.handlers, handler)
}

func (evt *UIRedraw) Trigger(payload UIRedrawPayload) {
	for _, handler := range evt.handlers {
		go handler(payload)
	}
}
