// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package engine

import (
	"encoding/json"
	"github.com/chain5j/chain5j-pkg/util/dateutil"
	log "github.com/chain5j/log15"
	"github.com/chain5j/sync_eth/apis/sync/chain"
	"github.com/chain5j/sync_eth/apis/sync/dao"
	"github.com/chain5j/sync_eth/apis/sync/models"
	"github.com/chain5j/sync_eth/params"
	"github.com/chain5j/sync_eth/params/global"
	"github.com/olivere/elastic/v7"
	"github.com/shopspring/decimal"
	"math/big"
)

func (s *Sync) parseTx(block *chain.Block, blockInfo *models.Block, tx *chain.Transaction) (transferValue decimal.Decimal, err error) {
	has, _ := global.LevelDb.Has([]byte(params.DBKEY_PRE_TX + tx.Hash))
	if !has {
		bytes, err := tx.Bytes()
		if err != nil {
			log.Error("chain.Transaction to bytes err", "err", err)
			return s.parseTx(block, blockInfo, tx)
		} else {
			global.LevelDb.Put([]byte(params.DBKEY_PRE_TX+tx.Hash), bytes)
		}
	}

	var (
		fee      = decimal.Zero
		txStatus = false

		isContract = false
		to         = tx.To
	)

	chainTx := &models.Transaction{
		TxHash:   tx.Hash,
		From:     tx.From.Hex(),
		Nonce:    decimal.NewFromBigInt(big.NewInt(int64(tx.Nonce)), 0),
		Value:    decimal.NewFromBigInt(tx.Value.ToInt(), 0),
		GasPrice: decimal.NewFromBigInt(tx.GasPrice.ToInt(), 0),
		GasLimit: uint64(tx.GasLimit),
		Input:    tx.Input.String(),

		BlockNumber:      decimal.NewFromBigInt(big.NewInt(int64(tx.BlockNumber)), 0),
		BlockTime:        blockInfo.BlockTime,
		TransactionIndex: uint64(tx.TransactionIndex),
	}

	// receipt
	startTime := dateutil.CurrentTime()
	var txReceipt *chain.TransactionReceipt
	{
		has, err := global.LevelDb.Has([]byte(params.DBKEY_PRE_TXRECEIPT + tx.Hash))
		if err == nil && has {
			bytes, err := global.LevelDb.Get([]byte(params.DBKEY_PRE_TXRECEIPT + tx.Hash))
			if err == nil && len(bytes) > 0 {
				err = json.Unmarshal(bytes, &txReceipt)
				if err != nil {
					log.Error("json.Unmarshal(bytes, &txReceipt)", "err", err)
				}
			}
		}
		if txReceipt == nil {
			err = s.eth.GetTransactionReceipt(tx.Hash, &txReceipt)
			if err != nil {
				log.Error("GetTransactionReceipt err", "hash", tx.Hash, "err", err)
				return decimal.Zero, err
			}
			if txReceipt != nil {
				bytes1, err := txReceipt.Bytes()
				if err != nil {
					log.Error("txReceipt.Bytes err", "err", err)
				} else {
					global.LevelDb.Put([]byte(params.DBKEY_PRE_TXRECEIPT+tx.Hash), bytes1)
				}
			}
		}
		if txReceipt != nil {
			fee = decimal.NewFromBigInt(new(big.Int).Mul(tx.GasPrice.ToInt(), big.NewInt(int64(txReceipt.GasUsed))), 0)

			chainTx.GasUsed = decimal.NewFromBigInt(big.NewInt(int64(txReceipt.GasUsed)), 0)
			chainTx.Fee = fee
			if txReceipt.ContractAddress != nil {
				chainTx.Contract = txReceipt.ContractAddress.Hex()
			}
		}

		if to != nil {
			has, _ := global.LevelDb.Has([]byte(params.DBKEY_PRE_CONTRACT + to.Hex()))
			if has {
				isContract = true
			} else {
				if len(tx.Input.String()) > 3 {
					contract, err := dao.ContractDao.Info(s.es, to.Hex())
					if err != nil {
						if elastic.IsNotFound(err) || contract == nil {
							if global.Config.ChainConfig.SyncStartBlock >= 1 {
								isContract, err = s.eth.IsContract(to.Hex())
								if err != nil {
									log.Error("isContract err", "err", err)
									return decimal.Decimal{}, err
								}
								if isContract {
									err = s.getContractInfo(chainTx, to.Hex(), false)
									if err != nil {
										log.Error("getContractInfo err", "err", err)
										return decimal.Decimal{}, err
									}
								}
							}
						} else {
							log.Error("dao.ContractDao.Info err", "err", err)
							return decimal.Decimal{}, err
						}
					}
				}
			}
		}
		txStatus = chain.TxStatus(txReceipt, isContract)
		chainTx.Status = txStatus
	}
	log.Trace("TransactionReceipt elapsed", "txHash", tx.Hash, "elapsed", dateutil.GetDistanceTime(dateutil.CurrentTime()-startTime))

	toStr := ""
	if to != nil {
		toStr = to.Hex()
	}
	chainTx.To = toStr

	if txStatus {
		transferValue = decimal.NewFromBigInt(tx.Value.ToInt(), 0)
		blockInfo.BlockAward = blockInfo.BlockAward.Add(fee)
	} else {
		blockInfo.BlockAward = blockInfo.BlockAward.Add(fee)
		transferValue = decimal.Zero
	}

	if to == nil {
		// deploy contract
		chainTx.TxType = models.TokenDeploy
		err = s.parseDeployTx(chainTx)
	} else {
		if isContract {
			// token tx
			chainTx.TxType = models.TokenContract
			err = s.parseContract(chainTx)
		} else {
			// common tx
			chainTx.TxType = models.TokenETH
			if txStatus {
				if s.isUseMQ {
					err = s.mqTxRaw.Produce(chainTx)
					if err != nil {
						log.Error("RabbitMq.Produce TransactionRaw", "txHash", chainTx.TxHash, "err", err)
						return decimal.Decimal{}, err
					}
				}
			}
		}
	}
	if err != nil {
		log.Error("parse tx err", "err", err)
		return
	}
	s.tempRawTxList = append(s.tempRawTxList, chainTx)
	return
}

func (s *Sync) parseDeployTx(chainTx *models.Transaction) (err error) {
	txType := models.TokenDeploy
	chainTx.TxType = txType

	if chainTx.Status {
		contractAddress := chainTx.Contract
		err = s.getContractInfo(chainTx, contractAddress, true)
		if err != nil {
			return err
		}
	}

	return
}

func (s *Sync) getContractInfo(chainTx *models.Transaction, contractAddress string, isDeploy bool) error {
	totalSupply := decimal.Zero
	supply, err := s.eth.GetTokenTotalSupply(contractAddress)
	if err != nil {
		log.Error("GetTokenTotalSupply err", "contractAddress", contractAddress, "err", err)
		return err
	} else {
		totalSupply = decimal.NewFromBigInt(&supply, 0)
	}

	decimals := int64(0)
	tokenDecimals, err := s.eth.GetTokenDecimals(contractAddress)
	if err != nil {
		log.Error("GetTokenName err", "contractAddress", contractAddress, "err", err)
		return err
	} else {
		decimals = tokenDecimals.Int64()
	}

	// token name
	tokenName := ""
	tokenName, err = s.eth.GetTokenName(contractAddress)
	if err != nil {
		log.Error("GetTokenName err", "contractAddress", contractAddress, "err", err)
		return err
	}
	log.Debug("tokenName", "tokenName", tokenName)
	// symbol
	symbol := ""
	symbol, err = s.eth.GetTokenSymbol(contractAddress)
	if err != nil {
		log.Error("GetTokenSymbol err", "contractAddress", contractAddress, "err", err)
		return err
	}
	log.Debug("symbol", "symbol", symbol)

	contractType := 0
	if totalSupply.Cmp(decimal.Zero) > 0 && tokenName != "" && symbol != "" {
		contractType = 1 //ERC20
	}

	// 1）save contract
	contract := &models.Contract{
		Address:      contractAddress,
		Name:         tokenName,
		Symbol:       symbol,
		Icon:         "",
		UnitLength:   decimals,
		Introduction: "",
		WebSite:      "",
		Total:        totalSupply,
		Type:         contractType,
		//Timestamp:    chainTx.BlockTime,
		//Block:        uint64(chainTx.BlockNumber.IntPart()),
		//Creator:      chainTx.From,
		//TxHash:       chainTx.TxHash,
	}
	if isDeploy {
		contract.Creator = chainTx.From
		contract.BlockNumber = uint64(chainTx.BlockNumber.IntPart())
		contract.BlockTime = chainTx.BlockTime
		contract.TxHash = chainTx.TxHash
	}
	bytes, _ := contract.Bytes()
	global.LevelDb.Put([]byte(params.DBKEY_PRE_CONTRACT+contractAddress), bytes)
	s.tempContractMapSet.Update(contract)
	return nil
}

// parse contract
func (s *Sync) parseContract(chainTx *models.Transaction) (err error) {
	txType := models.TokenContract
	chainTx.TxType = txType
	chainTx.Contract = chainTx.To

	contractType := 0 // 0：unknown，1：ERC20，2：ERC721
	contract, ok := s.tempContractMapSet.Get(chainTx.To)
	if !ok {
		contract, err = dao.ContractDao.Info(s.es, chainTx.To)
		if err != nil {
			log.Error("ContractDao.Info", "err", err)
			return
		}
	}
	if contract != nil {
		contractType = contract.Type
	}
	var tokenTransferParams *chain.TokenTransferParams
	if contractType == 1 {
		txType = models.TokenERC20

		tokenTransferInput := chain.ParseInput(chainTx.Input)
		if tokenTransferInput != nil {
			tokenTransferParams = chain.ParseTransferInput(tokenTransferInput)
		}
	}

	contract.TxNum = contract.TxNum + 1
	s.tempContractMapSet.Update(contract)

	switch contractType {
	case 1:
		// ERC20
		s.parseContractErc20(chainTx, contract, tokenTransferParams)
	//case 2:
	//	// ERC721
	default:
		// other
	}
	return
}

func (s *Sync) parseContractErc20(chainTx *models.Transaction, contract *models.Contract, tokenTransferParams *chain.TokenTransferParams) (err error) {
	if tokenTransferParams == nil {
		return
	}
	txType := models.TokenERC20
	chainTx.TxType = txType
	chainTx.To = tokenTransferParams.To
	chainTx.Value = tokenTransferParams.Value

	contract.TokenTxNum = contract.TokenTxNum + 1
	s.tempContractMapSet.Update(contract)

	if chainTx.Status {
		contract.TokenCirculation = contract.TokenCirculation.Add(tokenTransferParams.Value)
		s.tempContractMapSet.Update(contract)
		if s.isUseMQ {
			err = s.mqTxRaw.Produce(chainTx)
			if err != nil {
				log.Error("RabbitMq.Produce TransactionRaw", "txHash", chainTx.TxHash, "err", err)
				return err
			}
		}
	}

	return
}
