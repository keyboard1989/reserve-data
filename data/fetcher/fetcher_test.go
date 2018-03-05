package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
	"testing"
)

func TestUnchangedFunc(t *testing.T) {
	// test different len
	a1 := map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "mined", nil,
		},
	}
	b1 := map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "mined", nil,
		},
		common.ActivityID{2, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "mined", nil,
		},
	}
	if unchanged(a1, b1) != false {
		t.Fatalf("Expected unchanged() to return false, got true")
	}
	// test different id
	a1 = map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "mined", nil,
		},
	}
	b1 = map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{2, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "mined", nil,
		},
	}
	if unchanged(a1, b1) != false {
		t.Fatalf("Expected unchanged() to return false, got true")
	}
	// test different exchange status
	a1 = map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"", "0x123", 0, "mined", nil,
		},
	}
	b1 = map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "mined", nil,
		},
	}
	if unchanged(a1, b1) != false {
		t.Fatalf("Expected unchanged() to return false, got true")
	}
	// test different mining status
	a1 = map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "mined", nil,
		},
	}
	b1 = map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "", nil,
		},
	}
	if unchanged(a1, b1) != false {
		t.Fatalf("Expected unchanged() to return false, got true")
	}
	// test different tx
	a1 = map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "mined", nil,
		},
	}
	b1 = map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"done", "0x124", 0, "mined", nil,
		},
	}
	if unchanged(a1, b1) != false {
		t.Fatalf("Expected unchanged() to return false, got true")
	}
	// test identical statuses
	a1 = map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "mined", nil,
		},
	}
	b1 = map[common.ActivityID]common.ActivityStatus{
		common.ActivityID{1, "1"}: common.ActivityStatus{
			"done", "0x123", 0, "mined", nil,
		},
	}
	if unchanged(a1, b1) != true {
		t.Fatalf("Expected unchanged() to return true, got false")
	}
}
