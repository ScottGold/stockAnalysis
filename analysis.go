package main

import (
	"fmt"
	"sort"
)

type AnalysisResult struct {
	stockname     string
	totalCount    float64
	totalAmount   float64
	totalDividend float64
}

func ARHeaders() string {
	return `"证券名称","股票剩余数量","盈亏","股息"`
}

func (ar AnalysisResult) String() string {
	return fmt.Sprintf("%s,%f,%f,%f", ar.stockname, ar.totalCount, ar.totalAmount, ar.totalDividend)
}

type SortObj []AnalysisResult

func (a SortObj) Len() int      { return len(a) }
func (a SortObj) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortObj) Less(i, j int) bool {
	var av, bv float64
	if a[i].totalCount == 0 {
		av = 100000000
	}
	if a[j].totalCount == 0 {
		bv = 100000000
	}
	return a[i].totalAmount+av > a[j].totalAmount+bv
}

type BaseAnalysis struct {
	mapStockList map[string][]*StockItem
}

func (this *BaseAnalysis) Init() {
	this.mapStockList = make(map[string][]*StockItem)
}

func (this *BaseAnalysis) AddItem(item *StockItem) {
	if item.stockname == "" {
		//fmt.Println("AddItem stockname is empty", item)
		return
	}
	listvalue := this.mapStockList[item.stockname]
	listvalue = append(listvalue, item)
	this.mapStockList[item.stockname] = listvalue
}

//清仓盈亏
//这个历史数据对未来预测没什么用处
func (this *BaseAnalysis) Analysis() {
	results := make([]AnalysisResult, 0)
	for stockname, list := range this.mapStockList {
		rd := AnalysisResult{stockname: stockname}
		for _, d := range list {
			if d.businessName == "证券买入" || d.businessName == "证券卖出" || d.businessName == "新股申购确认缴款" {
				rd.totalCount += d.count
				rd.totalAmount += d.amount
			}
			if d.businessName == "股息红利税补缴" || d.businessName == "股息入账" {
				rd.totalDividend += d.amount
			}
		}
		if rd.totalCount != 0 || rd.totalAmount != 0 || rd.totalDividend != 0 {
			results = append(results, rd)
		}
	}
	sort.Sort(SortObj(results))
	fmt.Println(ARHeaders())
	for _, rd := range results {
		fmt.Println(rd)
	}
}
