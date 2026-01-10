"""
Test Scenario Definitions for 3D Bin Packing Benchmark

This module defines realistic test scenarios for evaluating the packing algorithm
performance across different complexity levels.

Scenarios are based on the methodology defined in the paper:
- S1: 10 items, Homogeneous (single product type)
- S2: 25 items, Light Heterogeneity (2 product types)
- S3: 50 items, Medium Heterogeneity (3 product types)
- S4: 100 items, High Heterogeneity (4 product types)
- S5: 200 items, Stress Test (4 product types)
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import TypedDict


class ProductSpec(TypedDict):
    """Product specification in millimeters and kilograms."""
    length: float  # mm
    width: float   # mm
    height: float  # mm
    weight: float  # kg


class ContainerSpec(TypedDict):
    """Container specification."""
    name: str
    length: float  # mm
    width: float   # mm
    height: float  # mm
    max_weight: float  # kg


# =============================================================================
# CONTAINER DEFINITIONS
# Standard shipping container dimensions (internal dimensions)
# =============================================================================

CONTAINERS: dict[str, ContainerSpec] = {
    "20ft_standard": {
        "name": "20ft Standard",
        "length": 5898,
        "width": 2352,
        "height": 2393,
        "max_weight": 21770,
    },
    "40ft_standard": {
        "name": "40ft Standard",
        "length": 12032,
        "width": 2352,
        "height": 2393,
        "max_weight": 26680,
    },
    "40ft_high_cube": {
        "name": "40ft High Cube",
        "length": 12032,
        "width": 2352,
        "height": 2698,
        "max_weight": 26460,
    },
}

# Default container for benchmarks
DEFAULT_CONTAINER = CONTAINERS["40ft_high_cube"]


# =============================================================================
# PRODUCT DEFINITIONS
# Realistic logistics product dimensions
# =============================================================================

PRODUCTS: dict[str, ProductSpec] = {
    "euro_pallet": {
        "length": 1200,
        "width": 800,
        "height": 144,
        "weight": 25.0,
    },
    "large_crate": {
        "length": 1000,
        "width": 600,
        "height": 500,
        "weight": 45.0,
    },
    "medium_box": {
        "length": 600,
        "width": 400,
        "height": 400,
        "weight": 15.0,
    },
    "small_box": {
        "length": 400,
        "width": 300,
        "height": 200,
        "weight": 5.0,
    },
}


# =============================================================================
# SCENARIO DEFINITIONS
# =============================================================================

@dataclass
class ScenarioItem:
    """Item definition within a scenario."""
    product_key: str
    quantity: int
    
    @property
    def spec(self) -> ProductSpec:
        return PRODUCTS[self.product_key]


@dataclass
class Scenario:
    """Test scenario definition."""
    id: str
    name: str
    description: str
    heterogeneity: str
    items: list[ScenarioItem]
    container_key: str = "40ft_high_cube"
    
    @property
    def container(self) -> ContainerSpec:
        return CONTAINERS[self.container_key]
    
    @property
    def total_items(self) -> int:
        return sum(item.quantity for item in self.items)
    
    @property
    def num_product_types(self) -> int:
        return len(self.items)
    
    def get_total_volume_mm3(self) -> float:
        """Calculate total volume of all items in mm³."""
        total = 0.0
        for item in self.items:
            spec = item.spec
            vol = spec["length"] * spec["width"] * spec["height"]
            total += vol * item.quantity
        return total
    
    def get_total_weight_kg(self) -> float:
        """Calculate total weight of all items in kg."""
        total = 0.0
        for item in self.items:
            total += item.spec["weight"] * item.quantity
        return total
    
    def get_container_volume_mm3(self) -> float:
        """Calculate container volume in mm³."""
        c = self.container
        return c["length"] * c["width"] * c["height"]


# =============================================================================
# PREDEFINED SCENARIOS
# =============================================================================

SCENARIOS: list[Scenario] = [
    Scenario(
        id="S1",
        name="Homogeneous Small",
        description="50 items, single product type for functional validation",
        heterogeneity="Homogen",
        items=[
            ScenarioItem("medium_box", 50),
        ],
    ),
    Scenario(
        id="S2",
        name="Light Heterogeneity",
        description="100 items, 2 product types for standard operational testing",
        heterogeneity="Heterogen Ringan",
        items=[
            ScenarioItem("medium_box", 60),
            ScenarioItem("small_box", 40),
        ],
    ),
    Scenario(
        id="S3",
        name="Medium Heterogeneity",
        description="150 items, 3 product types for medium load testing",
        heterogeneity="Heterogen Sedang",
        items=[
            ScenarioItem("large_crate", 40),
            ScenarioItem("medium_box", 60),
            ScenarioItem("small_box", 50),
        ],
    ),
    Scenario(
        id="S4",
        name="High Heterogeneity",
        description="200 items, 4 product types for high complexity testing",
        heterogeneity="Sangat Heterogen",
        items=[
            ScenarioItem("euro_pallet", 30),
            ScenarioItem("large_crate", 50),
            ScenarioItem("medium_box", 70),
            ScenarioItem("small_box", 50),
        ],
    ),
    Scenario(
        id="S5",
        name="Stress Test",
        description="300 items, 4 product types for stress testing",
        heterogeneity="Sangat Heterogen",
        items=[
            ScenarioItem("euro_pallet", 50),
            ScenarioItem("large_crate", 80),
            ScenarioItem("medium_box", 100),
            ScenarioItem("small_box", 70),
        ],
    ),
]


def get_scenario_by_id(scenario_id: str) -> Scenario | None:
    """Get a scenario by its ID."""
    for s in SCENARIOS:
        if s.id == scenario_id:
            return s
    return None


def build_pack_request(
    scenario: Scenario,
    *,
    bigger_first: bool = True,
    check_stable: bool = True,
    support_surface_ratio: float = 0.75,
) -> dict:
    """
    Build a packing request payload from a scenario definition.
    
    Args:
        scenario: The test scenario
        bigger_first: Whether to pack bigger items first
        check_stable: Whether to check stability constraints
        support_surface_ratio: Required support surface ratio (0-1)
    
    Returns:
        Dictionary ready to be passed to pack_request()
    """
    container = scenario.container
    
    items = []
    for idx, scenario_item in enumerate(scenario.items):
        spec = scenario_item.spec
        items.append({
            "item_id": f"{scenario_item.product_key}_{idx}",
            "label": scenario_item.product_key.replace("_", " ").title(),
            "length": spec["length"],
            "width": spec["width"],
            "height": spec["height"],
            "weight": spec["weight"],
            "quantity": scenario_item.quantity,
        })
    
    return {
        "units": "mm",
        "container": {
            "length": container["length"],
            "width": container["width"],
            "height": container["height"],
            "max_weight": container["max_weight"],
        },
        "items": items,
        "options": {
            "bigger_first": bigger_first,
            "fix_point": True,  # Always enable gravity
            "check_stable": check_stable,
            "support_surface_ratio": support_surface_ratio,
            "put_type": 1,
        },
    }


if __name__ == "__main__":
    # Print scenario summary
    print("=" * 70)
    print("BENCHMARK SCENARIOS SUMMARY")
    print("=" * 70)
    
    for s in SCENARIOS:
        print(f"\n{s.id}: {s.name}")
        print(f"  Description: {s.description}")
        print(f"  Heterogeneity: {s.heterogeneity}")
        print(f"  Total Items: {s.total_items}")
        print(f"  Product Types: {s.num_product_types}")
        print(f"  Total Volume: {s.get_total_volume_mm3() / 1e9:.2f} m³")
        print(f"  Total Weight: {s.get_total_weight_kg():.1f} kg")
        print(f"  Container: {s.container['name']}")
        print(f"  Container Volume: {s.get_container_volume_mm3() / 1e9:.2f} m³")
