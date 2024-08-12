package issclient

type MoexAPIResponseSecurityInfo struct {
	Description struct {
		Data [][]string `json:"data"`
	} `json:"description"`
	Boards struct {
		Data [][3]string `json:"data"`
	} `json:"boards"`
}

type MoexAPIResponseCurrentPrices struct {
	Securities struct {
		Data [][3]any `json:"data"`
	} `json:"securities"`
}
