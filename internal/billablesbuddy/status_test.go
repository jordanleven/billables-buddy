package billablesbuddy

import (
	"testing"
)

func TestGetCurrentStatus(t *testing.T) {
	t.Run("Returns the correct status when behind", func(t *testing.T) {

		actual := getCurrentStatus(10, 20, 40)
		expected := StatusBehind

		if actual != expected {
			t.Errorf("Received %d; want %d", actual, expected)
		}
	})

	t.Run("Returns the correct status when ahead", func(t *testing.T) {

		actual := getCurrentStatus(25, 20, 40)
		expected := StatusAhead

		if actual != expected {
			t.Errorf("Received %d; want %d", actual, expected)
		}
	})

	t.Run("Returns the correct status when over", func(t *testing.T) {

		actual := getCurrentStatus(41, 30, 40)
		expected := StatusOver

		if actual != expected {
			t.Errorf("Received %d; want %d", actual, expected)
		}
	})

	t.Run("Returns the correct status when on track", func(t *testing.T) {

		actual := getCurrentStatus(30, 30, 40)
		expected := StatusOnTrack

		if actual != expected {
			t.Errorf("Received %d; want %d", actual, expected)
		}
	})

	t.Run("Returns the correct status when in the grace period but ahead", func(t *testing.T) {

		actual := getCurrentStatus(30.25, 30, 40)
		expected := StatusOnTrack

		if actual != expected {
			t.Errorf("Received %d; want %d", actual, expected)
		}
	})

	t.Run("Returns the correct status when in the grace period but behind", func(t *testing.T) {

		actual := getCurrentStatus(29.75, 30, 40)
		expected := StatusOnTrack

		if actual != expected {
			t.Errorf("Received %d; want %d", actual, expected)
		}
	})
}
