package supermemo

import (
	"fmt"
	"testing"
)
import . "github.com/smartystreets/goconvey/convey"

func TestNextN(t *testing.T) {
	qs := []int{5, 4, 3, 2, 1, 0}
	results := []int{4, 4, 4, 0, 0, 0}

	for _i := 0; _i < len(qs); _i++ {
		i := _i
		Convey(fmt.Sprintf("nextN(%d, 3) = %d", qs[i], results[i]), t, func() {
			So(nextN(qs[i], 3), ShouldEqual, results[i])
		})
	}
}

func TestNextInterval(t *testing.T) {
	var qs []int
	var results []int
	var efs []float64

	Convey("With i = 0", t, func() {
		qs = []int{5, 4, 3, 2, 1, 0}
		results = []int{1, 1, 0, 0, 0, 0}

		for _i := 0; _i < len(qs); _i++ {
			i := _i
			Convey(fmt.Sprintf("nextInterval(%d, 2.5, 0, 0) = %d", qs[i], results[i]), func() {
				So(nextInterval(qs[i], 2.5, 0, 0), ShouldEqual, results[i])
			})
		}
	})

	Convey("With i = 1", t, func() {
		qs = []int{5, 4, 3, 2, 1, 0}
		results = []int{6, 6, 0, 0, 0, 0}

		for _i := 0; _i < len(qs); _i++ {
			i := _i
			Convey(fmt.Sprintf("nextInterval(%d, 2.5, 1, 1) = %d", qs[i], results[i]), func() {
				So(nextInterval(qs[i], 2.5, 1, 1), ShouldEqual, results[i])
			})
		}
	})

	Convey("With i > 1", t, func() {
		qs = []int{5, 4, 3, 2, 1, 0}
		results = []int{16, 16, 0, 0, 0, 0}

		for _i := 0; _i < len(qs); _i++ {
			i := _i
			Convey(fmt.Sprintf("nextInterval(%d, 2.7, 2, 6) = %d", qs[i], results[i]), func() {
				So(nextInterval(qs[i], 2.7, 2, 6), ShouldEqual, results[i])
			})
		}
	})

	Convey("With changing EFs", t, func() {
		efs = []float64{1.3, 2.0, 2.5, 2.6, 2.7}
		results = []int{8, 12, 15, 16, 16}

		for _i := 0; _i < len(efs); _i++ {
			i := _i
			Convey(fmt.Sprintf("nextInterval(5, %f, 2, 6) = %d", efs[i], results[i]), func() {
				So(nextInterval(5, efs[i], 2, 6), ShouldEqual, results[i])
			})
		}
	})

}

func TestNextEF(t *testing.T) {
	Convey("Normal case", t, func() {
		qs := []int{5, 4, 3, 2, 1, 0}
		results := []float64{2.6, 2.5, 2.36, 2.18, 1.96, 1.7}

		for _i := 0; _i < len(qs); _i++ {
			i := _i
			Convey(fmt.Sprintf("nextEF(%d, 2.5) = %f", qs[i], results[i]), func() {
				So(nextEF(qs[i], 2.5), ShouldAlmostEqual, results[i])
			})
		}
	})
	Convey("Minimum case", t, func() {
		qs := []int{5, 4, 3, 2, 1, 0}
		results := []float64{1.4, 1.3, 1.3, 1.3, 1.3, 1.3}

		for _i := 0; _i < len(qs); _i++ {
			i := _i
			Convey(fmt.Sprintf("nextEF(%d, 1.3) = %f", qs[i], results[i]), func() {
				So(nextEF(qs[i], 1.3), ShouldAlmostEqual, results[i])
			})
		}
	})
}

type TraverseTC struct {
	history  []int
	ef       float64
	n        int
	interval int
}

func TestTraverse(t *testing.T) {
	cases := []TraverseTC{
		{[]int{}, 2.5, 0, 0},
		{[]int{5}, 2.6, 1, 1},
		{[]int{5, 5}, 2.7, 2, 6},
		{[]int{5, 5, 5}, 2.8, 3, 16},
		{[]int{5, 5, 5, 5}, 2.9, 4, 45},
	}

	for _, _tc := range cases {
		tc := _tc
		Convey(fmt.Sprintf("traverse(%v) = %f, %d, %d", tc.history, tc.ef, tc.n, tc.interval), t, func() {
			f := NewFact()
			for _, q := range tc.history {
				f.Assess(q)
			}
			So(f.ef, ShouldAlmostEqual, tc.ef)
			So(f.n, ShouldEqual, tc.n)
			So(f.interval, ShouldEqual, tc.interval)
		})
	}
}

func ExampleFact() {
	type Question struct {
		q, a string
		f    *Fact
	}

	q1 := &Question{"Capital of Canada?", "Ottawa", NewFact()}
	q1.f.Assess(5)                  // correct answer; immediate recall
	fmt.Println(q1.f.UpForReview()) // false

	q2 := &Question{"Capital of Ontario?", "Toronto", NewFact()}
	q2.f.Assess(0)                  // incorrect answer; no idea whatsoever.
	fmt.Println(q2.f.UpForReview()) // true

}
