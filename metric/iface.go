package metric

type LabelsParser interface {
	ParseToLabels() map[string]string
}
