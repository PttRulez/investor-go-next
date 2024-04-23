import { SecurityType, MoexSecurityType } from '@/types/enums';

const map = {
  [MoexSecurityType.common_share]: SecurityType.SHARE,
  [MoexSecurityType.preferred_share]: SecurityType.SHARE,

  [MoexSecurityType.corporate_bond]: SecurityType.BOND,
  [MoexSecurityType.exchange_bond]: SecurityType.BOND,
  [MoexSecurityType.ofz_bond]: SecurityType.BOND,

  [MoexSecurityType.exchange_ppif]: SecurityType.PIF,
  [MoexSecurityType.public_ppif]: SecurityType.PIF,
  [MoexSecurityType.stock_index_if]: SecurityType.PIF,

  [MoexSecurityType.futures]: SecurityType.FUTURES,

  [MoexSecurityType.stock_index]: SecurityType.INDEX,
};

export const getSecurityTypeFromMoexSecType = (
  type: MoexSecurityType,
): SecurityType => {
  return map[type];
};

export const getNestedProp = (obj: Record<string, any>, path: string) =>
  path.split('.').reduce((acc, part) => acc && acc[part], obj);
