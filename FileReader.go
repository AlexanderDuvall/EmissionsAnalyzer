package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

var outlierList []dataMap
var siteLocation []string //latitude,longitude

/**
How data is grouped. May vary from year to year.
*/
type dataMap struct {
	date                         string
	source                       string
	siteId                       float64
	POC                          string
	Daily_Mean_PM_Concentrations float64
	UNITS                        string
	DAILY_AQI_VALUE              float64
	Site_Name                    string
	DAILY_OBS_COUNT              string
	PERCENT_COMPLETE             string
	AQS_PARAMETER_CODE           string
	AQS_PARAMETER_DESC           string
	CBSA_CODE                    string
	CBSA_NAME                    string
	STATE_CODE                   string
	STATE                        string
	COUNTY_CODE                  string
	COUNTY                       string
	SITE_LATITUDE                string
	SITE_LONGITUDE               string
}

/**
Will return an array of the file per line.
*/
func readFile(f string) (data []string) {
	file, err := os.Open(f)
	defer file.Close()
	if err == nil {
		//fmt.Printf("Reading File %s", file.Name())
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data = append(data, scanner.Text())
		}
		//fmt.Println(strings.Join(data, "\n"))
	} else {
		//fmt.Println(err)
	}
	return
}

/*
Will return a map of sensor values. The key is the siteId of the Sensor
*/
func separateData(data []string, discriminator string) map[float64][]dataMap {
	var array = make(map[float64][]dataMap)
	for i := 0; i < len(data); i++ {
		ar := strings.Split(data[i], "\",\"")
		ar[2] = strings.TrimSpace(ar[2])
		ar[2] = strings.ReplaceAll(ar[2], "\"", "")
		ar[4] = strings.TrimSpace(ar[4])
		ar[4] = strings.ReplaceAll(ar[4], "\"", "")
		ar[6] = strings.TrimSpace(ar[6])
		ar[6] = strings.ReplaceAll(ar[6], "\"", "")
		ar[4] = strings.TrimSpace(ar[4])
		ar[6] = strings.TrimSpace(ar[6])
		i2, err := strconv.ParseFloat(ar[2], 64)
		i4, err2 := strconv.ParseFloat(ar[4], 64)
		i6, err3 := strconv.ParseFloat(ar[6], 64)
		if err != nil || err2 != nil || err3 != nil {
			//fmt.Println(err)
			//fmt.Println(err2)
			//fmt.Println(err3)
		} else {
			if discriminator == "-1" || strings.ReplaceAll(ar[17], "\"", "") == discriminator {
				var a = dataMap{
					ar[0],
					ar[1],
					i2,
					ar[3],
					i4,
					ar[5],
					i6,
					ar[7],
					ar[8],
					ar[9],
					ar[10],
					ar[11],
					ar[12],
					ar[13],
					ar[14],
					ar[15],
					ar[16],
					ar[17],
					ar[18],
					ar[19]}
				array[a.siteId] = append(array[a.siteId], a)
			}
		}
	}
	//fmt.Println(array)

	return array
}

/*
returns IQR of a given sensors Data.
*/
func getMedian(list []float64) (med1, med3 float64) {
	c := len(list)
	if c%2 == 0 && c >= 7 {
		q1End := c/2 - 1
		q3Start := c/2 + 1
		slice := list[:q1End]
		slicelen := len(slice)
		if slicelen%2 == 0 && c != 0 {
			med1 = float64(slice[slicelen/2]+slice[slicelen/2-1]) / 2
			slice = list[q3Start:]
			med3 = float64(slice[slicelen/2]+slice[slicelen/2-1]) / 2
		} else if c == 0 {
			fmt.Println("ZERO")
		} else {
			med1 = slice[int(math.Ceil(float64(slicelen/2)))]
			slice = list[q3Start:]
			med3 = slice[int(math.Ceil(float64(slicelen/2)))]
		}
	} else if c < 7 {
		//fmt.Println("LESS THAN 7")
	} else {
		q1End := int(math.Floor(float64(c / 2)))
		q3Start := int(math.Floor(float64(c/2)) + 1)
		slice := list[:q1End]
		slicelen := len(slice)
		if slicelen%2 == 0 {
			med1 = float64(slice[slicelen/2]+slice[slicelen/2-1]) / 2
			slice = list[q3Start:]
			med3 = float64(slice[slicelen/2]+slice[slicelen/2-1]) / 2
		} else {
			med1 = slice[int(math.Ceil(float64(slicelen/2)))]
			slice = list[q3Start:]
			med3 = slice[int(math.Ceil(float64(slicelen/2)))]
		}
	}
	return
}
func IQR(data []dataMap) (med1, med3, medQ float64) {
	var DailyMeanPM []float64
	for _, v := range data {
		DailyMeanPM = append(DailyMeanPM, v.Daily_Mean_PM_Concentrations)
	}
	sort.Slice(DailyMeanPM, func(i, j int) bool {
		return DailyMeanPM[i] < DailyMeanPM[j]
	})
	med1, med3 = getMedian(DailyMeanPM)
	medQ = med3 - med1
	return
}

/**
Find outliers of a sensor according to the IQR
*/
func outliers(q1 float64, q3 float64, medq float64, data []dataMap, totOutlier *int) {
	var outlier []dataMap
	var upper float64 = q3 + 1.5*medq
	var lesser float64 = q1 - 1.5*medq
	for _, v := range data {
		if v.Daily_Mean_PM_Concentrations > upper || v.Daily_Mean_PM_Concentrations < lesser {
			outlierList = append(outlierList, v)
		}
	}
	for _, _ = range outlier {
		*totOutlier++
	}
}

/**
sets Up initial locations (2020) to be compared
*/
func getLocations(data map[float64][]dataMap) {
	for _, v1 := range data {
		for _, v := range v1 {
			var s string
			s += strings.ReplaceAll(v.SITE_LATITUDE, "\"", "") + ","
			s += strings.ReplaceAll(v.SITE_LONGITUDE, "\"", "")
			if !findElement(siteLocation, s) { // element not found so appending
				siteLocation = append(siteLocation, s)
			}
		}
	}
	for _, v := range siteLocation {
		fmt.Println(v)
	}
}
func compareLocations(data map[float64][]dataMap) {
	var date string
	var uncommons []string
	for _, v1 := range data {
		if strings.Compare(date, "") == 0 {
			date = strings.ReplaceAll(v1[0].date, "\"", "")
		} //get date
		for _, v := range v1 {
			var s string
			s += strings.ReplaceAll(v.SITE_LATITUDE, "\"", "") + ","
			s += strings.ReplaceAll(v.SITE_LONGITUDE, "\"", "")
			if !findElement(siteLocation, s) && !findElement(uncommons, s) { // element not found so appending
				uncommons = append(uncommons, s)
			}
		}
	}
	fmt.Printf("Found %v uncommon sites for %s\n", len(uncommons), strings.Split(date, "/")[2])
}

/**
return False if no match, else true
*/
func findElement(array []string, element string) bool {
	b := false
	for _, v := range array {
		if v == element { //dont match
			b = true
		}
	}
	return b
}
func checkConsistency(index int, data map[float64][]dataMap) {
	if index == 0 {
		getLocations(data)

	} else {
		compareLocations(data)
	}
}

/**
Give a county to find. If just general information put "-1"
*/
func setUpOutliers(d string) {
	file := []string{
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2020.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2019.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2018.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2017.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2016.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2015.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2014.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2013.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2012.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2011.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2010.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2009.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2008.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2007.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2006.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2005.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2004.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2003.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2002.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2001.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2000.csv",
		"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_1999.csv"}
	function := func(v2 string) map[float64][]dataMap {
		var totOutliers int = 0
		fileLines := readFile(v2)
		mappedSensors := separateData(fileLines, d)
		for _, v := range mappedSensors {
			q1, q3, medq := IQR(v)
			outliers(q1, q3, medq, v, &totOutliers)
		}
		//fmt.Println(totOutliers)
		//totalOutliers += totOutliers
		//fmt.Printf("Total Outliers for %v: %v\n", date, totOutliers)
		return mappedSensors
	}
	for i, v := range file {
		var mappedSensors map[float64][]dataMap = function(v)
		checkConsistency(i, mappedSensors)
	}
	if d == "-1" {
		writeOutliers("")
		outlierData("")
	} else if d != "" {
		writeOutliers(d)
		outlierData(d)
	}
}
func writeOutliers(ending string) {
	f, err := os.Create("C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\Outliers" + ending + ".txt")
	defer f.Close()
	defer fmt.Println("Finished writing data")
	if err != nil {
		fmt.Println(err)
	} else {
		f.WriteString("Date, siteID, Daily_Mean_PM_Concentrations, Site_Latitude, Site_Longitude\n")
		for _, v := range outlierList {
			var s string = fmt.Sprintf("%v,%f,%fs,%v,%v\n", v.date, v.siteId, v.Daily_Mean_PM_Concentrations, v.SITE_LATITUDE, v.SITE_LONGITUDE)
			f.WriteString(s)
		}
	}
}

/**
Will overwrite data if called
*/
func outlierData(ending string) {
	s := readOutliers("C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\Outliers" + ending + ".txt")
	outliers := mapData(s)
	mappedData := separatebyYears(outliers)
	commonDatesByYear(mappedData, ending)
}
func outliersBySensor(mappedData []Outlier, s string) {
	sensorData := make(map[float64]int)
	for _, v := range mappedData {
		sensorData[v.siteID]++
	}
	fmt.Println("Number of Outliers each sensor has.")
	file, err := os.Create("C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\OutliersBySensor" + s + ".csv")
	file.WriteString("SiteID, Outlier Count\n")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		for k, v := range sensorData {
			s := fmt.Sprintf("%f,%v\n", k, v)
			fmt.Printf("%f,%v\n", k, v)
			file.WriteString(s)
		}
	}
}
func main() {
	//setUpOutliers("-1")
	s := readOutliers("C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\Outliers" + "Harris" + ".txt")
	outliers := mapData(s)
	outliersBySensor(outliers, "Harris")
}
