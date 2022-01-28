package runner

type Type string

const (
	TypePnpm Type = "pnpm"
)

type Program struct {
	WorkspaceType Type
	Cwd           string
}

type Module struct {
	Name string
	Deps []string
	Path string
}

func (receiver Program) getWorkspaces() []Module {
	return []Module{}
}
