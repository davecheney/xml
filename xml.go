package xml

type Node interface{}

type Text string

type Name struct {
	Space, Prefix, Local string
}

type Element struct {
	Name
	Children []Node
}
