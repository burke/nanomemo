/*
package supermemo implements the supermemo SM-2 algorithm, described here, and
documented in more detail at http://www.supermemo.com/english/ol/sm2.htm.

=== VOCABULARY =======

EF: Easiness Factor. Higher means easier. This represents the ease of memorizing a given item.

I(n): Interval, in days, between repetitions of an item. n represents the number of times the item has been seen.

q: The quality of a response, from 0-5, where:

  5 - perfect response;
  4 - correct response after a hesitation;
  3 - correct response recalled with serious difficulty;
  2 - incorrect response where the correct one seemed easy to recall;
  1 - incorrect response, correct one remembered upon seeing answer;
  0 - complete blackout.

=== ALGORITHM SM2 =======

1. Split the knowledge into the smallest possible items. Think flash cards.

2. With each item, associate an initial EF of 2.5.

3. Repeat items using the following intervals, expressed in days:

  I(1)       := 1
  I(2)       := 6
  I(n | n>2) := I(n-1)*EF

4. After each response, assess the quality of the response (q) as described
above in the Vocabulary section.

5. After each response, modify the EF by the formula:

  EF':=EF+(0.1-(5-q)*(0.08+(5-q)*0.02))

6. If the most recent q > 3, reset n; that is, restart repetitions from I(0).

7. After all items are processed, repeat all items where q < 4, until all items
have at least 4.
*/
package supermemo

import "time"

// SM2.1

// type Fact contains metadata associated to a conceptual knowledge item, with
// which the SM2 algorithm can work. Calling code will generally embed a Fact in
// a type containing a Question and Answer.
type Fact struct {
	// Easiness Factor of the fact. Higher means the item is easier for the user
	// to remember.
	ef float64
	// Interval number of days to wait before presenting this item again after the
	// end of this session.
	interval int
	// last time the fact was reviewed. Interval counts days from here.
	intervalFrom time.Time
	// number of times this fact has been presented; reset to 0 on failed answer.
	n int
}

// set intervaFrom to the first second of the current date.
func (f *Fact) setIntervalFrom() {
	y, m, d := time.Now().Date()
	f.intervalFrom = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// `interval` days after the date represented by `intervalFrom`, the fact is up
// for review. Are we at or past that date?
func (f *Fact) UpForReview() bool {
	reviewDate := f.intervalFrom.AddDate(0, 0, f.interval)
	return time.Now().After(reviewDate)
}

// SM2.2
const initialEF = 2.5

// NewFact initializes a new Fact object with the correct default EF.
func NewFact() *Fact {
	return &Fact{ef: initialEF, n: 0, interval: 0, intervalFrom: time.Now()}
}

// SM2.3
func unguardedNextInterval(ef float64, n, interval int) int {
	switch n {
	case 0:
		return 1
	case 1:
		return 6
	default:
		return int(float64(interval)*ef + 0.5)
	}
}

// SM2.4

// Assess updates the Fact object with a new Interval and Easiness Factor based
// on the difficulty assessment provided by the user.
func (f *Fact) Assess(q int) {
	f.setIntervalFrom()

	ef := nextEF(q, f.ef)
	n := nextN(q, f.n)
	f.interval = nextInterval(q, f.ef, f.n, f.interval)

	f.ef = ef
	f.n = n
}

// SM2.5
func nextEF(q int, ef float64) float64 {
	nxt := ef + (0.1 - float64(5-q)*(0.08+float64(5-q)*0.02))
	if 1.3 > nxt {
		return 1.3
	}
	return nxt
}

// SM2.6
func nextN(q, i int) int {
	if q < 3 {
		return 0
	}
	return i + 1
}

// SM2.7
func nextInterval(q int, ef float64, n, interval int) int {
	if q < 4 {
		return 0
	}
	return unguardedNextInterval(ef, n, interval)
}
