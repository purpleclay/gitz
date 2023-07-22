/*
Copyright (c) 2023 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package git_test

import (
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	gittest.InitRepository(t)
	gittest.ConfigSet(t, "user.name", "joker", "user.email", "joker@dc.com")

	client, _ := git.NewClient()
	cfg, err := client.Config()

	require.NoError(t, err)
	assert.Equal(t, "joker", cfg["user.name"])
	assert.Equal(t, "joker@dc.com", cfg["user.email"])
}

func TestConfigOnlyLatestValues(t *testing.T) {
	gittest.InitRepository(t)
	gittest.ConfigSet(t, "user.name", "joker", "user.name", "scarecrow")

	client, _ := git.NewClient()
	cfg, err := client.Config()

	require.NoError(t, err)
	assert.Equal(t, "scarecrow", cfg["user.name"])
}

func TestConfigL(t *testing.T) {
	gittest.InitRepository(t)
	gittest.ConfigSet(t, "user.name", "alfred")

	client, _ := git.NewClient()
	cfg, err := client.ConfigL("user.name", "user.email")

	require.NoError(t, err)
	require.Len(t, cfg["user.name"], 2)
	assert.Equal(t, "alfred", cfg["user.name"][0])
	assert.Equal(t, gittest.DefaultAuthorName, cfg["user.name"][1])

	require.Len(t, cfg["user.email"], 1)
	assert.Equal(t, gittest.DefaultAuthorEmail, cfg["user.email"][0])
}

func configEquals(t *testing.T, path, expected string) {
	t.Helper()
	cfg, err := gittest.Exec(t, "git config --get "+path)

	require.NoError(t, err)
	assert.Equal(t, expected, cfg)
}

func TestConfigSetL(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	err := client.ConfigSetL("user.phobia", "bats", "user.birth.place", "gotham")

	require.NoError(t, err)
	configEquals(t, "user.phobia", "bats")
	configEquals(t, "user.birth.place", "gotham")
}

func TestConfigSetLMismatchedPairsError(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	err := client.ConfigSetL("user.hobbies")

	assert.EqualError(t, err, "config paths mismatch. path: user.hobbies is missing a corresponding value")
}

func TestConfigSetLNothingSetIfError(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	err := client.ConfigSetL("user.hobbies", "fighting crime", "user.arch.3nemy", "joker")

	require.Error(t, err)
	configMissing(t, "user.hobbies")
	configMissing(t, "user.4rch.enemy")
}

func configMissing(t *testing.T, path string) {
	t.Helper()
	cfg, err := gittest.Exec(t, "git config --get "+path)

	require.Error(t, err)
	require.Empty(t, cfg)
}

func TestCheckConfigPathError(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		errMsg string
	}{
		{
			name:   "InvalidMissingDot",
			path:   "nodot",
			errMsg: "path: nodot invalid as dot separator is missing or is the last character",
		},
		{
			name:   "InvalidJustSection",
			path:   "section.only.",
			errMsg: "path: section.only|.| invalid as dot separator is missing or is the last character",
		},
		{
			name:   "InvalidDigitAfterLastDot",
			path:   "a.bad.4pple",
			errMsg: "path: a.bad.|4|pple invalid as first character after final dot must be a letter [a-zA-Z]",
		},
		{
			name:   "InvalidContainsNonAlphanumeric",
			path:   "no.$symbol.allowed",
			errMsg: "path: no.|$|symbol.allowed invalid as non alphanumeric character detected [a-zA-Z0-9]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualError(t, git.CheckConfigPath(tt.path), tt.errMsg)
		})
	}
}

func TestToInlineConfig(t *testing.T) {
	cfg, err := git.ToInlineConfig("user.name", "penguin", "user.email", "penguin@dc.com")

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"-c user.name='penguin'", "-c user.email='penguin@dc.com'"}, cfg)
}
