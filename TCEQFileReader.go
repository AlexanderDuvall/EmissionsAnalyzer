package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const TCEQDirectoru string = "C:\\Users\\Alex\\Documents\\Summer 2020 Work\\482010307_Manchester Central"

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

func (t detailedPollutant) checkConsistency(isDay bool, array []detailedPollutant) {
	byDate := make(map[string][]detailedPollutant) //How many hours that day has charted
	po := ""
	for _, v := range array {
		po = v.Pollutant
		byDate[v.Date] = append(byDate[v.Date], v)
	}
	if _, err := os.Stat(TCEQDirectoru + "\\Details_By_Pollutant\\Complete"); os.IsNotExist(err) {
		err := os.Mkdir(TCEQDirectoru+"\\Details_By_Pollutant\\Complete", os.ModeDir)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	file, err := os.Create(TCEQDirectoru + "\\Details_By_Pollutant\\Complete\\" + strings.ReplaceAll(strings.Split(po, "<")[0], "/", "") + ".csv")
	defer file.Close()
	file.WriteString("Date, Start Hour, Value, Pollutant\n")
	if err != nil {
		fmt.Println(err)
	}
	if !isDay {
		for k, v := range byDate {
			for _, v2 := range v {
				if len(v) == 24 {
					s := fmt.Sprintf("%v,%v,%v,%v\n", k, v2.Start_Hour, v2.Value, v2.Pollutant)
					file.WriteString(s)
				}
			}
		}
	} else {
		for k, v := range byDate {
			for _, v2 := range v {
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
type TCEQFuncs interface {
	organize()
	TCEQConsistencyList()
}

//tin tsp
func organizeSemiColon() {

	file := TCEQDirectoru + "\\Total_1Hour.txt"
	ender := "1_Hour_"
	list := readFile(file)
	listOfData := make(map[string][]TCEQOrganizer)
	var organizer []string
	for i, v := range list {
		if i < 10 {
			continue
		} else if i == 10 {
			a := strings.Split(v, "SOC;")
			b := strings.Split(a[1], ";")
			organizer = b
			organizer[0] = strings.ReplaceAll(organizer[0], ";\"", "")
			for _, v2 := range organizer {
				listOfData[v2] = []TCEQOrganizer{}
			}
		} else {
			a := strings.Split(v, ";")
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
		fmt.Println(name, "DIO")
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
		//time.Sleep(40 * time.Millisecond)
	}
	println("fin")

}

/**
@DEPRECATED -- Not Feasible for TCEQ DATA "," delimitation. Use semicolons instead OrganizeSemiColon
Organize data from TCEQ by comma delimitation.
Separates files in Detail_By_Pollution Directory
*/
func organize() {
	file := TCEQDirectoru + "\\Total_24Hour.txt"
	ender := "24_Hour_"
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
		fmt.Println(name, "DIO")
		if strings.Contains(splitter[0], "Rubidium") {
			continue
		}
		f, err := os.Create(TCEQDirectoru + "\\" + "Details_By_Pollutant\\" + ender + strings.ReplaceAll(name, ";", "") + ".csv")
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
		//time.Sleep(40 * time.Millisecond)
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
			if v.IsDir() == false {

				log.Println("hello->", v.Name())

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
							Pollutant:  strings.Replace(v.Name(), ".csv", "", 1),
							Units:      a[7],
						}
						mappedData[v.Name()] = append(mappedData[v.Name()], p)
					}
				}
			}
		}
	}
	return mappedData
}

func TCEQConsistencyList() {
	fmt.Println("--.................................")
	pollutants := getData(TCEQDirectoru)
	for k, v := range pollutants {
		fmt.Println(k, "DIO")
		isDay := false
		if strings.Contains(k, "24_Hour") {
			isDay = true
			fmt.Println("gottem")
		}
		if len(v) > 0 {

			fmt.Println(v[0].Pollutant, "POLLY", len(v))
			v[0].checkConsistency(isDay, v)
		}
	}
	d := TCEQDirectoru + "\\Details_By_Pollutant\\Complete"
	errx := filepath.Walk(TCEQDirectoru+"\\Details_By_Pollutant\\Complete", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false {
			fmt.Println(path, "HAHA")
			s := readFile(path)
			if len(s) <= 1 {
				err := os.Remove(path)
				if err != nil {
					log.Fatal(err)
				} else {
					log.Println("purged")
				}
			}
			return err
		} else if d != path {
			return filepath.SkipDir
		}
		return err
	})
	if errx != nil {
		log.Fatal(errx.Error())
	}
	fmt.Println("flash")
}
func OrganizeCompleteData(isDay bool, data []string) (averagePollution []Pollutant) {
	counter := 0
	sum := 0.0
	reset := func() {
		counter = 0
		sum = 0
	}
	for i, v := range data {
		if i == 0 {
			continue
		}
		list := strings.Split(v, ",")
		name := list[3]
		if len(list) >= 4 {
			strings.Join(list[3:], ",")
		}
		p := Pollutant{
			pollutant: name,
			date:      list[0],
			value:     list[2],
			unit:      "",
		}
		s, err := strconv.ParseFloat(p.value, 64)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			sum += s
		}
		counter++
		if isDay {
			poll := Pollutant{
				pollutant: list[3],
				date:      list[0],
				value:     strconv.FormatFloat(sum, 'f', -1, 64),
				unit:      "",
			}
			averagePollution = append(averagePollution, poll)
		} else if counter == 24 {
			sum /= 24
			poll := Pollutant{
				pollutant: list[3],
				date:      list[0],
				value:     strconv.FormatFloat(sum, 'f', -1, 64),
				unit:      "",
			}
			averagePollution = append(averagePollution, poll)
			reset()
		}
	}
	return
}
func normalizeName(s string) string {
	splitter := strings.Split(s, "<")
	splitter[0] = strings.ReplaceAll(splitter[0], "/", "#")
	var name string = splitter[0]
	if len(name) > 20 {
		name = name[:21]
	}
	return name
}
func CycleThroughCompleteFiles() {
	d := TCEQDirectoru + "\\Details_By_Pollutant\\Complete\\"
	err := filepath.Walk(TCEQDirectoru+"\\Details_By_Pollutant\\Complete\\", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false {
			fmt.Println(path)
			fmt.Println(info.Name())
			fmt.Println("---------------------------------")
			// proceed bc not directory

			var is24 = strings.Contains(info.Name(), "24_Hour")
			if err != nil {
				log.Fatal(err.Error())
			} else {
				lines := readFile(path)
				p := OrganizeCompleteData(is24, lines) // get daily averages of pollutant
				if _, err := os.Stat(d + "\\DailyAverages"); os.IsNotExist(err) {
					os.Mkdir(d+"\\DailyAverages", os.ModeDir)
				}
				file, err := os.Create(filepath.Join(d+"\\DailyAverages", strings.ReplaceAll(info.Name(), ".csv", "")+".csv"))
				if err != nil {
					log.Fatal(err.Error())
				}
				file.WriteString("Date, Daily Avg\n")
				for _, v := range p {
					file.WriteString(fmt.Sprintf("%v,%v\n", v.date, v.value))
				}
				file.Close()
			}
		} else if d != path {
			return filepath.SkipDir
		}
		return err
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
