package goclitools

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type Router struct {
	routes   map[string]func(Route) error
	fallback func(Route) error
}

func (r *Router) Add(route string, handler func(Route) error) *Router {
	if r.routes == nil {
		r.routes = make(map[string]func(Route) error)
	}

	r.routes[route] = handler
	return r
}

func (r *Router) Fallback(fallback func(Route) error) *Router {
	r.fallback = fallback
	return r
}

func (r *Router) runFallBack() error {
	if r.fallback != nil {
		return r.fallback(Route{})
	}

	return nil
}

func (r *Router) Run() error {
	if len(os.Args) < 2 {
		return r.runFallBack()
	}

	handler, ok := r.routes[os.Args[1]]

	if ok {
		return handler(getRoute(os.Args[2:]))
	}

	if r.fallback != nil {
		return r.runFallBack()
	}

	return errors.New("undefined option : " + os.Args[1])
}

func NewRouter() *Router {
	return &Router{}
}

func getRoute(args []string) Route {
	var r = Route{}

	flag := regexp.MustCompile(`^-[A-z0-9\-_]+$`)
	opt := regexp.MustCompile(`^--[A-z0-9\-_]+=[A-z0-9\-_]+$`)

	for _, i := range args {
		switch true {
		case flag.MatchString(i):
			r.flags = append(r.flags, i)
		case opt.MatchString(i):
			if r.options == nil {
				r.options = make(map[string]string)
			}
			kv := strings.SplitN(i, "=", 2)
			r.options[kv[0]] = kv[1]
		default:
			r.args = append(r.args, i)
		}
	}

	return r
}

type Route struct {
	flags   []string
	args    []string
	options map[string]string
}

func (r *Route) HasFlag(f string) bool {
	for _, i := range r.flags {
		if i == f {
			return true
		}
	}

	return false
}

func (r *Route) HasOption(opt string) bool {
	_, ok := r.options[opt]

	return ok
}

func (r *Route) GetOption(opt string) string {
	v, ok := r.options[opt]

	if ok {
		return v
	}

	return ""
}

func (r *Route) GetArgs() []string {
	return r.args
}

func (r *Route) GetArg(index int) string {
	if len(r.args) < index {
		return ""
	}

	return r.args[index]
}
