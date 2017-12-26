package core

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type ActivityStorage interface {
	Record(
		action string,
		id common.ActivityID,
		destination string,
		params map[string]interface{},
		result map[string]interface{},
		estatus string,
		mstatus string,
		timepoint uint64) error
}
