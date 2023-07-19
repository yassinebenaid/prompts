package goclitools

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type Router struct {
	routes   map[string]func(Route) error
	groups   map[string]func(*Router) *Router
	fallback func(Route) error
	Args     []string
}

func (r *Router) Add(route string, handler func(Route) error) *Router {
	if r.routes == nil {
		r.routes = make(map[string]func(Route) error)
	}

	r.routes[route] = handler
	return r
}

func (r *Router) Group(prefix string, group func(*Router) *Router) *Router {
	if r.groups == nil {
		r.groups = make(map[string]func(*Router) *Router)
	}
	r.groups[prefix] = group
	return r
}

func (r *Router) Fallback(fallback func(Route) error) *Router {
	r.fallback = fallback
	return r
}

func (r *Router) runFallBack() error {
	if r.fallback != nil {
		return r.fallback(getRoute(r.Args))
	}

	return nil
}

func (r *Router) Run() error {
	if len(os.Args) < 2 {
		return r.runFallBack()
	}

	g, ok := r.groups[r.Args[0]]

	if ok {
		r := g(&Router{Args: r.Args[1:]})

		if r != nil {
			return r.Run()
		}

		return nil
	}

	handler, ok := r.routes[r.Args[0]]

	if ok {
		return handler(getRoute(r.Args[1:]))
	}

	if r.fallback != nil {
		return r.runFallBack()
	}

	return errors.New("undefined option : " + r.Args[0])
}

func NewRouter() *Router {
	return &Router{
		Args: os.Args[1:],
	}
}

func getRoute(Args []string) Route {
	var r = Route{}

	flag := regexp.MustCompile(`^-[A-z0-9\-_]+$`)
	opt := regexp.MustCompile(`^--[A-z0-9\-_]+=[A-z0-9\-_]+$`)

	for _, i := range Args {
		switch true {
		case flag.MatchString(i):
			r.Flags = append(r.Flags, i)
		case opt.MatchString(i):
			if r.Options == nil {
				r.Options = make(map[string]string)
			}
			kv := strings.SplitN(i, "=", 2)
			r.Options[kv[0]] = kv[1]
		default:
			r.Args = append(r.Args, i)
		}
	}

	return r
}

type Route struct {
	Flags   []string
	Args    []string
	Options map[string]string
}

func (r *Route) HasFlag(f string) bool {
	for _, i := range r.Flags {
		if i == f {
			return true
		}
	}

	return false
}

func (r *Route) HasOption(opt string) bool {
	_, ok := r.Options[opt]

	return ok
}

func (r *Route) GetOption(opt string) string {
	v, ok := r.Options[opt]

	if ok {
		return v
	}

	return ""
}

func (r *Route) GetArg(index int) string {

	if len(r.Args)-1 < index {
		return ""
	}

	return r.Args[index]
}
