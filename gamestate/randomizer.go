package gamestate

import "math/rand"

var _ RandomizerInterface = (*Randomizer)(nil)
var _ RandomizerInterface = (*TestRandomizer)(nil)

type RandomizerInterface interface {
	randomPercent() int
}

type Randomizer struct{}
func (r *Randomizer) randomPercent() int {
  return rand.Int() % 100
}

type TestRandomizer struct{}
func (r *TestRandomizer) randomPercent() int {
	return 100 // for the sake of testing all attacks are 100% accurate
}
