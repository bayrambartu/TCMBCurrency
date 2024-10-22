package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CurrencyDay struct {
	ID         string
	Date       time.Time
	DayNo      string
	Currencies []Currency
}

type Currency struct {
	Code           string
	CrossOrder     int
	Unit           int
	CurrencyNameTR string
	CurrencyName   string
	ForexBuying    float64
	ForexSelling   float64
	CrossRateUSD   float64
	CrossRateOther float64
}

type tarih_Date struct {
	XML       xml.Name `xml:"Tarih_Date"`
	Tarih     string   `xml:"Tarih,attr"`
	Date      string   `xml:"Date,attr"`
	Bulten_No string   `xml:"Bulten_No,attr"`
	Currency  []currency
}
type currency struct {
	Kod             string `xml:"Kod,attr"`
	CrossOrder      string `xml:"CrossOrder,attr"`
	CurrencyCode    string `xml:"CurrencyCode,attr"`
	Unit            string `xml:"Unit"`
	Isim            string `xml:"Isim"`
	CurrencyName    string `xml:"CurrencyName"`
	ForexBuying     string `xml:"ForexBuying"`
	ForexSelling    string `xml:"ForexSelling"`
	BanknoteBuying  string `xml:"BanknoteBuying"`
	BanknoteSelling string `xml:"BanknoteSelling"`

	CrossRateUSD   string `xml:"CrossRateUSD"`
	CrossRateOther string `xml:"CrossRateOther"`
}

func (c *CurrencyDay) GetData(CurrencyDate time.Time) {
	xDate := CurrencyDate
	t := new(tarih_Date)
	currDay := t.getDate(CurrencyDate, xDate)

	for {
		if currDay == nil {
			CurrencyDate = CurrencyDate.AddDate(0, 0, -1)
			currDay := t.getDate(CurrencyDate, xDate)
			if currDay != nil {
				break
			}
		} else {
			break
		}
	}
}
func (c *tarih_Date) getDate(CurrencyDate time.Time, XDate time.Time) *CurrencyDay {
	currDay := new(CurrencyDay)
	var resp *http.Response
	var err error
	var url string

	currDay = new(CurrencyDay)
	url = "https:// www.tcmb.gov.tr/kurlar/" + CurrencyDate.Format("200601") + "/" + CurrencyDate.Format("02012006") + ".xml"
	resp, err = http.Get(url)

	// fmt.Println(url)

	if err != nil {
		fmt.Println(err)
	} else {

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			tarih := new(tarih_Date)
			d := xml.NewDecoder(resp.Body)
			marshalErr := d.Decode(&tarih)

			if marshalErr != nil {
				log.Printf("error: %v", marshalErr)
			}
			c = &tarih_Date{}
			currDay.ID = XDate.Format("20060102")
			currDay.Date = XDate
			currDay.DayNo = tarih.Bulten_No
			currDay.Currencies = make([]Currency, len(tarih.Currency))
			for i, curr := range tarih.Currency {
				currDay.Currencies[i].Code = curr.CurrencyCode
				currDay.Currencies[i].CurrencyName = curr.CurrencyName
				currDay.Currencies[i].BanknoteBuying = strconv.ParseFloat(curr.BanknoteBuying, 64)
				currDay.Currencies[i].BanknoteSelling = curr.BanknoteSelling
				currDay.Currencies[i].CurrencyName = curr.CurrencyName
				currDay.Currencies[i].CurrencyName = curr.CurrencyName

			}
		}
	}
}
