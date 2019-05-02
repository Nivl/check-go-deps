package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/Nivl/check-go-deps/modutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasModules(t *testing.T) {
	t.Parallel()

	hasModule := true

	testCases := []struct {
		description string
		res         Results
		expected    bool
	}{
		{
			description: "no modules",
			res:         Results{},
			expected:    !hasModule,
		},
		{
			description: "updated modules",
			res: Results{
				Updated: []*modutil.Module{
					{},
				},
			},
			expected: hasModule,
		},
		{
			description: "replaced modules",
			res: Results{
				Replaced: []*modutil.Module{
					{},
				},
			},
			expected: hasModule,
		},
		{
			description: "old modules",
			res: Results{
				Old: []*modutil.Module{
					{},
				},
			},
			expected: hasModule,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.res.HasModules())
		})
	}
}

func TestPrint(t *testing.T) {
	now, err := time.Parse("2006/01/02", "2019/04/30")
	require.NoError(t, err, "time.Parse() was expected to succeed")
	now = now.UTC()

	OneYearAgo := now.Add(-366 * 24 * time.Hour)
	validModule := &modutil.Module{Path: "valid/pkg", Version: "0.0.1"}
	indirectModule := &modutil.Module{Path: "indirect/pkg", Indirect: true, Replace: validModule}
	replacedModule := &modutil.Module{Path: "replace/pkg", Replace: validModule}
	updatedToModule := &modutil.Module{Path: "updated/pkg", Version: "1.0.0", Time: &now}
	updatedModule := &modutil.Module{Path: "updated/pkg", Version: "0.0.1", Update: updatedToModule}
	oldModule := &modutil.Module{Path: "old/pkg", Version: "0.0.1", Time: &OneYearAgo}

	testCases := []struct {
		description        string
		res                Results
		expectedOutputFile string
	}{
		{
			description:        "no modules",
			res:                Results{},
			expectedOutputFile: "empty",
		},
		{
			description:        "updated modules",
			expectedOutputFile: "test-print-updated",
			res: Results{
				Updated: []*modutil.Module{
					updatedModule,
				},
			},
		},
		{
			description:        "replaced modules",
			expectedOutputFile: "test-print-replaced",
			res: Results{
				Replaced: []*modutil.Module{
					replacedModule,
					indirectModule,
				},
			},
		},
		{
			description:        "old modules",
			expectedOutputFile: "test-print-old",
			res: Results{
				Old: []*modutil.Module{
					oldModule,
				},
			},
		},
		{
			description:        "all",
			expectedOutputFile: "test-print-all",
			res: Results{
				Updated: []*modutil.Module{
					updatedModule,
				},
				Replaced: []*modutil.Module{
					replacedModule,
					indirectModule,
				},
				Old: []*modutil.Module{
					oldModule,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			expectedOutput, err := ioutil.ReadFile(filepath.Join("testdata", tc.expectedOutputFile))
			require.NoError(t, err, "ioutil.ReadFile() was expected to succeed")

			output := bytes.Buffer{}
			w := bufio.NewWriter(&output)
			tc.res.print(w)

			require.NoError(t, w.Flush(), "Flush() should have work")

			assert.Equal(t, string(expectedOutput), output.String())
		})
	}
}
