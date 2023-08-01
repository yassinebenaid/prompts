package wind

type Group struct {
	description string
	handler     func(*Router)
}
