package harvestclient

import (
	"net/url"
	"testing"
)

func TestToUrlValues(t *testing.T) {
	t.Run("Returns expected URL values when values are not set", func(t *testing.T) {
		args := Arguments{}
		expected := url.Values{}
		actual := args.ToUrlValues()

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

		actual := args.ToUrlValues()

		if actual.Encode() == expected.Encode() {
			t.Errorf("Received %s; want %s", actual, expected)
		}
	})
}
