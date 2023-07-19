package goclitools

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type Router struct {
	ran       bool
	routes    map[string]func(Route) error
	groups    map[string]func(*Router)
	fallback  func(Route) error
	arguments []string
	err       error
}

// Adds new route to the router ,
// prefix is the first value after the program name ,
// if Add called in a group , prefix is the first value after the group prefix
func (r *Router) Add(prefix string, handler func(Route) error) *Router {
	if r.routes == nil {
		r.routes = make(map[string]func(Route) error)
	}

	prefix = strings.TrimSpace(prefix)

	if !validPrefix(prefix) {
		r.err = errors.New("router error : invalid prefix [" + prefix + "] , it should match [A-z0-9\\-\\_]")
	} else {
		r.routes[prefix] = handler
	}

	return r
}

// Adds new route group to the router ,
// prefix is the first value after the program name ,
// if Group called in a group , prefix is the first value after the current group prefix
//
// this function useful to group multiple routes under a single prefix
//
// for example `$ git origin ` has many sub routes (add,remove ...)
//
// handler recieves a Router instance to register your sub routes
func (r *Router) Group(prefix string, handler func(*Router)) *Router {
	if r.groups == nil {
		r.groups = make(map[string]func(*Router))
	}

	prefix = strings.TrimSpace(prefix)

	if !validPrefix(prefix) {
		r.err = errors.New("router error : invalid group prefix [" + prefix + "] , it should match [A-z0-9\\-\\_]")
	} else {
		r.groups[prefix] = handler
	}

	return r
}

// Adds a fallback route to the router ,
//
// handler will be invoked if no route match the current option
func (r *Router) Fallback(fallback func(Route) error) *Router {
	r.fallback = fallback
	return r
}

func (r *Router) runFallBack() error {
	if r.fallback != nil {
		return r.fallback(getRoute(r.arguments))
	}

	return nil
}

// Dispatch the router, reads the process args and invoke the convenience handler
//
// if something went wrong , it returns why, or returns the error returned by the handler it self
func (r *Router) Dispatch() error {
	if r.err != nil {
		return r.err
	}

	if r.ran {
		return nil
	}
	r.ran = true

	if len(r.arguments) < 1 {
		return r.runFallBack()
	}

	g, ok := r.groups[r.arguments[0]]

	if ok {
		r2 := Router{arguments: r.arguments[1:]}
		g(&r2)
		return r2.Dispatch()
	}

	handler, ok := r.routes[r.arguments[0]]

	if ok {
		return handler(getRoute(r.arguments[1:]))
	}

	if r.fallback != nil {
		return r.runFallBack()
	}

	return errors.New("undefined option : " + r.arguments[0])
}

// create new router instance
func NewRouter() *Router {
	return &Router{
		arguments: os.Args[1:],
	}
}

func getRoute(Args []string) Route {
	var r = Route{}

	flag := regexp.MustCompile(`^-[A-z0-9\-_]+$`)
	opt := regexp.MustCompile(`^--[A-z0-9\-_]+=[A-z0-9\-_]+$`)

	for _, i := range Args {
		switch true {
		case flag.MatchString(i):
			i = strings.TrimPrefix(i, "-")
			fs := strings.Split(i, "")

			for _, f := range fs {
				r.Flags = append(r.Flags, "-"+f)
			}
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

// determine wether f flag is present or not
//
// for example HasFlag("-h")
func (r *Route) HasFlag(f string) bool {
	for _, i := range r.Flags {
		if i == f {
			return true
		}
	}

	return false
}

// determine wether opt option is present or not
//
// for example HasOption("--help")
func (r *Route) HasOption(opt string) bool {
	_, ok := r.Options[opt]

	return ok
}

// get an option value
//
// for example , in case of "--path=some/path/here" , GetOption("--path") returns "/some/path/here"
func (r *Route) GetOption(opt string) string {
	v, ok := r.Options[opt]

	if ok {
		return v
	}

	return ""
}

// Get an argument by its index , or "" if doesn't exists
//
// this functions returns all args that are not a flag , and not options
func (r *Route) GetArg(index int) string {

	if len(r.Args)-1 < index {
		return ""
	}

	return r.Args[index]
}

func validPrefix(p string) bool {
	rx := regexp.MustCompile(`^[A-z0-9\-\_]+$`)
	return rx.MatchString(p)
}
