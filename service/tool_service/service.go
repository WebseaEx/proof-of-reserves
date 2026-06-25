package tool_service

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"websea-zkmerkle-proof/config"
	"websea-zkmerkle-proof/global"
	prover_server "websea-zkmerkle-proof/service/prover_service"
	"websea-zkmerkle-proof/service/witness_service"
	"websea-zkmerkle-proof/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CleanKvrocks() {
	global.Cfg = &config.Config{}
	jsonFile, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("load config err : %s", err.Error()))
	}
	err = json.Unmarshal(jsonFile, global.Cfg)
	if err != nil {
		panic(err.Error())
	}

	dbtoolConfig := global.Cfg
	client := redis.NewClient(&redis.Options{
		Addr:            dbtoolConfig.TreeDB.Option.Addr,
		PoolSize:        500,
		MaxRetries:      5,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     10 * time.Second,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		PoolTimeout:     15 * time.Second,
		IdleTimeout:     5 * time.Minute,
	})
	client.FlushAll(context.Background())
	fmt.Println("kvrocks data drop successfully")
}

func CheckProverStatus() {
	global.Cfg = &config.Config{}
	jsonFile, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("load config err : %s", err.Error()))
	}
	err = json.Unmarshal(jsonFile, global.Cfg)
	if err != nil {
		panic(err.Error())
	}

	dbtoolConfig := global.Cfg
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             60 * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Silent,    // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,            // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dbtoolConfig.MysqlDataSource), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err.Error())
	}
	witnessModel := witness_service.NewWitnessModel(db, dbtoolConfig.DbSuffix)
	proofModel := prover_server.NewProofModel(db, dbtoolConfig.DbSuffix)

	witnessCounts, err := witnessModel.GetRowCounts()
	if err != nil {
		panic(err.Error())
	}
	proofCounts, err := proofModel.GetRowCounts()
	fmt.Printf("Total witness item %d, Published item %d, Pending item %d, Finished item %d\n", witnessCounts[0], witnessCounts[1], witnessCounts[2], witnessCounts[3])
	fmt.Println(witnessCounts[0] - proofCounts)
}

func QueryCexAssets() {
	global.Cfg = &config.Config{}
	jsonFile, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("load config err : %s", err.Error()))
	}
	err = json.Unmarshal(jsonFile, global.Cfg)
	if err != nil {
		panic(err.Error())
	}

	dbtoolConfig := global.Cfg
	db, err := gorm.Open(mysql.Open(dbtoolConfig.MysqlDataSource))
	if err != nil {
		panic(err.Error())
	}
	witnessModel := witness_service.NewWitnessModel(db, dbtoolConfig.DbSuffix)
	latestWitness, err := witnessModel.GetLatestBatchWitness()
	if err != nil {
		panic(err.Error())
	}
	witness := utils.DecodeBatchWitness(latestWitness.WitnessData)
	if witness == nil {
		panic("decode invalid witness data")
	}
	cexAssetsInfo := utils.RecoverAfterCexAssets(witness)
	var newAssetsInfo []utils.CexAssetInfo
	for i := 0; i < len(cexAssetsInfo); i++ {
		if cexAssetsInfo[i].BasePrice != 0 {
			newAssetsInfo = append(newAssetsInfo, cexAssetsInfo[i])
		}
	}
	cexAssetsInfoBytes, _ := json.Marshal(newAssetsInfo)
	fmt.Println(string(cexAssetsInfoBytes))
}

func CheckUserFiles() {
	global.Cfg = &config.Config{}
	jsonFile, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("load config err : %s", err.Error()))
	}
	err = json.Unmarshal(jsonFile, global.Cfg)
	if err != nil {
		panic(err.Error())
	}

	_, _, err = utils.ReadUserAssetsV1(global.Cfg.UserDataFile)
	if err != nil {
		panic(err.Error())
	}
}

func CompareAccountID(filePath string, showAll bool, limit int) {
	if filePath == "" {
		global.Cfg = &config.Config{}
		jsonFile, err := ioutil.ReadFile("./config/config.json")
		if err != nil {
			panic(fmt.Sprintf("load config err : %s", err.Error()))
		}
		err = json.Unmarshal(jsonFile, global.Cfg)
		if err != nil {
			panic(err.Error())
		}
		filePath = global.Cfg.UserDataFile
	}
	if limit <= 0 {
		limit = 20
	}

	f, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err.Error())
	}
	if len(rows) <= 1 {
		panic("csv file has no data rows")
	}

	printed := 0
	checked := 0
	caseInsensitiveMatches := 0
	exactMatches := 0
	mismatches := 0
	invalidRows := 0

	fmt.Println("checking csv file:", filePath)
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < 2 {
			panic(fmt.Sprintf("row %d has less than 2 columns", i+1))
		}

		uid := rows[i][1]
		accountID, err := normalizeAccountID(uid)
		if err != nil {
			invalidRows += 1
			if showAll || printed < limit {
				fmt.Printf("row=%d uid=%s error=%s\n", i+1, uid, err.Error())
				printed += 1
			}
			continue
		}

		checked += 1
		caseInsensitiveSame := strings.EqualFold(uid, accountID)
		exactSame := uid == accountID
		if caseInsensitiveSame {
			caseInsensitiveMatches += 1
		} else {
			mismatches += 1
		}
		if exactSame {
			exactMatches += 1
		}

		if showAll || (!caseInsensitiveSame && printed < limit) {
			fmt.Printf("row=%d uid=%s account_id=%s case_insensitive_same=%t exact_same=%t\n",
				i+1, uid, accountID, caseInsensitiveSame, exactSame)
			printed += 1
		}
	}

	fmt.Printf("summary: checked=%d case_insensitive_matches=%d exact_matches=%d mismatches=%d invalid_rows=%d\n",
		checked, caseInsensitiveMatches, exactMatches, mismatches, invalidRows)
	if mismatches == 0 && invalidRows == 0 {
		fmt.Println("all csv uid values are already consistent with generated account_id")
		return
	}
	if mismatches > 0 {
		fmt.Println("mismatch means the csv uid is not already in canonical field-element encoding")
	}
	if invalidRows > 0 {
		panic(fmt.Sprintf("found %d invalid uid rows", invalidRows))
	}
}

func normalizeAccountID(uid string) (string, error) {
	accountIDBytes, err := hex.DecodeString(uid)
	if err != nil || len(accountIDBytes) != 32 {
		return "", fmt.Errorf("uid must be a 64-char hex string decoded to 32 bytes")
	}
	return hex.EncodeToString(new(fr.Element).SetBytes(accountIDBytes).Marshal()), nil
}
