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
		ran:       false,
		routes:    make(map[string]*Route),
		groups:    make(map[string]func(*Router)),
		fallback:  func(*Context) {},
		arguments: os.Args[1:],
		err:       nil,
	}
}

// Add a new route ,
// schema is the command name followed by its arguments and flags .
//
// for example :
//
//	router.Add("copy",func(ctx *goclitools.Context){
//		// ...
//	})
//
// you can run this command using
//
//	$ <PRORAM_NAME> copy
//
// you can also define the flags :
//
//	router.Add("copy [-a --verbose]",func(ctx *goclitools.Context){
//		// ...
//	})
//
// to define the arguments , wrape them in a curly braces , you will use that name to retrieve them later
//
//	router.Add("copy [-a --verbose] {source} {destination}",func(ctx *goclitools.Context){
//		// ...
//	})
func (router *Router) Add(schema string, handler func(*Context)) *Router {
	route := Route{handler: handler}

	if err := route.parseSchema(schema); err != nil {
		router.err = err
	}

	router.routes[route.prefix] = &route
	return router
}

// Adds new route group to the router ,
// prefix will be used to differentiate the group,
// handler recieves a Router instance to register your sub routes
//
// this function useful to group multiple routes under a single prefix :
//
//	router.Group("copy",func(router *goclitools.Router){
//		router.Add("file",func(ctx *goclitools.Context){
//		    // ...
//	    })
//		router.Add("dir",func(ctx *goclitools.Context){
//		    // ...
//	    })
//	})
//
// these commands can be invoked like
//
//	$ <PROGRAM_NAME> copy file
//	$ <PROGRAM_NAME> copy dir
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
// handler will be invoked if no command match the current program args
func (r *Router) Fallback(fallback func(*Context)) *Router {
	r.fallback = fallback
	return r
}

// Dispatch the router, reads the process args and invoke the convenience handler
//
// if something went wrong , it returns why
//
// this function should be called lastely , after its first call, the router become useless
//
// this function is useless within groups , so it won't do anything if you run it within a group
func (r *Router) Dispatch() error {
	if r.err != nil {
		return r.err
	}

	if r.ran {
		return nil
	}
	r.ran = true

	if len(r.arguments) < 1 {
		r.fallback(getContext(r.arguments, []string{}))
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
		r.fallback(getContext(r.arguments, []string{}))
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
	ctx := getContext(r.arguments[1:], route.vars)
	route.path = strings.Join(r.arguments, " ")
	err := route.match(ctx)

	if err != nil {
		return err
	}

	route.handler(ctx)
	return nil
}

func validPrefix(p string) bool {
	rx := regexp.MustCompile(`^[A-z0-9\-]+$`)
	return rx.MatchString(p)
}
