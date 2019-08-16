package fmthtml

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestFile struct {
	n          int
	BeforePath string
	AfterPath  string
}

func readTestFiles(t *testing.T) []*TestFile {
	dir := "testfiles"
	files, err := ioutil.ReadDir(dir)
	require.NoError(t, err)
	testCasesMap := map[int]*TestFile{}
	for _, fi := range files {
		name := fi.Name()
		parts := strings.Split(name, ".")
		require.Equal(t, 2, len(parts))
		require.Equal(t, "html", parts[1])
		s := parts[0] // "00-before" or "00-after"
		parts = strings.Split(s, "-")
		require.Equal(t, 2, len(parts))
		n, err := strconv.Atoi(parts[0])
		require.NoError(t, err)
		tc := testCasesMap[n]
		if tc == nil {
			tc = &TestFile{
				n: n,
			}
			testCasesMap[n] = tc
		}
		path := filepath.Join(dir, fi.Name())
		switch parts[1] {
		case "before":
			tc.BeforePath = path
		case "after":
			tc.AfterPath = path
		default:
			require.True(t, false, "unexpected file name: '%s'\n", fi.Name())
		}
	}
	var testCases []*TestFile
	for _, tc := range testCasesMap {
		testCases = append(testCases, tc)
	}

	sort.Slice(testCases, func(i, j int) bool {
		return testCases[i].n < testCases[j].n
	})
	return testCases
}

func TestFiles(t *testing.T) {
	testCases := readTestFiles(t)
	for _, tc := range testCases {
		require.NotEmpty(t, tc.BeforePath)
		require.NotEmpty(t, tc.AfterPath)
		before, err := ioutil.ReadFile(tc.BeforePath)
		require.NoError(t, err)
		exp, err := ioutil.ReadFile(tc.AfterPath)
		require.NoError(t, err)
		got := Format(before)
		require.Equal(t, string(exp), string(got))
	}
}
