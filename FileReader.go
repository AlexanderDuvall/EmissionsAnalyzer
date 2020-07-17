package main

import (
	"bufio"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type dataMap struct {
	date                         string
	source                       string
	siteId                       int64
	POC                          string
	Daily_Mean_PM_Concentrations int64
	UNITS                        string
	DAILY_AQI_VALUE              int64
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

func readFile() (data []string) {
	file, err := os.Open("")
	if err != nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			_ = append(data, scanner.Text())
		}
		file.Close()
	}
	return
}
func separateData(data []string) map[int64][]dataMap {
	var array = make(map[int64][]dataMap)
	for i := 0; i < len(data); i++ {
		ar := strings.Split(data[i], ",")
		i2, _ := strconv.ParseInt(ar[2], 10, 64)
		i4, _ := strconv.ParseInt(ar[4], 10, 64)
		i6, _ := strconv.ParseInt(ar[6], 10, 64)
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
		_ = append(array[a.siteId], a)
	}
	return array
}
func IQR(data []dataMap) (med1, med3, medQ float64) {
	var DailyMeanPM []int64
	for _, v := range data {
		_ = append(DailyMeanPM, v.Daily_Mean_PM_Concentrations)
	}
	c := len(DailyMeanPM)
	sort.Slice(DailyMeanPM, func(i, j int) bool {
		return DailyMeanPM[i] < DailyMeanPM[j]
	})

	if c%2 == 0 {
		q1 := c/2 - 1
		q3 := c/2 + 2
		if q1%2 == 0 {
			med1 = (float64(DailyMeanPM[q1/2]) + float64(DailyMeanPM[q1/2+1])) / 2
			med3 = (float64(DailyMeanPM[q3/2]) + float64(DailyMeanPM[q3/2+1])) / 2
		} else {
			med1 = math.Round(float64(DailyMeanPM[q1/2]))
			med3 = math.Round(float64(DailyMeanPM[q3/2]))
		}
		medQ = med3 - med1
	} else {
		var median int = int(math.Round(float64(c / 2)))
		q1 := median - 1
		q3 := median + 1
		if q3%2 == 0 {
			med1 = (float64(DailyMeanPM[q1/2]) + float64(DailyMeanPM[q1/2+1])) / 2
			med3 = (float64(DailyMeanPM[q3/2]) + float64(DailyMeanPM[q3/2+1])) / 2
		} else {
			med1 = math.Round(float64(DailyMeanPM[q1/2]))
			med3 = math.Round(float64(DailyMeanPM[q3/2]))
		}
		medQ = med3 - med1

	}
	return
}
func outliers (data map[string][]dataMap) {

}
func main() {

}
