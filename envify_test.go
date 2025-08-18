package goenvify_test

import (
	"maps"
	"slices"
	"testing"

	goenvify "github.com/efeligne/GoEnvify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFile(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name    string
		file    string
		want    string
		wantErr bool
	}

	tests := []TestCase{
		{
			name:    "LoadFile should return the content of the file",
			file:    "./mocks/mock.file.env",
			want:    "ENV_VAR1=value1\nENV_VAR2=value2\nENV_VAR3=value3\nENV_VAR4=value4\n",
			wantErr: false,
		},
		{
			name:    "LoadFile should return an error if the file does not exist",
			file:    "./mocks/mock.not.exist",
			want:    "",
			wantErr: true,
		},
		{
			name:    "LoadFile should return an error if the file is empty",
			file:    "./mocks/mock.file.empty",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := goenvify.LoadFile(tt.file)
			if gotErr != nil && !tt.wantErr {
				require.NoError(t, gotErr, "LoadFile() should not return an error")
			}

			if tt.wantErr {
				require.Error(t, gotErr, "LoadFile() should return an error")
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSplitContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "SplitContent should split the content by line",
			content: "ENV_VAR1=value1\nENV_VAR2=value2\n\rENV_VAR3=value3",
			want: []string{
				"ENV_VAR1=value1",
				"ENV_VAR2=value2",
				"ENV_VAR3=value3",
			},
		},
		{
			name:    "SplitContent should remove comments",
			content: "# shellcheck disable=SC2034\nENV_VAR3=value3\nENV_VAR4=value4",
			want: []string{
				"ENV_VAR3=value3",
				"ENV_VAR4=value4",
			},
		},

		{
			name:    "SplitContent should remove empty lines",
			content: "\nENV_VAR5=value5\n\nENV_VAR6=value6",
			want: []string{
				"ENV_VAR5=value5",
				"ENV_VAR6=value6",
			},
		},
		{
			name:    "SplitContent should trim spaces",
			content: "  ENV_VAR7=value7  \n  ENV_VAR8=value8",
			want: []string{
				"ENV_VAR7=value7",
				"ENV_VAR8=value8",
			},
		},
		{
			name:    "SplitContent should remove inline comments",
			content: "ENV_VAR9=value9 # test comment\nENV_VAR10=value10#test #comment",
			want: []string{
				"ENV_VAR9=value9",
				"ENV_VAR10=value10#test",
			},
		},
		{
			name:    "Empty content should return an empty slice",
			content: "",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := goenvify.SplitContent(tt.content)

			if !slices.Equal(got, tt.want) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMapContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content []string
		want    map[string]string
	}{
		{
			name:    "MapContent should map the content to a map",
			content: []string{"ENV_VAR1=value1", "ENV_VAR2=value2", "ENV_VAR3="},
			want: map[string]string{
				"ENV_VAR1": "value1",
				"ENV_VAR2": "value2",
				"ENV_VAR3": "",
			},
		},
		{
			name:    "MapContent should trim spaces",
			content: []string{"ENV_VAR3= value3 ", " ENV_VAR4=   value4"},
			want: map[string]string{
				"ENV_VAR3": "value3",
				"ENV_VAR4": "value4",
			},
		},
		{
			name:    "Empty content should return an empty map",
			content: []string{},
			want:    map[string]string{},
		},
		{
			name:    "Invalid content should return an empty map",
			content: []string{"", "=value", "str"},
			want:    map[string]string{},
		},
		{
			name: "Its possible to use = and : as separators",
			content: []string{
				"ENV_VAR5=value5",
				"ENV_VAR6:value6",
				"DB_URL=postgres://user:pass@host:5432/db",
			},
			want: map[string]string{
				"ENV_VAR5": "value5",
				"ENV_VAR6": "value6",
				"DB_URL":   "postgres://user:pass@host:5432/db",
			},
		},
		{
			name:    "MapContent must accept quotation marks around the value ",
			content: []string{"ENV_VAR7=\"value7\"", "ENV_VAR8='value8'", "ENV_VAR9: \"value9\""},
			want: map[string]string{
				"ENV_VAR7": "value7",
				"ENV_VAR8": "value8",
				"ENV_VAR9": "value9",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := goenvify.MapContent(tt.content)

			if !maps.Equal(got, tt.want) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
