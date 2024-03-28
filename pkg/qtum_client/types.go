package qtum

type rpcError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

//	type resStruct struct {
//		Jsonrpc string      `json:"jsonrpc,omitempty"`
//		Result  interface{} `json:"result"`
//		Error   *rpcError   `json:"error,omitempty"`
//		Id      string      `json:"id"`
//	}
type getBlockCountResp struct {
	Jsonrpc string    `json:"jsonrpc,omitempty"`
	Result  float64   `json:"result"`
	Error   *rpcError `json:"error"`
	ID      string    `json:"id"`
}

type GetBlockHashResp struct {
	Jsonrpc string    `json:"jsonrpc,omitempty"`
	Result  string    `json:"result"`
	Error   *rpcError `json:"error"`
	ID      string    `json:"id"`
}

type GetBlockResp struct {
	Jsonrpc string    `json:"jsonrpc,omitempty"`
	Result  BlockInfo `json:"result"`
	Error   *rpcError `json:"error"`
	ID      string    `json:"id"`
}

type GetBlockVerboseResp struct {
	Jsonrpc string           `json:"jsonrpc,omitempty"`
	Result  BlockInfoVerbose `json:"result"`
	Error   *rpcError        `json:"error"`
	ID      string           `json:"id"`
}

type GetLogsResp struct {
	Jsonrpc string      `json:"jsonrpc,omitempty"`
	Result  []TxLogInfo `json:"result"`
	Error   *rpcError   `json:"error"`
	ID      string      `json:"id"`
}

type BlockInfo struct {
	Hash              string   `json:"hash"`
	Confirmations     int      `json:"confirmations"`
	Height            int      `json:"height"`
	Version           int      `json:"version"`
	VersionHex        string   `json:"versionHex"`
	Merkleroot        string   `json:"merkleroot"`
	Time              int64    `json:"time"`
	Mediantime        int64    `json:"mediantime"`
	Nonce             int      `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	Chainwork         string   `json:"chainwork"`
	NTx               int      `json:"nTx"`
	HashStateRoot     string   `json:"hashStateRoot"`
	HashUTXORoot      string   `json:"hashUTXORoot"`
	PrevoutStakeHash  string   `json:"prevoutStakeHash"`
	PrevoutStakeVoutN int      `json:"prevoutStakeVoutN"`
	Previousblockhash string   `json:"previousblockhash"`
	Nextblockhash     string   `json:"nextblockhash"`
	Flags             string   `json:"flags"`
	Proofhash         string   `json:"proofhash"`
	Modifier          string   `json:"modifier"`
	Signature         string   `json:"signature"`
	ProofOfDelegation string   `json:"proofOfDelegation"`
	Strippedsize      int      `json:"strippedsize"`
	Size              int      `json:"size"`
	Weight            int      `json:"weight"`
	Tx                []TxInfo `json:"tx"`
}

type BlockInfoVerbose struct {
	Hash              string          `json:"hash"`
	Confirmations     int             `json:"confirmations"`
	Height            int             `json:"height"`
	Version           int             `json:"version"`
	VersionHex        string          `json:"versionHex"`
	Merkleroot        string          `json:"merkleroot"`
	Time              int64           `json:"time"`
	Mediantime        int64           `json:"mediantime"`
	Nonce             int             `json:"nonce"`
	Bits              string          `json:"bits"`
	Difficulty        float64         `json:"difficulty"`
	Chainwork         string          `json:"chainwork"`
	NTx               int             `json:"nTx"`
	HashStateRoot     string          `json:"hashStateRoot"`
	HashUTXORoot      string          `json:"hashUTXORoot"`
	PrevoutStakeHash  string          `json:"prevoutStakeHash"`
	PrevoutStakeVoutN int             `json:"prevoutStakeVoutN"`
	Previousblockhash string          `json:"previousblockhash"`
	Nextblockhash     string          `json:"nextblockhash"`
	Flags             string          `json:"flags"`
	Proofhash         string          `json:"proofhash"`
	Modifier          string          `json:"modifier"`
	Signature         string          `json:"signature"`
	ProofOfDelegation string          `json:"proofOfDelegation"`
	Strippedsize      int             `json:"strippedsize"`
	Size              int             `json:"size"`
	Weight            int             `json:"weight"`
	Tx                []TxInfoVerbose `json:"tx"`
}

type GetTransactionInfoResp struct {
	Jsonrpc string    `json:"jsonrpc,omitempty"`
	Result  TxInfo    `json:"result"`
	Error   *rpcError `json:"error"`
	ID      string    `json:"id"`
}

type TxInfo struct {
	InActiveChain bool   `json:"in_active_chain,omitempty"`
	Txid          string `json:"txid"`
	Hash          string `json:"hash"`
	Version       int    `json:"version"`
	Size          int    `json:"size"`
	Vsize         int    `json:"vsize"`
	Weight        int    `json:"weight"`
	Locktime      int    `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	Hex           string `json:"hex"`
	Blockhash     string `json:"blockhash,omitempty"`
	Confirmations int    `json:"confirmations,omitempty"`
	Time          int    `json:"time,omitempty"`
	Blocktime     int    `json:"blocktime,omitempty"`
}

type TxInfoVerbose struct {
	InActiveChain bool         `json:"in_active_chain,omitempty"`
	Txid          string       `json:"txid"`
	Hash          string       `json:"hash"`
	Version       int          `json:"version"`
	Size          int          `json:"size"`
	Vsize         int          `json:"vsize"`
	Weight        int          `json:"weight"`
	Locktime      int64        `json:"locktime"`
	Vin           []VinVerbose `json:"vin"`
	Vout          []Vout       `json:"vout"`
	Hex           string       `json:"hex"`
	Blockhash     string       `json:"blockhash,omitempty"`
	Confirmations int          `json:"confirmations,omitempty"`
	Time          int64        `json:"time,omitempty"`
	Blocktime     int64        `json:"blocktime,omitempty"`
}

type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}
type Vin struct {
	Txid        string    `json:"txid"`
	Vout        int       `json:"vout"`
	Coinbase    string    `json:"coinbase"`
	ScriptSig   ScriptSig `json:"scriptSig"`
	Txinwitness []string  `json:"txinwitness"`
	Sequence    int64     `json:"sequence"`
}

type VinVerbose struct {
	Txid        string    `json:"txid"`
	Vout        int       `json:"vout"`
	Coinbase    string    `json:"coinbase"`
	ScriptSig   ScriptSig `json:"scriptSig"`
	PrevOut     PrevOut   `json:"prevout"`
	Txinwitness []string  `json:"txinwitness"`
	Sequence    int64     `json:"sequence"`
}

type ScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address,omitempty"`
	Type    string `json:"type"`
}

type PrevOut struct {
	Generated    bool         `json:"generated"`
	Height       int64        `json:"height"`
	Value        float64      `json:"value"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type FilterLogQuery struct {
	FromBlock int64
	ToBlock   int64
	Addresses []string
	Topics    []string
	MinConf   int64
}

type TxLogInfo struct {
	BlockHash         string    `json:"blockHash"`
	BlockNumber       uint64    `json:"blockNumber"`
	TransactionHash   string    `json:"transactionHash"`
	TransactionIndex  uint      `json:"transactionIndex"`
	OutputIndex       uint      `json:"outputIndex"`
	From              string    `json:"from"`
	To                string    `json:"to"`
	CumulativeGasUsed uint      `json:"cumulativeGasUsed"`
	GasUsed           uint      `json:"gasUsed"`
	ContractAddress   string    `json:"contractAddress"`
	Excepted          string    `json:"excepted"`
	ExceptedMessage   string    `json:"exceptedMessage"`
	Bloom             string    `json:"bloom"`
	StateRoot         string    `json:"stateRoot"`
	UtxoRoot          string    `json:"utxoRoot"`
	Log               []LogInfo `json:"log"`
}

type LogInfo struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}
