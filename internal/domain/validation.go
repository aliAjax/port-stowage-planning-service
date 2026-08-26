package domain

import "fmt"

func (c Container) Validate() error {
	if c.ID == "" || c.VoyageID == "" {
		return fmt.Errorf("container id and voyage_id are required")
	}
	if c.WeightKg <= 0 || c.WeightKg > 50000 {
		return fmt.Errorf("container %s weight out of range", c.ID)
	}
	if c.ISOSize != "20" && c.ISOSize != "40" && c.ISOSize != "45" {
		return fmt.Errorf("container %s unsupported size", c.ID)
	}
	if c.HazardClass != "" && (len(c.HazardClass) > 4 || c.HazardClass == "0") {
		return fmt.Errorf("container %s invalid hazard class", c.ID)
	}
	return nil
}

func (v Voyage) Validate() error {
	if v.ID == "" || v.VesselID == "" || v.PortID == "" {
		return fmt.Errorf("voyage references are required")
	}
	if !v.ETD.After(v.ETA) {
		return fmt.Errorf("etd must be after eta")
	}
	return nil
}
func (s Slot) Validate() error {
	if s.Bay <= 0 || s.Row <= 0 || s.Tier <= 0 {
		return fmt.Errorf("invalid slot coordinates")
	}
	return nil
}
