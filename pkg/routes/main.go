package routes

import "github.com/codersidprogrammer/gonotif/pkg/model"

type appRoute struct {
	apps []model.App
	name string
}

type Routes interface {
	Register(app model.App)
	Handler()
}

func NewRoute(name string) Routes {
	return &appRoute{
		apps: []model.App{},
		name: name,
	}
}

// Register implements Routes.
func (p *appRoute) Register(app model.App) {
	p.apps = append(p.apps, app)
}

func (p *appRoute) Handler() {
	for _, app := range p.apps {
		app.Route()
	}
}
