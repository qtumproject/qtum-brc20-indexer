package qtum

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/qtumproject/qtumsuite/base58"
	"golang.org/x/crypto/ripemd160"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type QtumClient struct {
	rpcUrl           string
	basicAccessToken string
	netType          string
}

func NewQtumClient(rpc string, basicAccessToken string, version string) *QtumClient {
	return &QtumClient{
		rpcUrl:           rpc,
		basicAccessToken: basicAccessToken,
		netType:          version,
	}
}

func (c QtumClient) GetUrl() string {
	return c.rpcUrl
}

func isP2PKHScript(scriptBytes []byte) bool {
	// P2PKH scripts typically have this format:
	// OP_DUP OP_HASH160 <20-byte pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG

	if len(scriptBytes) == 25 && scriptBytes[0] == 0x76 && scriptBytes[1] == 0xa9 &&
		scriptBytes[2] == 0x14 && scriptBytes[23] == 0x88 && scriptBytes[24] == 0xac {
		return true
	}

	return false
}

func hexToPublicKeyHash(hexStr string) ([]byte, error) {
	decodedHex, err := hex.DecodeString(hexStr)
	if err != nil {
		log.Fatal(err)
	}
	if !isP2PKHScript(decodedHex) {
		return []byte{}, fmt.Errorf("hexStr is not P2PKHScript type")
	}
	return decodedHex[3:23], nil
}

// 将公钥哈希转换为Qtum地址
func (c QtumClient) publicKeyHashToQtumAddress(publicKeyHash []byte) string {
	// Qtum地址以"Q"开头
	// Qtum 主网版本字节是 0x3a，测试网版本字节是 0x78
	// 使用主网版本字节
	versionByte := byte(0x3a)
	if c.netType == "test" {
		//使用测试网版本字节
		versionByte = byte(0x78)
	}

	// 添加版本字节到公钥哈希前面
	payload := append([]byte{versionByte}, publicKeyHash...)

	// 计算双 SHA256 校验和
	checksum := sha256.Sum256(payload)
	checksum = sha256.Sum256(checksum[:])
	// 创建地址
	addressBytes := append(payload, checksum[:4]...)

	// 使用Base58编码
	qtumAddress := base58.Encode(addressBytes)
	if !strings.HasPrefix(qtumAddress, "q") && !strings.HasPrefix(qtumAddress, "Q") {
		if c.netType == "main" {
			qtumAddress = "Q" + qtumAddress
		} else if c.netType == "test" {
			qtumAddress = "q" + qtumAddress
		}
	}
	return qtumAddress
}

func (c QtumClient) ExtractQtumAddressFromPubkey(script ScriptPubKey) (string, error) {
	publicKeyBytes, err := hex.DecodeString(script.Hex)
	if err != nil {
		fmt.Println("DecodeString failed:", err)
		return "", err
	}
	// SHA-256 hash
	sha256Hash := sha256.Sum256(publicKeyBytes)

	// RIPEMD-160 hash
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(sha256Hash[:])
	publicKeyHash := ripemd160Hasher.Sum(nil)

	// Prepend 0x3a for mainnet or 0x78 for testnet
	versionByte := byte(0x3a) // Mainnet
	//versionByte := byte(0x78) // Testnet
	payload := append([]byte{versionByte}, publicKeyHash...)

	// Double SHA-256 hash
	checksum := sha256.Sum256(payload)
	checksum = sha256.Sum256(checksum[:])

	// Append the first 4 bytes of the double hash as a checksum
	payload = append(payload, checksum[:4]...)

	// Encode as base58
	qtumAddress := base58.Encode(payload)
	if !strings.HasPrefix(qtumAddress, "q") && !strings.HasPrefix(qtumAddress, "Q") {
		if c.netType == "main" {
			qtumAddress = "Q" + qtumAddress
		} else if c.netType == "test" {
			qtumAddress = "q" + qtumAddress
		}
	}
	return qtumAddress, nil
}

func (c QtumClient) ExtractQtumAddressFromP2PKH(script ScriptPubKey) (string, error) {
	if len(script.Address) > 0 {
		return script.Address, nil
	}
	if strings.HasPrefix(script.Desc, "addr(q") || strings.HasPrefix(script.Desc, "addr(Q") {
		splitRes := strings.Split(script.Desc[5:], ")")
		return splitRes[0], nil
	}
	bytes, err := hexToPublicKeyHash(script.Hex)
	if err != nil {
		return "", nil
	}
	return c.publicKeyHashToQtumAddress(bytes), nil

}

// ExtractQtumAddress 提取Qtum地址
func (c QtumClient) ExtractQtumAddress(vout Vout) (string, error) {
	//hexToPublicKeyHash(vout.ScriptPubKey)
	// 根据不同的脚本类型执行适当的操作
	switch vout.ScriptPubKey.Type {
	case "pubkey":
		return c.ExtractQtumAddressFromPubkey(vout.ScriptPubKey)
	case "pubkeyhash":
		return c.ExtractQtumAddressFromP2PKH(vout.ScriptPubKey)

	case "witness_v1_taproot":
		if vout.ScriptPubKey.Address != "" {
			return vout.ScriptPubKey.Address, nil
		}
		return "", nil
	case "nonstandard":
		return "", nil
	default:
		return "", fmt.Errorf("未知的scriptPubKeyType: %s", vout.ScriptPubKey.Type)
	}
}

func (c QtumClient) getRpc(ctx context.Context, method string, url string, host string, params ...interface{}) (respBytes []byte, err error) {
	if params == nil {
		params = []interface{}{}
	}
	data, _ := json.Marshal(struct {
		Jsonrpc string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
		Id      string        `json:"id"`
	}{"1.0", method, params, strconv.Itoa(int(rand.Int31()))})
	//fmt.Println(string(data))
	respBytes, err = c.HttpRequest(ctx, "POST", url, host, bytes.NewReader(data))
	return
}

func (c QtumClient) HttpRequest(ctx context.Context, method string, requestUrl string, host string, body io.Reader) ([]byte, error) {
	var (
		req     *http.Request
		resp    *http.Response
		cancel  func()
		timeout context.Context
		err     error
	)
	for retries := 3; retries > 0; retries-- {
		timeout, cancel = context.WithTimeout(ctx, 20*time.Second)
		req, err = http.NewRequestWithContext(timeout, method, requestUrl, body)
		req.Header.Set("Content-Type", "text/plain")
		if c.basicAccessToken != "" {
			splitRes := strings.Split(c.basicAccessToken, ":")
			if len(splitRes) == 2 {
				req.SetBasicAuth(splitRes[0], splitRes[1])
			}
		}
		if err != nil {
			cancel()
			return nil, err
		} //generate timeout failed
		if host != "" {
			req.Host = host
		}
		//proxyUrl, err := url.Parse("http://127.0.0.1:7890")
		//client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		//resp, err = client.Do(req)
		resp, err = http.DefaultClient.Do(req) // request with timeout context
		if err == nil {                        // success
			break
		}
		fmt.Printf("HttpRequest to %s failed, err: %s. Retrying...\n", requestUrl, err.Error())
		cancel() // cancel the current timeout
	}
	defer cancel()
	if err != nil { // exceeds max retires
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBytes, err
}

// base58Encode 将字节数组编码为Base58字符串
func base58Encode(input []byte) string {
	var result []byte

	x := new(big.Int).SetBytes(input)

	base := big.NewInt(58)
	zero := big.NewInt(0)

	for x.Cmp(zero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, base, mod)
		result = append(result, base58Alphabet[mod.Int64()])
	}

	// 反转结果
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	// 添加前导的'1'
	for _, b := range input {
		if b == 0x00 {
			result = append([]byte{'1'}, result...)
		} else {
			break
		}
	}

	return string(result)
}

// ////////////////////////////////////////////////////
// ///////////////RPC Node Commands///////////////////
// ////////////////////////////////////////////////////

func (c QtumClient) BlockNumber(ctx context.Context) (num int64, err error) {
	respBytes, err := c.getRpc(ctx, "getblockcount", c.GetUrl(), "")
	rs := &getBlockCountResp{}

	if err = json.Unmarshal(respBytes, &rs); err != nil {
		newErr := fmt.Errorf("request response body: %s, err: %s", string(respBytes), err.Error())
		fmt.Println("【debug】" + newErr.Error())
		return num, newErr
	}
	if rs.Error != nil {
		return num, fmt.Errorf("jsonrpc call: %v", rs.Error.Message)
	}
	num = int64(rs.Result)
	return
}

func (c QtumClient) GetBlockByHeight(ctx context.Context, height int64) (*BlockInfo, error) {
	blk := &BlockInfo{}
	blkResp := &GetBlockResp{}

	respBytes, err := c.getRpc(ctx, "getblockhash", c.GetUrl(), "", height)
	rs := &GetBlockHashResp{}

	if err = json.Unmarshal(respBytes, &rs); err != nil {
		newErr := fmt.Errorf("request response body: %s, err: %s", string(respBytes), err.Error())
		fmt.Println("【debug】" + newErr.Error())
		return blk, newErr
	}
	if rs.Error != nil {
		return &blkResp.Result, fmt.Errorf("jsonrpc call: %v", rs.Error.Message)
	}
	blockHash := rs.Result

	return &blkResp.Result, c.GetBlockByHash(ctx, blkResp, blockHash, 2)
}

func (c QtumClient) GetBlockVerboseByHeight(ctx context.Context, height int64) (*BlockInfoVerbose, error) {
	blk := &BlockInfoVerbose{}
	blkResp := &GetBlockVerboseResp{}

	respBytes, err := c.getRpc(ctx, "getblockhash", c.GetUrl(), "", height)
	rs := &GetBlockHashResp{}

	if err = json.Unmarshal(respBytes, &rs); err != nil {
		newErr := fmt.Errorf("request response body: %s, err: %s", string(respBytes), err.Error())
		fmt.Println("【debug】" + newErr.Error())
		return blk, newErr
	}
	if rs.Error != nil {
		return &blkResp.Result, fmt.Errorf("jsonrpc call: %v", rs.Error.Message)
	}
	blockHash := rs.Result

	return &blkResp.Result, c.GetBlockByHash(ctx, blkResp, blockHash, 3)
}

func (c QtumClient) GetBlockByHash(ctx context.Context, getBlockResp interface{}, blockHash string, verbosity int64) error {
	respBytes, err := c.getRpc(ctx, "getblock", c.GetUrl(), "", blockHash, verbosity)
	//fmt.Println(string(respBytes))
	if err = json.Unmarshal(respBytes, &getBlockResp); err != nil {
		newErr := fmt.Errorf("request response body: %s, err: %s", string(respBytes), err.Error())
		fmt.Println("【debug】" + newErr.Error())
		return newErr
	}
	return nil
}

func (c QtumClient) GetTransactionInfo(ctx context.Context, txId, blockHash string) (*TxInfo, error) {
	blk := &TxInfo{}
	resp := &GetTransactionInfoResp{}
	var respBytes []byte
	var err error
	//blockHash = "afcb272cdcc55bd6ec9b1b1e252ef54116ac9289df6635f8da94720c516f4f4a"
	if len(blockHash) > 0 {
		respBytes, err = c.getRpc(
			ctx,
			"getrawtransaction",
			c.GetUrl(),
			"",
			txId,
			true,
			blockHash)
	} else {
		respBytes, err = c.getRpc(
			ctx,
			"getrawtransaction",
			c.GetUrl(),
			"",
			txId,
			true)
	}
	//fmt.Println(string(respBytes))
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		newErr := fmt.Errorf("request response body: %s, err: %s", string(respBytes), err.Error())
		fmt.Println("【debug】" + newErr.Error())
		return blk, newErr
	}
	return &resp.Result, nil
}

func (c QtumClient) GetLogsByQuery(ctx context.Context, logQuery FilterLogQuery) ([]TxLogInfo, error) {
	var blk []TxLogInfo
	logResp := &GetLogsResp{}
	respBytes, err := c.getRpc(
		ctx,
		"searchlogs",
		c.GetUrl(),
		"",
		logQuery.FromBlock,
		logQuery.ToBlock,
		map[string]interface{}{"addresses": logQuery.Addresses},
		map[string]interface{}{"topics": logQuery.Topics},
		logQuery.MinConf,
	)
	//fmt.Println(string(respBytes))
	if err = json.Unmarshal(respBytes, &logResp); err != nil {
		newErr := fmt.Errorf("request response body: %s, err: %s", string(respBytes), err.Error())
		fmt.Println("【debug】" + newErr.Error())
		return blk, newErr
	}
	return logResp.Result, nil
}
