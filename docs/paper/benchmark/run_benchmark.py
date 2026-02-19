#!/usr/bin/env python3
"""
3D Bin Packing Benchmark Runner

This script runs comprehensive benchmarks on the packing algorithm
and generates results for academic paper analysis.

Usage:
    python run_benchmark.py [--iterations N] [--scenarios S1,S2,...]
                            [--containers 20ft_standard,40ft_standard,40ft_high_cube]
                            [--variants baseline,no_stability,smaller_first]

Output:
    - results/results.json: Raw benchmark data
    - results/results.csv: Summary table
    - results/latex_tables.tex: LaTeX-formatted tables
    - results/figures/*.pdf: Publication-quality plots
"""

from __future__ import annotations

import argparse
import csv
import json
import os
import sys
import time
from dataclasses import dataclass, field
from datetime import datetime
from pathlib import Path
from statistics import mean, stdev
from typing import Any

# Add packing module to path
SCRIPT_DIR = Path(__file__).parent.resolve()
REPO_ROOT = SCRIPT_DIR.parents[2]
PACKING_DIR = REPO_ROOT / "cmd" / "packing"
sys.path.insert(0, str(PACKING_DIR))

from packing import pack_request
from scenarios import SCENARIOS, CONTAINERS, Scenario, build_pack_request, PRODUCTS


# =============================================================================
# DATA STRUCTURES
# =============================================================================

@dataclass
class SingleRunResult:
    """Result of a single benchmark run."""
    scenario_id: str
    container_key: str
    iteration: int
    variant: str  # e.g., "baseline", "no_stable", "bigger_first_off"

    # Packing results
    expanded_items: int
    fitted_count: int
    unfitted_count: int

    # Time metrics (ms)
    pack_time_ms: int
    total_time_ms: int

    # Calculated metrics
    volume_utilization_pct: float
    weight_utilization_pct: float
    fill_rate_pct: float

    # Raw placement data (for validation)
    placements: list[dict] = field(default_factory=list)


@dataclass
class AggregatedResult:
    """Aggregated statistics for a scenario+container+variant combination."""
    scenario_id: str
    container_key: str
    variant: str
    num_runs: int

    # Item counts
    total_items: int
    avg_fitted: float
    std_fitted: float

    # Utilization stats
    avg_volume_util: float
    std_volume_util: float
    min_volume_util: float
    max_volume_util: float

    avg_weight_util: float
    std_weight_util: float

    # Fill rate stats
    avg_fill_rate: float
    std_fill_rate: float

    # Time stats (ms)
    avg_time: float
    std_time: float
    min_time: float
    max_time: float


# =============================================================================
# BENCHMARK RUNNER
# =============================================================================

def calculate_volume_utilization(
    placements: list[dict],
    scenario: Scenario,
    request_items: list[dict],
) -> float:
    """
    Calculate volume utilization percentage.

    Volume utilization = (sum of fitted item volumes) / container_volume * 100
    """
    # Build item dimension lookup (item_id -> dimensions in mm)
    item_dims: dict[str, tuple[float, float, float, float]] = {}
    for item in request_items:
        item_dims[item["item_id"]] = (
            item["length"],
            item["width"],
            item["height"],
            item["weight"],
        )

    # Calculate fitted volume
    fitted_volume_mm3 = 0.0
    for p in placements:
        item_id = p["item_id"]
        if item_id in item_dims:
            l, w, h, _ = item_dims[item_id]
            fitted_volume_mm3 += l * w * h

    container_volume_mm3 = scenario.get_container_volume_mm3()

    if container_volume_mm3 == 0:
        return 0.0

    return (fitted_volume_mm3 / container_volume_mm3) * 100


def calculate_weight_utilization(
    placements: list[dict],
    scenario: Scenario,
    request_items: list[dict],
) -> float:
    """
    Calculate weight utilization percentage.

    Weight utilization = (sum of fitted item weights) / max_weight * 100
    """
    # Build item weight lookup
    item_weights: dict[str, float] = {}
    for item in request_items:
        item_weights[item["item_id"]] = item["weight"]

    # Calculate fitted weight
    fitted_weight_kg = 0.0
    for p in placements:
        item_id = p["item_id"]
        if item_id in item_weights:
            fitted_weight_kg += item_weights[item_id]

    max_weight = scenario.container["max_weight"]

    if max_weight == 0:
        return 0.0

    return (fitted_weight_kg / max_weight) * 100


def run_single_benchmark(
    scenario: Scenario,
    iteration: int,
    variant: str,
    *,
    bigger_first: bool = True,
    check_stable: bool = True,
    support_surface_ratio: float = 0.75,
) -> SingleRunResult:
    """Run a single benchmark iteration."""

    # Build request (uses scenario.container_key internally)
    request = build_pack_request(
        scenario,
        bigger_first=bigger_first,
        check_stable=check_stable,
        support_surface_ratio=support_surface_ratio,
    )

    # Run packing
    start_time = time.perf_counter()
    result = pack_request(request)
    total_time_ms = int((time.perf_counter() - start_time) * 1000)

    # Extract stats
    data = result["data"]
    stats = data["stats"]
    placements = data["placements"]

    # Calculate utilization metrics
    volume_util = calculate_volume_utilization(placements, scenario, request["items"])
    weight_util = calculate_weight_utilization(placements, scenario, request["items"])

    # Calculate fill rate
    expanded = stats["expanded_items"]
    fitted = stats["fitted_count"]
    fill_rate = (fitted / expanded * 100) if expanded > 0 else 0.0

    return SingleRunResult(
        scenario_id=scenario.id,
        container_key=scenario.container_key,
        iteration=iteration,
        variant=variant,
        expanded_items=expanded,
        fitted_count=fitted,
        unfitted_count=stats["unfitted_count"],
        pack_time_ms=stats["pack_time_ms"],
        total_time_ms=total_time_ms,
        volume_utilization_pct=volume_util,
        weight_utilization_pct=weight_util,
        fill_rate_pct=fill_rate,
        placements=placements,
    )


def aggregate_results(runs: list[SingleRunResult]) -> AggregatedResult:
    """Aggregate multiple runs into statistics."""
    if not runs:
        raise ValueError("Cannot aggregate empty list of runs")

    scenario_id = runs[0].scenario_id
    container_key = runs[0].container_key
    variant = runs[0].variant

    fitted_counts = [r.fitted_count for r in runs]
    volume_utils = [r.volume_utilization_pct for r in runs]
    weight_utils = [r.weight_utilization_pct for r in runs]
    fill_rates = [r.fill_rate_pct for r in runs]
    times = [r.pack_time_ms for r in runs]

    n = len(runs)

    return AggregatedResult(
        scenario_id=scenario_id,
        container_key=container_key,
        variant=variant,
        num_runs=n,
        total_items=runs[0].expanded_items,
        avg_fitted=mean(fitted_counts),
        std_fitted=stdev(fitted_counts) if n > 1 else 0.0,
        avg_volume_util=mean(volume_utils),
        std_volume_util=stdev(volume_utils) if n > 1 else 0.0,
        min_volume_util=min(volume_utils),
        max_volume_util=max(volume_utils),
        avg_weight_util=mean(weight_utils),
        std_weight_util=stdev(weight_utils) if n > 1 else 0.0,
        avg_fill_rate=mean(fill_rates),
        std_fill_rate=stdev(fill_rates) if n > 1 else 0.0,
        avg_time=mean(times),
        std_time=stdev(times) if n > 1 else 0.0,
        min_time=min(times),
        max_time=max(times),
    )


# =============================================================================
# BENCHMARK VARIANTS
# =============================================================================

VARIANTS = {
    "baseline": {
        "bigger_first": True,
        "check_stable": True,
        "support_surface_ratio": 0.75,
    },
    "no_stability": {
        "bigger_first": True,
        "check_stable": False,
        "support_surface_ratio": 0.75,
    },
    "smaller_first": {
        "bigger_first": False,
        "check_stable": True,
        "support_surface_ratio": 0.75,
    },
}

# All available container keys
ALL_CONTAINERS = list(CONTAINERS.keys())  # ["20ft_standard", "40ft_standard", "40ft_high_cube"]


# =============================================================================
# OUTPUT GENERATORS
# =============================================================================

def save_json_results(
    all_runs: list[SingleRunResult],
    aggregated: list[AggregatedResult],
    output_dir: Path,
) -> None:
    """Save results to JSON file."""

    # Convert to serializable format
    runs_data = []
    for r in all_runs:
        runs_data.append({
            "scenario_id": r.scenario_id,
            "container_key": r.container_key,
            "iteration": r.iteration,
            "variant": r.variant,
            "expanded_items": r.expanded_items,
            "fitted_count": r.fitted_count,
            "unfitted_count": r.unfitted_count,
            "pack_time_ms": r.pack_time_ms,
            "total_time_ms": r.total_time_ms,
            "volume_utilization_pct": round(r.volume_utilization_pct, 2),
            "weight_utilization_pct": round(r.weight_utilization_pct, 2),
            "fill_rate_pct": round(r.fill_rate_pct, 2),
        })

    agg_data = []
    for a in aggregated:
        agg_data.append({
            "scenario_id": a.scenario_id,
            "container_key": a.container_key,
            "variant": a.variant,
            "num_runs": a.num_runs,
            "total_items": a.total_items,
            "avg_fitted": round(a.avg_fitted, 1),
            "std_fitted": round(a.std_fitted, 2),
            "avg_volume_util": round(a.avg_volume_util, 2),
            "std_volume_util": round(a.std_volume_util, 2),
            "min_volume_util": round(a.min_volume_util, 2),
            "max_volume_util": round(a.max_volume_util, 2),
            "avg_weight_util": round(a.avg_weight_util, 2),
            "std_weight_util": round(a.std_weight_util, 2),
            "avg_fill_rate": round(a.avg_fill_rate, 2),
            "std_fill_rate": round(a.std_fill_rate, 2),
            "avg_time_ms": round(a.avg_time, 1),
            "std_time_ms": round(a.std_time, 2),
            "min_time_ms": round(a.min_time, 1),
            "max_time_ms": round(a.max_time, 1),
        })

    output = {
        "generated_at": datetime.now().isoformat(),
        "total_runs": len(all_runs),
        "raw_results": runs_data,
        "aggregated_results": agg_data,
    }

    output_path = output_dir / "results.json"
    with open(output_path, "w") as f:
        json.dump(output, f, indent=2)

    print(f"  Saved JSON results to {output_path}")


def save_csv_results(aggregated: list[AggregatedResult], output_dir: Path) -> None:
    """Save aggregated results to CSV file."""

    output_path = output_dir / "results.csv"

    with open(output_path, "w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow([
            "Scenario",
            "Container",
            "Variant",
            "Total Items",
            "Avg Fitted",
            "Avg Volume Util (%)",
            "Std Volume Util",
            "Avg Weight Util (%)",
            "Avg Fill Rate (%)",
            "Avg Time (ms)",
            "Std Time (ms)",
        ])

        for a in aggregated:
            container_name = CONTAINERS[a.container_key]["name"]
            writer.writerow([
                a.scenario_id,
                container_name,
                a.variant,
                a.total_items,
                f"{a.avg_fitted:.1f}",
                f"{a.avg_volume_util:.2f}",
                f"{a.std_volume_util:.2f}",
                f"{a.avg_weight_util:.2f}",
                f"{a.avg_fill_rate:.2f}",
                f"{a.avg_time:.1f}",
                f"{a.std_time:.2f}",
            ])

    print(f"  Saved CSV results to {output_path}")


def generate_latex_tables(aggregated: list[AggregatedResult], output_dir: Path) -> None:
    """Generate LaTeX-formatted tables for the paper."""

    output_path = output_dir / "latex_tables.tex"

    # Helpers
    baseline_hc = [
        a for a in aggregated
        if a.variant == "baseline" and a.container_key == "40ft_high_cube"
    ]
    baseline_all_containers = [
        a for a in aggregated if a.variant == "baseline"
    ]

    lines = []
    lines.append("% Auto-generated LaTeX tables from benchmark results")
    lines.append(f"% Generated at: {datetime.now().isoformat()}")
    lines.append("")

    # -------------------------------------------------------------------------
    # Table 1: Main Results (Baseline, 40ft HC) — primary results
    # -------------------------------------------------------------------------
    lines.append("% =============================================================================")
    lines.append("% TABLE 1: Main Performance Results (Baseline, 40ft High Cube)")
    lines.append("% =============================================================================")
    lines.append(r"\begin{table}[H]")
    lines.append(r"    \centering")
    lines.append(r"    \caption{Hasil Pengujian Performa Algoritma pada Kontainer 40ft High Cube}")
    lines.append(r"    \label{tab:results}")
    lines.append(r"    \small")
    lines.append(r"    \begin{tabular}{lrrrrr}")
    lines.append(r"        \toprule")
    lines.append(r"        \textbf{Skenario} & \textbf{Item} & \textbf{Vol. Util. (\%)} & \textbf{Berat Util. (\%)} & \textbf{Fill Rate (\%)} & \textbf{Waktu (ms)} \\")
    lines.append(r"        \midrule")

    for a in baseline_hc:
        lines.append(
            f"        {a.scenario_id} & {a.total_items} & "
            f"{a.avg_volume_util:.2f} $\\pm$ {a.std_volume_util:.2f} & "
            f"{a.avg_weight_util:.2f} $\\pm$ {a.std_weight_util:.2f} & "
            f"{a.avg_fill_rate:.2f} & "
            f"{a.avg_time:.1f} $\\pm$ {a.std_time:.1f} \\\\"
        )

    lines.append(r"        \bottomrule")
    lines.append(r"    \end{tabular}")
    lines.append(r"\end{table}")
    lines.append("")

    # -------------------------------------------------------------------------
    # Table 2: Algorithm Variant Comparison (S3, 40ft HC)
    # -------------------------------------------------------------------------
    lines.append("% =============================================================================")
    lines.append("% TABLE 2: Algorithm Variant Comparison (S3, 40ft High Cube)")
    lines.append("% =============================================================================")
    lines.append(r"\begin{table}[H]")
    lines.append(r"    \centering")
    lines.append(r"    \caption{Perbandingan Varian Algoritma pada Skenario S3 (150 item)}")
    lines.append(r"    \label{tab:variants}")
    lines.append(r"    \small")
    lines.append(r"    \begin{tabular}{lrrrr}")
    lines.append(r"        \toprule")
    lines.append(r"        \textbf{Varian} & \textbf{Vol. Util. (\%)} & \textbf{Fill Rate (\%)} & \textbf{Waktu (ms)} & \textbf{Item Termuat} \\")
    lines.append(r"        \midrule")

    s3_hc_results = [
        a for a in aggregated
        if a.scenario_id == "S3" and a.container_key == "40ft_high_cube"
    ]
    variant_labels = {
        "baseline": "Baseline (Bigger First + Stability)",
        "no_stability": "Tanpa Pemeriksaan Stabilitas",
        "smaller_first": "Smaller First + Stability",
    }

    for a in s3_hc_results:
        label = variant_labels.get(a.variant, a.variant)
        lines.append(
            f"        {label} & "
            f"{a.avg_volume_util:.2f} & "
            f"{a.avg_fill_rate:.2f} & "
            f"{a.avg_time:.1f} & "
            f"{a.avg_fitted:.0f} \\\\"
        )

    lines.append(r"        \bottomrule")
    lines.append(r"    \end{tabular}")
    lines.append(r"\end{table}")
    lines.append("")

    # -------------------------------------------------------------------------
    # Table 3: Scalability Analysis (Baseline, 40ft HC)
    # -------------------------------------------------------------------------
    lines.append("% =============================================================================")
    lines.append("% TABLE 3: Scalability Analysis (Baseline, 40ft High Cube)")
    lines.append("% =============================================================================")
    lines.append(r"\begin{table}[H]")
    lines.append(r"    \centering")
    lines.append(r"    \caption{Analisis Skalabilitas Waktu Komputasi}")
    lines.append(r"    \label{tab:scalability}")
    lines.append(r"    \small")
    lines.append(r"    \begin{tabular}{lrrrr}")
    lines.append(r"        \toprule")
    lines.append(r"        \textbf{Skenario} & \textbf{Jumlah Item} & \textbf{Min (ms)} & \textbf{Rata-rata (ms)} & \textbf{Max (ms)} \\")
    lines.append(r"        \midrule")

    for a in baseline_hc:
        lines.append(
            f"        {a.scenario_id} & {a.total_items} & "
            f"{a.min_time:.0f} & {a.avg_time:.1f} & {a.max_time:.0f} \\\\"
        )

    lines.append(r"        \bottomrule")
    lines.append(r"    \end{tabular}")
    lines.append(r"\end{table}")
    lines.append("")

    # -------------------------------------------------------------------------
    # Table 4: Multi-Container Comparison (Baseline, all containers)
    # -------------------------------------------------------------------------
    container_order = ["20ft_standard", "40ft_standard", "40ft_high_cube"]
    container_labels = {
        "20ft_standard": "20ft Standard",
        "40ft_standard": "40ft Standard",
        "40ft_high_cube": "40ft High Cube",
    }

    lines.append("% =============================================================================")
    lines.append("% TABLE 4: Multi-Container Comparison (Baseline variant, all containers)")
    lines.append("% =============================================================================")
    lines.append(r"\begin{table}[H]")
    lines.append(r"    \centering")
    lines.append(r"    \caption{Perbandingan Performa Algoritma pada Berbagai Tipe Kontainer}")
    lines.append(r"    \label{tab:multicontainer}")
    lines.append(r"    \small")
    lines.append(r"    \begin{tabular}{llrrrr}")
    lines.append(r"        \toprule")
    lines.append(r"        \textbf{Skenario} & \textbf{Kontainer} & \textbf{Vol. Util. (\%)} & \textbf{Fill Rate (\%)} & \textbf{Item Termuat} & \textbf{Waktu (ms)} \\")
    lines.append(r"        \midrule")

    scenario_ids = [a.scenario_id for a in baseline_hc]  # use HC list to preserve scenario order
    for i, sid in enumerate(scenario_ids):
        if i > 0:
            lines.append(r"        \midrule")
        for ck in container_order:
            match = next(
                (a for a in baseline_all_containers if a.scenario_id == sid and a.container_key == ck),
                None,
            )
            if match:
                clabel = container_labels.get(ck, ck)
                lines.append(
                    f"        {sid if container_order.index(ck) == 0 else ''} & {clabel} & "
                    f"{match.avg_volume_util:.2f} & "
                    f"{match.avg_fill_rate:.2f} & "
                    f"{match.avg_fitted:.0f} & "
                    f"{match.avg_time:.1f} \\\\"
                )

    lines.append(r"        \bottomrule")
    lines.append(r"    \end{tabular}")
    lines.append(r"\end{table}")

    with open(output_path, "w") as f:
        f.write("\n".join(lines))

    print(f"  Saved LaTeX tables to {output_path}")


# =============================================================================
# MAIN BENCHMARK RUNNER
# =============================================================================

def run_benchmark(
    scenarios: list[Scenario],
    iterations: int = 10,
    variants: list[str] | None = None,
    containers: list[str] | None = None,
    output_dir: Path | None = None,
) -> tuple[list[SingleRunResult], list[AggregatedResult]]:
    """
    Run the complete benchmark suite.

    Args:
        scenarios: List of scenarios to test
        iterations: Number of iterations per scenario/container/variant
        variants: List of variant names to test (default: all)
        containers: List of container keys to test (default: all)
        output_dir: Directory for output files

    Returns:
        Tuple of (all_runs, aggregated_results)
    """
    if variants is None:
        variants = list(VARIANTS.keys())

    if containers is None:
        containers = ALL_CONTAINERS

    if output_dir is None:
        output_dir = SCRIPT_DIR / "results"

    output_dir.mkdir(parents=True, exist_ok=True)

    all_runs: list[SingleRunResult] = []
    aggregated: list[AggregatedResult] = []

    total_tests = len(containers) * len(scenarios) * len(variants) * iterations
    current_test = 0

    print("=" * 70)
    print("3D BIN PACKING BENCHMARK")
    print("=" * 70)
    print(f"Containers: {containers}")
    print(f"Scenarios:  {[s.id for s in scenarios]}")
    print(f"Variants:   {variants}")
    print(f"Iterations per combination: {iterations}")
    print(f"Total test runs: {total_tests}")
    print("=" * 70)

    for container_key in containers:
        container_name = CONTAINERS[container_key]["name"]
        print(f"\n{'=' * 70}")
        print(f"CONTAINER: {container_name}")
        print(f"{'=' * 70}")

        for scenario in scenarios:
            # Override the container on a shallow copy so original is unchanged
            import copy
            scenario_copy = copy.copy(scenario)
            scenario_copy.container_key = container_key

            print(f"\n  [{scenario_copy.id}] {scenario_copy.name} ({scenario_copy.total_items} items)")

            for variant_name in variants:
                variant_config = VARIANTS[variant_name]
                variant_runs: list[SingleRunResult] = []

                print(f"    Variant: {variant_name}", end=" ")
                sys.stdout.flush()

                for i in range(iterations):
                    current_test += 1

                    result = run_single_benchmark(
                        scenario_copy,
                        iteration=i + 1,
                        variant=variant_name,
                        **variant_config,
                    )

                    variant_runs.append(result)
                    all_runs.append(result)

                    # Progress indicator
                    print(".", end="")
                    sys.stdout.flush()

                # Aggregate variant results
                agg = aggregate_results(variant_runs)
                aggregated.append(agg)

                print(f" Vol: {agg.avg_volume_util:.1f}%, Fill: {agg.avg_fill_rate:.1f}%, Time: {agg.avg_time:.0f}ms")

    print("\n" + "=" * 70)
    print("SAVING RESULTS")
    print("=" * 70)

    save_json_results(all_runs, aggregated, output_dir)
    save_csv_results(aggregated, output_dir)
    generate_latex_tables(aggregated, output_dir)

    return all_runs, aggregated


# =============================================================================
# CLI
# =============================================================================

def main():
    parser = argparse.ArgumentParser(
        description="Run 3D Bin Packing Benchmark",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )

    parser.add_argument(
        "-i", "--iterations",
        type=int,
        default=10,
        help="Number of iterations per scenario/container/variant (default: 10)",
    )

    parser.add_argument(
        "-s", "--scenarios",
        type=str,
        default=None,
        help="Comma-separated list of scenario IDs (default: all)",
    )

    parser.add_argument(
        "-c", "--containers",
        type=str,
        default=None,
        help=(
            "Comma-separated list of container keys to test "
            "(default: all). Available: 20ft_standard, 40ft_standard, 40ft_high_cube"
        ),
    )

    parser.add_argument(
        "-v", "--variants",
        type=str,
        default=None,
        help="Comma-separated list of variant names (default: all)",
    )

    parser.add_argument(
        "-o", "--output",
        type=str,
        default=None,
        help="Output directory (default: results/)",
    )

    parser.add_argument(
        "--no-plots",
        action="store_true",
        help="Skip plot generation",
    )

    args = parser.parse_args()

    # Parse scenarios
    if args.scenarios:
        scenario_ids = [s.strip().upper() for s in args.scenarios.split(",")]
        scenarios = [s for s in SCENARIOS if s.id in scenario_ids]
        if not scenarios:
            print(f"Error: No valid scenarios found. Available: {[s.id for s in SCENARIOS]}")
            sys.exit(1)
    else:
        scenarios = SCENARIOS

    # Parse containers
    if args.containers:
        container_keys = [c.strip() for c in args.containers.split(",")]
        invalid = [c for c in container_keys if c not in CONTAINERS]
        if invalid:
            print(f"Error: Invalid containers: {invalid}. Available: {list(CONTAINERS.keys())}")
            sys.exit(1)
        containers = container_keys
    else:
        containers = ALL_CONTAINERS

    # Parse variants
    if args.variants:
        variants = [v.strip() for v in args.variants.split(",")]
        invalid = [v for v in variants if v not in VARIANTS]
        if invalid:
            print(f"Error: Invalid variants: {invalid}. Available: {list(VARIANTS.keys())}")
            sys.exit(1)
    else:
        variants = None

    # Output directory
    output_dir = Path(args.output) if args.output else SCRIPT_DIR / "results"

    # Run benchmark
    all_runs, aggregated = run_benchmark(
        scenarios=scenarios,
        iterations=args.iterations,
        variants=variants,
        containers=containers,
        output_dir=output_dir,
    )

    # Generate plots
    if not args.no_plots:
        print("\n" + "=" * 70)
        print("GENERATING PLOTS")
        print("=" * 70)

        try:
            from plot_results import generate_all_plots
            generate_all_plots(aggregated, output_dir / "figures")
        except ImportError as e:
            print(f"  Warning: Could not import plot_results: {e}")
            print("  Skipping plot generation. Run plot_results.py separately.")

    print("\n" + "=" * 70)
    print("BENCHMARK COMPLETE")
    print("=" * 70)
    print(f"Results saved to: {output_dir}")


if __name__ == "__main__":
    main()
