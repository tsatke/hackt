package workspace

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNoProjectDefinitionFile Error = "no project definition file could be found (" + ProjectDefinitionFileName + ")"
)
