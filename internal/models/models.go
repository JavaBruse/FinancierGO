package models

type Payment struct {
	ID       int64
	CreditID int64
	Amount   float64
	Date     string
	Status   string
}

type Stats struct {
	TotalBalance    float64
	MonthlyIncome   float64
	MonthlyExpenses float64
	ActiveCredits   int
}

type CreditLoad struct {
	TotalDebt      float64
	MonthlyPayment float64
	DebtToIncome   float64
}

type Prediction struct {
	NextMonthBalance float64
	Trend            string
	Confidence       float64
}

type KeyRate struct {
	Rate float64
	Date string
}
