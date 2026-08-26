package cargo

import "fmt"

// Dimensions are physical container dimensions in millimetres.
type Dimensions struct {
	LengthMM int
	WidthMM  int
	HeightMM int
}

// SlotDimensions are the physical clearances of a cell slot.
type SlotDimensions struct {
	LengthMM int
	WidthMM  int
	HeightMM int
}

// OOGResult reports how a container exceeds its slot.
type OOGResult struct {
	IsOOG          bool
	Directions     []string
	OverLengthMM   int
	OverWidthMM    int
	OverHeightMM   int
}

// Standard sizes used when slot dimensions are not provided.
var StandardDims = map[string]Dimensions{
	"20": {LengthMM: 6058, WidthMM: 2438, HeightMM: 2591},
	"40": {LengthMM: 12192, WidthMM: 2438, HeightMM: 2591},
	"45": {LengthMM: 13716, WidthMM: 2438, HeightMM: 2896},
}

// Classify decides whether the container is out of gauge for the slot.
func Classify(iso string, dim Dimensions, slot SlotDimensions) OOGResult {
	res := OOGResult{}
	if slot.LengthMM == 0 && slot.WidthMM == 0 && slot.HeightMM == 0 {
		if base, ok := StandardDims[iso]; ok {
			slot = SlotDimensions{LengthMM: base.LengthMM, WidthMM: base.WidthMM, HeightMM: base.HeightMM}
		} else {
			return res
		}
	}
	if dim.LengthMM > slot.LengthMM {
		res.Directions = append(res.Directions, "length")
		res.OverLengthMM = dim.LengthMM - slot.LengthMM
	}
	if dim.WidthMM > slot.WidthMM {
		res.Directions = append(res.Directions, "width")
		res.OverWidthMM = dim.WidthMM - slot.WidthMM
	}
	if dim.HeightMM > slot.HeightMM {
		res.Directions = append(res.Directions, "height")
		res.OverHeightMM = dim.HeightMM - slot.HeightMM
	}
	res.IsOOG = len(res.Directions) > 0
	return res
}

// Stowable validates that an OOG container can be placed in the given
// position. OOG cargo must sit on deck and never under a hatch cover.
func Stowable(res OOGResult, onDeck, underHatch bool) error {
	if !res.IsOOG {
		return nil
	}
	if !onDeck {
		return fmt.Errorf("out-of-gauge container must be stowed on deck")
	}
	if underHatch {
		return fmt.Errorf("out-of-gauge container cannot be under a hatch cover")
	}
	if res.OverHeightMM > 1500 {
		return fmt.Errorf("overheight %dmm exceeds lashing clearance", res.OverHeightMM)
	}
	return nil
}

// OOGCount returns the number of over-gauge containers in a list.
func OOGCount(dims []Dimensions, iso []string, slot SlotDimensions) int {
	n := 0
	for i, d := range dims {
		if Classify(iso[i], d, slot).IsOOG {
			n++
		}
	}
	return n
}
