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

func (group *Group) dispatch(args []string) (suggestions []string, err error) {
	subrouter := &Router{
		ran:       true,
		routes:    make(map[string]*Route),
		groups:    make(map[string]*Group),
		arguments: args,
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
