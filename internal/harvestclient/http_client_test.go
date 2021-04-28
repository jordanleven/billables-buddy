package harvestclient

import (
	"net/url"
	"testing"
)

func TestToURLValues(t *testing.T) {
	t.Run("Returns expected URL values when values are not set", func(t *testing.T) {
		args := Arguments{}
		expected := url.Values{}
		actual := args.ToURLValues()

		if len(actual) != 0 {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})

	t.Run("Returns expected URL values when values are set", func(t *testing.T) {
		args := Arguments{
			"foo": "bar",
			"baz": "bax",
		}
		expected := url.Values{}
		expected.Set("foo", "bar")

		actual := args.ToURLValues()

		if actual.Encode() == expected.Encode() {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})
}
