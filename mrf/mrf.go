package mrf

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/britojr/kbn/factor"
	"github.com/britojr/utl/conv"
)

// Mrf markov random field
type Mrf struct {
	cardin     []int
	potentials []*factor.Factor
}

// LoadFromUAI creates a mrf from a reader in uai format
func LoadFromUAI(r io.Reader) *Mrf {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	scanner.Scan()
	// numvar := utils.Atoi(scanner.Text())
	scanner.Scan()
	cardin := conv.Satoi(strings.Fields(scanner.Text()))
	scanner.Scan()
	potentials := make([]*factor.Factor, conv.Atoi(scanner.Text()))
	for i := range potentials {
		scanner.Scan()
		varlist := conv.Satoi(strings.Fields(scanner.Text()))
		potentials[i] = factor.NewFactor(varlist[1:], cardin)
	}
	// here we have problem with different UAI formats
	scanner.Scan()
	if len(scanner.Text()) == 0 {
		for i := range potentials {
			scanner.Scan()
			scanner.Scan()
			potentials[i].SetValues(conv.Satof(strings.Fields(scanner.Text())))
			scanner.Scan()
		}
	} else {
		for i := range potentials {
			for j := range potentials[i].Values() {
				scanner.Scan()
				potentials[i].Values()[j] = conv.Atof(scanner.Text())
			}
			scanner.Scan()
			scanner.Scan()
		}
	}
	return &Mrf{cardin, potentials}
}

// UnnormalizedProb returns the "unnormalized probability" of given evidence
func (m *Mrf) UnnormalizedProb(evid []int) float64 {
	q := float64(1)
	for _, f := range m.potentials {
		q *= f.GetEvidValue(evid)
	}
	return q
}

// UnnormLogProb returns the log of the "unnormalized probability" of evidence
func (m *Mrf) UnnormLogProb(evid []int) (q float64) {
	for _, f := range m.potentials {
		q += math.Log(f.GetEvidValue(evid))
	}
	return
}

// SaveOnLibdaiFormat saves a Mrf in libDAI factor graph format on the given writer
func (m *Mrf) SaveOnLibdaiFormat(w io.Writer) {
	// number of potentials
	fmt.Fprintf(w, "%d\n", len(m.potentials))
	fmt.Fprintln(w)
	for _, p := range m.potentials {
		// number of variables
		fmt.Fprintf(w, "%d\n", len(p.Variables()))
		// variables
		for _, v := range p.Variables() {
			fmt.Fprintf(w, "%d ", v)
		}
		fmt.Fprintln(w)
		// cardinalities
		for _, v := range p.Variables() {
			fmt.Fprintf(w, "%d ", p.Cardinality()[v])
		}
		fmt.Fprintln(w)
		// number of factor values
		fmt.Fprintf(w, "%d\n", len(p.Values()))
		// factor values
		for j, v := range p.Values() {
			fmt.Fprintf(w, "%d     %v\n", j, v)
		}
		fmt.Fprintln(w)
	}
}

// Print prints all mrf values
func (m *Mrf) Print() {
	fmt.Println(m.cardin)
	for _, f := range m.potentials {
		fmt.Println(f.Variables(), f.Values())
	}
}
