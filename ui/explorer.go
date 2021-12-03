package ui

import (
	"path/filepath"
	"sort"
	"strings"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/tsatke/hackt/event"
	"github.com/tsatke/hackt/workspace"
)

var _ cview.Primitive = (*Explorer)(nil)

type Explorer struct {
	log    zerolog.Logger
	events *event.Bus

	tree *cview.TreeView
	root *cview.TreeNode
}

func NewExplorer(log zerolog.Logger, events *event.Bus) *Explorer {
	tree := cview.NewTreeView()
	tree.SetBorder(true)
	tree.SetTitle("Workspace")
	root := cview.NewTreeNode("Projects")
	tree.SetRoot(root)

	explorer := &Explorer{
		log:    log,
		events: events,

		tree: tree,
		root: root,
	}
	events.Project.ProjectLoaded.Register(explorer.projectLoaded)

	return explorer
}

func (e *Explorer) projectLoaded(event event.ProjectLoadedPayload) {
	project := event.Project
	projectNode, err := e.createTreeNodeFromProject(project)
	if err != nil {
		e.log.Error().
			Err(err).
			Str("project", project.Name()).
			Msg("create project tree")
		return
	}

	e.log.Debug().
		Str("name", project.Name()).
		Msg("add project to explorer")
	e.root.AddChild(projectNode)
}

func (e *Explorer) createTreeNodeFromProject(p *workspace.Project) (*cview.TreeNode, error) {
	node := cview.NewTreeNode(p.Name())

	var addRecursive func(afero.Fs, string, *cview.TreeNode)
	depth := 0
	addRecursive = func(root afero.Fs, path string, parent *cview.TreeNode) {
		depth++
		defer func() { depth-- }()

		if depth > 5 { // introduce this to prevent an infinite hang
			e.log.Debug().
				Str("path", path).
				Msg("depth limit exceeded")
			return
		}

		projectDir, err := root.Open(path)
		if err != nil {
			e.log.Error().
				Err(err).
				Msg("open")
			return
		}
		defer func() { _ = projectDir.Close() }()

		res, err := projectDir.Readdir(0)

		sort.Slice(res, func(i, j int) bool {
			if res[i].IsDir() && !res[j].IsDir() {
				return true
			} else if !res[i].IsDir() && res[j].IsDir() {
				return false
			}

			return strings.Compare(res[i].Name(), res[j].Name()) == -1
		})

		for _, info := range res {
			if strings.HasPrefix(info.Name(), ".") {
				continue // skip dot files
			}

			child := cview.NewTreeNode(info.Name())
			if info.IsDir() {
				e.makeCollapsible(child)
				addRecursive(root, filepath.Join(path, info.Name()), child)
			} else {
				e.makeOpenable(child, p, filepath.Join(path, info.Name()))
				child.SetColor(tcell.ColorCadetBlue)
			}
			parent.AddChild(child)
		}
	}

	projectFs := p.Fs()
	addRecursive(projectFs, ".", node)
	node.ExpandAll()

	return node, nil
}

func (e *Explorer) makeOpenable(node *cview.TreeNode, p *workspace.Project, path string) {
	node.SetSelectedFunc(func() {
		e.events.Project.ProjectFileOpen.Trigger(event.ProjectFileOpenPayload{
			Project: p,
			Path:    path,
		})
	})
}

func (e *Explorer) makeCollapsible(node *cview.TreeNode) {
	node.SetSelectedFunc(func() {
		if node.IsExpanded() {
			node.Collapse()
		} else {
			node.Expand()
		}
	})
}

// implement cview.Primitive, but delegate everything to the actual underlying TreeView

func (e Explorer) Draw(screen tcell.Screen) {
	e.tree.Draw(screen)
}

func (e Explorer) GetRect() (int, int, int, int) {
	return e.tree.GetRect()
}

func (e Explorer) SetRect(x, y, width, height int) {
	e.tree.SetRect(x, y, width, height)
}

func (e Explorer) GetVisible() bool {
	return e.tree.GetVisible()
}

func (e Explorer) SetVisible(v bool) {
	e.tree.SetVisible(v)
}

func (e Explorer) InputHandler() func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
	return e.tree.InputHandler()
}

func (e Explorer) Focus(delegate func(p cview.Primitive)) {
	e.tree.Focus(delegate)
}

func (e Explorer) Blur() {
	e.tree.Blur()
}

func (e Explorer) GetFocusable() cview.Focusable {
	return e.tree.GetFocusable()
}

func (e Explorer) MouseHandler() func(action cview.MouseAction, event *tcell.EventMouse, setFocus func(p cview.Primitive)) (consumed bool, capture cview.Primitive) {
	return e.tree.MouseHandler()
}
