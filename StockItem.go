package main

import (
	"fmt"
	"strconv"
)

type StockItem struct {
	stockname    string  //证券名称
	buyDirection string  //买卖方向
	businessName string  //业务名称
	date         string  //成交日期
	time         string  //成交时间
	count        float64 //成交数量
	price        float64 //成交价格
	amount       float64 //交易金额
}

func (this *StockItem) CSVString() string {
	return fmt.Sprintf(`"%s","%s","%s","%s","%s","%f","%f","%f"`, this.stockname,
		this.buyDirection,
		this.businessName,
		this.date,
		this.time,
		this.count,
		this.price,
		this.amount)
}

func (this *StockItem) CSVRecord() []string {
	record := make([]string, 0)
	record = append(record, this.stockname,
		this.buyDirection,
		this.businessName,
		this.date,
		this.time,
		strconv.FormatFloat(this.count, 'f', -1, 64),
		strconv.FormatFloat(this.price, 'f', -1, 64),
		strconv.FormatFloat(this.amount, 'f', -1, 64))
	return record
}

func (this *StockItem) String() string {
	return this.CSVString()
}
