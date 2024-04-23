export enum Exchange {
  MOEX = 'MOEX',
}

export enum DealType {
  BUY = 'BUY',
  SELL = 'SELL'
}

export enum OpinionType {
  FLAT = 'FLAT',
  GENERAL = 'GENERAL',
  GROWTH = 'GROWTH',
  REDUCTION = 'REDUCTION',
}

export enum Role {
	ADMIN = 'ADMIN',
	INVESTOR = 'INVESTOR'
}

export enum SecurityType {
  BOND = 'BOND',
  CURRENCY = 'CURRENCY',
  FUTURES = 'FUTURES',
  INDEX = 'INDEX',
  PIF = 'PIF',
  SHARE = 'SHARE'
}

export enum TransactionType {
  CASHOUT = 'CASHOUT',
  DEPOSIT = 'DEPOSIT'
}


//  MOEX

export enum MoexEngine {
  stock = 'stock',
  currency = 'currency'
}

export enum MoexMarket {
  shares = 'shares',
  bonds = 'bonds',
  index = 'index',
  selt = 'selt'       // Валюта: Биржевые сделки с ЦК
}

export enum MoexBoard {
  TQBR = 'TQBR',												// Т+: Акции и ДР - безадрес.
  CETS = 'CETS'                         // Системные сделки - безадрес.
}

export enum MoexSecurityType {
  common_share = 'common_share', 				// 'акция обыкновенная'
  preferred_share = 'preferred_share', 	// 'акция привелигированная'
  
	corporate_bond = 'corporate_bond',		// 'корпоративная облигация'
  exchange_bond = 'exchange_bond', 			// 'облигация'
  ofz_bond = 'ofz_bond',								// 'ОФЗ'

  exchange_ppif = 'exchange_ppif', 			// 'биржевой ПИФ'
  public_ppif = 'public_ppif', 					// 'публичный ПИФ'
  stock_index_if = 'stock_index_if', 		// 'iNAV облигаций'

  futures = 'futures', 									// 'фьючерс'

  stock_index = 'stock_index',					// 'индекс'
}