package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/beevik/etree"
)

type CBRService struct{}

// Метод вызываемый из handler
func (s *CBRService) GetKeyRate() (float64, error) {
	soapRequest := buildSOAPRequest()
	rawBody, err := sendRequest(soapRequest)
	if err != nil {
		return 0, err
	}
	rate, err := parseXMLResponse(rawBody)
	if err != nil {
		return 0, err
	}
	return rate + 5.0, nil // +5% маржа банка
}

// Генерация SOAP-запроса с диапазоном дат
func buildSOAPRequest() string {
	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
  <soap12:Body>
    <KeyRate xmlns="http://web.cbr.ru/">
      <fromDate>%s</fromDate>
      <ToDate>%s</ToDate>
    </KeyRate>
  </soap12:Body>
</soap12:Envelope>`, fromDate, toDate)
}

// Отправка HTTP POST-запроса
func sendRequest(soapRequest string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(
		"POST",
		"https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx",
		bytes.NewBuffer([]byte(soapRequest)),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://web.cbr.ru/KeyRate")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}
	return rawBody, nil
}

// Парсинг XML-ответа и извлечение последней ставки
func parseXMLResponse(rawBody []byte) (float64, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(rawBody); err != nil {
		return 0, fmt.Errorf("ошибка парсинга XML: %v", err)
	}

	// Поиск всех KR (KeyRate) и выбор последнего
	krElements := doc.FindElements("//diffgram/KeyRate/KR")
	if len(krElements) == 0 {
		return 0, errors.New("ставка не найдена")
	}

	// Последний KR — первый в списке (если упорядочен по убыванию)
	latest := krElements[0]
	rateEl := latest.FindElement("./Rate")
	if rateEl == nil {
		return 0, errors.New("тег <Rate> не найден")
	}

	rateStr := rateEl.Text()
	rateStr = normalizeRate(rateStr)

	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		return 0, fmt.Errorf("ошибка конвертации ставки: %v", err)
	}
	return rate, nil
}

// Заменяем запятую на точку
func normalizeRate(s string) string {
	return bytes.NewBufferString(s).String()
}
