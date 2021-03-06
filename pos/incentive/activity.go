package incentive

import (
	"github.com/wanchain/go-wanchain/common"
	"github.com/wanchain/go-wanchain/consensus"
	"github.com/wanchain/go-wanchain/core/vm"
	"github.com/wanchain/go-wanchain/crypto"
	"github.com/wanchain/go-wanchain/log"
	"github.com/wanchain/go-wanchain/pos/util"
	"github.com/wanchain/go-wanchain/pos/util/convert"
)

func getEpochLeaderActivity(stateDb vm.StateDB, epochID uint64) ([]common.Address, []int) {
	epochLeaders := util.GetEpocherInst().GetEpochLeaders(epochID)
	if epochLeaders == nil || len(epochLeaders) == 0 {
		log.Error("incentive activity GetEpochLeaders error", "epochID", epochID)
		return []common.Address{}, []int{}
	}

	addrs := make([]common.Address, len(epochLeaders))
	activity := make([]int, len(addrs))
	for i := 0; i < len(addrs); i++ {
		addrs[i] = crypto.PubkeyToAddress(*crypto.ToECDSAPub(epochLeaders[i]))
		activity[i] = 0
	}

	for i := 0; i < len(addrs); i++ {
		epochIDBuf := convert.Uint64ToBytes(epochID)
		selfIndexBuf := convert.Uint64ToBytes(uint64(i))
		keyHash := vm.GetSlotLeaderStage2KeyHash(epochIDBuf, selfIndexBuf)

		data := stateDb.GetStateByteArray(vm.GetSlotLeaderSCAddress(), keyHash)
		if data == nil {
			continue
		}

		epID, slfIndex, selfPk, _, _, err := vm.RlpUnpackStage2DataForTx(data)
		if err != nil {
			continue
		}

		if epID != epochID || uint64(i) != slfIndex {
			continue
		}
		//TODO: CHECK
		addr := crypto.PubkeyToAddress(*selfPk)
		if addr.Hex() == addrs[i].Hex() {
			activity[i] = 1
		}
	}

	return addrs, activity
}

func getRandomProposerActivity(stateDb vm.StateDB, epochID uint64) ([]common.Address, []int) {
	if getRandomProposerAddress == nil {
		log.Error("incentive activity getRandomProposerAddress == nil", "epochID", epochID)
		return []common.Address{}, []int{}
	}

	leaders := getRandomProposerAddress(epochID)
	addrs := make([]common.Address, len(leaders))
	for i := 0; i < len(leaders); i++ {
		addrs[i] = leaders[i].SecAddr
	}

	if (addrs == nil) || (len(addrs) == 0) {
		log.Error("incentive activity getRandomProposerAddress error", "epochID", epochID)
		return []common.Address{}, []int{}
	}

	activity := make([]int, len(addrs))
	for i := 0; i < len(addrs); i++ {
		if vm.IsRBActive(stateDb, epochID, uint32(i)) {
			activity[i] = 1
		} else {
			activity[i] = 0
		}
	}
	return addrs, activity
}

func getSlotLeaderActivity(chain consensus.ChainReader, epochID uint64, slotCount int) ([]common.Address, []int, float64) {
	currentNumber := chain.CurrentHeader().Number.Uint64()
	miners := make(map[common.Address]int)
	for i := currentNumber - 1; i > 0; i-- {
		header := chain.GetHeaderByNumber(i)
		if header == nil {
			continue
		}

		epID := getEpochIDFromDifficulty(header.Difficulty)
		if epID == epochID {
			cnt, ok := miners[header.Coinbase]
			if ok {
				cnt++
				miners[header.Coinbase] = cnt
			} else {
				miners[header.Coinbase] = 1
			}
		}

		if epID < epochID {
			break
		}
	}

	addrs := make([]common.Address, 0)
	blocks := make([]int, 0)

	for k, v := range miners {
		addrs = append(addrs, k)
		blocks = append(blocks, v)
	}

	epochBlockCnt := sumIntArray(blocks)
	if epochBlockCnt > slotCount {
		epochBlockCnt = slotCount
	}
	activePercent := float64(epochBlockCnt) / float64(slotCount)
	return addrs, blocks, activePercent
}
