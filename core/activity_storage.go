package core

type ActivityStorage interface {
	Record(
		action string,
		id string,
		destination string,
		params map[string]interface{},
		result map[string]interface{},
		status string,
		timepoint uint64) error
}
