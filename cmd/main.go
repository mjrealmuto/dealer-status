package main

import (
	"dealer-status/internal/dbclient"
	"dealer-status/internal/sshclient"
	"dealer-status/pkg/fileinteraction"
	"fmt"

	"net"
	"os"
	"strings"
	"time"
)

type row struct {
	ID int
	site string
	CCID int
}

type siteStatus struct {
	url string
	cname string
	available bool
	message string
}

func main () {
	start := time.Now()
	fmt.Println("This process has started at: ", start)

	ex, _ := os.Executable()

	fmt.Println(ex);



	if _, exists := os.Stat(fmt.Sprintf("%v/assets/", ex)); os.IsNotExist(exists) {
		fmt.Println("file does not exist")
		os.Exit(1);
	}

	client := sshclient.TunnelIn(sshclient.SetCredentials())

	defer client.Close()

	db := dbclient.OpenDatabaseAndTestConnection(dbclient.ConnectToDatabase())	

	defer db.Close()

	rows, dmErr := db.Query(dbclient.GetWpQuery())

	handleError(dmErr)

	defer rows.Close()

	fmt.Println("Creating InActive Sites CSV. This file will be stored in the assets directory of this project.")

	csvFile := fileinteraction.CreateFile("/app/assets/InactiveSites.csv")

	defer csvFile.Close()

	fileinteraction.CreateCsvWriter(csvFile)
	
	fileinteraction.WriteCsvReportHeader()

	siteMessages := make(chan siteStatus)

	var dmLink string

	for rows.Next() {
		
		thisrow := new(row)
		
		rows.Scan(&thisrow.ID, &thisrow.site, &thisrow.CCID)
		
		go func(siteAddr string, c chan siteStatus) {

			cname, err := net.LookupCNAME(siteAddr)

			if err != nil {
				c <- siteStatus{
					siteAddr,
					cname,
					false,
					err.Error(),
				}
				return
			}

			if strings.Contains(cname, ".dealerinspire.com") {
				c <- siteStatus{
					siteAddr,
					cname,
					true,
					"Looks like our site",
				}
				return
			}
		
			c <- siteStatus{
				siteAddr,
				cname,
				false,
				"Looks like this is not our site",
			}
		}( normalizeSiteAddr(thisrow.site), siteMessages)
		
		status := <- siteMessages
		
		if !status.available {
				fmt.Println(status)
				dmLink = fmt.Sprintf("https://dev.dealerinspire.com/wp-admin/post.php?post=%v&action=edit", thisrow.ID)
				fileinteraction.WriteToCsv([]string{status.url, status.cname, status.message, dmLink, fmt.Sprintf("%v", thisrow.CCID)})
		}
		thisrow = nil
	}

	fileinteraction.FlushCsvWriter()
	duration := time.Since(start)

	fmt.Println("Process Took: ", duration)
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func normalizeSiteAddr(url string) string {

	/**
	seems that some urls are returning UPPERCASE
	**/
	url = strings.ToLower(url)

	/**
	Going to check here if the URL begins with https or http 
	with a trailing `www` then we will replace with just `www`
	**/
	if strings.HasPrefix(url, "https://") {
		url = strings.ReplaceAll(url, "https://", "")
	} else if strings.HasPrefix(url, "http://") {
		url = strings.ReplaceAll(url, "http://", "")
	}

	/**
	Here we're going to check if the URL starts with `www`
	If it doesn't we're going to add it to the beginning
	of the url
	**/
	if !strings.HasPrefix(url, "www") {
		url = "www." + url
	}

	if strings.HasSuffix(url, "/") {
		url = strings.Replace(url, "/", "", 1)
	}

	return url
}
