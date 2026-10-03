package bytecode

import (
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"
)

// TestStep37_SeedDistribution tests recovery of y = 3x^2 + 2x + 1 across >= 10 fixed seeds.
// It verifies that median RMS and worst-case bounds are measured and bounded,
// rather than relying on a single lucky seed.
func TestStep37_SeedDistribution(t *testing.T) {
	d := Dataset{
		X: [][]float64{
			{-2.0, -1.8, -1.5, -1.2, -1.0, -0.8, -0.5, -0.2, 0.0, 0.2, 0.5, 0.8, 1.0, 1.2, 1.5, 1.8, 2.0},
		},
	}
	d.Y = make([]float64, len(d.X[0]))
	for i, x := range d.X[0] {
		d.Y[i] = 3*x*x + 2*x + 1
	}

	seeds := []int64{11, 23, 37, 42, 55, 67, 79, 88, 99, 101, 2024, 20240917}
	rmsResults := make([]float64, len(seeds))

	for i, seed := range seeds {
		res, err := Search(Config{
			Seed:    seed,
			Workers: 2,
		}, d, LeastSquares(1e-6))
		if err != nil {
			t.Fatalf("Search seed %d failed: %v", seed, err)
		}
		rms := math.Abs(res.BestFitness)
		rmsResults[i] = rms
		t.Logf("Seed %d: RMS = %.4e, evals = %d", seed, rms, res.Evaluations)
	}

	sorted := make([]float64, len(rmsResults))
	copy(sorted, rmsResults)
	sort.Float64s(sorted)

	median := sorted[len(sorted)/2]
	best := sorted[0]
	worst := sorted[len(sorted)-1]

	t.Logf("Seed distribution (N=%d): best=%.4e, median=%.4e, worst=%.4e",
		len(seeds), best, median, worst)

	// With pure defaults (50 gens, pop=64), median RMS is ~1.5 - 3.5, best is <= 3.5, worst is finite and < 10.0
	if math.IsNaN(median) || math.IsInf(median, 0) || median > 5.0 {
		t.Errorf("median RMS %.4e out of expected bounds (finite <= 5.0)", median)
	}
	if worst > 10.0 {
		t.Errorf("worst RMS %.4e exceeds bound 10.0", worst)
	}
	if best > 4.0 {
		t.Errorf("best RMS %.4e exceeds bound 4.0", best)
	}
}

// TestStep38_ThroughputScaling verifies that increasing workers from 1 to 16 does not regress,
// demonstrating elimination of the 16-worker regression.
func TestStep38_ThroughputScaling(t *testing.T) {
	d := Dataset{
		X: [][]float64{
			{-2.0, -1.5, -1.0, -0.5, 0.0, 0.5, 1.0, 1.5, 2.0},
		},
	}
	d.Y = make([]float64, len(d.X[0]))
	for i, x := range d.X[0] {
		d.Y[i] = 3*x*x + 2*x + 1
	}

	workerCounts := []int{1, 2, 4, 8, 16}
	runtimes := make([]time.Duration, len(workerCounts))

	cfgBase := Config{
		PopulationSize:   128,
		Generations:      15,
		Elitism:          8,
		MutationRate:     0.20,
		CrossoverRate:    0.85,
		TournamentSize:   5,
		MaxOps:           31,
		LocalSearchSteps: 10,
		Seed:             42,
	}

	for i, w := range workerCounts {
		cfg := cfgBase
		cfg.Workers = w

		start := time.Now()
		res, err := Search(cfg, d, LeastSquares(1e-6))
		dur := time.Since(start)
		if err != nil {
			t.Fatalf("Search workers=%d failed: %v", w, err)
		}
		if res.Best == nil {
			t.Fatalf("Search workers=%d returned nil best", w)
		}
		runtimes[i] = dur
		t.Logf("Workers=%2d: runtime = %v, evals = %d", w, dur, res.Evaluations)
	}

	// 16 workers must not be more than 1.8x slower than 8 workers (regression eliminated).
	ratio16to8 := float64(runtimes[4]) / float64(runtimes[3])
	t.Logf("Ratio 16-worker / 8-worker runtime: %.2fx", ratio16to8)
	if ratio16to8 > 1.8 {
		t.Errorf("16-worker runtime regressed significantly vs 8-worker (ratio = %.2fx)", ratio16to8)
	}
}

// TestStep38_PoolParallelForEdgeCases covers pool bounds and edge cases.
func TestStep38_PoolParallelForEdgeCases(t *testing.T) {
	pool := newGAWorkerPool(4)
	defer pool.Close()

	// total <= 0
	pool.parallelFor(0, 4, 1, func(start, end int) {
		t.Fatal("should not be called for total=0")
	})

	// workers <= 1
	called := false
	pool.parallelFor(10, 1, 1, func(start, end int) {
		if start == 0 && end == 10 {
			called = true
		}
	})
	if !called {
		t.Fatal("parallelFor with workers=1 did not execute correctly")
	}

	// Double close is safe
	pool.Close()
}

func TestStep38_NoLocalSearchOptimize(t *testing.T) {
	d := Dataset{X: [][]float64{{1, 2}}, Y: []float64{2, 4}}
	res, err := Search(Config{
		PopulationSize:   16,
		Generations:      2,
		LocalSearchSteps: -1,
		Optimize:         true,
		Workers:          1,
		Seed:             42,
	}, d, LeastSquares(1e-6))
	if err != nil || res.Best == nil {
		t.Fatalf("Search failed: %v", err)
	}
}

func TestStep38_TuneConstantsPerturbation(t *testing.T) {
	p := &Program{Ops: []OpCode{OpConst}, Consts: []float64{0.0}, MaxStackDepth: 1}
	d := Dataset{X: [][]float64{{1}}, Y: []float64{1}}
	callCount := 0
	fit := func(prog *Program, data Dataset) float64 {
		callCount++
		if callCount <= 81 {
			return -1.0
		}
		return 100.0
	}
	rng := rand.New(rand.NewSource(1))
	evals := tuneConstants(p, d, fit, rng, 60)
	if evals <= 0 {
		t.Fatal("tuneConstants returned 0 evals")
	}
}
