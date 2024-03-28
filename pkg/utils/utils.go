package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/jmcvetta/randutil"
	"github.com/oschwald/geoip2-golang"
	"math"
	"math/big"
	mathRand "math/rand"
	"net"
	"sort"
	"strconv"
	"strings"
)

const (
	activeTokenLen    = 32
	tokenElementTable = "*+=.0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	tokenElementLen   = byte(len(tokenElementTable))

	Wei   = 1
	GWei  = 1e9
	Ether = 1e18
)

func GenActiveToken() string {
	keyBytes := make([]byte, activeTokenLen)
	rand.Read(keyBytes)
	for i := 0; i < activeTokenLen; i++ {
		keyBytes[i] = tokenElementTable[keyBytes[i]%tokenElementLen]
	}
	return string(keyBytes)
}

func StringInSlice(arr []string, str string) bool {
	for _, element := range arr {
		if str == element {
			return true
		}
	}
	return false
}

func Int64InSlice(arr []int64, element int64) bool {
	for _, e := range arr {
		if e == element {
			return true
		}
	}
	return false
}

func IntInSlice(arr []int, element int) bool {
	for _, e := range arr {
		if e == element {
			return true
		}
	}
	return false
}

func StringDeduplicate(arr []string) []string {
	if len(arr) <= 1 {
		return arr
	}
	var temp []string
	for _, element := range arr {
		if !StringInSlice(temp, element) {
			temp = append(temp, element)
		}
	}
	return temp
}

func RemoveSelectedString(arr []string, selector string) []string {
	var r []string
	for _, str := range arr {
		if str != selector {
			r = append(r, str)
		}
	}
	return r
}

func RemoveSelectedIntFromSlice(arr []int, selector int) []int {
	var res []int
	for _, element := range arr {
		if element != selector {
			res = append(res, element)
		}
	}
	return res
}

func RollingInSlice(arr []string) string {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
	return arr[r.Int64()]
}

func WeiToEther(wei *big.Int, decimals int) *big.Float {
	dec := fmt.Sprintf("1e%d", decimals)
	deci, _ := strconv.ParseFloat(dec, 64)
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(deci))
}

func BigIntToString(ints []*big.Int) string {
	var res []string
	for _, i := range ints {
		res = append(res, i.String())
	}
	return strings.Join(res, ",")
}

func RankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func IpCountryDetect(reader *geoip2.Reader, ipStr string, countryIsoCode string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		fmt.Printf("IpCountryDetect failed, ip %s parse fail\n", ipStr)
		return false
	}
	record, err := reader.Country(ip)
	if err != nil {
		fmt.Printf("IpCountryDetect failed, get country fail, reason: %s\n", err.Error())
		return false
	}
	return record.Country.IsoCode == countryIsoCode
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

//--------------------- math related--------------------------------------------

func AlmostEqual(a, b float64) bool {
	epsilon := 1e-9
	diff := math.Abs(a - b)
	return diff < epsilon
}

func IsZero(f float64) bool {
	return math.Float64bits(f) == 0
}

func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

func FormatPrice(price float64) string {
	str := strconv.FormatFloat(price, 'f', -1, 64) // Convert float to string

	parts := strings.Split(str, ".")
	if len(parts) < 2 {
		return str
	}

	countZero := 0
	for _, char := range parts[1] {
		if char == '0' {
			countZero++
		} else {
			break
		}
	}

	if countZero == 0 {
		return fmt.Sprintf("%.2f", price)
	}

	return parts[0] + "." + strings.Repeat("0", countZero) + string(parts[1][countZero])
}

// HelperPaginate returns total page number and indexes
// startIndex = -1 indicates error
func HelperPaginate(totalNum, pageNum, pageSize int) (totalPageNum, startIndex, endIndex int) {
	if pageSize < 1 {
		pageSize = 1
	}
	totalPageNum = totalNum / pageSize
	if totalNum%pageSize != 0 {
		totalPageNum++
	}

	startIndex = (pageNum - 1) * pageSize
	endIndex = startIndex + pageSize

	if startIndex >= totalNum {
		startIndex = -1
		endIndex = -1
	} else if endIndex > totalNum {
		endIndex = totalNum
	}
	return
}

func ReadHttpParamToInt64(param string, defaultValue, minValue int64) int64 {
	intValue, err := strconv.ParseInt(param, 10, 32)
	if err != nil || intValue < minValue {
		return defaultValue
	}
	return intValue
}

/*
func findCeil(sortedArr []int, target, low, high int) int {
	for low < high {
		mid := low + ((high - low) >> 1) // Same as mid = (low+high)/2
		if target > sortedArr[mid] {
			low = mid + 1
		} else {
			high = mid
		}
	}

	if sortedArr[low] >= target {
		return low
	} else {
		return -1
	}
}

func RandomGenerator(values, frequencies []int) int {
	n := len(values)

	// Create and fill prefix array
	prefix := make([]int, n)
	prefix[0] = frequencies[0]
	for i := 1; i < n; i++ {
		prefix[i] = prefix[i-1] + frequencies[i]
	}

	// prefix[n-1] is the sum of all frequencies.
	// Generate a random number with a value from 1 to this sum.
	randMath.Seed(time.Now().UnixNano())
	r := randMath.Intn(prefix[n-1]) + 1

	// Find the index of the ceiling of r in the prefix array
	selectedIndex := findCeil(prefix, r, 0, n-1)
	return values[selectedIndex]
}

*/

func SortMapKeys(hash map[int]int) []int {
	values := make([]int, 0)
	for val := range hash {
		values = append(values, val)
	}
	sort.Ints(values)
	return values
}

func GenerateWeightedChoice(values []int, frequencies []int) (int, error) {
	if len(values) != len(frequencies) {
		return 0, errors.New("lottery configuration error")
	}
	listChoice := make([]randutil.Choice, 0)
	for i := 0; i < len(values); i++ {
		listChoice = append(listChoice, randutil.Choice{
			Weight: frequencies[i],
			Item:   values[i],
		})
	}
	choice, err := randutil.WeightedChoice(listChoice)
	if err != nil {
		return 0, err
	}
	result, ok := choice.Item.(int)
	if !ok {
		return 0, errors.New("assertion failed: item is not an int")
	}
	return result, err
}

// DiceWithProbability 给定总数和概率，抽取结果
func DiceWithProbability(total int64, probability float64) (hit int64) {
	if total == 0 {
		return total
	}
	for i := int64(0); i < total; i++ {
		if mathRand.Float64() < probability {
			hit += 1
		}
	}
	return
}

func AverageAmount(num1, num2 string) (string, error) {
	num1Int, err1 := strconv.Atoi(num1)
	num2Int, err2 := strconv.Atoi(num2)
	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("conversion error")
	}
	return strconv.FormatFloat((float64(num1Int)+float64(num2Int))/2, 'f', -1, 64), nil
}

func FindStringWithMaxCount(counts map[string]int) (string, int) {
	fmt.Printf("%+v\n", counts)
	var maxString string
	var maxCount int

	for str, count := range counts {
		if count > maxCount {
			maxString = str
			maxCount = count
		}
	}

	return maxString, maxCount
}

func EscapeDoubleQuotes(input string) string {
	// 使用 strings.ReplaceAll 函数将双引号替换为带有转义符号的双引号
	escaped := strings.ReplaceAll(input, "\"", "\\\"")
	return escaped
}
