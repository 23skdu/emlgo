package gpu

// Device represents a GPU device with its properties.
type Device struct {
	ID                 int
	Name               string
	ComputeMajor       int
	ComputeMinor       int
	MemoryBytes        int64
	MaxThreadsPerBlock int
	WarpSize           int
	ClockRateKHz       int
}

// Stream represents an asynchronous CUDA stream for overlapping
// kernel execution with data transfers.
type Stream struct {
	handle uintptr
}

// BatchConfig configures a GPU batch operation.
type BatchConfig struct {
	// BlockSize is threads per block (0 = default 256).
	BlockSize int
	// UsePinned enables zero-copy pinned memory for host<->device transfers.
	UsePinned bool
	// Stream enables async execution on a non-default stream.
	Stream *Stream
}

// DefaultBlockSize is the default threads per block.
const DefaultBlockSize = 256

// GridSize computes the grid size (number of blocks) for a given
// number of elements and block size.
func GridSize(n, blockSize int) int {
	if blockSize <= 0 {
		blockSize = DefaultBlockSize
	}
	return (n + blockSize - 1) / blockSize
}

// LaunchConfig holds pre-computed kernel launch parameters.
type LaunchConfig struct {
	GridDimX  int
	GridDimY  int
	GridDimZ  int
	BlockDimX int
	BlockDimY int
	BlockDimZ int
}

// DefaultLaunchConfig returns a 1D launch config for n elements.
func DefaultLaunchConfig(n, blockSize int) LaunchConfig {
	if blockSize <= 0 || blockSize > 1024 {
		blockSize = DefaultBlockSize
	}
	return LaunchConfig{
		GridDimX:  GridSize(n, blockSize),
		BlockDimX: blockSize,
		BlockDimY: 1,
		BlockDimZ: 1,
		GridDimY:  1,
		GridDimZ:  1,
	}
}

// MatMulLaunchConfig returns a 2D launch config for MxN matrix multiply.
func MatMulLaunchConfig(M, N, tileSize int) LaunchConfig {
	if tileSize <= 0 {
		tileSize = 16
	}
	return LaunchConfig{
		GridDimX:  (N + tileSize - 1) / tileSize,
		GridDimY:  (M + tileSize - 1) / tileSize,
		GridDimZ:  1,
		BlockDimX: tileSize,
		BlockDimY: tileSize,
		BlockDimZ: 1,
	}
}

// Error codes returned by the C API.
const (
	ErrSuccess       = 0
	ErrNoDevice      = 1
	ErrNotInit       = 2
	ErrAlloc         = 3
	ErrCopy          = 4
	ErrLaunch        = 5
	ErrInvalidConfig = 6
)
