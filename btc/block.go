package btc

import (
	pb "github.com/Appscrunch/Multy-back/node-streamer/btc"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// ProcessTransaction from block
func (c *Client) blockTransactions(hash *chainhash.Hash) {
	log.Debugf("New block connected %s", hash.String())
	// Block Height
	blockVerbose, err := c.RPCClient.GetBlockVerbose(hash)
	if err != nil {
		log.Errorf("parseNewBlock:GetBlockVerbose: %s", err.Error())
		return
	}
	blockHeight := blockVerbose.Height

	// Parse all block transactions
	rawBlock, err := c.RPCClient.GetBlock(hash)
	allBlockTransactions, err := rawBlock.TxHashes()
	if err != nil {
		log.Errorf("parseNewBlock:rawBlock.TxHashes: %s", err.Error())
	}

	// Broadcast to client to delete mempool
	for _, hash := range allBlockTransactions {
		c.DeleteMempool <- pb.MempoolToDelete{
			Hash: hash.String(),
		}
	}

	for _, txHash := range allBlockTransactions {

		blockTxVerbose, err := c.RPCClient.GetRawTransactionVerbose(&txHash)
		if err != nil {
			log.Errorf("parseNewBlock:RPCClient.GetRawTransactionVerbose: %s", err.Error())
			continue
		}

		c.ProcessTransaction(blockHeight, blockTxVerbose, false)
	}
}
