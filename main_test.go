package main

import (
	"bufio"
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/Nivl/check-deps/modutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	testCases := []struct {
		description    string
		argv           []string
		expectedResult Flags
		expectedError  error
	}{
		{
			description: "default flags",
			argv:        []string{"bin"},
			expectedResult: Flags{
				CheckOldPkgs:  false,
				CheckIndirect: false,
				IgnoredPkgs:   []string{},
			},
			expectedError: nil,
		},
		{
			description: "set all",
			argv: []string{
				"bin",
				"--old",
				"--indirect",
				"-i",
				"pkg1,pkg2",
				"--ignore=pkg3,pkg4",
			},
			expectedResult: Flags{
				CheckOldPkgs:  true,
				CheckIndirect: true,
				IgnoredPkgs:   []string{"pkg1", "pkg2", "pkg3", "pkg4"},
			},
			expectedError: nil,
		},
		{
			description: "invalid flag",
			argv: []string{
				"bin",
				"--nope",
			},
			expectedResult: Flags{},
			expectedError:  errors.New("unknown flag: --nope"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			flags, err := parseFlags(tc.argv)
			if tc.expectedError != nil {
				require.Error(t, err, "parseFlags should have failed")
				require.Equal(t, tc.expectedError, err, "parseFlags failed with an unexpected error")
				return
			}

			require.NoError(t, err, "parseFlags should have succeed")
			assert.Equal(t, tc.expectedResult, *flags)

		})
	}
}

func TestParseModules(t *testing.T) {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	OneYearAgo := now.Add(-366 * 24 * time.Hour)
	validModule := &modutil.Module{Path: "valid/pkg", Version: "0.0.1"}
	indirectModule := &modutil.Module{Path: "indirect/pkg", Indirect: true, Replace: validModule}
	replacedModule := &modutil.Module{Path: "replace/pkg", Replace: validModule}
	updatedToModule := &modutil.Module{Path: "updated/pkg", Version: "1.0.0", Time: &now}
	updatedModule := &modutil.Module{Path: "updated/pkg", Version: "0.0.1", Update: updatedToModule}
	updatedModuleInvalid := &modutil.Module{Path: "updated/pkg", Version: "0.0.1", Update: updatedToModule, Time: &tomorrow}
	oldModule := &modutil.Module{Path: "old/pkg", Version: "0.0.1", Time: &OneYearAgo}

	testCases := []struct {
		description string
		flags       Flags
		modules     []*modutil.Module
		expected    Results
	}{
		{description: "no modules"},
		{
			description: "one ignored module",
			flags: Flags{
				IgnoredPkgs: []string{validModule.Path},
			},
			modules: []*modutil.Module{
				validModule,
			},
		},
		{
			description: "one skipped indirect module",
			flags: Flags{
				CheckIndirect: false,
			},
			modules: []*modutil.Module{
				indirectModule,
			},
		},
		{
			description: "one replaced module",
			modules: []*modutil.Module{
				replacedModule,
			},
			expected: Results{
				Replaced: []*modutil.Module{
					replacedModule,
				},
			},
		},
		{
			description: "one replaced indirect module",
			flags: Flags{
				CheckIndirect: true,
			},
			modules: []*modutil.Module{
				indirectModule,
			},
			expected: Results{
				Replaced: []*modutil.Module{
					indirectModule,
				},
			},
		},
		{
			description: "one updated module",
			modules: []*modutil.Module{
				updatedModule,
			},
			expected: Results{
				Updated: []*modutil.Module{
					updatedModule,
				},
			},
		},
		{
			description: "one invalid updated module",
			modules: []*modutil.Module{
				updatedModuleInvalid,
			},
		},
		{
			description: "one skipped old module",
			modules: []*modutil.Module{
				oldModule,
			},
		},
		{
			description: "one old module",
			flags: Flags{
				CheckOldPkgs: true,
			},
			modules: []*modutil.Module{
				oldModule,
			},
			expected: Results{
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

			res := parseModules(&tc.flags, tc.modules)
			assert.Equal(t, tc.expected, *res)
		})
	}
}

func TestRun(t *testing.T) {
	testCases := []struct {
		description  string
		args         []string
		expectedBuf  string
		expectedCode int
	}{
		{
			description:  "invalid flags",
			args:         []string{"bin", "--nope"},
			expectedCode: ExitFailure,
		},
		{
			description:  "happy path",
			args:         []string{"bin"},
			expectedCode: ExitSuccess,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			buf := bytes.Buffer{}
			w := bufio.NewWriter(&buf)
			exitStatus := run(tc.args, w)
			require.NoError(t, w.Flush(), "Flush() should have work")
			require.Equal(t, tc.expectedCode, exitStatus)
			if tc.expectedBuf != "" {
				require.Equal(t, tc.expectedBuf, buf.String())
			}
		})
	}
}
