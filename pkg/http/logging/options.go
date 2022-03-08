package logging

import "io"

type Option func(m *middleware)

func WithIOWriter(w io.Writer) Option {
	return func(m *middleware) {
		m.output = w
	}
}
