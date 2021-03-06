package posconfig

import (
	"math/big"

	"github.com/wanchain/go-wanchain/accounts/keystore"
	"github.com/wanchain/go-wanchain/common"
	"github.com/wanchain/go-wanchain/node"
	bn256 "github.com/wanchain/go-wanchain/crypto/bn256"
)

var (
	// EpochBaseTime is the pos start time such as: 2018-12-12 00:00:00 == 1544544000
	EpochBaseTime = uint64(0)
	// SelfTestMode config whether it is in a simlate tese mode
	SelfTestMode = false
)

const (
	RbLocalDB  = "rblocaldb"
	EpLocalDB  = "eplocaldb"
	PosLocalDB = "pos"
)

const (
	// EpochLeaderCount is count of pk in epoch leader group which is select by stake
	EpochLeaderCount = 50
	// RandomProperCount is count of pk in random leader group which is select by stake
	RandomProperCount = 21
	// SlotTime is the time span of a slot in second, So it's 1 hours for a epoch
	SlotTime = 10
	// GenesisPK is the epoch0 pk


	//Incentive should perform delay some epochs.
	IncentiveDelayEpochs = 1
	IncentiveStartStage  = Stage2K

	// K count of each epoch
	KCount = 12
	K      = 10
	// SlotCount is slot count in an epoch
	SlotCount = K * KCount

	// Stage1K is divde a epoch into 10 pieces
	Stage1K  = uint64(K)
	Stage2K  = Stage1K * 2
	Stage3K  = Stage1K * 3
	Stage4K  = Stage1K * 4
	Stage5K  = Stage1K * 5
	Stage6K  = Stage1K * 6
	Stage7K  = Stage1K * 7
	Stage8K  = Stage1K * 8
	Stage9K  = Stage1K * 9
	Stage10K = Stage1K * 10
	Stage11K = Stage1K * 11
	Stage12K = Stage1K * 12

	Sma1Start = 0
	Sma1End   = Stage3K
	Sma2Start = Stage6K
	Sma2End   = Stage8K
	Sma3Start = Stage10K
	Sma3End   = Stage12K
)
var GenesisPK = "04dc40d03866f7335e40084e39c3446fe676b021d1fcead11f2e2715e10a399b498e8875d348ee40358545e262994318e4dcadbc865bcf9aac1fc330f22ae2c786"
type Config struct {
	PolymDegree   uint
	K             uint
	RBThres       uint
	EpochInterval uint64
	PosStartTime  int64
	MinerKey      *keystore.Key
	Dbpath        string
	NodeCfg       *node.Config
	Dkg1End       uint64
	Dkg2Begin     uint64
	Dkg2End       uint64
	SignBegin     uint64
	SignEnd       uint64
}

var DefaultConfig = Config{
	10,
	K,
	11,
	0,
	0,
	nil,
	"",
	nil,
	Stage2K - 1,
	Stage4K,
	Stage6K - 1,
	Stage8K,
	Stage10K - 1,
}

func Cfg() *Config {
	return &DefaultConfig
}

func (c *Config) GetMinerAddr() common.Address {
	if c.MinerKey == nil {
		return common.Address{}
	}

	return c.MinerKey.Address
}

func (c *Config) GetMinerBn256PK() *bn256.G1 {
	if c.MinerKey == nil {
		return nil
	}

	return new(bn256.G1).Set(c.MinerKey.PrivateKey3.PublicKeyBn256.G1)
}

func (c *Config) GetMinerBn256SK() *big.Int {
	if c.MinerKey == nil {
		return nil
	}

	return new(big.Int).Set(c.MinerKey.PrivateKey3.D)
}
