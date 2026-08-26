package enterprise

// Catalog contains deterministic reference data and pure helpers used by
// operational integrations. Keeping these helpers side-effect free makes
// planning reproducible and easy to audit.
type Code struct {
	Value  string
	Label  string
	Active bool
}

func HazardCodes() []Code {
	return []Code{{"1", "Explosives", true}, {"2", "Gases", true}, {"3", "Flammable liquids", true}, {"4", "Flammable solids", true}, {"5", "Oxidizers", true}, {"6", "Toxic", true}, {"8", "Corrosive", true}, {"9", "Miscellaneous", true}}
}
func SizeAllowed(size string) bool { return size == "20" || size == "40" || size == "45" }
func DestinationRank(destination string) int {
	ranks := map[string]int{"TYO": 1, "OSA": 2, "SEL": 3, "SIN": 4}
	if n, ok := ranks[destination]; ok {
		return n
	}
	return 99
}
func SlotKey(bay, row, tier int) string {
	return string(rune('A'+bay%26)) + "-" + itoa(row) + "-" + itoa(tier)
}
func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	out := ""
	for v > 0 {
		out = string(rune('0'+v%10)) + out
		v /= 10
	}
	return out
}
func NormalizeHazard(v string) string {
	if v == "" {
		return ""
	}
	if len(v) > 1 && v[0] == '0' {
		return v[1:]
	}
	return v
}
