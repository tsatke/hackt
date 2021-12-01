package workspace

type Workspace struct {
	// projects are the projects that are contained within this workspace.
	// The projects don't have any relation to each other. This also means,
	// that not all projects in a workspace live in the same folder or even
	// on the same drive.
	projects []*Project
}

func (ws *Workspace) AddProject(p *Project) {
	ws.projects = append(ws.projects, p)
}
