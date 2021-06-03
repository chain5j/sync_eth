// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/chain5j/chain5j-pkg/math"
	"github.com/chain5j/chain5j-pkg/util/dateutil"
	log "github.com/chain5j/log15"
	"github.com/chain5j/sync_eth/apis/sync/chain"
	"github.com/chain5j/sync_eth/apis/sync/dao"
	"github.com/chain5j/sync_eth/apis/sync/models"
	"github.com/chain5j/sync_eth/apis/sync/mqImpl"
	"github.com/chain5j/sync_eth/params"
	"github.com/chain5j/sync_eth/params/global"
	"github.com/chain5j/sync_eth/pkg/database/es"
	"github.com/chain5j/sync_eth/pkg/database/mq"
	"github.com/chain5j/sync_eth/pkg/util/convutil"
	"github.com/olivere/elastic/v7"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

type Sync struct {
	dbBlockNumber      int64
	eth                *chain.Eth
	wg                 sync.WaitGroup
	tempContractMapSet *models.ContractMapSet
	isRawTxInsert      bool
	tempRawTxList      []*models.Transaction

	es      *es.ES
	mqTxRaw *mqImpl.TransactionRawMq
	isUseMQ bool
}

func NewSync(esHosts []string) (*Sync, error) {
	newES, err := es.NewES(esHosts)
	if err != nil {
		return nil, err
	}

	isUseMQ := global.Config.Database.Mq.IsUse
	var mqTxRaw *mqImpl.TransactionRawMq
	if isUseMQ {
		chainName := global.Config.ChainConfig.ChainName
		mqTxRaw1 := mq.NewDefaultMq(
			global.Config.Database.Mq,
			chainName+"."+mqImpl.MqTxRawEx,
			chainName+"."+mqImpl.MqTxRawQueue,
			chainName+"."+mqImpl.MqTxRawRouteKey)
		mqTxRaw = mqImpl.NewTxRawMq(mqTxRaw1, newES)
	}

	return &Sync{
		eth:                global.RpcClient,
		tempContractMapSet: models.NewContractMapSet(),
		tempRawTxList:      make([]*models.Transaction, 0),
		es:                 newES,
		isRawTxInsert:      true,
		mqTxRaw:            mqTxRaw,
	}, nil
}

// Start Start
func (s *Sync) Start() error {
	err := s.initEsDb()
	if err != nil {
		return err
	}
	if s.isUseMQ {
		go s.mqTxRaw.Consumer()
	}
	go s.listen()
	return nil
}

func (s *Sync) initEsDb() error {
	ctx := context.Background()
	blockModel := models.Block{}
	exist, err := s.es.Client().IndexExists(blockModel.TableName()).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exist {
		// create
		do, err := s.es.Client().CreateIndex(blockModel.TableName()).Body(es.Mapping(blockModel)).Do(context.Background())
		if err != nil {
			log.Error("blockModel CreateIndex err", "err", err)
			return err
		}
		if !do.Acknowledged {
			log.Error("blockModel createIndex.Acknowledged", "acknowledged", do.Acknowledged)
			return errors.New(fmt.Sprintf("blockModel createIndex.Acknowledged:%t", do.Acknowledged))
		}
	}

	transactionModel := models.Transaction{}
	exist, err = s.es.Client().IndexExists(transactionModel.TableName()).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exist {
		// create
		do, err := s.es.Client().CreateIndex(transactionModel.TableName()).Body(es.Mapping(transactionModel)).Do(context.Background())
		if err != nil {
			log.Error("transactionModel CreateIndex err", "err", err)
			return err
		}
		if !do.Acknowledged {
			log.Error("transactionModel createIndex.Acknowledged", "acknowledged", do.Acknowledged)
			return errors.New(fmt.Sprintf("tansactionModel createIndex.Acknowledged:%t", do.Acknowledged))
		}
	}

	contractModel := models.Contract{}
	exist, err = s.es.Client().IndexExists(contractModel.TableName()).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exist {
		// 创建
		do, err := s.es.Client().CreateIndex(contractModel.TableName()).Body(es.Mapping(contractModel)).Do(context.Background())
		if err != nil {
			log.Error("contractModel CreateIndex err", "err", err)
			return err
		}
		if !do.Acknowledged {
			log.Error("contractModel createIndex.Acknowledged", "acknowledged", do.Acknowledged)
			return errors.New(fmt.Sprintf("contractModel createIndex.Acknowledged:%t", do.Acknowledged))
		}
	}

	listenAddressModel := models.ListenAddress{}
	exist, err = s.es.Client().IndexExists(listenAddressModel.TableName()).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exist {
		// 创建
		do, err := s.es.Client().CreateIndex(listenAddressModel.TableName()).Body(es.Mapping(listenAddressModel)).Do(context.Background())
		if err != nil {
			log.Error("listenAddressModel CreateIndex err", "err", err)
			return err
		}
		if !do.Acknowledged {
			log.Error("listenAddressModel createIndex.Acknowledged", "acknowledged", do.Acknowledged)
			return errors.New(fmt.Sprintf("listenAddressModel createIndex.Acknowledged:%t", do.Acknowledged))
		}
	}
	return nil
}

// init account
func (s *Sync) initAccount() bool {
	if global.Config.ChainConfig.Alloc != nil {
		for k, v := range global.Config.ChainConfig.Alloc {
			account, err := dao.AccountsDao.GetAccount(k, "")
			if err == sql.ErrNoRows || (err == nil && account == nil) {
				bigInt := decimal.NewFromBigInt(math.MustParseBig256(v), 0)
				a := &models.Accounts{
					Address:          k,
					Balance:          bigInt,
					BalanceInput:     bigInt,
					AccountType:      models.TypeAccountNormal,
					Contract:         "",
					TransactionCount: 0,
				}
				err1 := dao.AccountsDao.Insert(a)
				if err1 != nil {
					log.Error("init account err", "err", err1)
					return false
				}
			} else if err != nil {
				log.Error("dao.AccountsDao.GetAccount err", "err", err)
				return false
			}
		}
	}

	return true
}

func (s *Sync) listen() {
	log.Info("syncing......")
	// 获取数据库中的区块高度
	latestBlock, err := dao.BlockDao.LatestBlock(s.es)
	if err != nil {
		log.Error("BlockDao.DbLatestBlock get err", "err", err)
		return
	}
	if latestBlock == nil {
		s.dbBlockNumber = global.Config.ChainConfig.SyncStartBlock - 1
	} else {
		latestBlock, err = dao.BlockDao.VerifyBlock(s.es, 0, latestBlock)
		if err != nil {
			log.Error("verify block err", "err", err)
			return
		}
		s.dbBlockNumber = int64(latestBlock.BlockNumber)
	}
	// start sync
	s.wg.Add(1)
	go s.sync()

	s.wg.Wait()
}

// sync
func (s *Sync) sync() {
	for {
		var isNeedWait = true
		var latestBlock chain.BlockTxHashes
		err := s.eth.GetLatestBlock(false, &latestBlock)
		if err != nil {
			log.Error("GetLatestBlock error", "err", err)
			time.Sleep(10 * time.Second)
			continue
		}
		// Synchronize only data that is considered unchangeable
		latestBlockNumber := uint64(latestBlock.Number) - global.Config.ChainConfig.MinConfirms

		if int64(latestBlockNumber) > s.dbBlockNumber {
			for i := s.dbBlockNumber + 1; i <= int64(latestBlockNumber); i++ {
				// 循环开始获取区块内容
				var block *chain.Block
				err := s.eth.GetBlockByNumber(uint64(i), true, &block)
				if err != nil {
					log.Error("get block from chain err", "err", err)
					break
				}
				if block == nil {
					log.Error("get block from chain nil")
					break
				}

				{
					preBlockNumber := uint64(i - 1)
					if i > (global.Config.ChainConfig.SyncStartBlock) {
						dbBlock, err := dao.BlockDao.GetBlockByHeight(s.es, preBlockNumber)
						if err != nil {
							if elastic.IsNotFound(err) {
								s.dbBlockNumber = int64(preBlockNumber - 1)
								isNeedWait = false
								break
							} else {
								log.Error("get block from db err", "blockNumber", preBlockNumber, "err", err)
								break
							}
						}
						if dbBlock == nil {
							log.Error("get block from db nil", "blockNumber", preBlockNumber)
							break
						}
						if dbBlock.BlockHash != block.ParentHash {
							log.Info("rollback data", "preBlockNumber", preBlockNumber)
							err := dao.BlockDao.Rollback(s.es, preBlockNumber)
							if err != nil {
								log.Error("rollback failed", "err", err)
								time.Sleep(5 * time.Second)
								continue
							}
							s.dbBlockNumber = int64(preBlockNumber - 1)
							isNeedWait = false
							break
						}
					}
				}

				log.Info("block processing ......", "block", i, "txsLen", len(block.Transactions))
				startTime := dateutil.CurrentTime()
				err = s.parseBlock(block)
				if err != nil {
					log.Error("parse block err", "err", err)
					time.Sleep(10 * time.Second)
					break
				}
				log.Info("block update success", "block", i, "txsLen", len(block.Transactions), "elapsed", dateutil.GetDistanceTime(dateutil.CurrentTime()-startTime))
				s.dbBlockNumber++
			}
		}
		if isNeedWait {
			log.Info("wait block", "waitBlockNumber", s.dbBlockNumber+1)
			time.Sleep(10 * time.Second)
		}
	}
}

func (s *Sync) parseBlock(block *chain.Block) error {
	var err error
	if block == nil {
		return errors.New("chain.Block is empty")
	}
	blockTxHashes := chain.Block2BlockTxHashes(block)
	if blockTxHashes != nil {
		bytes, err := blockTxHashes.Bytes()
		if err == nil && bytes != nil {
			err := global.LevelDb.Put([]byte(params.DBKEY_PRE_BLOCK+convutil.ToString(uint64(block.Number))), bytes)
			if err != nil {
				log.Error("LevelDB save block err", "blockNumber", block.Number, "err", err)
			}
		}
	}
	//
	blockInfo := &models.Block{
		BlockNumber: uint64(block.Number),
		BlockHash:   block.Hash,
		Miner:       block.Coinbase,
		BlockTime:   uint64(block.Timestamp),
		TxCount:     len(block.Transactions),
		BlockSize:   uint64(block.Size),
		ParentHash:  block.ParentHash,
		Txs:         blockTxHashes.Transactions,
	}

	log.Debug(fmt.Sprintf("===============================%d===========================", blockInfo.BlockNumber))
	baseBalanceTransferValue := decimal.Zero
	if block.Transactions != nil && len(block.Transactions) > 0 {
		parseTxStartTime := dateutil.CurrentTime()
		for _, tx := range block.Transactions {
			transferValue, err := s.parseTx(block, blockInfo, tx)
			if err != nil {
				log.Error("parseTx err", "err", err)
				s.tempContractMapSet = models.NewContractMapSet()
				s.tempRawTxList = make([]*models.Transaction, 0)
				return errors.New("parseTx err")
			}
			baseBalanceTransferValue = baseBalanceTransferValue.Add(transferValue)
		}
		log.Debug("parse block elapsed", "blockNumber", blockInfo.BlockNumber, "txLen", blockInfo.TxCount, "elapsed", dateutil.GetDistanceTime(dateutil.CurrentTime()-parseTxStartTime))
	}

	startTime := dateutil.CurrentTime()
	err = dao.BlockDao.BatchSave(s.es, blockInfo, s.tempContractMapSet, s.tempRawTxList)
	if err != nil {
		log.Error("BlockDao.BatchSave err", "err", err)
		s.tempContractMapSet = models.NewContractMapSet()
		s.tempRawTxList = make([]*models.Transaction, 0)
		return err
	}
	log.Debug("batch save elapsed", "blockNumber", uint64(block.Number), "txLen", blockInfo.TxCount, "elapsed", dateutil.GetDistanceTime(dateutil.CurrentTime()-startTime))
	s.tempContractMapSet = models.NewContractMapSet()
	s.tempRawTxList = make([]*models.Transaction, 0)
	return nil
}
