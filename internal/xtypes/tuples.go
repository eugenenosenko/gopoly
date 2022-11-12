package xtypes

import (
	"fmt"
)

type Tuple[F, S any] struct {
	First  F
	Second S
}

func (t *Tuple[F, S]) String() string {
	return fmt.Sprintf("('%v','%v')", t.First, t.Second)
}

var _ fmt.Stringer = (*Tuple[string, string])(nil)
