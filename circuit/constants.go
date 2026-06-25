package circuit

import "math/big"

var (
	//	is poseidon hash(empty account info)
	//
	// EmptyAccountLeafNodeHash, _ = new(big.Int).SetString("2853aadbbd06deb5d6a1389d23a89ff47658aaf4a45b287ecfd62df192bd91a4", 16)
	// EmptyAccountLeafNodeHash, _ = new(big.Int).SetString("251a46d446ae444bdce9ff893361bba7cf80a5d43a320449ca7a9de0f0e48e74", 16) // 618
	EmptyAccountLeafNodeHash, _ = new(big.Int).SetString("17109aabd55f1571098efd581899989c4340e103843a11eb21a356400c3fcf62", 16) // 64
// EmptyAccountLeafNodeHash, _ = new(big.Int).SetString("26507dd7020c1919ce0048c710ea5087be0a75359010315c17d58d5370882545", 16) //621
)
