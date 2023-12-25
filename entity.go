package entity

type Product struct {
	ProductName   string
	PortionHigh   int
	PortionMedium int
	PortionLow    int
	PortionSize   string
	CarbId        []int
	Stage         int
	CategoryId    int
}
