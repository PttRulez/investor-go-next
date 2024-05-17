package iss_client

type MoexApiResponseSecurityInfo struct {
	Description struct {
		Data [][]string `json:"data"`
	}
	Boards struct {
		Data [][4]any `json:"data"`
	}
}

type MoexApiResponseCurrentPrices struct {
	Securities struct {
		Data [][3]any
	}
}
