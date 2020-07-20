package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Outlier struct {
	Date                         string
	siteID                       float64
	Daily_Mean_PM_Concentrations float64
	Site_Latitude                string
	Site_Longitude               string
}
type Dates struct {
	Jan uint16
	Mar uint16
	Feb uint16
	Apr uint16
	May uint16
	Jun uint16
	Jul uint16
	Aug uint16
	Sep uint16
	Oct uint16
	Nov uint16
	Dec uint16
}

func (d Dates) totalCount() (x uint16) {
	x = d.Jan + d.Feb + d.Mar + d.Apr + d.May + d.Jun + d.Jul + d.Aug + d.Sep + d.Oct + d.Nov + d.Dec
	return
}

func (d Dates) common() {
	fmt.Printf("%+v", d)

}
func readOutliers(f string) []string {
	file, err := os.Open(f)
	defer file.Close()
	var data []string
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data = append(data, scanner.Text())
		}
	} else {
		fmt.Println(err)
	}
	return data
}
func mapData(data []string) (mappedData []Outlier) {
	for i, v := range data {
		if i != 0 {
			a := strings.ReplaceAll(v, "\"", "") //rid of quotes
			array := strings.Split(a, ",")
			siteId, _ := strconv.ParseFloat(array[1], 64)
			pm, _ := strconv.ParseFloat(array[2], 64)
			outLier := Outlier{array[0], siteId, pm, array[3], array[4]}
			mappedData = append(mappedData, outLier)
		}
	}

	return
}
func separatebyYears(mappedData []Outlier) (data map[string][]Outlier) {
	data = make(map[string][]Outlier)
	for _, v := range mappedData {
		year := strings.Split(v.Date, "/")[2]
		data[year] = append(data[year], v)
	}
	return
}

func commonDatesByYear(data map[string][]Outlier, d string) {
	var datemap = make(map[string]Dates)
	for i := 1999; i <= 2020; i++ {
		var array []string
		v := data[strconv.FormatInt(int64(i), 10)]
		for _, v2 := range v {
			array = append(array, v2.Date)
		}
		dates := commonDates(array)
		datemap[strconv.FormatInt(int64(i), 10)] = dates
		fmt.Printf("%v:%v ---- %+v\n", i, dates.totalCount(), dates)
	}
	writeOutlierYears(datemap, d)

}
func commonDates(data []string) Dates {
	var recurringDates = Dates{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	var dates []string
	for _, v := range data {
		a := strings.ReplaceAll(v, "\"", "")
		dates = append(dates, strings.Split(a, ",")[0])
	}
	for _, v := range dates {
		num, _ := strconv.ParseFloat(strings.Split(v, "/")[0], 64)
		switch num {
		case 0:
			recurringDates.Jan++
			break
		case 1:
			recurringDates.Feb++
			break
		case 2:
			recurringDates.Mar++
			break
		case 3:
			recurringDates.Apr++
			break
		case 4:
			recurringDates.May++
			break
		case 5:
			recurringDates.Jun++
			break
		case 6:
			recurringDates.Jul++
			break
		case 7:
			recurringDates.Aug++
			break
		case 8:
			recurringDates.Sep++
			break
		case 9:
			recurringDates.Oct++
			break
		case 10:
			recurringDates.Nov++
			break
		case 11:
			recurringDates.Dec++
			break
		}
	}
	return recurringDates
}

/*
Warning... Will overwrite Graph Data if called
*/
func writeOutlierYears(dates map[string]Dates, d string) {
	file, err := os.Create("C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\OutliersYears-" + d + ".csv")
	defer file.Close()
	if err == nil {
		file.WriteString("Date, OutlierCount\n")
		for k, v := range dates {
			file.WriteString(fmt.Sprintf("%v,%v\n", k, strconv.FormatInt(int64(v.totalCount()), 10)))
		}
		fmt.Println("fin writing csv")
	} else {
		fmt.Println(err)
	}
}
