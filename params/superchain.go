package params

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"

	"github.com/jeffcogswell/golembase-op-geth/common"
	"github.com/jeffcogswell/golembase-op-geth/superchain"
)

var OPStackSupport = ProtocolVersionV0{Build: [8]byte{}, Major: 9, Minor: 0, Patch: 0, PreRelease: 0}.Encode()

func init() {
	for id, ch := range superchain.Chains {
		NetworkNames[fmt.Sprintf("%d", id)] = ch.Name
	}
}

// uint64ptr is a weird helper to allow 1-line constant pointer creation.
func uint64ptr(n uint64) *uint64 {
	return &n
}

func LoadOPStackChainConfig(chConfig *superchain.ChainConfig) (*ChainConfig, error) {
	hardforks := chConfig.Hardforks
	genesisActivation := uint64(0)
	out := &ChainConfig{
		ChainID:                 new(big.Int).SetUint64(chConfig.ChainID),
		HomesteadBlock:          common.Big0,
		DAOForkBlock:            nil,
		DAOForkSupport:          false,
		EIP150Block:             common.Big0,
		EIP155Block:             common.Big0,
		EIP158Block:             common.Big0,
		ByzantiumBlock:          common.Big0,
		ConstantinopleBlock:     common.Big0,
		PetersburgBlock:         common.Big0,
		IstanbulBlock:           common.Big0,
		MuirGlacierBlock:        common.Big0,
		BerlinBlock:             common.Big0,
		LondonBlock:             common.Big0,
		ArrowGlacierBlock:       common.Big0,
		GrayGlacierBlock:        common.Big0,
		MergeNetsplitBlock:      common.Big0,
		ShanghaiTime:            hardforks.CanyonTime,  // Shanghai activates with Canyon
		CancunTime:              hardforks.EcotoneTime, // Cancun activates with Ecotone
		PragueTime:              hardforks.IsthmusTime, // Prague activates with Isthmus
		BedrockBlock:            common.Big0,
		RegolithTime:            &genesisActivation,
		CanyonTime:              hardforks.CanyonTime,
		EcotoneTime:             hardforks.EcotoneTime,
		FjordTime:               hardforks.FjordTime,
		GraniteTime:             hardforks.GraniteTime,
		HoloceneTime:            hardforks.HoloceneTime,
		IsthmusTime:             hardforks.IsthmusTime,
		JovianTime:              hardforks.JovianTime,
		TerminalTotalDifficulty: common.Big0,
		Ethash:                  nil,
		Clique:                  nil,
	}

	if chConfig.Optimism != nil {
		out.Optimism = &OptimismConfig{
			EIP1559Elasticity:  chConfig.Optimism.EIP1559Elasticity,
			EIP1559Denominator: chConfig.Optimism.EIP1559Denominator,
		}
		if chConfig.Optimism.EIP1559DenominatorCanyon != nil {
			out.Optimism.EIP1559DenominatorCanyon = uint64ptr(*chConfig.Optimism.EIP1559DenominatorCanyon)
		}
	}

	// special overrides for OP-Stack chains with pre-Regolith upgrade history
	switch chConfig.ChainID {
	case OPMainnetChainID:
		out.BerlinBlock = big.NewInt(3950000)
		out.LondonBlock = big.NewInt(105235063)
		out.ArrowGlacierBlock = big.NewInt(105235063)
		out.GrayGlacierBlock = big.NewInt(105235063)
		out.MergeNetsplitBlock = big.NewInt(105235063)
		out.BedrockBlock = big.NewInt(105235063)
	}

	return out, nil
}

// ProtocolVersion encodes the OP-Stack protocol version. See OP-Stack superchain-upgrade specification.
type ProtocolVersion [32]byte

func (p ProtocolVersion) MarshalText() ([]byte, error) {
	return common.Hash(p).MarshalText()
}

func (p *ProtocolVersion) UnmarshalText(input []byte) error {
	return (*common.Hash)(p).UnmarshalText(input)
}

func (p ProtocolVersion) Parse() (versionType uint8, build [8]byte, major, minor, patch, preRelease uint32) {
	versionType = p[0]
	if versionType != 0 {
		return
	}
	// bytes 1:8 reserved for future use
	copy(build[:], p[8:16])                        // differentiates forks and custom-builds of standard protocol
	major = binary.BigEndian.Uint32(p[16:20])      // incompatible API changes
	minor = binary.BigEndian.Uint32(p[20:24])      // identifies additional functionality in backwards compatible manner
	patch = binary.BigEndian.Uint32(p[24:28])      // identifies backward-compatible bug-fixes
	preRelease = binary.BigEndian.Uint32(p[28:32]) // identifies unstable versions that may not satisfy the above
	return
}

func (p ProtocolVersion) String() string {
	versionType, build, major, minor, patch, preRelease := p.Parse()
	if versionType != 0 {
		return "v0.0.0-unknown." + common.Hash(p).String()
	}
	ver := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	if preRelease != 0 {
		ver += fmt.Sprintf("-%d", preRelease)
	}
	if build != ([8]byte{}) {
		if humanBuildTag(build) {
			ver += fmt.Sprintf("+%s", strings.TrimRight(string(build[:]), "\x00"))
		} else {
			ver += fmt.Sprintf("+0x%x", build)
		}
	}
	return ver
}

// humanBuildTag identifies which build tag we can stringify for human-readable versions
func humanBuildTag(v [8]byte) bool {
	for i, c := range v { // following semver.org advertised regex, alphanumeric with '-' and '.', except leading '.'.
		if c == 0 { // trailing zeroed are allowed
			for _, d := range v[i:] {
				if d != 0 {
					return false
				}
			}
			return true
		}
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || (c == '.' && i > 0)) {
			return false
		}
	}
	return true
}

// ProtocolVersionComparison is used to identify how far ahead/outdated a protocol version is relative to another.
// This value is used in metrics and switch comparisons, to easily identify each type of version difference.
// Negative values mean the version is outdated.
// Positive values mean the version is up-to-date.
// Matching versions have a 0.
type ProtocolVersionComparison int

const (
	AheadMajor         ProtocolVersionComparison = 4
	OutdatedMajor      ProtocolVersionComparison = -4
	AheadMinor         ProtocolVersionComparison = 3
	OutdatedMinor      ProtocolVersionComparison = -3
	AheadPatch         ProtocolVersionComparison = 2
	OutdatedPatch      ProtocolVersionComparison = -2
	AheadPrerelease    ProtocolVersionComparison = 1
	OutdatedPrerelease ProtocolVersionComparison = -1
	Matching           ProtocolVersionComparison = 0
	DiffVersionType    ProtocolVersionComparison = 100
	DiffBuild          ProtocolVersionComparison = 101
	EmptyVersion       ProtocolVersionComparison = 102
	InvalidVersion     ProtocolVersionComparison = 103
)

func (p ProtocolVersion) Compare(other ProtocolVersion) (cmp ProtocolVersionComparison) {
	if p == (ProtocolVersion{}) || (other == (ProtocolVersion{})) {
		return EmptyVersion
	}
	aVersionType, aBuild, aMajor, aMinor, aPatch, aPreRelease := p.Parse()
	bVersionType, bBuild, bMajor, bMinor, bPatch, bPreRelease := other.Parse()
	if aVersionType != bVersionType {
		return DiffVersionType
	}
	if aBuild != bBuild {
		return DiffBuild
	}
	// max values are reserved, consider versions with these values invalid
	if aMajor == ^uint32(0) || aMinor == ^uint32(0) || aPatch == ^uint32(0) || aPreRelease == ^uint32(0) ||
		bMajor == ^uint32(0) || bMinor == ^uint32(0) || bPatch == ^uint32(0) || bPreRelease == ^uint32(0) {
		return InvalidVersion
	}
	fn := func(a, b uint32, ahead, outdated ProtocolVersionComparison) ProtocolVersionComparison {
		if a == b {
			return Matching
		}
		if a > b {
			return ahead
		}
		return outdated
	}
	if aPreRelease != 0 { // if A is a pre-release, then decrement the version before comparison
		if aPatch != 0 {
			aPatch -= 1
		} else if aMinor != 0 {
			aMinor -= 1
			aPatch = ^uint32(0) // max value
		} else if aMajor != 0 {
			aMajor -= 1
			aMinor = ^uint32(0) // max value
		}
	}
	if bPreRelease != 0 { // if B is a pre-release, then decrement the version before comparison
		if bPatch != 0 {
			bPatch -= 1
		} else if bMinor != 0 {
			bMinor -= 1
			bPatch = ^uint32(0) // max value
		} else if bMajor != 0 {
			bMajor -= 1
			bMinor = ^uint32(0) // max value
		}
	}
	if c := fn(aMajor, bMajor, AheadMajor, OutdatedMajor); c != Matching {
		return c
	}
	if c := fn(aMinor, bMinor, AheadMinor, OutdatedMinor); c != Matching {
		return c
	}
	if c := fn(aPatch, bPatch, AheadPatch, OutdatedPatch); c != Matching {
		return c
	}
	return fn(aPreRelease, bPreRelease, AheadPrerelease, OutdatedPrerelease)
}

type ProtocolVersionV0 struct {
	Build                           [8]byte
	Major, Minor, Patch, PreRelease uint32
}

func (v ProtocolVersionV0) Encode() (out ProtocolVersion) {
	copy(out[8:16], v.Build[:])
	binary.BigEndian.PutUint32(out[16:20], v.Major)
	binary.BigEndian.PutUint32(out[20:24], v.Minor)
	binary.BigEndian.PutUint32(out[24:28], v.Patch)
	binary.BigEndian.PutUint32(out[28:32], v.PreRelease)
	return
}
