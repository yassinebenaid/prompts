package wind

import (
	"fmt"
	"os"
	"regexp"
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
	prefix = strings.TrimSpace(prefix)
	group := &Group{}

	if !validPrefix(prefix) {
		router.err = router.error("router error : invalid group prefix [%s] , it should match [A-z0-9\\-\\_]", prefix)
	} else {
		group.handler = handler
		router.groups[prefix] = group
	}

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
		return nil, router.err
	}

	if router.ran {
		return nil, nil
	}
	router.ran = true

	if len(router.arguments) < 1 {
		if router.fallback != nil {
			router.fallback(getContext(router.arguments, []string{}))
		}

		return nil, nil
	}

	g, ok := router.groups[router.arguments[0]]

	if ok {
		return router.dispatchGroup(g)
	}

	route, ok := router.routes[router.arguments[0]]

	if ok {
		return nil, route.dispatch(router.arguments[1:])
	}

	if router.fallback != nil {
		router.fallback(getContext(router.arguments, []string{}))
		return nil, nil
	}

	return router.getSuggestions(router.arguments[0]), RouteErr{message: "undefined command : " + router.arguments[0]}
}

func (router *Router) dispatchGroup(group *Group) (suggestions []string, err error) {
	subrouter := &Router{
		ran:       true,
		routes:    make(map[string]*Route),
		groups:    make(map[string]*Group),
		arguments: router.arguments[1:],
	}
	group.handler(subrouter)
	subrouter.ran = false
	return subrouter.Dispatch()
}

func (router *Router) error(err string, args ...any) error {
	return RouterErr{
		message: fmt.Sprintf(err, args...),
	}
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

func validPrefix(p string) bool {
	rx := regexp.MustCompile(`^[A-z0-9\:]+$`)
	return rx.MatchString(p)
}

func (router *Router) isUniquePrefix(p string) bool {
	_, ok := router.routes[p]

	if ok {
		return false
	}

	_, ok = router.groups[p]

	return !ok
}

// this function helps you test the  router with hardcoded command string
//
// the cmd represents the command string you write in a terminal :
//
//	router.Add("model <path> <name>", func(ctx *wind.Context) {
//		fmt.Println(ctx.GetArg("path"))
//		fmt.Println(ctx.GetArg("name"))
//	})
//
//	err := router.Test("model foo bar")
//	if err != nil{
//		panic(err)
//	}
func (router *Router) Test(cmd string) (suggestions []string, err error) {
	router.arguments = strings.Fields(cmd)
	router.ran = false
	return router.Dispatch()
}
