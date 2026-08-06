package data

import "testing"

func TestEmojisLoaded(t *testing.T) {
	if len(Emojis) == 0 {
		t.Fatal("Emojis is empty; embedded JSON failed to load")
	}
	first := Emojis[0]
	if first.Emoji == "" || first.Description == "" {
		t.Fatalf("first emoji entry incomplete: %+v", first)
	}
}
