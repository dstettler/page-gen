package pagegencore

import (
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetHTMLTagFromStringFor(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	var tag HTMLTag
	tag.TagName = "for"

	var items map[string]interface{} = make(map[string]interface{})
	items["var"] = "apps"
	items["refname"] = "app"
	items["unnecessary-int"] = 2

	tag.TagItems = items

	tag.StartPos = 0

	testStr := "<for var=\"apps\" refname =\"app\"  unnecessary-int = \"2\">"
	tag.EndPos = len(testStr) - 1

	genTag, found := GetHTMLTagFromString(testStr)
	if !found {
		t.Fatalf(`GetHTMLTagFromString returned not found!`)
	}
	if !cmp.Equal(tag, genTag) {
		t.Log("Given:", tag, "\n")
		t.Log("Generated:", genTag, "\n")
		t.Fatalf(`Tags unequal!`)
	}
}

func TestGetHTMLTagFromStringIfInt(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	var tag HTMLTag
	tag.TagName = "if"

	var items map[string]interface{} = make(map[string]interface{})
	items["var"] = "apps"
	items["val"] = 3
	items["condition"] = "lt"

	tag.TagItems = items

	tag.StartPos = 3

	testStr := "\n  <if var=\"apps\"   val =\"3\" condition=\"lt\" >  "
	tag.EndPos = len(testStr) - 3

	genTag, found := GetHTMLTagFromString(testStr)
	if !found {
		t.Fatalf(`GetHTMLTagFromString returned not found!`)
	}
	if !cmp.Equal(tag, genTag) {
		t.Log("Given:", tag, "\n")
		t.Log("Generated:", genTag, "\n")
		t.Fatalf(`Tags unequal!`)
	}
}

func TestGetHTMLTagFromStringIfFloat(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	var tag HTMLTag
	tag.TagName = "if"

	var items map[string]interface{} = make(map[string]interface{})
	items["var"] = "apps"
	items["val"] = 3.0
	items["condition"] = "lt"

	tag.TagItems = items

	tag.StartPos = 3

	testStr := "\n  <if var=\"apps\"   val =\"3.0\" condition=\"lt\" >  "
	tag.EndPos = len(testStr) - 3

	genTag, found := GetHTMLTagFromString(testStr)
	if !found {
		t.Fatalf(`GetHTMLTagFromString returned not found!`)
	}
	if !cmp.Equal(tag, genTag) {
		t.Log("Given:", tag, "\n")
		t.Log("Generated:", genTag, "\n")
		t.Fatalf(`Tags unequal!`)
	}
}
