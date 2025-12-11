package detector

// GetAllDetectors returns all available detectors
func GetAllDetectors() []Detector {
	return []Detector{
		&JavaDetector{},
		&PythonDetector{},
		&NodeDetector{},
	}
}
