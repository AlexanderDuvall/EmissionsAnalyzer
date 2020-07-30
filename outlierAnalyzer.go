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
	Feb uint16
	Mar uint16
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
	//fmt.Printf("%+v", d)
}

/** reads outlier file and returns each line as a string
 */
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

/**
Organizes the Outlier Data
*/
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

/**
Separates the Outliers by years in a map
*/
func separatebyYears(mappedData []Outlier) (data map[string][]Outlier) {
	data = make(map[string][]Outlier)
	for _, v := range mappedData {
		year := strings.Split(v.Date, "/")[2]
		data[year] = append(data[year], v)
	}
	return
}

/**
Counts Number of Outliers by yearand writes them to a file
*/
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
		//fmt.Printf("%v:%v ---- %+v\n", i, dates.totalCount(), dates)
	}
	writeOutlierYears(datemap, d)

}

/**
returns outlier count by Month for a given year.
*/
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
		case 1:
			recurringDates.Jan++
			break
		case 2:
			recurringDates.Feb++
			break
		case 3:
			recurringDates.Mar++
			break
		case 4:
			recurringDates.Apr++
			break
		case 5:
			recurringDates.May++
			break
		case 6:
			recurringDates.Jun++
			break
		case 7:
			recurringDates.Jul++
			break
		case 8:
			recurringDates.Aug++
			break
		case 9:
			recurringDates.Sep++
			break
		case 10:
			recurringDates.Oct++
			break
		case 11:
			recurringDates.Nov++
			break
		case 12:
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
	file, err := os.Create("C:\\Users\\Alex\\Documents\\Summer 2020 Work\\California_2019\\" + pollutant + "OutliersYears-" + d + ".csv")
	defer file.Close()
	if err == nil {
		file.WriteString("Date, OutlierCount\n")
		for k, v := range dates {
			file.WriteString(fmt.Sprintf("%v, %v\n", k, strconv.FormatInt(int64(v.totalCount()), 10)))
		}
		fmt.Println("fin writing csv")
	} else {
		fmt.Println(err)
	}
}

/**
Will overwrite data if called
*/
func outlierData(ending string) {
	s := readOutliers("C:\\Users\\Alex\\Documents\\Summer 2020 Work\\California_2019\\" + pollutant + "Outliers" + ending + ".csv")
	outliers := mapData(s)
	mappedData := separatebyYears(outliers)
	commonDatesByYear(mappedData, ending)
}

/**
Counts Outliers by sensor and writes it to a file
*/
func outliersBySensor(mappedData []Outlier, s string) {
	sensorData := make(map[float64]int)
	for _, v := range mappedData {
		sensorData[v.siteID]++
	}
	fmt.Println("Number of Outliers each sensor has.")
	file, err := os.Create("C:\\Users\\Alex\\Documents\\Summer 2020 Work\\California_2019\\" + pollutant + "OutliersBySensor" + s + ".csv")
	file.WriteString("SiteID, Outlier Count\n")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		for k, v := range sensorData {
			s := fmt.Sprintf("%f,%v\n", k, v)
			//	fmt.Printf("%f,%v\n", k, v)
			file.WriteString(s)
		}
	}
}
func OutliersBySensorDetailed(mappedData []Outlier, s string) {
	sensorData := make(map[string]int)
	for _, v := range mappedData {
		a := fmt.Sprintf("%s:%f", strings.Split(v.Date, "/")[2], v.siteID)
		sensorData[a]++
	}
	fmt.Println("Number of Outliers each sensor has by Year.")
	file, err := os.Create("C:\\Users\\Alex\\Documents\\Summer 2020 Work\\California_2019\\" + pollutant + "OutliersBySensorDetailed" + s + ".csv")
	defer file.Close()
	file.WriteString("Date-SiteId, Outlier Count\n")
	if err != nil {
		fmt.Println(err)
	} else {
		for k, v := range sensorData {
			s := fmt.Sprintf("%s,%v\n", k, v)
			file.WriteString(s)
		}
	}
}

/**
Function to get outliers by sensor
isDetailed == true will give date+siteid and create separate file
false will give just Site Id
*/
func OutliersBySensorFile(d string, isDetailed bool) {
	s := readOutliers("C:\\Users\\Alex\\Documents\\Summer 2020 Work\\California_2019\\" + pollutant + "Outliers" + d + ".csv")
	outliers := mapData(s)
	if isDetailed {
		OutliersBySensorDetailed(outliers, "General")
	} else {
		outliersBySensor(outliers, "General")
	}
}
