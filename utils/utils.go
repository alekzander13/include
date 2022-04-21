package utils

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ChkErrFatal(err error) {
	if err != nil {
		AddToLog(GetProgramPath()+"-error.txt", err)
		os.Exit(1)
	}
	/*
		f, err := os.OpenFile("errors_"+os.Args[0], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			os.Exit(1)
		}
		defer f.Close()

		fmt.Fprintln(f, v)
		os.Exit(1)
	*/
}

//Exists exist file or folder
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//AddToLog add to log file
func AddToLog(name string, info interface{}) {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return
	}
	defer f.Close()

	date := time.Now().Local().Format("02.01.06 15:04:05 ")

	fmt.Fprintln(f, date, info)
}

func GetPathWhereExe() string {
	p, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Dir(p)
}

//GetProgramPath return full path to file with file name
func GetProgramPath() string {
	path := os.Args[0]
	p, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	path = filepath.Dir(p)
	ext := filepath.Ext(filepath.Base(p))
	p = strings.TrimSuffix(filepath.Base(p), ext)
	return filepath.Join(path, p)
}

//Get2SymbolsinYear 1996 -> 96
func Get2SymbolsinYear(year string) string {
	return year[2:]
}

//MakePortsFromSlice {"1-5","9","11:13"} = ["1" "2" "3" "4" "5" "9" "11" "12" "13"]
func MakePortsFromSlice(ps []string) ([]string, error) {
	var res []string
	for _, p := range ps {
		if strings.Contains(p, "-") {
			r, err := makeSlicePort(strings.Split(p, "-"))
			if err != nil {
				return nil, err
			}
			res = append(res, r...)
		}
		if strings.Contains(p, ":") {
			r, err := makeSlicePort(strings.Split(p, ":"))
			if err != nil {
				return nil, err
			}
			res = append(res, r...)
		}
		port, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		res = append(res, strconv.Itoa(port))
	}
	return res, nil
}

func makeSlicePort(s []string) ([]string, error) {
	var res []string
	if len(s) == 2 {
		start, err := strconv.Atoi(s[0])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}
		if start > end {
			return nil, errors.New("bad slice " + strings.Join(s, "-"))
		}
		for i := start; i <= end; i++ {
			res = append(res, strconv.Itoa(i))
		}
		return res, nil
	}
	return nil, errors.New("bad slice")
}

func Formatdatetime(datetime time.Time, withtime bool) string {
	if withtime {
		return fmt.Sprintf("%02d.%02d.%02d %02d:%02d.%02d", datetime.Day(), datetime.Month(), datetime.Year(), datetime.Hour(), datetime.Minute(), datetime.Second())
	}
	return fmt.Sprintf("%02d.%02d.%02d", datetime.Day(), datetime.Month(), datetime.Year())
}

//Gpspathdate Формат для каталогов ЖПС данных 2019_12_21
func Gpspathdate(datetime time.Time) string {
	return fmt.Sprintf("%02d_%02d_%02d", datetime.Year(), datetime.Month(), datetime.Day())
}

//toFixedFloat округление до кол-ва после запятой
func toFixedFloat(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(roundFloat(num*output)) / output
}

func roundFloat(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ConvertCoordToFloat(str string) float64 {
	res, err := strconv.ParseFloat(str, 64)
	//res, err := strconv.ParseFloat(strings.ReplaceAll(str, ",", "."), 64)
	if err != nil {
		return -1
	}

	//return res
	return converCoord(res)
}

func converCoord(coord float64) float64 {
	/*
		temp := []rune(fmt.Sprintf("%f", coord/100))
		gr := string(temp[:2])
		min := string(temp[2:])
		minfl, err := strconv.ParseFloat(min, 10)
		if err != nil {
			return -1
		}
		grfl, err := strconv.ParseFloat(gr, 10)
		if err != nil {
			return -1
		}
		return grfl + (minfl * 100 / 60)
	*/
	coord /= 100
	gr := int(coord)
	min := coord - float64(gr)
	min = min * 100 / 60
	//return float64(gr) + min
	return toFixedFloat(float64(gr)+min, 7)
}
