package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestUtils(t *testing.T) {
	arr := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for i := 0; i < 30; i++ {
		fmt.Println(RollingInSlice(arr))
	}
}

func TestPaginate(t *testing.T) {
	d := make([]string, 0)
	a, b, c := HelperPaginate(0, 1, 20)
	fmt.Printf("%d %d %d \n", a, b, c)
	fmt.Printf("%+v\n", d[b:c])

}

func TestSortMapKeys(t *testing.T) {
	hash := map[int]int{10: 1, 20: 2, 30: 3, 40: 4}
	resp := SortMapKeys(hash)
	fmt.Println(resp)
}

func TestRandomGenerator(t *testing.T) {
	points := []int{10, 20, 30, 40, 50, 100, 300}
	probs := []int{600, 200, 100, 60, 20, 15, 5}
	startTime := time.Now()
	resp, err := GenerateWeightedChoice(points, probs)
	elapsedTime := time.Since(startTime)
	fmt.Println("Elapsed time: ", elapsedTime)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	fmt.Println(resp)
	return
}

func TestUTCdate(t *testing.T) {
	t1 := time.Now().UTC()
	time.Sleep(1 * time.Second)
	t2 := time.Now().UTC()
	fmt.Printf("t1: %+v, t2: %+v\n", t1, t2)

	fmt.Println(IsSameUTCDay(t1.Unix(), t2.Unix()))
}

func TestDic(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(DiceWithProbability(100, float64(24)/float64(24)))
	}
}

func TestAverageAmount(t *testing.T) {
	fmt.Println(AverageAmount("30", "4000"))
}
