// Package em implements expectation-maximization algorithm
package em

import (
	"fmt"
	"math"

	"github.com/britojr/kbn/cliquetree"
	"github.com/britojr/kbn/factor"
	"github.com/britojr/kbn/filehandler"
	"github.com/britojr/kbn/utils"
)

const epslon = 1e-14

// ExpectationMaximization ..
func ExpectationMaximization(ct *cliquetree.CliqueTree, ds *filehandler.DataSet, norm bool) {
	diff := epslon * 10
	var err error
	for i := 1; diff >= epslon; i++ {
		fmt.Printf("Iteration: %v\n", i)
		newpot := expectationStep(ct, ds)
		if norm {
			for j := range newpot {
				newpot[j].Normalize()
			}
		}
		fmt.Printf("Count param: %v (%v)=0\n", newpot[0].Values()[0], newpot[0].Variables())
		diff, _, _, err = factor.MaxDifference(ct.BkpPotentialList(), newpot)
		utils.ErrCheck(err, "")
		fmt.Printf("current diff: %v\n", diff)
		ct.SetAllPotentials(newpot)
	}
}

// expectationStep ..
func expectationStep(ct *cliquetree.CliqueTree, ds *filehandler.DataSet) []*factor.Factor {
	// initialize counter
	count := make([]*factor.Factor, ct.Size())
	for i := range count {
		count[i] = ct.CurrPotential(i).ClearCopy()
	}

	// calculate probability of every instance
	for _, m := range ds.Data() {
		ct.ReduceByEvidence(m)
		ct.UpDownCalibration()
		// ct.LoadCalibration()
		checkCliqueTree(ct)
		for i := range count {
			ct.Calibrated(i).Normalize()
			for j, v := range ct.Calibrated(i).Values() {
				count[i].Values()[j] += v
				if math.IsNaN(count[i].Values()[j]) {
					panic(fmt.Sprintf("count %v, index %v is NaN", i, j))
				}
			}
		}
	}

	return count
}

func checkCliqueTree(ct *cliquetree.CliqueTree) {
	for i := range ct.BkpPotentialList() {
		f := ct.Calibrated(i)
		sum := 0.0
		for _, v := range f.Values() {
			sum += v
		}
		if sum == 0 {
			fmt.Printf("(%v)\n", f.Variables())
			panic("original zero factor")
		}
	}
}
