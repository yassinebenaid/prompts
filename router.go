package wind

import (
	"os"
	"strings"
)

type Router struct {
	ran       bool
	routes    map[string]*Route
	groups    map[string]*Group
	fallback  func(*Context)
	arguments []string
	err       error
}

type RouterErr struct {
	message string
}

func (e RouterErr) Error() string {
	return e.message
}

// create new router instance
func NewRouter() *Router {
	return &Router{
		ran:       false,
		routes:    make(map[string]*Route),
		groups:    make(map[string]*Group),
		arguments: os.Args[1:],
		err:       nil,
	}
}

// Add a new route ,
// schema is the command name followed by its arguments and flags .
//
// for example :
//
//	router.Add("copy",handler)
//
// you can run this command using
//
//	$ <PRORAM_NAME> copy
//
// you can also define the flags :
//
//	router.Add("copy [-a --verbose]",handler)
//
// here is how to define the arguments ,  you will use those names to retrieve them later
//
//	router.Add("copy [-a --verbose] <source> <destination>",handler)
//
// you can make the arguments optional by adding "?" question mark at the end :
//
//	router.Add("copy [-a --verbose] <source> <destination?>",handler)
//
// there is a rule here, all optional arguments must be at the en after the required arguments and after the flags
//
// also , keep in mind that the order you choose in the schema , is the order will be used to invoke the command
// so this will throw an error :
//
//	router.Add("copy [-a --verbose] <source> <destination?>",handler)
//
//	$<PROGRAM_NAME> copy somesource -a  // RouteErr
//
// this throws an error because the flag used after the source, but the order is important,
//
//	$ <PROGRAM_NAME> copy -a somesource  // Works!
func (router *Router) Add(schema string, handler func(*Context)) *Route {
	route := &Route{handler: handler}

	if err := route.parseSchema(schema); err != nil {
		router.err = err
	}

	if !router.isUniquePrefix(route.prefix) {
		router.err = RouterErr{"command " + route.prefix + " is not unique!"}
	}

	router.routes[route.prefix] = route
	return route
}

// Adds new route group to the router ,
// prefix will be used to differentiate the group,
// handler recieves a Router instance to register your sub routes
//
// this function useful to group multiple routes under a single prefix :
//
//	router.Group("copy",func(router *wind.Router){
//		router.Add("file",func(ctx *wind.Context){
//		    // ...
//	    })
//		router.Add("dir",func(ctx *wind.Context){
//		    // ...
//	    })
//	})
//
// these commands can be invoked like
//
//	$ <PROGRAM_NAME> copy file
//	$ <PROGRAM_NAME> copy dir
func (router *Router) Group(prefix string, handler func(*Router)) *Group {
	group := &Group{handler: handler}

	if err := group.parsePrefix(prefix); err != nil {
		router.err = err
		return group
	}

	if !router.isUniquePrefix(group.prefix) {
		router.err = RouterErr{"command " + group.prefix + " is not unique!"}
	}

	router.groups[group.prefix] = group

	return group
}

// Adds a fallback route to the router ,
//
// handler will be invoked if no command match the current program args
func (r *Router) Fallback(fallback func(*Context)) *Router {
	r.fallback = fallback
	return r
}

// Dispatch the router, reads the process args and invoke the handler accordingly
//
// if something went wrong , it returns why , and if no command match , suggestions will contain the most similar commands:
//
//		router.Add("print",handler)
//		router.Add("sprint",handler)
//
//		sugestions,err := router.Add("sprint",handler)
//
//		$ <PROGRAM_NAME> pr
//
//	 // err == RouteErr : undefined command pr
//	 // suggestions == [print, sprint]
//
// this function should be called after all commands registered .
// also it won't do anything if you run it within a group
func (router *Router) Dispatch() (suggestions []string, err error) {
	if router.err != nil {
		return []string{}, router.err
	}

	if router.ran {
		return []string{}, nil
	}
	router.ran = true

	if len(router.arguments) < 1 {
		if router.fallback != nil {
			router.fallback(getContext(router.arguments, []string{}))
		}

		return []string{}, nil
	}

	group, ok := router.groups[router.arguments[0]]

	if ok {
		return group.dispatch(router.arguments[1:])
	}

	route, ok := router.routes[router.arguments[0]]

	if ok {
		return []string{}, route.dispatch(router.arguments[1:])
	}

	if router.fallback != nil {
		router.fallback(getContext(router.arguments, []string{}))
		return []string{}, nil
	}

	return router.getSuggestions(router.arguments[0]),
		RouteErr{message: "undefined command : " + router.arguments[0]}
}

func (router *Router) getSuggestions(cmd string) []string {
	var sgs []string

	for key := range router.routes {
		if strings.Contains(key, cmd) {
			sgs = append(sgs, key)
		}
	}

	return sgs
}

func (router *Router) isUniquePrefix(p string) bool {
	_, ok := router.routes[p]

	if ok {
		return false
	}

	_, ok = router.groups[p]

	return !ok
}
