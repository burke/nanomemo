package supermemo

type FactSet []*Fact
type Fact struct {
	Question string
	Answer   string

	*FactMetadata
}

func NewFact(q, a string) *Fact {
	return &Fact{q, a, newFactMetadata()}
}

func (f *Fact) Dump() (q, a string, ef float64, n, interval int, intervalFrom string) {
	ts := f.intervalFrom.Format("2006-01-02")
	return f.Question, f.Answer, f.ef, f.n, f.interval, ts
}

func LoadFact(q, a string, ef float64, n, interval int, intervalFrom string) (*Fact, error) {
	md, err := loadFactMetadata(ef, n, interval, intervalFrom)
	if err != nil {
		return nil, err
	}
	return &Fact{q, a, md}, nil
}

func (s FactSet) ForReview() FactSet {
	var subset FactSet
	for _, fact := range s {
		if fact.UpForReview() {
			subset = append(subset, fact)
		}
	}
	return subset
}
