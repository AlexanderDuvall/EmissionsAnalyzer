package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const TCEQDirectoru string = "C:\\Users\\Alex\\Documents\\Summer 2020 Work\\482011035 Clinton"

//const TCEQDirectoru string = "C:\\Users\\Alex\\Documents\\Summer 2020 Work\\482010057 Galena Park"

//const TCEQDirectoru string = "C:\\Users\\Alex\\Documents\\Summer 2020 Work\\482010307_Manchester Central"
//const TCEQDirectoru string = "C:\\Users\\Alex\\Documents\\Summer 2020 Work\\482011049_Pasadena North"
type detailedPollutant struct {
	Date       string
	AQS_Code   string
	Latitude   string
	Longitude  string
	Start_Hour string
	Value      string
	Pollutant  string
	Units      string
}

func (t detailedPollutant) checkConsistency(array []detailedPollutant) {

	byDate := make(map[string][]detailedPollutant) //How many hours that day has charted
	po := ""
	os.Stat("")
	for _, v := range array {
		po = v.Pollutant
		byDate[v.Date] = append(byDate[v.Date], v)
	}
	if _, err := os.Stat(TCEQDirectoru + "\\Details_By_Pollutant\\Complete"); os.IsNotExist(err) {
		err := os.Mkdir(TCEQDirectoru+"\\Details_By_Pollutant\\Complete", os.ModeDir)
		if err != nil {
			println(err)
		}
	}
	file, err := os.Create(TCEQDirectoru + "\\Details_By_Pollutant\\Complete\\" + strings.ReplaceAll(strings.Split(po, "<")[0], "/", "") + ".csv")
	defer file.Close()
	file.WriteString("Date, Start Hour, Value, Pollutant\n")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range byDate {
		for _, v2 := range v {
			if len(v) == 24 {
				s := fmt.Sprintf("%v,%v,%v,%v\n", k, v2.Start_Hour, v2.Value, v2.Pollutant)
				file.WriteString(s)
			}
		}
	}
}

type TCEQOrganizer struct {
	State        string
	Region       string
	County       string
	City         string
	AQS_Code     string
	Site_Name    string
	Latitude     string
	Longitude    string
	Year         string
	Month        string
	Day          string
	Date         string
	Start_Hour   string
	Start_Minute string
	Start_Time   string
	Duration     string
	Sampler_Type string
	SOC          string
	Data         Pollutant
}
type Pollutant struct {
	pollutant string
	date      string
	value     string
	unit      string
}

/**
Organize data from TCEQ by comma delimitation.
Separates files in Detail_By_Pollution Directory
*/
func organize() {
	file := TCEQDirectoru + "\\Total_24Hour.txt"
	ender := "24_Hour"
	list := readFile(file)
	listOfData := make(map[string][]TCEQOrganizer)
	var organizer []string
	for i, v := range list {
		if i < 10 {
			continue
		} else if i == 10 {
			a := strings.Split(v, "SOC")
			b := strings.Split(a[1], "\",\"")
			organizer = b
			organizer[0] = strings.ReplaceAll(organizer[0], ",\"", "")
			for _, v2 := range organizer {
				listOfData[v2] = []TCEQOrganizer{}
			}
		} else {
			a := strings.Split(v, ",")
			if len(a) <= 1 {
				break
			}
			splice := a[18:]
			for i2, v2 := range splice {
				if v2 != "?" {
					data := TCEQOrganizer{
						State:        a[0],
						Region:       a[1],
						County:       a[2],
						City:         a[3],
						AQS_Code:     a[4],
						Site_Name:    a[5],
						Latitude:     a[6],
						Longitude:    a[7],
						Year:         a[8],
						Month:        a[9],
						Day:          a[10],
						Date:         a[11],
						Start_Hour:   a[12],
						Start_Minute: a[13],
						Start_Time:   a[14],
						Duration:     a[15],
						Sampler_Type: a[16],
						SOC:          a[17],
						Data: Pollutant{
							pollutant: organizer[i2],
							date:      a[11],
							value:     v2,
							unit:      "",
						},
					}
					key := organizer[i2]
					listOfData[key] = append(listOfData[key], data)
				}
			}
		}
	}
	if _, ar := os.Stat(TCEQDirectoru + "\\" + "Details_By_Pollutant"); os.IsNotExist(ar) {
		os.Mkdir(TCEQDirectoru+"\\"+"Details_By_Pollutant", os.ModeDir)
		fmt.Println("Making Directory")
	}

	for k, v := range listOfData {
		splitter := strings.Split(k, "<")
		splitter[0] = strings.ReplaceAll(splitter[0], "/", "#")
		var name string = splitter[0]
		if len(name) > 20 {
			name = name[:21]
		}
		fmt.Println(name)
		if strings.Contains(splitter[0], "Rubidium") {
			continue
		}
		f, err := os.Create(TCEQDirectoru + "\\" + "Details_By_Pollutant\\" + ender + name + ".csv")
		_, err2 := f.WriteString("Date,AQS_Code,Latitude,Longitude,Start_Hour,Value,Pollutant,Units\n")

		if err != nil || err2 != nil {
			fmt.Println(err)
			fmt.Println(err2)
		} else {
			//20200101
			for _, v2 := range v {
				year := v2.Date[:4]
				month := v2.Date[4:6]
				day := v2.Date[6:]
				date := fmt.Sprintf("%v-%v-%v", month, day, year)
				s := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v\n", date, v2.AQS_Code, v2.Latitude, v2.Longitude, v2.Start_Hour, v2.Data.value, v2.Data.pollutant, v2.Data.unit)
				f.WriteString(s)
			}
		}
		if err3 := f.Close(); err3 != nil {
			fmt.Println(err3)
			os.Exit(2)
		}
		time.Sleep(40 * time.Millisecond)
	}
	println("fin")

}

func getData(dir string) map[string][]detailedPollutant {
	dir = dir + "\\" + "Details_By_Pollutant"
	fileInfo, err := ioutil.ReadDir(dir)
	mappedData := make(map[string][]detailedPollutant)
	if err != nil {
		println(err)
	} else {
		for _, v := range fileInfo {
			file, err2 := os.Open(dir + "\\" + v.Name())
			mappedData[v.Name()] = []detailedPollutant{}
			if err2 != nil {
				println(err2)
			} else {
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					asdf := scanner.Text()
					a := strings.Split(asdf, ",")
					if a[0] == "Date" {
						continue
					}
					p := detailedPollutant{
						Date:       a[0],
						AQS_Code:   a[1],
						Latitude:   a[2],
						Longitude:  a[3],
						Start_Hour: a[4],
						Value:      a[5],
						Pollutant:  a[6],
						Units:      a[7],
					}
					mappedData[v.Name()] = append(mappedData[v.Name()], p)
				}
			}
		}
	}
	return mappedData
}

func TCEQConsistencyList() {
	pollutants := getData(TCEQDirectoru)
	for _, v := range pollutants {
		if len(v) > 0 {
			fmt.Println(v[0].Pollutant)
			v[0].checkConsistency(v)
		}
	}
}
