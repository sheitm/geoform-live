package competitions

import (
	"sync"
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

func Test_impl_getAll(t *testing.T) {
	// Arrange
	comps := []*comp {
		&comp{Series: "GEoFOrm", Season: "2020", Number: 1, Name: "1"},
		&comp{Series: "geoform", Season: "2020", Number: 2, Name: "2"},
		&comp{Series: "geoform", Season: "2019", Number: 12, Name: "12"},
	}

	i := &impl{
		comps: map[string]*comp{},
		mux:   &sync.Mutex{},
	}
	for _, c := range comps {
		i.add(c)
	}

	// Act
	cs := i.getAll("GEOFORM", "2020")

	// Assert
	if len(cs) != 2 {
		t.Errorf("expected 2 comps, got %d", len(cs))
	}
}