package config

import (
	"math/big"
	"websea-zkmerkle-proof/utils"
)

type Config struct {
	MysqlDataSource string
	UserDataFile    string
	DbSuffix        string
	TreeDB          struct {
		Driver string
		Option struct {
			Addr string
		}
	}
	Redis struct {
		Host     string
		Type     string
		Password string
	}
	ZkKeyName string
}

type CexConfig struct {
	ProofCsv                  string
	ZkKeyVKDirectoryAndPrefix string
	CexAssetsInfo             []utils.CexAssetInfo
}

type UserConfig struct {
	Arrangement          uint32
	UniqueIdentification string
	TotalAssetEquity     *big.Int
	TotalAssetDebt       *big.Int
	AssetDetails         []utils.AccountAsset
	TreeRootHash         string
	MerkleProofEncode    []string
}
