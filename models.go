package main

import (
	"sort"
	"time"
)

var taxCostIncreasedAt, _ = time.Parse(time.RFC3339, "2019-10-01T00:00:00+00:00")

var accountMap = map[string]string{
	"家賃":      "地代家賃",
	"WEBサービス": "雑費",
	"通信":      "通信費",
	"買い物":     "雑費",
	"税金":      "租税公課",
	"水道":      "水道光熱費",
	"電気":      "水道光熱費",
	"交際費":     "交際費",
	"本・雑誌":    "新聞図書費",
	"タクシー":    "旅費交通費",
	"ガス":      "水道光熱費",
	"オフィス設備":  "修繕費",
	"携帯電話":    "通信費",
	"電車":      "旅費交通費",
	"電化製品":    "消耗品費",
	"インターネット": "通信費",
	"消耗品":     "消耗品費",
	"交通":      "旅費交通費",
	"住民税":     "事業主貸",
}

type (
	// DateTime struct
	DateTime struct {
		time.Time
	}

	// MoneytreeBankAccountHistory struct
	MoneytreeBankAccountHistory struct {
		OccuredAt           *DateTime `csv:"日付"`
		Amount              int       `csv:"金額"`
		BankAccountCurrency string    `csv:"口座通貨"`
		Summary             string    `csv:"ご利用先・摘要"`
		Memo                string    `csv:"メモ"`
		Receipt             string    `csv:"領収書"`
		Balance             int       `csv:"取引後残高"`
		Category            string    `csv:"カテゴリ"`
		ExpenseCategory     string    `csv:"経費"`
	}

	// MoneytreeBankAccountHistories slice
	MoneytreeBankAccountHistories []MoneytreeBankAccountHistory

	// MoneytreeExpense struct
	MoneytreeExpense struct {
		OccuredAt           *DateTime `csv:"日付"`
		Amount              int       `csv:"金額"`
		BankAccountCurrency string    `csv:"口座通貨"`
		Summary             string    `csv:"ご利用先・摘要"`
		Memo                string    `csv:"メモ"`
		Receipt             string    `csv:"領収書"`
		Category            string    `csv:"カテゴリ"`
		BankAccountName     string    `csv:"口座名"`
		BankAccountNumber   string    `csv:"口座番号"`
		LocalCurrency       string    `csv:"現地通貨"`
		LocalCurrencyAmount int       `csv:"現地通貨金額"`
	}

	// MoneytreeExpenses slice
	MoneytreeExpenses []MoneytreeExpense

	// FreeeDeal struct
	FreeeDeal struct {
		Category              string    `csv:"収支区分"`
		ManagementNumber      string    `csv:"管理番号"`
		OccuredAt             *DateTime `csv:"発生日"`
		SettlementDeadlineAt  *DateTime `csv:"決済期日"`
		SupplierCode          string    `csv:"取引先コード"`
		Supplier              string    `csv:"取引先"`
		Account               string    `csv:"勘定科目"`
		TaxCategory           string    `csv:"税区分"`
		Amount                int       `csv:"金額"`
		TaxCalcCategory       string    `csv:"税計算区分"`
		Taxcost               string    `csv:"税額"`
		Remarks               string    `csv:"備考"`
		Item                  string    `csv:"品目"`
		Department            string    `csv:"部門"`
		Memotag               string    `csv:"メモタグ（複数指定可、カンマ区切り）"`
		Segment1              string    `csv:"セグメント1"`
		Segment2              string    `csv:"セグメント2"`
		Segment3              string    `csv:"セグメント3"`
		SettlementedAt        *DateTime `csv:"決済日"`
		SettlementBankAccount string    `csv:"決済口座"`
		Settlement            int       `csv:"決済金額"`
	}

	// FreeeDeals slice
	FreeeDeals []FreeeDeal
)

// toFreeeDeals MoneytreeBankAccountHistories -> FreeeDeals
func (m *MoneytreeBankAccountHistories) toFreeeDeals() *FreeeDeals {
	deals := FreeeDeals{}
	for _, v := range *m {
		freeeDeal := v.toFreeeDeal()
		if nil != freeeDeal {
			deals = append(deals, *freeeDeal)
		}
	}
	sort.SliceStable(deals, func(i, j int) bool { return deals[i].OccuredAt.Before(deals[j].OccuredAt.Time) })

	return &deals
}

// toFreeeDeals MoneytreeExpenses -> FreeeDeals
func (m *MoneytreeExpenses) toFreeeDeals() *FreeeDeals {
	deals := FreeeDeals{}
	for _, v := range *m {
		freeeDeal := v.toFreeeDeal()
		if nil != freeeDeal {
			deals = append(deals, *freeeDeal)
		}
	}
	sort.SliceStable(deals, func(i, j int) bool { return deals[i].OccuredAt.Before(deals[j].OccuredAt.Time) })

	return &deals
}

// toFreeeDeal MoneytreeBankAccountHistory -> FreeeDeal (incomes only)
func (m *MoneytreeBankAccountHistory) toFreeeDeal() *FreeeDeal {
	if m.Amount <= 0 {
		// ignore expense
		return nil
	}

	taxCategory := "課税売上10%"
	if m.OccuredAt.Before(taxCostIncreasedAt) {
		taxCategory = "課税売上8%"
	}

	return &FreeeDeal{
		Category:              "収入",
		ManagementNumber:      "",
		OccuredAt:             m.OccuredAt,
		SettlementDeadlineAt:  nil,
		SupplierCode:          "",
		Supplier:              "",
		Account:               "売上高",
		TaxCategory:           taxCategory,
		Amount:                m.Amount,
		TaxCalcCategory:       "税込",
		Taxcost:               "",
		Remarks:               m.Memo,
		Item:                  "",
		Department:            "",
		Memotag:               "",
		Segment1:              "",
		Segment2:              "",
		Segment3:              "",
		SettlementedAt:        m.OccuredAt,
		SettlementBankAccount: "事業主貸",
		Settlement:            m.Amount,
	}
}

// toFreeeDeal MoneytreeExpense -> FreeeDeal
func (m *MoneytreeExpense) toFreeeDeal() *FreeeDeal {
	taxCategory := "課対仕入10%"
	if m.OccuredAt.Before(taxCostIncreasedAt) {
		taxCategory = "課対仕入8%"
	}

	account, _ := accountMap[m.Category]
	if account == "" {
		account = "雑費"
	}

	return &FreeeDeal{
		Category:              "支出",
		ManagementNumber:      "",
		OccuredAt:             m.OccuredAt,
		SettlementDeadlineAt:  nil,
		SupplierCode:          "",
		Supplier:              "",
		Account:               account,
		TaxCategory:           taxCategory,
		Amount:                -m.Amount,
		TaxCalcCategory:       "税込",
		Taxcost:               "",
		Remarks:               m.Memo,
		Item:                  "",
		Department:            "",
		Memotag:               "",
		Segment1:              "",
		Segment2:              "",
		Segment3:              "",
		SettlementedAt:        m.OccuredAt,
		SettlementBankAccount: "事業主借",
		Settlement:            -m.Amount,
	}
}

// MarshalCSV convert the internal date as CSV string
func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Format("2006/01/02"), nil
}

// UnmarshalCSV convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2006/01/02", csv)
	return err
}
