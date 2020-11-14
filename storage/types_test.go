package storage

import (
	"fmt"
	"testing"
)

func Test_computedAthlete_computePoints(t *testing.T) {
	// Arrange
	cs := &computedAthlete{
		Results:           []athleteResult{
			{Points:       140.44},
			{Points:       136.86},
			{Points:       139.07},
			{Points:       135.37},
			{Points:       142.82},
			{Points:       138.46},
			{Points:       143.23},
			{Points:       145.18},
			{Points:       146.15},
			{Points:       143.32},
			{Points:       138.41},
		},
	}

	// Act
	cs.computePoints(1)
	if cs.PointsOfficial != 146.15 {
		t.Errorf("unexpected official points after 1 event, got %f", cs.PointsOfficial)
	}

	cs.computePoints(2)
	fs := fmt.Sprintf("%f", cs.PointsOfficial)
	if fs != "291.330000" {
		t.Errorf("unexpected official points after 2 events, got %s", fs)
	}

	cs.computePoints(999999)
	fs = fmt.Sprintf("%f", cs.PointsOfficial)
	if fs != "1549.310000" {
		t.Errorf("unexpected official points after all events, got %s", fs)
	}
	//""
}
