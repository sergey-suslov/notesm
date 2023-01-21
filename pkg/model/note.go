package model

import "time"

type Note struct {
	name    string
	ModTime time.Time
}

func (n Note) Title() string       { return n.name }
func (n Note) Description() string { return n.ModTime.Format(time.RFC822) }
func (n Note) FilterValue() string { return n.name }
