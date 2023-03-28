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

package git

import (
	"fmt"
	"strings"
)

// SortKey represents a structured [field name] that can be used as a sort key
// when analysing referenced objects such as tags
//
// [field name]: https://git-scm.com/docs/git-for-each-ref#_field_names
type SortKey string

const (
	// CreatorDate sorts the reference in ascending order by the creation date
	// of the underlying commit
	CreatorDate SortKey = "creatordate"

	// CreatorDateDesc sorts the reference in descending order by the creation date
	// of the underlying commit
	CreatorDateDesc SortKey = "-creatordate"

	// RefName sorts the reference by its name in ascending lexicographic order
	RefName SortKey = "refname"

	// RefNameDesc sorts the reference by its name in descending lexicographic order
	RefNameDesc SortKey = "-refname"

	// TaggerDate sorts the reference in ascending order by its tag creation date
	TaggerDate SortKey = "taggerdate"

	// TaggerDateDesc sorts the reference in descending order by its tag
	// creation date
	TaggerDateDesc SortKey = "-taggerdate"

	// Version interpolates the references as a version number and sorts in
	// ascending order
	Version SortKey = "version:refname"

	// VersionDesc interpolates the references as a version number and sorts in
	// descending order
	VersionDesc SortKey = "-version:refname"
)

// String converts the sort key from an enum into its string counterpart
func (k SortKey) String() string {
	return string(k)
}

// CreateTagOption provides a way for setting specific options during a tag
// creation operation. Each supported option can customize the way the tag is
// created against the current repository (working directory)
type CreateTagOption func(*createTagOptions)

type createTagOptions struct {
	Annotation string
}

// WithAnnotation ensures the created tag is annotated with the provided
// message. This ultimately converts the standard lightweight tag into
// an annotated tag which is stored as a full object within the git
// database. Any leading and trailing whitespace will automatically be
// trimmed from the message. This allows empty messages to be ignored
func WithAnnotation(message string) CreateTagOption {
	return func(opts *createTagOptions) {
		opts.Annotation = strings.TrimSpace(message)
	}
}

// Tag a specific point within a repositories history and push it to the
// configured remote. Tagging comes in two flavours:
//   - A lightweight tag, which points to a specific commit within
//     the history and marks a specific point in time
//   - An annotated tag, which is treated as a full object within
//     git, and must include a tagging message (or annotation)
//
// By default, a lightweight tag will be created, unless specific tag
// options are provided
func (c *Client) Tag(tag string, opts ...CreateTagOption) (string, error) {
	options := &createTagOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Build command based on the provided options
	var tagCmd strings.Builder
	tagCmd.WriteString("git tag ")

	if options.Annotation != "" {
		tagCmd.WriteString(fmt.Sprintf("-m '%s' -a ", options.Annotation))
	}
	tagCmd.WriteString(fmt.Sprintf("'%s'", tag))

	if out, err := exec(tagCmd.String()); err != nil {
		return out, err
	}

	return exec(fmt.Sprintf("git push origin '%s'", tag))
}

// DeleteTag a tag both locally and from the remote origin
func (c *Client) DeleteTag(tag string) (string, error) {
	if out, err := exec(fmt.Sprintf(`git tag -d "%s"`, tag)); err != nil {
		return out, err
	}

	return exec(fmt.Sprintf("git push --delete origin '%s'", tag))
}

// ListTagsOption provides a way for setting specific options during a list
// tags operation. Each supported option can customize the way in which the
// tags are queried and returned from the current repository (workng directory)
type ListTagsOption func(*listTagsOptions)

type listTagsOptions struct {
	Count        int
	Filters      []TagFilter
	ShellGlobs   []string
	SemanticSort bool
	SortBy       []string
}

// TagFilter allows a tag to be filtered based on any user-defined
// criteria. If the filter returns true, the tag will be included
// within the filtered results:
//
//	componentFilter := func(tag string) bool {
//		return strings.HasPrefix(tag, "component/")
//	}
type TagFilter func(tag string) bool

// WithCount limits the number of tags that are returned after all
// processing and filtering has been applied the retrieved list
func WithCount(n int) ListTagsOption {
	return func(opts *listTagsOptions) {
		opts.Count = n
	}
}

// WithFilters allows the retrieved list of tags to be processed
// with a set of user-defined filters. Each filter is applied in
// turn to the working set. Nil filters are ignored
func WithFilters(filters ...TagFilter) ListTagsOption {
	return func(opts *listTagsOptions) {
		opts.Filters = make([]TagFilter, 0, len(filters))
		for _, filter := range filters {
			if filter == nil {
				continue
			}

			opts.Filters = append(opts.Filters, filter)
		}
	}
}

// WithShellGlob limits the number of tags that will be retrieved, by only
// returning tags that match a given [Shell Glob] pattern. If multiple
// patterns are provided, tags will be retrieved if they match against
// a single pattern. All leading and trailing whitespace will be trimmed,
// allowing empty patterns to be ignored
//
// [Shell Glob]: https://tldp.org/LDP/GNU-Linux-Tools-Summary/html/x11655.htm
func WithShellGlob(patterns ...string) ListTagsOption {
	return func(opts *listTagsOptions) {
		opts.ShellGlobs = TrimAndPrefix("refs/tags/", patterns...)
	}
}

// WithSortBy allows the retrieved order of tags to be changed by sorting
// against a reserved [field name]. By default, sorting will always be in
// ascending order. To change this behaviour, prefix a field name with a
// hyphen (-<fieldname>). You can sort tags against multiple fields, but
// this does change the expected behavior. The last field name is treated
// as the primary key for the entire sort. All leading and trailing whitespace
// will be trimmed, allowing empty field names to be ignored
//
// [field name]: https://git-scm.com/docs/git-for-each-ref#_field_names
func WithSortBy(keys ...SortKey) ListTagsOption {
	return func(opts *listTagsOptions) {
		converted := make([]string, len(keys))
		for _, key := range keys {
			if key == Version || key == VersionDesc {
				// Ensure semantic versioning tags are going to be sorted correctly
				opts.SemanticSort = true
			}

			converted = append(converted, key.String())
		}

		opts.SortBy = TrimAndPrefix("--sort=", converted...)
	}
}

// Tags retrieves all local tags from the current repository (working directory).
// By default, all tags are retrieved in ascending lexicographic order as implied
// through the [RefName] sort key. Options can be provided to customize retrieval
func (c *Client) Tags(opts ...ListTagsOption) ([]string, error) {
	options := &listTagsOptions{
		Count: disabledNumericOption,
	}
	for _, opt := range opts {
		opt(options)
	}

	if len(options.ShellGlobs) == 0 {
		options.ShellGlobs = append(options.ShellGlobs, "refs/tags/**")
	}

	var config string
	if options.SemanticSort {
		config = "-c versionsort.suffix=-"
	}

	tags, err := exec(fmt.Sprintf("git %s for-each-ref %s --format='%%(refname:lstrip=2)' %s --color=never",
		config,
		strings.Join(options.SortBy, " "),
		strings.Join(options.ShellGlobs, " ")))
	if err != nil {
		return nil, err
	}

	if tags == "" {
		return nil, nil
	}

	splitTags := strings.Split(tags, "\n")
	splitTags = filterTags(splitTags, options.Filters)

	if options.Count > disabledNumericOption && options.Count <= len(splitTags) {
		return splitTags[:options.Count], nil
	}

	return splitTags, nil
}

func filterTags(tags []string, filters []TagFilter) []string {
	filtered := tags
	for _, filter := range filters {
		keep := make([]string, 0, len(filtered))
		for _, tag := range filtered {
			if filter(tag) {
				keep = append(keep, tag)
			}
		}

		filtered = keep
	}

	return filtered
}
