#!/usr/bin/env python3
"""Fail only when repeated same-runner benchmark samples show a clear slowdown."""

import argparse
import random
import re
import statistics
import sys
from pathlib import Path

BENCH_LINE = re.compile(
    r"^(\S+)-\d+\s+\d+\s+([0-9.]+)\s+ns/op(?:\s+[0-9.]+\s+B/op)?"
    r"(?:\s+[0-9.]+\s+allocs/op)?\s*$"
)
BOOTSTRAP_RESAMPLES = 20_000
MIN_SAMPLES = 5


def read_samples(path: Path) -> dict[str, list[float]]:
    samples: dict[str, list[float]] = {}
    for line in path.read_text(encoding="utf-8").splitlines():
        match = BENCH_LINE.match(line)
        if match:
            samples.setdefault(match.group(1), []).append(float(match.group(2)))
    return samples


def lower_ratio_bound(baseline: list[float], current: list[float], seed: int) -> float:
    rng = random.Random(seed)
    ratios = []
    for _ in range(BOOTSTRAP_RESAMPLES):
        old = [rng.choice(baseline) for _ in baseline]
        new = [rng.choice(current) for _ in current]
        ratios.append(statistics.median(new) / statistics.median(old))
    ratios.sort()
    return ratios[int(0.025 * len(ratios))]


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("baseline", type=Path)
    parser.add_argument("current", type=Path)
    parser.add_argument("--max-regression", type=float, default=0.10)
    args = parser.parse_args()
    if not 0 <= args.max_regression < 1:
        parser.error("--max-regression must be in [0, 1)")

    baseline = read_samples(args.baseline)
    current = read_samples(args.current)
    names = sorted(baseline.keys() & current.keys())
    if not names:
        print("no matching benchmark samples found", file=sys.stderr)
        return 2

    failed = False
    evaluated = 0
    for index, name in enumerate(names):
        old, new = baseline[name], current[name]
        if len(old) < MIN_SAMPLES or len(new) < MIN_SAMPLES:
            print(f"SKIP {name}: need at least {MIN_SAMPLES} samples per revision")
            continue
        evaluated += 1
        old_median = statistics.median(old)
        new_median = statistics.median(new)
        lower = lower_ratio_bound(old, new, seed=index)
        print(
            f"{'FAIL' if lower > 1 + args.max_regression else 'PASS'} {name}: "
            f"median ratio={new_median / old_median:.3f}, "
            f"95% lower bound={lower:.3f}, threshold={1 + args.max_regression:.3f}"
        )
        failed |= lower > 1 + args.max_regression
    if evaluated == 0:
        print("no benchmarks had enough samples to evaluate", file=sys.stderr)
        return 2
    return int(failed)


if __name__ == "__main__":
    raise SystemExit(main())
