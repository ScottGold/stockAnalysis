package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	help    bool
	infile  string
	outfile string
)

var (
	mapStockCode2Name map[string]string
)

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.StringVar(&infile, "i", "", "input file")
	flag.StringVar(&outfile, "o", "", "output file")
}

func usage() {
	fmt.Fprintf(os.Stderr,
		`app version: app/1.10.0 Usage:  [-h] 
Options:
`)
	flag.PrintDefaults()
}

func readline(filename string, cline chan string) error {
	ifile, err := os.Open(filename)
	if err != nil {
		return err
	}

	pargesize := os.Getpagesize()
	fmt.Println("Getpagesize", pargesize)
	go func() {
		defer ifile.Close()
		line := ""
		bytes := make([]byte, pargesize)
		count, _ := ifile.Read(bytes)
		for count > 0 {
			//fmt.Println("read", count, err)
			lastlinestart := 0
			for i := 0; i < count; i++ {
				if bytes[i] == '\n' {
					line = line + string(bytes[lastlinestart:i])
					cline <- line
					line = ""
					lastlinestart = i + 1
				}
				if i == count-1 && i > lastlinestart {
					line = string(bytes[lastlinestart : i+1])
				}
			}
			count, _ = ifile.Read(bytes)
		}
		if line != "" {
			cline <- line
		}
		close(cline)
	}()

	return nil
}

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

func Headers() []string {
	h := []string{"证券名称", "买卖方向", "业务名称", "成交日期", "成交时间", "成交数量", "成交价格", "交易金额"}
	return h
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

//成交日期,成交时间,证券名称,证券代码,买卖方向,业务名称,市场,成交数量,成交价格,交易金额,成交编号
func parseData(row string) StockItem {
	//row = strings.Trim(row, "\" ")
	fields := strings.Split(row, ",")

	date := strings.TrimSpace(strings.Trim(fields[0], "\" "))
	time := strings.TrimSpace(strings.Trim(fields[1], "\" "))
	stockname := strings.TrimSpace(strings.Trim(fields[2], "\" "))
	stockcode := strings.TrimSpace(strings.Trim(fields[3], `" `))
	name, iscontain := mapStockCode2Name[stockcode]
	if iscontain {
		stockname = name
	} else {
		mapStockCode2Name[stockcode] = stockname
	}

	buyDirection := strings.TrimSpace(strings.Trim(fields[4], "\" "))
	businessName := strings.TrimSpace(strings.Trim(fields[5], "\" "))
	//市场
	countStr := strings.TrimSpace(strings.Trim(fields[7], "\" ")) //成交数量 strconv.ParseInt(, 10, 64)
	priceStr := strings.TrimSpace(strings.Trim(fields[8], "\" ")) //
	amountStr := strings.TrimSpace(strings.Trim(fields[9], "\" "))
	count, err := strconv.ParseFloat(countStr, 64)
	if err != nil {
		fmt.Println("ParseInt countStr", err, countStr)
	}
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Println("ParseInt priceStr", err, priceStr)
	}
	amount, err := strconv.ParseFloat(strings.Trim(amountStr, `"`), 64)
	if err != nil {
		fmt.Println("ParseInt amountStr", err, amountStr)
	}
	if count > 0 {
		amount = amount * -1
	}

	//fmt.Println(date, time, stockname, buyDirection, businessName, count, price, amount)
	return StockItem{stockname,
		buyDirection,
		businessName,
		date,
		time,
		count,
		price,
		amount}
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	fmt.Println("需要手动将输出文件转成 utf-8 bom 格式")

	mapStockCode2Name = make(map[string]string)
	chRLine := make(chan string, 100)
	err := readline(infile, chRLine)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	finished := make(chan bool)
	if outfile == "" {
		outfile = infile[:strings.LastIndex(infile, ".csv")] + strconv.FormatInt(time.Now().Unix(), 10) + ".csv"
		fmt.Println("output file", outfile)
	}
	ofile, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	defer ofile.Close()

	writer := csv.NewWriter(ofile)
	defer writer.Flush()
	writer.Write(Headers())
	go func() {
		ln := 0
		l, ok := <-chRLine
		for ; ok == true; l, ok = <-chRLine {
			//fmt.Println(l)
			ln++
			if ln == 1 {
				continue
			}
			data := parseData(l)
			if data.amount == 0 { //|| data.count == 0 || data.price == 0
				continue
			}
			if data.businessName == "新股申购确认缴款" ||
				data.businessName == "货币基金赎回" ||
				data.businessName == "货币基金申购" {
				continue
			}
			//fmt.Println(data.amount)
			writer.Write(data.CSVRecord())
		}
		finished <- true
	}()
	<-finished
}
