package tmoex

type MoexApiResponseSecurityInfo struct {
	Description struct {
		Data [][]string `json:"data"`
	}
	Boards struct {
		Data [][4]any `json:"data"`
	}
}

type Security struct {
	Board     Board  `json:"board" db:"board"`
	Engine    Engine `json:"engine" db:"engine"`
	Id        int    `json:"id" db:"id"`
	Market    Market `json:"market" db:"market"`
	Name      string `json:"name" db:"name"`
	ShortName string `json:"shortname" db:"shortname"`
	Ticker    string `json:"ticker" db:"ticker"`
}

type Share struct {
	Security
}

type Bond struct {
	Security
}
