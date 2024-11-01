package issclient

type MoexAPIResponseSecurityInfo struct {
	Description struct {
		Data [][]string `json:"data"`
	} `json:"description"`
	Boards struct {
		Data [][4]any `json:"data"`
	} `json:"boards"`
}

type MoexAPIResponseCurrentPrices struct {
	Securities struct {
		Data [][3]any `json:"data"`
	} `json:"securities"`
}

type MoexAPIResponseCurrency struct {
	Cbrf struct {
		RateNames []string `json:"columns"`
		Data      [][]any  `json:"data"`
	} `json:"cbrf"`
}
