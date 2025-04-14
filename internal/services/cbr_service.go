package services

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

type CBRService struct{}

func (s *CBRService) GetKeyRate() (float64, error) {
	soapBody := `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ws="http://web.cbr.ru/">
   <soapenv:Header/>
   <soapenv:Body>
      <ws:KeyRateXML/>
   </soapenv:Body>
</soapenv:Envelope>`

	resp, err := http.Post("https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx",
		"text/xml; charset=utf-8", bytes.NewBuffer([]byte(soapBody)))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("failed to get data from server: " + resp.Status)
	}

	return parseKeyRate(resp.Body)
}

func parseKeyRate(body io.Reader) (float64, error) {
	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(body); err != nil {
		return 0, err
	}

	rateElem := doc.FindElement("//KeyRateXMLResult/KeyRate/Rate")
	if rateElem == nil {
		return 0, errors.New("key rate not found")
	}

	rate, err := parseFloat(rateElem.Text())
	return rate, err
}

func parseFloat(s string) (float64, error) {
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}
