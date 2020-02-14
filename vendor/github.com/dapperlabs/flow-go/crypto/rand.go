package crypto

// Rand is a pseudo random number generator
type Rand struct {
	states     []uint64
	stateIndex int
}

// Next generates a new state given previous state
// Xorshift is used with these parameters (13, 9, 15)
func (x *Rand) next(index int) {
	x.states[index] ^= (x.states[index] << 13)
	x.states[index] ^= (x.states[index] >> 9)
	x.states[index] ^= (x.states[index] << 15)
}

// IntN returns an int random number between 0 and "to" (exclusive)
func (x *Rand) IntN(to int) int {
	res := x.states[x.stateIndex] % uint64(to)
	x.next(x.stateIndex)
	x.stateIndex = (x.stateIndex + 1) % len(x.states)
	return int(res)
}

// NewRand returns a new random generator (Rand)
func NewRand(seed []uint64) (*Rand, error) {
	size := len(seed)
	// seed slice can not include zeros
	for i := range seed {
		if seed[i] == uint64(0) {
			return nil, &cryptoError{"seed slice can not have a zero value"}
		}
	}
	// stateIndex initialized with first element of seed % rand size
	states := make([]uint64, len(seed))
	copy(states, seed)
	rand := &Rand{states: states, stateIndex: int(seed[0] % uint64(size))}
	// initial next
	for i := range rand.states {
		rand.next(i)
	}
	return rand, nil
}

// RandomPermutationSubset implements Fisher-Yates Shuffle (inside-outside) on a slice of ints with
// size n using the given seed s and returns a subset m of it
func RandomPermutationSubset(n int, m int, s []uint64) ([]int, error) {
	items, err := RandomPermutation(n, s)
	if err != nil {
		return nil, err
	}
	return items[:m], nil
}

// RandomPermutation implements Fisher-Yates Shuffle (inside-outside) on a slice of ints with
// size n using the given seed s and returns the slice
func RandomPermutation(n int, s []uint64) ([]int, error) {
	items := make([]int, n)
	rand, err := NewRand(s)
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		j := rand.IntN(i + 1)
		if j != i {
			items[i] = items[j]
		}
		items[j] = i
	}
	return items, nil
}
