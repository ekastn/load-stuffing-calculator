#!/usr/bin/env python3
"""
Plot Generation for 3D Bin Packing Benchmark Results

Generates publication-quality figures for academic paper.

Usage:
    python plot_results.py [--input results/] [--output results/figures/]
"""

from __future__ import annotations

import argparse
import json
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any

import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
import numpy as np

# Use a clean academic style
plt.style.use('seaborn-v0_8-whitegrid')

# Set global font settings for academic papers
plt.rcParams.update({
    'font.family': 'serif',
    'font.serif': ['Times New Roman', 'DejaVu Serif', 'serif'],
    'font.size': 10,
    'axes.titlesize': 11,
    'axes.labelsize': 10,
    'xtick.labelsize': 9,
    'ytick.labelsize': 9,
    'legend.fontsize': 9,
    'figure.titlesize': 12,
    'figure.dpi': 150,
    'savefig.dpi': 300,
    'savefig.bbox': 'tight',
    'savefig.pad_inches': 0.1,
})

# Color palette - professional academic colors
COLORS = {
    'primary': '#2E86AB',      # Blue
    'secondary': '#A23B72',    # Magenta
    'tertiary': '#F18F01',     # Orange
    'quaternary': '#C73E1D',   # Red
    'success': '#3A7D44',      # Green
    'neutral': '#6C757D',      # Gray
}

VARIANT_COLORS = {
    'baseline': COLORS['primary'],
    'no_stability': COLORS['tertiary'],
    'smaller_first': COLORS['secondary'],
}

VARIANT_LABELS = {
    'baseline': 'Bigger First + Stability',
    'no_stability': 'No Stability Check',
    'smaller_first': 'Smaller First',
}


@dataclass
class AggregatedResult:
    """Aggregated statistics for a scenario+variant combination."""
    scenario_id: str
    variant: str
    num_runs: int
    total_items: int
    avg_fitted: float
    std_fitted: float
    avg_volume_util: float
    std_volume_util: float
    min_volume_util: float
    max_volume_util: float
    avg_weight_util: float
    std_weight_util: float
    avg_fill_rate: float
    std_fill_rate: float
    avg_time: float
    std_time: float
    min_time: float
    max_time: float


def load_results(results_dir: Path) -> list[AggregatedResult]:
    """Load aggregated results from JSON file."""
    json_path = results_dir / "results.json"
    
    with open(json_path, "r") as f:
        data = json.load(f)
    
    results = []
    for item in data["aggregated_results"]:
        results.append(AggregatedResult(
            scenario_id=item["scenario_id"],
            variant=item["variant"],
            num_runs=item["num_runs"],
            total_items=item["total_items"],
            avg_fitted=item["avg_fitted"],
            std_fitted=item["std_fitted"],
            avg_volume_util=item["avg_volume_util"],
            std_volume_util=item["std_volume_util"],
            min_volume_util=item["min_volume_util"],
            max_volume_util=item["max_volume_util"],
            avg_weight_util=item["avg_weight_util"],
            std_weight_util=item["std_weight_util"],
            avg_fill_rate=item["avg_fill_rate"],
            std_fill_rate=item["std_fill_rate"],
            avg_time=item["avg_time_ms"],
            std_time=item["std_time_ms"],
            min_time=item["min_time_ms"],
            max_time=item["max_time_ms"],
        ))
    
    return results


def plot_utilization_comparison(
    results: list[AggregatedResult],
    output_dir: Path,
) -> None:
    """
    Generate bar chart comparing volume and weight utilization across scenarios.
    """
    # Filter baseline results
    baseline = [r for r in results if r.variant == "baseline"]
    baseline.sort(key=lambda x: x.scenario_id)
    
    scenarios = [r.scenario_id for r in baseline]
    volume_utils = [r.avg_volume_util for r in baseline]
    volume_errs = [r.std_volume_util for r in baseline]
    weight_utils = [r.avg_weight_util for r in baseline]
    weight_errs = [r.std_weight_util for r in baseline]
    
    x = np.arange(len(scenarios))
    width = 0.35
    
    fig, ax = plt.subplots(figsize=(8, 5))
    
    bars1 = ax.bar(
        x - width/2, volume_utils, width,
        yerr=volume_errs, capsize=3,
        label='Volume Utilization',
        color=COLORS['primary'],
        edgecolor='white',
        linewidth=0.5,
    )
    
    bars2 = ax.bar(
        x + width/2, weight_utils, width,
        yerr=weight_errs, capsize=3,
        label='Weight Utilization',
        color=COLORS['secondary'],
        edgecolor='white',
        linewidth=0.5,
    )
    
    ax.set_xlabel('Scenario')
    ax.set_ylabel('Utilization (%)')
    ax.set_title('Volume and Weight Utilization by Scenario')
    ax.set_xticks(x)
    ax.set_xticklabels(scenarios)
    ax.legend(loc='upper left')
    ax.set_ylim(0, 100)
    
    # Add value labels on bars
    for bar, val in zip(bars1, volume_utils):
        ax.annotate(
            f'{val:.1f}%',
            xy=(bar.get_x() + bar.get_width()/2, bar.get_height()),
            xytext=(0, 3),
            textcoords='offset points',
            ha='center', va='bottom',
            fontsize=8,
        )
    
    for bar, val in zip(bars2, weight_utils):
        ax.annotate(
            f'{val:.1f}%',
            xy=(bar.get_x() + bar.get_width()/2, bar.get_height()),
            xytext=(0, 3),
            textcoords='offset points',
            ha='center', va='bottom',
            fontsize=8,
        )
    
    plt.tight_layout()
    
    # Save in multiple formats
    output_dir.mkdir(parents=True, exist_ok=True)
    fig.savefig(output_dir / "utilization_comparison.pdf")
    fig.savefig(output_dir / "utilization_comparison.png")
    plt.close(fig)
    
    print(f"  Generated: utilization_comparison.pdf")


def plot_computation_time(
    results: list[AggregatedResult],
    output_dir: Path,
) -> None:
    """
    Generate line chart showing computation time scalability.
    """
    # Filter baseline results
    baseline = [r for r in results if r.variant == "baseline"]
    baseline.sort(key=lambda x: x.total_items)
    
    items = [r.total_items for r in baseline]
    times = [r.avg_time for r in baseline]
    time_errs = [r.std_time for r in baseline]
    
    fig, ax = plt.subplots(figsize=(8, 5))
    
    # Plot with error band
    ax.errorbar(
        items, times,
        yerr=time_errs,
        marker='o',
        markersize=8,
        linewidth=2,
        capsize=4,
        capthick=1.5,
        color=COLORS['primary'],
        label='Computation Time',
    )
    
    # Fill error band
    times_arr = np.array(times)
    errs_arr = np.array(time_errs)
    ax.fill_between(
        items,
        times_arr - errs_arr,
        times_arr + errs_arr,
        alpha=0.2,
        color=COLORS['primary'],
    )
    
    # Add trend line (polynomial fit)
    z = np.polyfit(items, times, 2)
    p = np.poly1d(z)
    x_smooth = np.linspace(min(items), max(items), 100)
    ax.plot(
        x_smooth, p(x_smooth),
        '--',
        color=COLORS['neutral'],
        alpha=0.7,
        label='Quadratic Trend',
    )
    
    ax.set_xlabel('Number of Items')
    ax.set_ylabel('Computation Time (ms)')
    ax.set_title('Algorithm Scalability: Computation Time vs Item Count')
    ax.legend(loc='upper left')
    ax.set_xlim(0, max(items) * 1.1)
    ax.set_ylim(0, max(times) * 1.3)
    
    # Add data labels
    for x, y in zip(items, times):
        ax.annotate(
            f'{y:.0f}ms',
            xy=(x, y),
            xytext=(5, 10),
            textcoords='offset points',
            fontsize=8,
        )
    
    plt.tight_layout()
    
    output_dir.mkdir(parents=True, exist_ok=True)
    fig.savefig(output_dir / "computation_time.pdf")
    fig.savefig(output_dir / "computation_time.png")
    plt.close(fig)
    
    print(f"  Generated: computation_time.pdf")


def plot_fill_rate(
    results: list[AggregatedResult],
    output_dir: Path,
) -> None:
    """
    Generate stacked bar chart showing fitted vs unfitted items.
    """
    # Filter baseline results
    baseline = [r for r in results if r.variant == "baseline"]
    baseline.sort(key=lambda x: x.scenario_id)
    
    scenarios = [r.scenario_id for r in baseline]
    fitted = [r.avg_fitted for r in baseline]
    unfitted = [r.total_items - r.avg_fitted for r in baseline]
    fill_rates = [r.avg_fill_rate for r in baseline]
    
    x = np.arange(len(scenarios))
    width = 0.6
    
    fig, ax = plt.subplots(figsize=(8, 5))
    
    bars1 = ax.bar(
        x, fitted, width,
        label='Fitted Items',
        color=COLORS['success'],
        edgecolor='white',
        linewidth=0.5,
    )
    
    bars2 = ax.bar(
        x, unfitted, width,
        bottom=fitted,
        label='Unfitted Items',
        color=COLORS['quaternary'],
        edgecolor='white',
        linewidth=0.5,
    )
    
    ax.set_xlabel('Scenario')
    ax.set_ylabel('Number of Items')
    ax.set_title('Item Placement Success Rate by Scenario')
    ax.set_xticks(x)
    ax.set_xticklabels(scenarios)
    ax.legend(loc='upper left')
    
    # Add fill rate labels
    for i, (bar, rate) in enumerate(zip(bars1, fill_rates)):
        total = fitted[i] + unfitted[i]
        ax.annotate(
            f'{rate:.1f}%',
            xy=(bar.get_x() + bar.get_width()/2, total),
            xytext=(0, 5),
            textcoords='offset points',
            ha='center', va='bottom',
            fontsize=9,
            fontweight='bold',
        )
    
    plt.tight_layout()
    
    output_dir.mkdir(parents=True, exist_ok=True)
    fig.savefig(output_dir / "fill_rate.pdf")
    fig.savefig(output_dir / "fill_rate.png")
    plt.close(fig)
    
    print(f"  Generated: fill_rate.pdf")


def plot_variant_comparison(
    results: list[AggregatedResult],
    output_dir: Path,
) -> None:
    """
    Generate grouped bar chart comparing algorithm variants on S4.
    """
    # Filter S4 results (or largest scenario)
    s4_results = [r for r in results if r.scenario_id == "S4"]
    
    if not s4_results:
        # Use largest available scenario
        all_scenarios = set(r.scenario_id for r in results)
        target = sorted(all_scenarios)[-1]
        s4_results = [r for r in results if r.scenario_id == target]
    
    variants = [r.variant for r in s4_results]
    volume_utils = [r.avg_volume_util for r in s4_results]
    fill_rates = [r.avg_fill_rate for r in s4_results]
    times = [r.avg_time for r in s4_results]
    
    x = np.arange(len(variants))
    width = 0.25
    
    fig, ax1 = plt.subplots(figsize=(10, 5))
    
    # Volume utilization bars
    bars1 = ax1.bar(
        x - width, volume_utils, width,
        label='Volume Utilization (%)',
        color=COLORS['primary'],
        edgecolor='white',
    )
    
    # Fill rate bars
    bars2 = ax1.bar(
        x, fill_rates, width,
        label='Fill Rate (%)',
        color=COLORS['success'],
        edgecolor='white',
    )
    
    ax1.set_xlabel('Algorithm Variant')
    ax1.set_ylabel('Percentage (%)')
    ax1.set_ylim(0, 110)
    
    # Secondary axis for time
    ax2 = ax1.twinx()
    bars3 = ax2.bar(
        x + width, times, width,
        label='Time (ms)',
        color=COLORS['tertiary'],
        edgecolor='white',
    )
    ax2.set_ylabel('Computation Time (ms)')
    
    # X-axis labels
    variant_labels = [VARIANT_LABELS.get(v, v) for v in variants]
    ax1.set_xticks(x)
    ax1.set_xticklabels(variant_labels, rotation=15, ha='right')
    
    ax1.set_title('Algorithm Variant Comparison (Scenario S4: 100 items)')
    
    # Combined legend
    lines1, labels1 = ax1.get_legend_handles_labels()
    lines2, labels2 = ax2.get_legend_handles_labels()
    ax1.legend(lines1 + lines2, labels1 + labels2, loc='upper right')
    
    plt.tight_layout()
    
    output_dir.mkdir(parents=True, exist_ok=True)
    fig.savefig(output_dir / "variant_comparison.pdf")
    fig.savefig(output_dir / "variant_comparison.png")
    plt.close(fig)
    
    print(f"  Generated: variant_comparison.pdf")


def plot_detailed_metrics(
    results: list[AggregatedResult],
    output_dir: Path,
) -> None:
    """
    Generate a 2x2 subplot with all key metrics.
    """
    # Filter baseline results
    baseline = [r for r in results if r.variant == "baseline"]
    baseline.sort(key=lambda x: x.scenario_id)
    
    scenarios = [r.scenario_id for r in baseline]
    items = [r.total_items for r in baseline]
    volume_utils = [r.avg_volume_util for r in baseline]
    weight_utils = [r.avg_weight_util for r in baseline]
    fill_rates = [r.avg_fill_rate for r in baseline]
    times = [r.avg_time for r in baseline]
    time_errs = [r.std_time for r in baseline]
    
    fig, axes = plt.subplots(2, 2, figsize=(12, 10))
    
    # 1. Volume Utilization
    ax1 = axes[0, 0]
    bars = ax1.bar(scenarios, volume_utils, color=COLORS['primary'], edgecolor='white')
    ax1.set_xlabel('Scenario')
    ax1.set_ylabel('Volume Utilization (%)')
    ax1.set_title('(a) Volume Utilization by Scenario')
    ax1.set_ylim(0, 100)
    for bar, val in zip(bars, volume_utils):
        ax1.annotate(f'{val:.1f}%', xy=(bar.get_x() + bar.get_width()/2, val),
                     xytext=(0, 3), textcoords='offset points', ha='center', fontsize=8)
    
    # 2. Fill Rate
    ax2 = axes[0, 1]
    bars = ax2.bar(scenarios, fill_rates, color=COLORS['success'], edgecolor='white')
    ax2.set_xlabel('Scenario')
    ax2.set_ylabel('Fill Rate (%)')
    ax2.set_title('(b) Item Placement Success Rate')
    ax2.set_ylim(0, 110)
    for bar, val in zip(bars, fill_rates):
        ax2.annotate(f'{val:.1f}%', xy=(bar.get_x() + bar.get_width()/2, val),
                     xytext=(0, 3), textcoords='offset points', ha='center', fontsize=8)
    
    # 3. Computation Time (line)
    ax3 = axes[1, 0]
    ax3.errorbar(items, times, yerr=time_errs, marker='o', markersize=8,
                 linewidth=2, capsize=4, color=COLORS['primary'])
    ax3.fill_between(items, np.array(times) - np.array(time_errs),
                     np.array(times) + np.array(time_errs), alpha=0.2, color=COLORS['primary'])
    ax3.set_xlabel('Number of Items')
    ax3.set_ylabel('Computation Time (ms)')
    ax3.set_title('(c) Computation Time Scalability')
    for x, y in zip(items, times):
        ax3.annotate(f'{y:.0f}ms', xy=(x, y), xytext=(5, 10),
                     textcoords='offset points', fontsize=8)
    
    # 4. Weight vs Volume Utilization (scatter)
    ax4 = axes[1, 1]
    scatter = ax4.scatter(weight_utils, volume_utils, c=items, cmap='viridis',
                          s=100, edgecolor='white', linewidth=1)
    ax4.set_xlabel('Weight Utilization (%)')
    ax4.set_ylabel('Volume Utilization (%)')
    ax4.set_title('(d) Weight vs Volume Utilization')
    ax4.set_xlim(0, max(weight_utils) * 1.2)
    ax4.set_ylim(0, 100)
    cbar = plt.colorbar(scatter, ax=ax4)
    cbar.set_label('Number of Items')
    
    # Add scenario labels to scatter
    for i, (w, v, s) in enumerate(zip(weight_utils, volume_utils, scenarios)):
        ax4.annotate(s, xy=(w, v), xytext=(5, 5), textcoords='offset points', fontsize=9)
    
    plt.tight_layout()
    
    output_dir.mkdir(parents=True, exist_ok=True)
    fig.savefig(output_dir / "detailed_metrics.pdf")
    fig.savefig(output_dir / "detailed_metrics.png")
    plt.close(fig)
    
    print(f"  Generated: detailed_metrics.pdf")


def generate_all_plots(
    results: list[AggregatedResult],
    output_dir: Path,
) -> None:
    """Generate all plots for the paper."""
    
    output_dir.mkdir(parents=True, exist_ok=True)
    
    print("Generating plots...")
    
    plot_utilization_comparison(results, output_dir)
    plot_computation_time(results, output_dir)
    plot_fill_rate(results, output_dir)
    plot_variant_comparison(results, output_dir)
    plot_detailed_metrics(results, output_dir)
    
    print(f"\nAll plots saved to: {output_dir}")


def main():
    parser = argparse.ArgumentParser(
        description="Generate plots from benchmark results",
    )
    
    parser.add_argument(
        "-i", "--input",
        type=str,
        default="results",
        help="Input directory containing results.json (default: results/)",
    )
    
    parser.add_argument(
        "-o", "--output",
        type=str,
        default=None,
        help="Output directory for figures (default: <input>/figures/)",
    )
    
    args = parser.parse_args()
    
    script_dir = Path(__file__).parent.resolve()
    input_dir = Path(args.input)
    if not input_dir.is_absolute():
        input_dir = script_dir / input_dir
    
    if args.output:
        output_dir = Path(args.output)
        if not output_dir.is_absolute():
            output_dir = script_dir / output_dir
    else:
        output_dir = input_dir / "figures"
    
    if not (input_dir / "results.json").exists():
        print(f"Error: results.json not found in {input_dir}")
        print("Run run_benchmark.py first to generate results.")
        sys.exit(1)
    
    results = load_results(input_dir)
    generate_all_plots(results, output_dir)


if __name__ == "__main__":
    main()
