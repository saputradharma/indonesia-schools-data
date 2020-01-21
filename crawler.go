package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type location struct {
	KodeWilayah               string `json:"kode_wilayah"`
	Nama                      string `json:"nama"`
	Level                     int    `json:"id_level_wilayah"`
	KodeWilayahIndukProvinsi  string `json:"kode_wilayah_induk_provinsi"`
	KodeWilayahIndukKabupaten string `json:"kode_wilayah_induk_kabupaten"`
}

// BaseLocationsAPIURL Base url format for getting Provinces, Cities, and Districts data.
const BaseLocationsAPIURL = "https://dapo.dikdasmen.kemdikbud.go.id/rekap/dataSekolah?id_level_wilayah=%d&kode_wilayah=%s"
const outputPath = "./crawler-data"

func syncProvinces() {
	locations := callHTTP(fmt.Sprintf(BaseLocationsAPIURL, 0, "000000"))

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, os.ModePerm)
	}

	outfile, err := os.Create(outputPath + "/provinces.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer outfile.Close()
	csvWriter := csv.NewWriter(outfile)

	// write header
	csvWriter.Write([]string{"kode_wilayah", "nama", "level"})
	for _, loc := range locations {
		csvWriter.Write([]string{strings.TrimSpace(loc.KodeWilayah), loc.Nama, strconv.Itoa(loc.Level)})
	}

	csvWriter.Flush()

	fmt.Printf("Provinces data saved to %s/%s\n", outputPath, "provinces.csv")
}

func callHTTP(url string) []location {
	locations := []location{}

	client := http.Client{
		Timeout: time.Second * 2,
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(responseBody, &locations)
	if err != nil {
		log.Fatal(err)
	}

	return locations
}
