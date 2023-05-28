package test

func LabelsEquals(got map[string]string, want map[string]string) bool {
	if len(got) != len(want) {
		return false
	}
	for kg, vg := range got {
		if want[kg] != vg {
			return false
		}
	}
	return true
}
