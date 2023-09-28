package prompts

import (
	"regexp"
	"strings"
)

type Group struct {
	prefix      string
	description string
	handler     *Router
}

func newGroup(cfg RouterConfig) *Group {
	return &Group{
		handler: &Router{
			config: cfg,
			routes: make(map[string]*Route),
			groups: make(map[string]*Group),
		},
	}
}

// Add a new route to the group,
// schema is the command name followed by its arguments and flags .
//
// for example :
//
//	group = router.Group("copy")
//	group.Add("file",handler)
//
// you can run this command using
//
//	$ <PRORAM_NAME> copy file
//
// this function has the same usage as *prompts.Router.Add()
func (group *Group) Add(schema string, handler func(*Context)) *Route {
	return group.handler.Add(schema, handler)
}

// Adds nested group to the current group ,
// prefix will be used to differentiate the group,
//
// this function has the same usage as *prompts.Router.Group():
//
//	group1 := router.Group("do")
//	group2 := group.Group("print")
//	group2.Add("hello",handler})
//
// the hello command can be invoked like
//
//	$ <PROGRAM_NAME> do print hello
func (group *Group) Group(prefix string) *Group {
	return group.handler.Group(prefix)
}

// Adds a fallback route to the group ,
//
// handler will be invoked if no command match the current program args
func (group *Group) Fallback(fallback func(*Context)) *Group {
	group.handler.Fallback(fallback)
	return group
}

func (group *Group) dispatch(args []string) (suggestions []string, err error) {
	group.handler.arguments = args
	return group.handler.Dispatch()
}

func (group *Group) parsePrefix(prefix string) error {
	prefix = strings.TrimSpace(prefix)

	if !regexp.MustCompile(`^[A-z0-9\:]+$`).MatchString(prefix) {
		return RouterErr{"router error : invalid group prefix [" + prefix + "] , it should match [A-z0-9\\-\\_]"}
	}

	group.prefix = prefix
	return nil
}

func (group *Group) Description(d string) *Group {
	group.description = d
	return group
}
