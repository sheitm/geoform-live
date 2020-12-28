package competitions

import (
	"testing"
	"time"
)

func Test_getElapsedTimeInfo(t *testing.T) {
	elapsed := 3600 + (5 * 60) + 29
	elapsedDuration := time.Duration(elapsed) * time.Second
	type args struct {
		elapsed time.Duration
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 string
	}{
		{"1", args{elapsedDuration }, elapsed, "1:05:29"},
		{"2", args{time.Duration(53*60) * time.Second}, 53*60, "0:53:00"},
		{"2", args{time.Duration(1) * time.Second}, 1, "0:00:01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getElapsedTimeInfo(tt.args.elapsed)
			if got != tt.want {
				t.Errorf("getElapsedTimeInfo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getElapsedTimeInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}