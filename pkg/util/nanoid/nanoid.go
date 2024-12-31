package nanoid

import gonanoid "github.com/matoous/go-nanoid/v2"

func New() (string, error) {
	return gonanoid.New()
}

func Must() string {
	return gonanoid.Must()
}
