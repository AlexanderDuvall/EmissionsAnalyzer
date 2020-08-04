package main

import (
	"fmt"
	"os"
	"strings"
)

//const TCEQDirectoru string = "C:\\Users\\Alex\\Documents\\Summer 2020 Work\\482011035 Clinton"
//const TCEQDirectoru string = "C:\\Users\\Alex\\Documents\\Summer 2020 Work\\482010057 Galena Park"
//const TCEQDirectoru string = "C:\\Users\\Alex\\Documents\\Summer 2020 Work\\482010307_Manchester Central"
const TCEQDirectoru string = "C:\\Users\\Alex\\Documents\\Summer 2020 Work\\482011049_Pasadena North"

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

func start() {
	file := TCEQDirectoru + "\\Totals_24HourPart1.txt"
	ender := "24_Hour"
	list := readFile(file)
	listOfData := make(map[string][]TCEQOrganizer)
	var organizer []string
	for i, v := range list {
		if i < 10 {
			continue
		} else if i == 10 {
			a := strings.Split(v, "SOC")
			a = strings.Split(a[1], "\",\"")
			if len(a) <= 1 {
				a = strings.Split(v, "\t")
			}
			organizer = a
			organizer[0] = strings.ReplaceAll(organizer[0], ",\"", "")
			for _, v2 := range organizer {
				listOfData[v2] = []TCEQOrganizer{}
			}
		} else {
			a := strings.Split(v, ",")
			if len(a) <= 1 {
				a = strings.Split(v, "\t")
			}

			if len(a) <= 1 {
				break
			}
			splice := a[18:]

			for i2, v2 := range splice {
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
				//fmt.Printf("%+v\n", data)
				listOfData[organizer[i2]] = append(listOfData[organizer[i2]], data)
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
		f.WriteString("Date,AQS_Code,Latitude,Longitude,Start_Hour,Value,Units\n")
		defer f.Close()

		if err != nil {
			fmt.Println(err)
		} else {
			for _, v2 := range v {
				s := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v\n", v2.Date, v2.AQS_Code, v2.Latitude, v2.Longitude, v2.Start_Hour, v2.Data.value, v2.Data.unit)
				f.WriteString(s)
			}
		}
	}

}
