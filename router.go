package goclitools

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type Router struct {
	ran       bool
	routes    map[string]Route
	groups    map[string]func(*Router)
	fallback  func(*Context)
	arguments []string
	err       error
}

// Adds new route to the router ,
// prefix is the first value after the program name ,
// if Add called in a group , prefix is the first value after the group prefix
func (r *Router) Add(prefix string, route Route) *Router {
	if r.routes == nil {
		r.routes = make(map[string]Route)
	}

	prefix = strings.TrimSpace(prefix)

	if !validPrefix(prefix) {
		r.err = errors.New("router error : invalid prefix [" + prefix + "] , it should match [A-z0-9\\-\\_]")
	} else {
		if route.Handler == nil {
			r.err = errors.New("router error : handler cannot be nil for route [" + prefix + "]")
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
		if r.fallback != nil {
			r.fallback(getContext(r.arguments))
		}
		return nil
	}

	g, ok := r.groups[r.arguments[0]]

	if ok {
		r2 := Router{arguments: r.arguments[1:], ran: true}
		g(&r2)
		r2.ran = false
		return r2.Dispatch()
	}

	route, ok := r.routes[r.arguments[0]]

	if ok {
		ctx := getContext(r.arguments[1:])
		err := route.match(ctx)

		if err != nil {
			return err
		}

		route.Handler(ctx)
		return nil
	}

	if r.fallback != nil {
		r.fallback(getContext(r.arguments))
		return nil
	}

	return errors.New("undefined option : " + r.arguments[0])
}

// create new router instance
func NewRouter() *Router {
	return &Router{
		arguments: os.Args[1:],
	}
}

func validPrefix(p string) bool {
	rx := regexp.MustCompile(`^[A-z0-9\-\_]+$`)
	return rx.MatchString(p)
}
