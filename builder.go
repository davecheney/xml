package xml

import "strings"

type builder struct {
	parent             *builder
	name               string
	content            []Node
	declaredNamespaces map[string]string
}

func (b *builder) doElementStart(name string) *builder {
	return &builder{
		parent: b,
		name:   name,
	}
}

func (b *builder) build() *Element {
	return &Element{
		Name:     b.qname(),
		Children: b.content,
	}
}

func (b *builder) doCharacters(text string) {
	t := Text(text)
	b.addContent(&t)
}

func (b *builder) addContent(n Node) {
	b.content = append(b.content, n)
}

func (b *builder) qname() Name {
	if index := strings.Index(b.name, ":"); index > -1 {
		prefix := b.name[0:index]
		local := b.name[index+1:]
		space := b.declaredNamespaceForPrefix(prefix)
		return Name{Space: space, Prefix: prefix, Local: local}
	} else {
		// TODO, should call declaredNamespace to get the default namespace
		space := b.declaredNamespaceForPrefix("")
		return Name{Space: space, Local: b.name}
	}
}

func (b *builder) declaredNamespaceForPrefix(prefix string) string {
	if ns, ok := b.declaredNamespaces[prefix]; ok {
		return ns
	} else if b.parent != nil {
		return b.parent.declaredNamespaceForPrefix(prefix)
	} else {
		return ""
	}
}
