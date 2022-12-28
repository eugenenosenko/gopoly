package source

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFromPackage(t *testing.T) {
	t.Run("should correctly load source from package", func(t *testing.T) {
		got, err := LoadFromPackage("./testdata", "./testdata/b", "./testdata/c")

		require.NoError(t, err)
		require.Len(t, got, 3)

		require.Equal(t, "testdata", got[0].Name)
		require.Equal(t, "github.com/eugenenosenko/gopoly/internal/source/testdata", got[0].Path)
		require.Len(t, got[0].Files, 1)

		require.Equal(t, "b", got[1].Name)
		require.Equal(t, "github.com/eugenenosenko/gopoly/internal/source/testdata/b", got[1].Path)
		require.Len(t, got[1].Files, 1)

		require.Equal(t, "c", got[2].Name)
		require.Equal(t, "github.com/eugenenosenko/gopoly/internal/source/testdata/c", got[2].Path)
		require.Len(t, got[2].Files, 1)
	})
}
