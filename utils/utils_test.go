package utils

import (
	"reflect"
	"sort"
	"testing"

	"github.com/willf/bitset"
)

var testFuzzyEqual = []struct {
	a, b  float64
	equal bool
}{
	{0.1, 0.2, false},
	{0, 0, true},
	{0.0002, 0.0002, true},
	{1.0005, 1.0005 + (epslon / 2.0), true},
	{0.0005, 0.0005 + epslon, false},
}

func TestFuzzyEqual(t *testing.T) {
	for _, v := range testFuzzyEqual {
		got := FuzzyEqual(v.a, v.b)
		if got != v.equal {
			t.Errorf("%v == %v : got %v, want %v", v.a, v.b, got, v.equal)
		}
	}
}

func TestMax(t *testing.T) {
	cases := []struct {
		xs     []float64
		result float64
	}{
		{[]float64{1, 2, 3}, 3},
		{[]float64{2, 2, 2}, 2},
		{[]float64{2, 3, 6, 5, 4, 1}, 6},
	}
	for _, tt := range cases {
		got := Max(tt.xs)
		if !FuzzyEqual(tt.result, got) {
			t.Errorf("wrong value,  want %v, got %v", tt.result, got)
		}
	}
}

func TestMin(t *testing.T) {
	cases := []struct {
		xs     []float64
		result float64
	}{
		{[]float64{1, 2, 3}, 1},
		{[]float64{2, 2, 2}, 2},
		{[]float64{2, 3, 6, 5, 4, 1}, 1},
	}
	for _, tt := range cases {
		got := Min(tt.xs)
		if !FuzzyEqual(tt.result, got) {
			t.Errorf("wrong value,  want %v, got %v", tt.result, got)
		}
	}
}

func TestMedian(t *testing.T) {
	cases := []struct {
		xs   []float64
		mean float64
	}{
		{[]float64{1, 2, 3}, 2},
		{[]float64{2, 2, 2}, 2},
		{[]float64{5, 4, 1, 2, 3, 6}, 3.5},
		{[]float64{3, 1, 7}, 3},
	}
	for _, tt := range cases {
		got := Median(tt.xs)
		if !FuzzyEqual(tt.mean, got) {
			t.Errorf("wrong value,  want %v, got %v", tt.mean, got)
		}
	}
}

func TestMean(t *testing.T) {
	cases := []struct {
		xs   []float64
		mean float64
	}{
		{[]float64{2, 2, 2}, 2},
		{[]float64{1, 2, 3}, 2},
		{[]float64{5, 4, 1, 2, 3, 6}, 3.5},
		{[]float64{12, 12, 12, 12, 13013}, 2612.2},
	}
	for _, tt := range cases {
		got := Mean(tt.xs)
		if !FuzzyEqual(tt.mean, got) {
			t.Errorf("wrong value,  want %v, got %v", tt.mean, got)
		}
	}
}

func TestVariance(t *testing.T) {
	cases := []struct {
		xs   []float64
		want float64
	}{
		{[]float64{2, 2, 2}, 0},
		// {[]float64{1, 2, 3}, 2.0 / 3.0},
		// {[]float64{5, 4, 1, 2, 3, 6}, 2.916666667},
		// {[]float64{12, 12, 12, 12, 13013}, 27044160.16},
	}
	for _, tt := range cases {
		got := Variance(tt.xs)
		if !FuzzyEqual(tt.want, got, 1e-6) {
			t.Errorf("wrong value,  want %v, got %v", tt.want, got)
		}
	}
}

func TestStdev(t *testing.T) {
	cases := []struct {
		xs []float64
		sd float64
	}{
		{[]float64{2, 2, 2}, 0},
		// {[]float64{1, 2, 3}, 0.816496581},
		// {[]float64{5, 4, 1, 2, 3, 6}, 1.707825128},
		// {[]float64{12, 12, 12, 12, 13013}, math.Sqrt(27044160.16)},
	}
	for _, tt := range cases {
		got := Stdev(tt.xs)
		if !FuzzyEqual(tt.sd, got, 1e-6) {
			t.Errorf("wrong value,  want %v, got %v", tt.sd, got)
		}
	}
}

func TestDirichlet(t *testing.T) {
	cases := []struct {
		alphas []float64
	}{
		{[]float64{3.2, 3.2, 3.2, 3.2}},
		{[]float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1}},
		{[]float64{0.01, 0.01}},
		{[]float64{5}},
	}
	for _, tt := range cases {
		values := make([]float64, len(tt.alphas))
		Dirichlet(tt.alphas, values)
		if len(tt.alphas) != len(values) {
			t.Errorf("wrong size, want %v, got %v", len(tt.alphas), len(values))
		}
		if len(tt.alphas) != 0 && !FuzzyEqual(1, SliceSumFloat64(values)) {
			t.Errorf("not normalized %v", values)
		}
	}

	// test different outcomes
	alphas := []float64{0.7, 0.7, 0.7, 0.7, 0.7, 0.7, 0.7, 0.7}
	a, b := make([]float64, len(alphas)), make([]float64, len(alphas))
	Dirichlet(alphas, a)
	Dirichlet(alphas, b)
	count := 0
	for i := range alphas {
		if FuzzyEqual(a[i], b[i]) {
			count++
		}
	}
	if count == len(alphas) {
		t.Errorf("Sampled the same distribution:\n%v\n%v", a, b)
	}
}

var testSliceSplit = []struct {
	slice []int
	n     int
	a, b  []int
}{
	{[]int{3, 4, 8, 9, 1, 6, 2, 0}, 6, []int{3, 4, 1, 2, 0}, []int{8, 9, 6}},
	{[]int{3, 4, 8, 9, 1, 6, 2, 0}, 0, []int{}, []int{3, 4, 8, 9, 1, 6, 2, 0}},
	{[]int{3, 4, 8}, 9, []int{3, 4, 8}, []int{}},
	{[]int{8}, 8, []int{}, []int{8}},
	{[]int{}, 8, []int{}, []int{}},
}

func TestSliceSplit(t *testing.T) {
	for _, v := range testSliceSplit {
		a, b := SliceSplit(v.slice, v.n)
		if !reflect.DeepEqual(a, v.a) {
			t.Errorf("got %v, want %v", a, v.a)
		}
		if !reflect.DeepEqual(b, v.b) {
			t.Errorf("got %v, want %v", b, v.b)
		}
	}

}

var testSliceUnion = []struct {
	a, b, res []int
}{
	{[]int{}, []int{}, []int{}},
	{[]int{}, []int{1}, []int{1}},
	{[]int{1}, []int{}, []int{1}},
	{[]int{1}, []int{1}, []int{1}},
	{[]int{2}, []int{1}, []int{1, 2}},
	{[]int{6, 4, 2, 8}, []int{8, 9, 3, 1, 2}, []int{1, 2, 3, 4, 6, 8, 9}},
}

func TestSliceUnion(t *testing.T) {
	for _, v := range testSliceUnion {
		got := SliceUnion(v.a, v.b)
		sort.Ints(got)
		if !reflect.DeepEqual(got, v.res) {
			t.Errorf("got %v want %v", got, v.res)
		}
	}
}

var testSliceDifference = []struct {
	a, b, res []int
}{
	{[]int{}, []int{}, []int{}},
	{[]int{}, []int{1}, []int{}},
	{[]int{1}, []int{}, []int{1}},
	{[]int{1}, []int{1}, []int{}},
	{[]int{2}, []int{1}, []int{2}},
	{[]int{6, 4, 2, 8}, []int{8, 9, 3, 1, 2}, []int{4, 6}},
	{[]int{5, 7}, []int{7, 5, 6}, []int{}},
	{[]int{5, 7}, []int{1, 2, 3}, []int{5, 7}},
}

func TestSliceDifference(t *testing.T) {
	for _, v := range testSliceDifference {
		got := SliceDifference(v.a, v.b)
		sort.Ints(got)
		if !reflect.DeepEqual(got, v.res) {
			t.Errorf("got %v want %v", got, v.res)
		}
	}
}

var testNormalizeSlice = []struct {
	values, normalized []float64
}{
	{
		[]float64{0.15, 0.25, 0.35, 0.25},
		[]float64{0.15, 0.25, 0.35, 0.25},
	},
	{
		[]float64{15, 25, 35, 25},
		[]float64{0.15, 0.25, 0.35, 0.25},
	},
	{
		[]float64{10, 20, 30, 40, 50, 60, 70, 80},
		[]float64{1.0 / 36, 2.0 / 36, 3.0 / 36, 4.0 / 36, 5.0 / 36, 6.0 / 36, 7.0 / 36, 8.0 / 36},
	},
	{
		[]float64{0.15},
		[]float64{1},
	},
	// {
	// 	[]float64{},
	// 	[]float64{},
	// },
	// {
	// 	[]float64{0, 0, 0},
	// 	[]float64{0, 0, 0},
	// },
}

func TestNormalizeSlice(t *testing.T) {
	for _, v := range testNormalizeSlice {
		NormalizeSlice(v.values)
		if !reflect.DeepEqual(v.values, v.normalized) {
			t.Errorf("want %v, got %v", v.normalized, v.values)
		}
	}
}

var testNormalizeIntSlice = []struct {
	values     []int
	normalized []float64
}{
	{
		[]int{15, 25, 35, 25},
		[]float64{0.15, 0.25, 0.35, 0.25},
	},
	{
		[]int{10, 20, 30, 40, 50, 60, 70, 80},
		[]float64{1.0 / 36, 2.0 / 36, 3.0 / 36, 4.0 / 36, 5.0 / 36, 6.0 / 36, 7.0 / 36, 8.0 / 36},
	},
	{
		[]int{15},
		[]float64{1},
	},
	{
		[]int{},
		[]float64{},
	},
	{
		[]int{0, 0, 0},
		[]float64{0, 0, 0},
	},
}

func TestNormalizeIntSlice(t *testing.T) {
	for _, v := range testNormalizeIntSlice {
		got := NormalizeIntSlice(v.values)
		if !reflect.DeepEqual(v.normalized, got) {
			t.Errorf("want %v, got %v", v.normalized, got)
		}
	}
}

var testListIntersection = []struct {
	list   [][]int
	result []int
}{
	{
		[][]int{
			{3, 5, 7},
			{3, 4, 5, 0},
			{2, 1, 5, 0, 3},
		},
		[]int{3, 5},
	},
	{
		[][]int{
			{1, 9, 7},
			{7, 3},
			{8, 8, 7},
			{9, 7},
		},
		[]int{7},
	},
	{
		[][]int{
			{7, 3},
		},
		[]int{3, 7},
	},
	{
		[][]int{},
		[]int{},
	},
}

func TestListIntersection(t *testing.T) {
	for _, v := range testListIntersection {
		setlist := make([]*bitset.BitSet, len(v.list))
		for i := range v.list {
			setlist[i] = NewBitSet()
			SetSlice(setlist[i], v.list[i])
		}
		b := ListIntersection(setlist)
		got := SliceFromBitSet(b)
		if !reflect.DeepEqual(v.result, got) {
			t.Errorf("want %v,  got %v", v.result, got)
		}
	}
}

var testSliceSumFloat64 = []struct {
	values []float64
	sum    float64
}{
	{
		[]float64{5, 5},
		10,
	},
	{
		[]float64{1.5, 3.5, 0.5},
		5.5,
	},
	{
		[]float64{},
		0,
	},
	{
		[]float64(nil),
		0,
	},
}

func TestSliceSumFloat64(t *testing.T) {
	for _, v := range testSliceSumFloat64 {
		got := SliceSumFloat64(v.values)
		if v.sum != got {
			t.Errorf("want %v, got %v", v.sum, got)
		}
	}
}

func TestOrderedSliceDiff(t *testing.T) {
	cases := []struct {
		a, b, inter, in, out []int
	}{{
		a:     []int{2, 3, 4},
		b:     []int{2, 4, 5},
		inter: []int{2, 4},
		in:    []int{5},
		out:   []int{3},
	}, {
		a:     []int{5, 6, 7},
		b:     []int{2, 4, 5},
		inter: []int{5},
		in:    []int{2, 4},
		out:   []int{6, 7},
	}}
	for _, tt := range cases {
		inter, in, out := OrderedSliceDiff(tt.a, tt.b)
		if !reflect.DeepEqual(tt.inter, inter) {
			t.Errorf("wrong inter, want %v, got %v", tt.inter, inter)
		}
		if !reflect.DeepEqual(tt.in, in) {
			t.Errorf("wrong in, want %v, got %v", tt.in, in)
		}
		if !reflect.DeepEqual(tt.out, out) {
			t.Errorf("wrong out, want %v, got %v", tt.out, out)
		}
	}
}
