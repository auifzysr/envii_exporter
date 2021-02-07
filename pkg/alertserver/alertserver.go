package alertserver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type alertHandleEntry struct {
	h       func()
	pattern string
}

type AlertMux struct {
	handlers map[string]alertHandleEntry
	httpMux  *http.ServeMux
}

func NewAlertMux() *AlertMux {
	return &AlertMux{}
}

func (mux *AlertMux) AlertHandle(pattern string, handler func()) {
	if pattern == "" {
		panic("AlertMux: invalid pattern")
	}
	if handler == nil {
		panic("AlertMux: nil handler")
	}
	if _, exist := mux.handlers[pattern]; exist {
		panic("AlertMux: multiple registrations for " + pattern)
	}
	if mux.handlers == nil {
		mux.handlers = make(map[string]alertHandleEntry)
	}
	e := alertHandleEntry{
		h:       handler,
		pattern: pattern,
	}
	mux.handlers[pattern] = e
}

type AlertHandler interface {
	AlertHandle(string, func())
	http.Handler
}

func ListenAndServe(addr string, handler AlertHandler) error {
	mux := http.NewServeMux()
	mux.Handle("/", handler)

	return http.ListenAndServe(addr, mux)
}

func (mux *AlertMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(mux.handlers) < 1 {
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("body: ", string(body))

	type _handler struct {
		State string `json:"state"`
		Tags  struct {
			Threshold string `json:"threshold"`
		} `json:"tags"`
	}
	_body := &_handler{}
	err = json.Unmarshal([]byte(body), _body)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, entry := range mux.handlers {
		log.Println(entry.pattern)
		_entry := &_handler{}
		err := json.Unmarshal([]byte(entry.pattern), _entry)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("%+v", _entry.Tags)
		if _body.Tags.Threshold == _entry.Tags.Threshold {
			entry.h()
			return
		}
	}
}
