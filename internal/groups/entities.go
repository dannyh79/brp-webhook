package groups

type Group struct {
	Id string
}

func NewGroup(id string) *Group {
	return &Group{id}
}
