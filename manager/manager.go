package manager

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GameType interface {
	Update() error
	Draw(screen *ebiten.Image)
	Init()
}

type Scene struct {
	Name        string
	Update      func(func()) error
	Draw        func(screen *ebiten.Image)
	Init        func(func())
	UpdateReady bool
	DrawReady   bool
}

type Manager struct {
	Scenes map[string]*Scene
	Data   map[string]interface{}
	Active *Scene
}

func MakeManager(scenes map[string]*Scene, start string) *Manager {
	m := &Manager{
		Scenes: scenes,
	}
	m.Active = m.Scenes[start]
	m.Active.Init(m.Active.SetUpdateReady)
	return m
}

func (s *Scene) SetDrawReady() {
	if !s.DrawReady {
		s.DrawReady = true
	}
}

func (s *Scene) SetUpdateReady() {
	if !s.UpdateReady {
		s.UpdateReady = true
	}
}

func (m *Manager) SwitchTo(name string) {
	m.Active = m.Scenes[name]
	m.Active.DrawReady = false
	m.Active.UpdateReady = false
	m.Active.Init(m.Active.SetUpdateReady)
}

func (m *Manager) Update() error {
	if m.Active.UpdateReady {
		return m.Active.Update(m.Active.SetDrawReady)
	}
	return nil
}

func (m *Manager) Draw(screen *ebiten.Image) {
	if m.Active.DrawReady {
		m.Active.Draw(screen)
	}
}

func (m *Manager) Init() {
	m.Active.Init(m.Active.SetUpdateReady)
}
