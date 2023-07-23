package goclitools

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type Router struct {
	ran       bool
	routes    map[string]*Route
	groups    map[string]func(*Router)
	fallback  func(*Context)
	arguments []string
	err       error
}

// create new router instance
func NewRouter() *Router {
	return &Router{
		ran:    false,
		routes: make(map[string]*Route),
		groups: make(map[string]func(*Router)),
		fallback: func(*Context) {
		},
		arguments: os.Args[1:],
		err:       nil,
	}
}

// Adds new handler to the router ,
// prefix is the first value after the program name ,
// if Add called in a group , prefix is the first value after the group prefix
func (r *Router) Add(schema string, handler func(*Context)) *Router {
	route := Route{}

	if err := route.parseSchema(schema); err != nil {
		r.err = err
	}

	return r
	// return r.AddRoute(schema, &Route{prefix: prefix, Handler: handler})
}

// Adds new route to the router ,
// prefix is the first value after the program name ,
// if Add called in a group , prefix is the first value after the group prefix
func (r *Router) AddRoute(prefix string, route *Route) *Router {
	prefix = strings.TrimSpace(prefix)

	if !validPrefix(prefix) {
		r.err = errors.New(" invalid prefix [" + prefix + "] , it should match [A-z0-9\\-\\_]")
	} else {
		if route.Handler == nil {
			r.err = errors.New("handler cannot be nil for route [" + prefix + "]")
		} else {
			route.prefix = prefix
			r.routes[prefix] = route
		}
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
func (r *Router) Fallback(fallback func(*Context)) *Router {
	r.fallback = fallback
	return r
}

// Dispatch the router, reads the process args and invoke the convenience handler
//
// if something went wrong , it returns why, or returns the error returned by the handler it self
//
// this function should be called lastely , after its first call, the router become useless
//
// if this function invoked within a group, it does nothing
// but returns early, this gives the router the chance to read all routes correctly
func (r *Router) Dispatch() error {
	if r.err != nil {
		return r.err
	}

	if r.ran {
		return nil
	}
	r.ran = true

	if len(r.arguments) < 1 {
		r.fallback(getContext(r.arguments))
		return nil
	}

	g, ok := r.groups[r.arguments[0]]

	if ok {
		return r.dispatchGroup(g)
	}

	route, ok := r.routes[r.arguments[0]]

	if ok {
		return r.dispatchHandler(route)
	}

	if r.fallback != nil {
		r.fallback(getContext(r.arguments))
		return nil
	}

	return errors.New("undefined command : " + r.arguments[0])
}

func (r *Router) dispatchGroup(g func(*Router)) error {
	r2 := &Router{
		ran:       true,
		routes:    make(map[string]*Route),
		groups:    make(map[string]func(*Router)),
		fallback:  func(*Context) {},
		arguments: r.arguments[1:],
	}
	g(r2)
	r2.ran = false
	return r2.Dispatch()
}

func (r *Router) dispatchHandler(route *Route) error {
	ctx := getContext(r.arguments[1:])
	err := route.match(ctx)

	if err != nil {
		return err
	}

	route.Handler(ctx)
	return nil
}

func validPrefix(p string) bool {
	rx := regexp.MustCompile(`^[A-z0-9\-]+$`)
	return rx.MatchString(p)
}
