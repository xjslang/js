package js_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xjslang/js"
)

func TestTryCatch(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		input := `try {opendb()}`
		_, err := js.Parse([]byte(input))
		require.ErrorContains(t, err, "missing catch or finally after try")
	})
}

func TestLanguageFeatures(t *testing.T) {
	pattern := filepath.Join("testdata", "*.js")
	files, err := filepath.Glob(pattern)
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		testName := strings.TrimSuffix(filepath.Base(file), ".js")
		t.Run(testName, func(t *testing.T) {
			// read file
			dat, err := os.ReadFile(file)
			require.NoError(t, err)
			// parse data
			result, err := js.Parse(dat)
			require.NoError(t, err)
			// print result
			out, err := js.Print(result)
			require.NoError(t, err)
			// re-parse the output
			result, err = js.Parse([]byte(out))
			require.NoError(t, err)
			// re-print the result
			out, err = js.Print(result)
			require.NoError(t, err)
			// the original must match the final printed result
			require.Equal(t, string(dat), out)
		})
	}
}
