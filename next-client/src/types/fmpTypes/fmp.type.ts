export interface historicalPrice {
  date : string,
  open : number,
  low : number,
  high : number,
  close : number,
  volume : number
}

export interface searchStock {
  symbol : string,
  name : string,
  currency : string,
  stockExchange : string,
  exchangeShortName : string
}

export interface historicalDaily {
  date: string,
  open: number,
  high: number,
  low: number,
  close: number,
  adjClose: number,
  volume: number,
  unadjustedVolume: number,
  change: number,
  changePercent: number,
  vwap: number,
  label: string,
  changeOverTime: number
}

export type requestPeriodType = 'year' | 'quarter'