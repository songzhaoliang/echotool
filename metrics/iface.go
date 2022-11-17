package metrics

type LabelsParser interface {
	ParseToLabels() map[string]string
}
