package supermemo

import "time"

// SM2.1

// type FactMetadata contains metadata associated to a conceptual knowledge item, with
// which the SM2 algorithm can work. Calling code will generally embed a FactMetadata in
// a type containing a Question and Answer.
type FactMetadata struct {
	// Easiness FactMetadataor of the fact. Higher means the item is easier for the user
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

// set intervalFrom to the first second of the current date.
func (f *FactMetadata) setIntervalFrom() {
	y, m, d := time.Now().Date()
	f.intervalFrom = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// `interval` days after the date represented by `intervalFrom`, the fact is up
// for review. Are we at or past that date?
func (f *FactMetadata) UpForReview() bool {
	reviewDate := f.intervalFrom.AddDate(0, 0, f.interval)
	return time.Now().After(reviewDate)
}

// SM2.2
const initialEF = 2.5

// NewFactMetadata initializes a new FactMetadata object with the correct default EF.
func newFactMetadata() *FactMetadata {
	f := &FactMetadata{ef: initialEF, n: 0, interval: 0}
	f.setIntervalFrom()
	return f
}

func loadFactMetadata(ef float64, n, interval int, intervalFrom string) (*FactMetadata, error) {
	t, err := time.Parse("2006-01-02", intervalFrom)
	if err != nil {
		return nil, err
	}
	return &FactMetadata{ef: ef, interval: interval, n: n, intervalFrom: t}, nil
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

// Assess updates the FactMetadata object with a new Interval and Easiness FactMetadataor based
// on the difficulty assessment provided by the user.
func (f *FactMetadata) Assess(q int) {
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
