package model

type Note struct {
	name string
}

func (n Note) Title() string       { return n.name }
func (n Note) Description() string { return "" }
func (n Note) FilterValue() string { return n.name }
