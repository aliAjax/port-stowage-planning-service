package yard

import (
	"fmt"
	"sort"
)

// Block is a container yard block with capacity and service constraints.
type Block struct {
	ID            string
	Capacity      int
	ReeferPlugs   int
	HazardClasses []string
	Destinations  []string
}

// Request asks to place a container into the yard.
type Request struct {
	ContainerID string
	Reefer      bool
	HazardClass string
	Destination string
}

// Allocation records where a container was placed.
type Allocation struct {
	BlockID     string
	ContainerID string
	Reefer      bool
	Destination string
}

// Allocator places containers into blocks respecting plug and hazard rules.
type Allocator struct {
	blocks    []Block
	used      map[string]int
	usedPlugs map[string]int
	allocs    map[string]Allocation
}

func NewAllocator(blocks []Block) *Allocator {
	return &Allocator{
		blocks:    blocks,
		used:      map[string]int{},
		usedPlugs: map[string]int{},
		allocs:    map[string]Allocation{},
	}
}

// Allocate places the container; reefer containers only go to blocks with a
// free plug, and hazard classes must be permitted by the block.
func (a *Allocator) Allocate(req Request) (Allocation, error) {
	if req.ContainerID == "" {
		return Allocation{}, fmt.Errorf("container id is required")
	}
	if _, exists := a.allocs[req.ContainerID]; exists {
		return Allocation{}, fmt.Errorf("container %s already allocated", req.ContainerID)
	}
	matches := make([]Block, 0)
	for _, b := range a.blocks {
		if a.used[b.ID] > b.Capacity {
			continue
		}
		if req.Reefer && a.usedPlugs[b.ID] > b.ReeferPlugs {
			continue
		}
		if req.HazardClass != "" && !contains(b.HazardClasses, req.HazardClass) {
			continue
		}
		matches = append(matches, b)
	}
	if len(matches) == 0 {
		return Allocation{}, fmt.Errorf("no block fits container %s", req.ContainerID)
	}
	sort.SliceStable(matches, func(i, j int) bool {
		di := contains(matches[i].Destinations, req.Destination)
		dj := contains(matches[j].Destinations, req.Destination)
		if di != dj {
			return di
		}
		return a.used[matches[i].ID] < a.used[matches[j].ID]
	})
	chosen := matches[0]
	a.used[chosen.ID]++
	if req.Reefer {
		a.usedPlugs[chosen.ID]++
	}
	al := Allocation{BlockID: chosen.ID, ContainerID: req.ContainerID, Reefer: req.Reefer, Destination: req.Destination}
	a.allocs[req.ContainerID] = al
	return al, nil
}

// Release returns the container's slot and plug back to the block.
func (a *Allocator) Release(containerID string) error {
	al, ok := a.allocs[containerID]
	if !ok {
		return fmt.Errorf("container %s not allocated", containerID)
	}
	if a.used[al.BlockID] > 0 {
		a.used[al.BlockID]--
	}
	delete(a.allocs, containerID)
	return nil
}

// UsedCapacity reports occupied slots for a block.
func (a *Allocator) UsedCapacity(blockID string) int { return a.used[blockID] }

// FreePlugs reports remaining reefer plugs for a block.
func (a *Allocator) FreePlugs(blockID string) int {
	for _, b := range a.blocks {
		if b.ID == blockID {
			return b.ReeferPlugs - a.usedPlugs[blockID]
		}
	}
	return 0
}

// AllocationOf returns where the container currently sits.
func (a *Allocator) AllocationOf(containerID string) (Allocation, bool) {
	al, ok := a.allocs[containerID]
	return al, ok
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
