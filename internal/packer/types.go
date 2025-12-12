package packer

// ContainerInput represents the dimensions and weight capacity of a container.
type ContainerInput struct {
	ID        string
	Length    float64 // mm
	Width     float64 // mm
	Height    float64 // mm
	MaxWeight float64 // kg
}

// ItemInput represents an item to be packed.
// Dimensions are in mm, Weight in kg per unit.
type ItemInput struct {
	ID            string
	Label         string
	Length        float64 // mm
	Width         float64 // mm
	Height        float64 // mm
	Weight        float64 // kg
	Quantity      int
	AllowRotation bool
	Color         string // Hex color code
	ProductSKU    string
}

// PackedItem represents a single instance of an item successfully placed in the container.
type PackedItem struct {
	ItemID        string // Original ItemInput ID
	InstanceID    string // Unique ID for this placed instance
	Label         string
	ProductSKU    string
	RotatedLength float64 // mm
	RotatedWidth  float64 // mm
	RotatedHeight float64 // mm
	Position      Position
	RotationType  int // 0-5 orientation code
}

// Position represents 3D coordinates in mm from the container origin.
type Position struct {
	X, Y, Z float64
}

// PackingResult contains the outcome of a packing operation.
type PackingResult struct {
	ContainerID          string
	PackedItems          []PackedItem
	UnfitItems           []ItemInput
	TotalPackedItems     int
	TotalVolumePackedM3  float64
	TotalWeightPackedKG  float64
	VolumeUtilisationPct float64 // 0-100
	WeightUtilisationPct float64 // 0-100
	IsFeasible           bool    // True if all requested items fit
	Algorithm            string
	DurationMs           int64
}
