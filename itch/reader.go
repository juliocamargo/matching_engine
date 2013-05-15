package itch

import (
	"bufio"
	"github.com/fmstephe/matching_engine/trade"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type ItchReader struct {
	lineCount uint
	maxBuy    int64
	minSell   int64
	r         *bufio.Reader
}

func NewItchReader(fName string) *ItchReader {
	f, err := os.Open(fName)
	if err != nil {
		panic(err.Error())
	}
	r := bufio.NewReader(f)
	// Clear column headers
	if _, err := r.ReadString('\n'); err != nil {
		panic(err.Error())
	}
	return &ItchReader{lineCount: 1, minSell: math.MaxInt32, r: r}
}

func (i *ItchReader) ReadOrderData() (o *trade.OrderData, line string, err error) {
	i.lineCount++
	for {
		line, err = i.r.ReadString('\n')
		if err != nil {
			return
		}
		if line != "" {
			break
		}
	}
	o, err = mkOrderData(line)
	if o != nil && o.Kind == trade.BUY && o.Price > i.maxBuy {
		i.maxBuy = o.Price
	}
	if o != nil && o.Kind == trade.SELL && o.Price < i.minSell {
		i.minSell = o.Price
	}
	return
}

func (i *ItchReader) ReadAll() (orders []*trade.OrderData, err error) {
	orders = make([]*trade.OrderData, 0)
	var o *trade.OrderData
	for err == nil {
		o, _, err = i.ReadOrderData()
		if o != nil {
			orders = append(orders, o)
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func (i *ItchReader) LineCount() uint {
	return i.lineCount
}

func (i *ItchReader) MaxBuy() int64 {
	return i.maxBuy
}

func (i *ItchReader) MinSell() int64 {
	return i.minSell
}

func mkOrderData(line string) (o *trade.OrderData, err error) {
	ss := strings.Split(line, " ")
	var useful []string
	for _, w := range ss {
		if w != "" && w != "\n" {
			useful = append(useful, w)
		}
	}
	cd, td, err := mkData(useful)
	if err != nil {
		return
	}
	switch useful[3] {
	case "B":
		o.WriteBuy(cd, td)
		return
	case "S":
		o.WriteSell(cd, td)
		return
	case "D":
		o.WriteCancel(td)
		return
	default:
		return
	}
	panic("Unreachable")
}

func mkData(useful []string) (cd trade.CostData, td trade.TradeData, err error) {
	//      print("ID: ", useful[2], " Type: ", useful[3], " Price: ",  useful[4], " Amount: ", useful[5])
	//      println()
	var price, amount, traderId, tradeId int
	amount, err = strconv.Atoi(useful[4])
	price, err = strconv.Atoi(useful[5])
	traderId, err = strconv.Atoi(useful[2])
	tradeId, err = strconv.Atoi(useful[2])
	if err != nil {
		return
	}
	cd = trade.CostData{Price: int64(price), Amount: uint32(amount)}
	td = trade.TradeData{TraderId: uint32(traderId), TradeId: uint32(tradeId), StockId: uint32(1)}
	return
}
