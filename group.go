package wind

import (
	"regexp"
	"strings"
)

type Group struct {
	prefix      string
	description string
	handler     func(*Router)
}

func (group *Group) dispatch(args []string, withHelp bool) (suggestions []string, err error) {
	subrouter := &Router{
		config:    RouterConfig{WithHelp: withHelp},
		ran:       true,
		routes:    make(map[string]*Route),
		groups:    make(map[string]*Group),
		arguments: args,
	}

	if withHelp {
		subrouter.Fallback(func(ctx *Context) {
			Println(group.description + "\n")
			subrouter.displayHelp()
		})
	}

	group.handler(subrouter)
	subrouter.ran = false

	return subrouter.Dispatch()
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
