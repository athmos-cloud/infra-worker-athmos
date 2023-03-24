package kubernetes

func labelsToString(labelsEntry map[string]string) string {
	labels := ""
	for key, value := range labelsEntry {
		labels += key + "=" + value + "\n"
	}
	return labels
}
