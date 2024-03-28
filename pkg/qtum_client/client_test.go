package qtum

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPPP(t *testing.T) {
	//c := NewQtumClient("https://qtum-rpc.foxnb.net/main/", "foxwallet:p455w0rdisfoxnb", "main")
	c := NewQtumClient("https://qtum-rpc.foxnb.net/test/", "foxwallet:foxwallet_test_password", "test")

	res, _ := c.BlockNumber(context.Background())
	//res, _ := c.GetTransactionInfo(context.Background(), "65b2156aaee571dd6a9ec6639165ec9b570b320e74d90193a7e6f940d88a08a3", "")
	fmt.Println(res)
}

func isCompressedPublicKey(publicKeyBytes []byte) bool {
	if len(publicKeyBytes) == 33 {
		firstByte := publicKeyBytes[0]
		return firstByte == 0x02 || firstByte == 0x03
	}
	return false
}

func TestProgram(t *testing.T) {
	//ss, _ := hex.DecodeString("20bed122fa1dd431dfe84fe94c7918a13a5b2930439e1dbcf98cd5114090d3b232ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800467b2270223a226272632d3230222c226f70223a226465706c6f79222c227469636b223a22494d4442222c226d6178223a22323130303030222c226c696d223a2231303030227d68")
	//fmt.Println((string(ss[69:139])))
	//sss, _ := hex.DecodeString("20bed122fa1dd431dfe84fe94c7918a13a5b2930439e1dbcf98cd5114090d3b232ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800357b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a22494d4442222c22616d74223a2232303030227d68")

	//fmt.Println(string(sss))
	//fmt.Println(len(ss))
	//fmt.Println(len(sss))
	//fmt.Println(string(ss[67:69]))
	//strconv.Atoi(string(ss[67:69]))
	//data := ss[67:69]
	//val := binary.BigEndian.Uint16(data)
	//fmt.Println(val) // 输出: 70

	//maintest()
	//token := ""
	//blkResp := &GetBlockVerboseResp{}
	c := NewQtumClient("https://qtum-rpc.foxnb.net/main/", "foxwallet:p455w0rdisfoxnb", "main") // 正式链
	//c := NewQtumClient("https://qtum-rpc.foxnb.net/test/", "foxwallet:foxwallet_test_password") // 测试链
	//_ = c.GetBlockByHash(context.Background(), blkResp, "e0e985df063b6f56594030a9011d0450f86e2ff0b42b0506c2c9a04ec13a6ae4", 3)
	//fmt.Println(blkResp.Result.Tx[2].Vin[0].PrevOut.ScriptPubKey)
	//num, _ := c.BlockNumber(context.Background())
	//fmt.Println(num)
	//res, _ := c.GetTransactionInfo(
	//	context.Background(), "1824539d6138b9eace838193f0e6c847bc7b97bb546685041f11d151d6007662", "")
	//res, _ := c.GetBlockByHash(context.Background(), "c4588d9bd7b8e66bc56884d3b37813f11b59e145423fd4cfdecc393012a28b13")

	//address, _ := ExtractQtumAddress(res.Vout[0])
	//fmt.Println(address)

	//resp, _ := c.GetBlockByHeight(context.Background(), 3518744)
	//fmt.Println(resp)

	query := FilterLogQuery{
		FromBlock: 3458443,
		ToBlock:   3458443,
		MinConf:   0,
	}
	//query.Addresses = []string{"e7e5caae57b34b93c57af9478a5130f62e3d2827"}
	query.Topics = []string{"7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65"}
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := c.GetLogsByQuery(timeout, query)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(resp))
	//fmt.Println(len(resp[0].Log))
	fmt.Println(resp)
}

type InscriptionBRC20Content struct {
	Proto        string `json:"p,omitempty"`
	Operation    string `json:"op,omitempty"`
	BRC20Tick    string `json:"tick,omitempty"`
	BRC20Max     string `json:"max,omitempty"`
	BRC20Amount  string `json:"amt,omitempty"`
	BRC20Limit   string `json:"lim,omitempty"` // option
	BRC20Decimal string `json:"dec,omitempty"` // option
}

func maintest() {
	//content, _ := hex.DecodeString("c1bed122fa1dd431dfe84fe94c7918a13a5b2930439e1dbcf98cd5114090d3b232")
	//fmt.Println(string(content))
	//fmt.Println("end")
	//hex.DecodeString(fields[7])
	//body := new(InscriptionBRC20Content)
	//if err := body.Unmarshal(data.ContentBody); err != nil {
	//	continue
	//}
}

func TestHeight(t *testing.T) {
	// 您的 scriptPubKey 的 hex 字段
	// SegWit地址
	// scriptPubKey 的 hex 字段
	//scriptPubKeyHex := "51206bfa8c8846c0d407707b71ea628247e91e50d6c8baf24b81129b6f053cac2998"
	//
	//// 解码 hex 字段
	//scriptPubKeyBytes, err := hex.DecodeString(scriptPubKeyHex)
	//if err != nil {
	//	fmt.Println("解码 hex 字段失败:", err)
	//	return
	//}
	//
	//// 使用 bech32 编码转换地址
	//address, err := bech32.SegWitAddressEncode("tq", scriptPubKeyBytes)
	//if err != nil {
	//	fmt.Println("编码地址失败:", err)
	//	return
	//}
	//
	//fmt.Println("Qtum 地址:", address)

}

func conv111() {

}
