package fileinteraction

import (
	"encoding/csv"
	"fmt"
	"os"
)

var csvHeader = []string{"Site URL", "Cname", "Message", "DM Link", "CCID"}

var csvWriter *csv.Writer

func CheckFileExists(filename string) bool {
	_, err := os.Stat(filename)
	
	return err == nil 
}

func CreateFile(filename string) (*os.File){

	createdfile, err := os.Create(filename)

	_ = os.Chmod(filename, 0777)

	if err != nil {
		fmt.Printf("Error Creating file [ %v ]. Error: %+v ", filename, err.Error())
		os.Exit(1)
	}

	return createdfile
}

func CreateTmpFile(value string) string{
	tmp, err := os.CreateTemp("", "tmp-id-rsa")

	if err != nil {
		panic(err.Error())
	}

	WriteCredentials(value, tmp)


	fmt.Println(tmp.Name())
	return tmp.Name()
}

func ReadFile(filename string) string {
	data, err := os.ReadFile(filename)
	
	if err != nil {
		fmt.Printf("Error Reading file [ %v ]. Error: %+v ", filename, err.Error())
		os.Exit(1)
	}

	return string(data)
}

func WriteCredentials(creds string, file *os.File) {
	defer file.Close()

	file.Write([]byte(creds))
}

func CreateCsvWriter(file *os.File) {
	csvWriter = csv.NewWriter(file)
}

func WriteCsvReportHeader() {
	 err := csvWriter.Write(csvHeader)

	 if err != nil {

		fmt.Println(err.Error())
	 }
}

func WriteToCsv(row []string) {
	err := csvWriter.Write(row)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func FlushCsvWriter() {
	csvWriter.Flush()
}

func GetHomeDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Can't Determine Home Directory.")
		return "", err
	}

	return homeDir, nil

}