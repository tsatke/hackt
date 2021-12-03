package hackt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/tsatke/hackt/event"
	"github.com/tsatke/hackt/ui"
	"github.com/tsatke/hackt/workspace"
)

type Application struct {
	log  zerolog.Logger
	view *ui.ApplicationView

	events *event.Bus

	workspace *workspace.Workspace
}

func NewApplication(log zerolog.Logger) *Application {
	events := event.NewBus(log)

	app := &Application{
		log:    log,
		events: events,
		view:   ui.NewApplicationView(log, events),
	}

	events.MenuBar.MenuBarExit.Register(func(_ event.MenuBarExitPayload) {
		app.log.Debug().
			Msg("menu bar requested exit")
		app.Stop()
	})
	return app
}

func (app Application) Run() {
	app.log.Info().
		Msg("start hackt ui")
	app.view.Run()
}

func (app Application) Stop() {
	app.view.Stop()
}

func (app Application) Wait() error {
	return <-app.view.Done()
}

// LoadProjectFromPath loads a project from the given directory path and displays it in the UI.
func (app Application) LoadProjectFromPath(path string) error {
	realpath := path
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get cwd: %w", err)
		}
		realpath = filepath.Join(cwd, path)
	}

	app.log.Info().Str("path", realpath).Msg("load project from disk")

	fs := afero.NewBasePathFs(afero.NewOsFs(), realpath)
	project, err := workspace.LoadProject(fs)
	if err != nil {
		app.log.Error().
			Err(err).
			Msg("can't load project")
		return err
	}

	return app.LoadProject(project)
}

func (app Application) LoadProject(p *workspace.Project) error {
	app.log.Info().
		Str("name", p.Name()).
		Msg("load project")

	app.events.Project.ProjectLoaded.Trigger(event.ProjectLoadedPayload{Project: p})

	return nil
}
