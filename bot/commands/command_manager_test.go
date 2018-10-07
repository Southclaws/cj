package commands

import (
	"fmt"
	"testing"
	"time"
)

func Test_percentageToClock(t *testing.T) {
	tests := []struct {
		since      time.Duration
		cooldown   time.Duration
		wantResult string
	}{
		{time.Second, time.Minute, "🕐"},
		{time.Second * 15, time.Minute, "🕒"},
		{time.Second * 30, time.Minute, "🕕"},
		{time.Second * 59, time.Minute, "🕛"},
	}
	for ii, tt := range tests {
		t.Run(fmt.Sprint(ii), func(t *testing.T) {
			if gotResult := pcd(tt.since, tt.cooldown); gotResult != tt.wantResult {
				t.Errorf("percentageToClock() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
