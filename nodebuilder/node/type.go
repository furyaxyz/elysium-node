package node

// Type defines the Node type (e.g. `light`, `bridge`) for identity purposes.
// The zero value for Type is invalid.
type Type uint8

// StorePath is an alias used in order to pass the base path of the node store to nodebuilder
// modules.
type StorePath string

const (
	// Bridge is a Elysium Node that bridges the Elysium consensus network and data availability
	// network. It maintains a trusted channel/connection to a Elysium Core node via the core.Client
	// API.
	Bridge Type = iota + 1
	// Light is a stripped-down Elysium Node which aims to be lightweight while preserving the highest
	// possible security guarantees.
	Light
	// Full is a Elysium Node that stores blocks in their entirety.
	Full
)

// String converts Type to its string representation.
func (t Type) String() string {
	if !t.IsValid() {
		return "unknown"
	}
	return typeToString[t]
}

// IsValid reports whether the Type is valid.
func (t Type) IsValid() bool {
	_, ok := typeToString[t]
	return ok
}

// ParseType converts string in a type if possible.
func ParseType(str string) Type {
	tp, ok := stringToType[str]
	if !ok {
		return 0
	}

	return tp
}

// typeToString keeps string representations of all valid Types.
var typeToString = map[Type]string{
	Bridge: "Bridge",
	Light:  "Light",
	Full:   "Full",
}

// typeToString maps strings representations of all valid Types.
var stringToType = map[string]Type{
	"Bridge": Bridge,
	"Light":  Light,
	"Full":   Full,
}
