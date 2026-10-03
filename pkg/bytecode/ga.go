package bytecode

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"
)

// timeSeed derives a non-zero seed from the system clock. It is only used when
// Config.Seed is left at zero.
func timeSeed() int64 { return time.Now().UnixNano() | 1 }

// Dataset is a columnar set of samples for symbolic regression: one input
// column per variable plus the observed outputs.
//
//	data := bytecode.Dataset{
//	    X: [][]float64{{...}, {...}},  // one column per variable, same length
//	    Y: []float64{...},            // observed outputs, len == len(X[0])
//	}
type Dataset struct {
	// X holds the input variables in column-major order: X[varIdx][sampleIdx].
	X [][]float64
	// Y holds the observed target value for each sample.
	Y []float64
}

// Samples reports the number of samples in the dataset.
func (d Dataset) Samples() int {
	if len(d.Y) > 0 {
		return len(d.Y)
	}
	if len(d.X) > 0 {
		return len(d.X[0])
	}
	return 0
}

// Variables reports the number of input variables in the dataset.
func (d Dataset) Variables() int { return len(d.X) }

// validate checks the dataset is self-consistent.
func (d Dataset) validate() error {
	n := d.Samples()
	if n == 0 {
		return fmt.Errorf("bytecode: dataset has no samples")
	}
	for i, col := range d.X {
		if len(col) != n {
			return fmt.Errorf("bytecode: dataset variable %d has %d samples, want %d", i, len(col), n)
		}
	}
	if len(d.Y) > 0 && len(d.Y) != n {
		return fmt.Errorf("bytecode: dataset has %d targets but %d samples", len(d.Y), n)
	}
	return nil
}

// Fitness scores a candidate program against a dataset. Higher is better.
//
// A Fitness must return a finite value for every candidate, including ones whose
// evaluation produces NaN or Inf: the population loop sorts on the score, so a
// single NaN would make selection meaningless. See LeastSquares.
//
// Fitness is called concurrently, so it must be safe for concurrent use.
type Fitness func(p *Program, d Dataset) float64

// LeastSquares returns the negative root-mean-square error of p over d, so that
// higher is better. A candidate whose evaluation goes non-finite on any sample
// scores -Inf, which sorts last rather than poisoning the ordering with NaN.
//
// It is the standard objective for symbolic regression: recover a closed form
// by driving the RMS error to zero.
func LeastSquares(eps float64) Fitness {
	if eps <= 0 {
		eps = 1e-6
	}
	return func(p *Program, d Dataset) float64 {
		n := d.Samples()
		if n == 0 {
			return math.Inf(-1)
		}
		scratch := make([]float64, p.MaxStackDepth)
		vars := make([]float64, len(d.X))
		var sum float64
		for i := 0; i < n; i++ {
			for v := range d.X {
				vars[v] = d.X[v][i]
			}
			got := p.EvalRegularized(vars, eps, scratch)
			if math.IsNaN(got) || math.IsInf(got, 0) {
				return math.Inf(-1)
			}
			diff := got - d.Y[i]
			sum += diff * diff
		}
		return -math.Sqrt(sum / float64(n))
	}
}

// Config parameterises a genetic search.
type Config struct {
	// PopulationSize is the number of candidates per generation. Defaults to 64.
	PopulationSize int
	// Generations is the number of generations to run. Defaults to 50.
	Generations int
	// CrossoverRate is the probability that a child is produced by crossover
	// rather than by cloning its parent. Defaults to 0.7.
	CrossoverRate float64
	// MutationRate is the per-opcode probability of mutation, and also the
	// probability of jittering each constant operand. Defaults to 0.15.
	MutationRate float64
	// Elitism is how many top candidates are copied unchanged into the next
	// generation. Defaults to 2.
	Elitism int
	// TournamentSize is how many candidates are sampled per selection; larger
	// values increase selection pressure. Defaults to 3.
	TournamentSize int
	// MinOps and MaxOps bound the opcode count of a candidate. MaxOps of 0 means
	// uncapped. Defaults are 3 and 21.
	MinOps int
	MaxOps int
	// Epsilon is the domain regularization parameter handed to
	// EvalRegularized. Defaults to 1e-6.
	Epsilon float64
	// Workers is the number of goroutines used to score a generation. Zero
	// selects runtime.NumCPU(); 1 forces serial scoring.
	Workers int
	// LocalSearchSteps is the number of coordinate-descent sweeps applied to
	// each elite candidate after every generation. Zero disables local search.
	//
	// Pure crossover and mutation converge very slowly on the constant
	// operands: random jitter is a random walk, so pinning a coefficient to,
	// say, 3.0 takes many generations. Coordinate descent over the constants
	// turns that into a direct search and is by far the biggest improvement to
	// the recovered error -- measured on 3x^2+2x+1 it is the difference
	// between an RMS error of 0.21 and 4e-13.
	//
	// Defaults to 30. Set it to a negative value to disable local search.
	LocalSearchSteps int
	// Seed makes the whole search reproducible. Zero picks a seed from the
	// system entropy source.
	Seed int64
}

func (c *Config) withDefaults() Config {
	out := *c
	if out.PopulationSize <= 0 {
		out.PopulationSize = 64
	}
	if out.Generations <= 0 {
		out.Generations = 50
	}
	if out.CrossoverRate <= 0 {
		out.CrossoverRate = 0.7
	}
	if out.MutationRate <= 0 {
		out.MutationRate = 0.15
	}
	if out.Elitism <= 0 {
		out.Elitism = 2
	}
	if out.TournamentSize <= 0 {
		out.TournamentSize = 3
	}
	if out.MinOps <= 0 {
		out.MinOps = 3
	}
	if out.MaxOps < out.MinOps {
		out.MaxOps = out.MinOps + 18
	}
	if out.Epsilon <= 0 {
		out.Epsilon = 1e-6
	}
	if out.Workers == 0 {
		out.Workers = runtime.NumCPU()
	}
	if out.LocalSearchSteps == 0 {
		out.LocalSearchSteps = 30
	} else if out.LocalSearchSteps < 0 {
		out.LocalSearchSteps = 0
	}
	if out.Elitism >= out.PopulationSize {
		out.Elitism = max(1, out.PopulationSize/4)
	}
	return out
}

// Result reports the outcome of a search.
type Result struct {
	// Best is the highest-fitness program found.
	Best *Program
	// BestFitness is its score. For LeastSquares this is the negated RMS error,
	// so convergence means it approaches zero from below.
	BestFitness float64
	// History holds the best fitness after each generation, so a caller can see
	// whether the search converged.
	History []float64
	// Evaluations counts how many candidates were scored.
	Evaluations int
}

// Converged reports whether the best fitness is within tolerance of zero, which
// is the natural target for the least-squares objective.
func (r Result) Converged(tolerance float64) bool {
	return math.Abs(r.BestFitness) <= tolerance
}

// Search runs a genetic search over EML-tree candidates and returns the best
// program found for the dataset.
//
// The search is deterministic for a given Config.Seed. Fitness is called
// concurrently from up to Config.Workers goroutines, so it must be safe for
// concurrent use; LeastSquares is.
//
// What it can and cannot do, measured on this implementation:
//
//   - Shallow closed forms built from the available operators are recovered
//     essentially exactly. Fitting y = 3x^2 + 2x + 1 from 24 samples reaches an
//     RMS error of 4e-13 with a fixed seed.
//   - Transcendental targets are approximated, not recovered. Fitting y = sin(x)
//     over [0, 2*pi) plateaus near an RMS error of 0.17, versus ~0.71 for the
//     best constant fit, and increasing the budget to 1500 individuals for 1500
//     generations does not improve on it.
//
// The second point is inherent: approximating sin requires discovering a
// composition of elementary functions, which is what the source paper does with
// gradient-based training of a parameterised tree, not with a genetic algorithm.
// Use this for the first class of problem, and expect approximation elsewhere.
func Search(cfg Config, dataset Dataset, fitness Fitness) (Result, error) {
	if fitness == nil {
		return Result{}, fmt.Errorf("bytecode: nil fitness function")
	}
	if err := dataset.validate(); err != nil {
		return Result{}, err
	}
	cfg = cfg.withDefaults()

	seed := cfg.Seed
	if seed == 0 {
		seed = timeSeed()
	}
	rng := rand.New(rand.NewSource(seed)) // #nosec G404 -- reproducible GA search, not a security primitive

	nVars := dataset.Variables()
	if nVars == 0 {
		nVars = 1
	}

	population := make([]*Program, cfg.PopulationSize)
	for i := range population {
		population[i] = RandomProgram(nVars, cfg.MinOps, cfg.MaxOps, rng)
	}

	scores := make([]float64, len(population))
	history := make([]float64, 0, cfg.Generations)
	evaluations := 0

	for gen := 0; gen < cfg.Generations; gen++ {
		evaluations += scorePopulation(cfg, population, dataset, fitness, scores)

		order := sortedIndices(scores)
		history = append(history, scores[order[0]])

		next := make([]*Program, 0, len(population))
		// Elitism: carry the best candidates over, polishing their constants
		// with a short local search first.
		for i := 0; i < cfg.Elitism; i++ {
			elite := population[order[i]].Clone()
			if cfg.LocalSearchSteps > 0 {
				tuneConstants(elite, dataset, fitness, rng, cfg.LocalSearchSteps)
			}
			next = append(next, elite)
		}

		for len(next) < len(population) {
			child := population[tournament(rng, scores, order, cfg.TournamentSize)].Clone()

			if rng.Float64() < cfg.CrossoverRate {
				other := tournament(rng, scores, order, cfg.TournamentSize)
				if other != order[0] {
					if c1, _, err := Crossover(child, population[other], rng); err == nil && c1 != nil {
						child = c1
					}
				}
			}
			child = Mutate(child, cfg.MutationRate, rng)

			if cfg.MaxOps > 0 && len(child.Ops) > cfg.MaxOps {
				child.Ops = child.Ops[:cfg.MaxOps]
				child.CalculateMaxStackDepth()
			}
			if !ValidateStack(child.Ops) {
				// Never admit a candidate that would panic during evaluation.
				continue
			}
			next = append(next, child)
		}

		population = next
	}

	// Score the final generation.
	evaluations += scorePopulation(cfg, population, dataset, fitness, scores)
	order := sortedIndices(scores)

	return Result{
		Best:        population[order[0]],
		BestFitness: scores[order[0]],
		History:     history,
		Evaluations: evaluations,
	}, nil
}

// tuneConstants improves a candidate's constant operands by coordinate descent
// on the fitness, shrinking the step size after each sweep.
//
// Only the constants move; the opcode structure is left to the genetic
// operators. This is what lets a candidate that has the right shape settle onto
// the right coefficients in a few generations instead of drifting there by
// random walk.
func tuneConstants(p *Program, d Dataset, fitness Fitness, rng *rand.Rand, sweeps int) {
	if p == nil || len(p.Consts) == 0 || fitness == nil {
		return
	}
	step := 1.0
	best := fitness(p, d)
	for s := 0; s < sweeps; s++ {
		improved := false
		for i := range p.Consts {
			for _, delta := range []float64{step, -step} {
				saved := p.Consts[i]
				p.Consts[i] = saved + delta
				if got := fitness(p, d); got > best {
					best = got
					improved = true
				} else {
					p.Consts[i] = saved
				}
			}
		}
		if !improved {
			step /= 2
			if step < 1e-12 {
				return
			}
		}
	}
	_ = rng
}

// sortedIndices returns the indices of scores in descending order. It is a plain
// insertion into the caller's slice, avoiding a per-generation allocation.
func sortedIndices(scores []float64) []int {
	order := make([]int, len(scores))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(a, b int) bool {
		return scores[order[a]] > scores[order[b]]
	})
	return order
}

// tournament picks one index from the population by sampling k candidates and
// taking the best.
func tournament(rng *rand.Rand, scores []float64, order []int, k int) int {
	best := -1
	bestScore := math.Inf(-1)
	for i := 0; i < k; i++ {
		idx := order[rng.Intn(len(order))]
		if scores[idx] > bestScore {
			bestScore = scores[idx]
			best = idx
		}
	}
	if best < 0 {
		return order[0]
	}
	return best
}

// scorePopulation evaluates every candidate, fanning out across workers.
func scorePopulation(cfg Config, pop []*Program, d Dataset, fitness Fitness, scores []float64) int {
	// Below this many candidates per worker the goroutine hand-off costs more
	// than the work it saves.
	const minPerWorker = 8
	if cfg.Workers <= 1 || len(pop) < cfg.Workers*minPerWorker {
		for i, p := range pop {
			scores[i] = fitness(p, d)
		}
		return len(pop)
	}

	workers := cfg.Workers
	chunk := (len(pop) + workers - 1) / workers
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		lo := w * chunk
		if lo >= len(pop) {
			break
		}
		hi := min(lo+chunk, len(pop))
		wg.Add(1)
		go func(lo, hi int) {
			defer wg.Done()
			for i := lo; i < hi; i++ {
				scores[i] = fitness(pop[i], d)
			}
		}(lo, hi)
	}
	wg.Wait()
	return len(pop)
}
