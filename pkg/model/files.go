package model

func (m *TeaModel) createNote(name, body string) error {
	return m.fr.SaveNote(name, body)
}
