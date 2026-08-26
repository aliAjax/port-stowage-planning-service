package repository

import (
	"context"
	"fmt"
	"github.com/example/port-stowage-planner/internal/domain"
	"sync"
)

type Store struct {
	mu         sync.RWMutex
	ports      map[string]domain.Port
	vessels    map[string]domain.Vessel
	voyages    map[string]domain.Voyage
	containers map[string]domain.Container
	cranes     map[string]domain.QuayCrane
	plans      map[string]domain.Plan
	acks       []domain.Ack
}

func New() *Store {
	return &Store{ports: map[string]domain.Port{}, vessels: map[string]domain.Vessel{}, voyages: map[string]domain.Voyage{}, containers: map[string]domain.Container{}, cranes: map[string]domain.QuayCrane{}, plans: map[string]domain.Plan{}}
}
func (s *Store) PutPort(_ context.Context, v domain.Port) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ports[v.ID] = v
}
func (s *Store) PutVessel(_ context.Context, v domain.Vessel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vessels[v.ID] = v
}
func (s *Store) PutVoyage(_ context.Context, v domain.Voyage) error {
	if err := v.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.voyages[v.ID] = v
	return nil
}
func (s *Store) PutContainer(_ context.Context, c domain.Container) error {
	if err := c.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.containers[c.ID] = c
	return nil
}
func (s *Store) PutCrane(_ context.Context, c domain.QuayCrane) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cranes[c.ID] = c
}
func (s *Store) Voyage(_ context.Context, id string) (domain.Voyage, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.voyages[id]
	return v, ok
}
func (s *Store) Containers(_ context.Context, voyage string) []domain.Container {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Container{}
	for _, c := range s.containers {
		if c.VoyageID == voyage {
			out = append(out, c)
		}
	}
	return out
}
func (s *Store) Cranes(_ context.Context, port string) []domain.QuayCrane {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.QuayCrane{}
	for _, c := range s.cranes {
		if c.PortID == port {
			out = append(out, c)
		}
	}
	return out
}
func (s *Store) SavePlan(_ context.Context, p domain.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if old, ok := s.plans[p.ID]; ok && old.Revision != p.Revision-1 {
		return fmt.Errorf("optimistic lock conflict")
	}
	s.plans[p.ID] = p
	return nil
}
func (s *Store) Plan(_ context.Context, id string) (domain.Plan, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.plans[id]
	return p, ok
}
func (s *Store) ListPlans(_ context.Context) []domain.Plan {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Plan{}
	for _, p := range s.plans {
		out = append(out, p)
	}
	return out
}
func (s *Store) AddAck(_ context.Context, a domain.Ack) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.acks = append(s.acks, a)
}
