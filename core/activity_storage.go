package core

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type ActivityStorage interface {
	Record(action string, params map[string]interface{}, result interface{}) error
	GetAllRecords() ([]common.ActivityRecord, error)
}
