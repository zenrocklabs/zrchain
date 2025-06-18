package main

import (
	"testing"
	"time"

	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
)

func TestShouldSkipEventBasedOnTimestamp(t *testing.T) {
	// Create a test oracle with a specific startup time
	startupTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	o := &Oracle{
		startupTimestamp: startupTime,
		Config: sidecartypes.Config{
			FilterEventsAfterStartup: true,
		},
		DebugMode: true,
	}

	// Test cases
	testCases := []struct {
		name          string
		blockTimeUnix int64
		expectedSkip  bool
		description   string
	}{
		{
			name:          "Event before startup should be skipped",
			blockTimeUnix: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC).Unix(),
			expectedSkip:  true,
			description:   "Event at 11:00 should be skipped when startup was at 12:00",
		},
		{
			name:          "Event after startup should not be skipped",
			blockTimeUnix: time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC).Unix(),
			expectedSkip:  false,
			description:   "Event at 13:00 should not be skipped when startup was at 12:00",
		},
		{
			name:          "Event at startup time should not be skipped",
			blockTimeUnix: startupTime.Unix(),
			expectedSkip:  false,
			description:   "Event at startup time should not be skipped",
		},
		{
			name:          "Zero block time should not be skipped",
			blockTimeUnix: 0,
			expectedSkip:  false,
			description:   "Events with zero block time should not be skipped",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := o.shouldSkipEventBasedOnTimestamp(tc.blockTimeUnix, "test-tx-sig", "test-event")
			if result != tc.expectedSkip {
				t.Errorf("%s: expected skip=%v, got skip=%v", tc.description, tc.expectedSkip, result)
			}
		})
	}
}

func TestShouldSkipEventBasedOnTimestampDisabled(t *testing.T) {
	// Create a test oracle with filtering disabled
	startupTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	o := &Oracle{
		startupTimestamp: startupTime,
		Config: sidecartypes.Config{
			FilterEventsAfterStartup: false,
		},
		DebugMode: true,
	}

	// Test that events are never skipped when filtering is disabled
	oldEventTime := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC).Unix()

	result := o.shouldSkipEventBasedOnTimestamp(oldEventTime, "test-tx-sig", "test-event")
	if result {
		t.Errorf("Event should not be skipped when FilterEventsAfterStartup is disabled")
	}
}
