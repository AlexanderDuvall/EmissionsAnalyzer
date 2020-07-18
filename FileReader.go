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

var totalOutliers int32 = 0

/**
Will return an array of the file per line.
*/
func readFile(f string) (data []string) {
	file, err := os.Open(f)
	defer file.Close()
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			data = append(data, scanner.Text())
		}
		//fmt.Println(strings.Join(data, "\n"))
	} else {
		fmt.Println(err)
	}
	return
}

/*
Will return a map of sensor values. The key is the siteId of the Sensor
*/
func separateData(data []string) map[float64][]dataMap {
	var array = make(map[float64][]dataMap)
	for i := 0; i < len(data); i++ {
		ar := strings.Split(data[i], ",")
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
			fmt.Println(err)
			fmt.Println(err2)
			fmt.Println(err3)
		} else {
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
	//fmt.Println(array)

	return array
}

/*
returns IQR of a given sensors Data.
*/
func getMedian(list []float64) (med1, med3 float64) {
	c := len(list)
	if c%2 == 0 {
		q1End := c/2 - 1
		q3Start := c/2 + 1
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
func outliers(q1 float64, q3 float64, medq float64, data []dataMap) {
	var outlier []dataMap
	var upper float64 = q3 + 1.5*medq
	var lesser float64 = q1 - 1.5*medq
	for _, v := range data {
		if v.Daily_Mean_PM_Concentrations > upper || v.Daily_Mean_PM_Concentrations < lesser {
			outlier = append(outlier, v)
		}
	}
	fmt.Println("outliers---------------------------------------------------------")
	for _, v := range outlier {
		fmt.Printf("date: %s, siteId:%f\n", v.date, v.siteId)
		totalOutliers++
	}

}

func main() {
	file := []string{"C:\\Users\\Alex\\Documents\\Summer 2020 Work\\PM2.5\\pm2.5_2020.csv"}
	function := func(v2 string) {
		fileLines := readFile(v2)
		mappedSensors := separateData(fileLines)
		for _, v := range mappedSensors {
			q1, q3, medq := IQR(v)
			outliers(q1, q3, medq, v)
		}

	}

	for _, v := range file {
		fmt.Println("scanning")
		function(v)
	}

	fmt.Printf("Total Outliers: %v", totalOutliers)
	//med1, med3 := getMedian([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11,12,13,14,15})
	//fmt.Println(med1, med3)
}
