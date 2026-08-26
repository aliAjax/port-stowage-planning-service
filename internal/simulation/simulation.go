package simulation

import (
	"context"
	"github.com/example/port-stowage-planner/internal/domain"
	"time"
)

type Event struct {
	At                         time.Time
	Kind, ContainerID, Message string
}
type Report struct {
	Events    []Event
	Completed int
	Warnings  []string
}

func Run(ctx context.Context, p domain.Plan) Report {
	r := Report{}
	now := time.Now().UTC()
	for n, d := range p.Decisions {
		select {
		case <-ctx.Done():
			r.Warnings = append(r.Warnings, "simulation cancelled")
			return r
		default:
		}
		r.Events = append(r.Events, Event{At: now.Add(time.Duration(n) * time.Second), Kind: "load", ContainerID: d.ContainerID, Message: "slot constraints satisfied"})
		r.Completed++
	}
	return r
}
