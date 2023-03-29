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

	client, _ := git.NewClient()
	cfg, err := client.Config("user.name")

	require.NoError(t, err)
	assert.Equal(t, gittest.DefaultAuthorName, cfg)
}

func TestConfigL(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	cfg, err := client.ConfigL("user.name", "user.email")

	require.NoError(t, err)
	require.Len(t, cfg, 2)
	assert.Equal(t, gittest.DefaultAuthorName, cfg["user.name"])
	assert.Equal(t, gittest.DefaultAuthorEmail, cfg["user.email"])
}

func TestConfigSet(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	err := client.ConfigSet("user.name", "")

	require.NoError(t, err)
}

func TestConfigSetL(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	err := client.ConfigSetL("user.name", "user.email", "user.email", "")

	require.NoError(t, err)
	// verify that it has been set
}

func TestConfigSetLMismatchedPairsError(t *testing.T) {
	client, _ := git.NewClient()

	err := client.ConfigSetL("user.name")
	require.Error(t, err)
	// TODO: verify error message
}

func TestConfigSetLInvalidPathError(t *testing.T) {
	client, _ := git.NewClient()

	err := client.ConfigSetL("user.email", "jdoe@gmail.com", "user.1ame", "jdoe")
	require.Error(t, err)
	// TODO: verify error message
}

func TestValidConfigPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		isValid bool
	}{
		{
			name:    "ValidSingleDotPath",
			path:    "gr8.path",
			isValid: true,
		},
		{
			name:    "ValidMultiDotPath",
			path:    "a.gr8.path",
			isValid: true,
		},
		{
			name:    "InvalidMissingDot",
			path:    "nodot",
			isValid: false,
		},
		{
			name:    "InvalidJustSection",
			path:    "section.only.",
			isValid: false,
		},
		{
			name:    "InvalidDigitAfterLastDot",
			path:    "a.bad.4pple",
			isValid: false,
		},
		{
			name:    "InvalidContainsNonAlphanumeric",
			path:    "no.$symbol.allowed",
			isValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, git.ValidConfigPath(tt.path))
		})
	}
}
