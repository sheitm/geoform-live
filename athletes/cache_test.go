package athletes

import (
	"sync"
	"testing"
)

func Test_sha(t *testing.T) {
	type args struct {
		name string
		club string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{"Ståle Heitmann", "Fossum IF"}, "pOD9F-Ft3jHKmz6JDDCVQKjrYZ8="},
		{"2", args{" Ståle Heitmann", "Fossum IF "}, "pOD9F-Ft3jHKmz6JDDCVQKjrYZ8="},
		{"3", args{"Ståle S. Heitmann", "Fossum IF "}, "yNrDVPtt6qCeQbWZQhdcXXONiF8="},
		{"4", args{"Ståle S. Heitmann", "O-entusiastene"}, "bud9H6gua4lphxEL4fF0x4IY2yk="},
		{"5", args{"Benny Carlsson", "Sturla"}, "p3WN_8Jvg2V-ZHgGYGas19Wrvlg="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sha(tt.args.name, tt.args.club); got != tt.want {
				t.Errorf("id() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_competitor(t *testing.T) {
	// Arrange
	name := "Jens Petter Hansen"
	club := "OK Linne"
	c := &cache{
		competitorsBySHA:  map[string]*athleteWithID{},
		competitorsByGuid: map[string]*athleteWithID{},
		mux:               sync.Mutex{},
	}

	// Act
	a1, e1 := c.competitor(name, club)
	a2, e2 := c.competitor(name, club)

	// Assert
	if e1 {
		t.Error("athlete should not exist first time")
	}
	if !e2 {
		t.Error("athlete should already exist the second time")
	}
	if a1.ID != a2.ID {
		t.Errorf("expected same athlete, og %s and %s", a1.ID, a2.ID)
	}

}