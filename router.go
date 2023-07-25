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
	groups    map[string]func(*Router)
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
// there is a rule here, all optional arguments must be at the en after the required arguments and after the options
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
func (router *Router) Group(prefix string, handler func(*Router)) *Router {
	prefix = strings.TrimSpace(prefix)

	if !validPrefix(prefix) {
		router.err = router.error("router error : invalid group prefix [%s] , it should match [A-z0-9\\-\\_]", prefix)
	} else {
		router.groups[prefix] = handler
	}

	return router
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
func (router *Router) Dispatch() error {
	if router.err != nil {
		return router.err
	}

	if router.ran {
		return nil
	}
	router.ran = true

	if len(router.arguments) < 1 {
		router.fallback(getContext(router.arguments, []string{}))
		return nil
	}

	g, ok := router.groups[router.arguments[0]]

	if ok {
		return router.dispatchGroup(g)
	}

	route, ok := router.routes[router.arguments[0]]

	if ok {
		return router.dispatchHandler(route)
	}

	if router.fallback != nil {
		router.fallback(getContext(router.arguments, []string{}))
		return nil
	}

	return RouteErr{message: "undefined command : " + router.arguments[0]}
}

func (router *Router) dispatchGroup(group func(*Router)) error {
	subrouter := &Router{
		ran:       true,
		routes:    make(map[string]*Route),
		groups:    make(map[string]func(*Router)),
		fallback:  func(*Context) {},
		arguments: router.arguments[1:],
	}
	group(subrouter)
	subrouter.ran = false
	return subrouter.Dispatch()
}

func (r *Router) dispatchHandler(route *Route) error {
	ctx := getContext(r.arguments[1:], route.vars)
	err := route.match(ctx)

	if err != nil {
		return err
	}

	route.handler(ctx)
	return nil
}

func (router *Router) error(err string, args ...any) error {
	return RouterErr{
		message: fmt.Sprintf(err, args...),
	}
}

func validPrefix(p string) bool {
	rx := regexp.MustCompile(`^[A-z0-9\:]+$`)
	return rx.MatchString(p)
}

func formatFields(f string) string {
	fs := strings.Fields(f)
	var fields string

	for _, v := range fs {
		fields += " " + regexp.MustCompile(`^-{1}[A-z]+$`).ReplaceAllStringFunc(v, func(s string) string {
			tmp := strings.TrimLeft(s, "-")
			s = ""
			for _, el := range tmp {
				s += " -" + string(el)
			}

			return strings.TrimSpace(s)
		})
	}

	return strings.TrimSpace(fields)
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
func (router *Router) Test(cmd string) error {
	router.arguments = strings.Fields(cmd)
	router.ran = false
	router.err = nil
	return router.Dispatch()
}
