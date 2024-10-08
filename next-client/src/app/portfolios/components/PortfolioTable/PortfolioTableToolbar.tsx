import { SelectInput } from '@pttrulez/mui-based-ui';
import { SelectChangeEvent } from '@mui/material';

type Props = {
  onChooseTransaction: (e: SelectChangeEvent) => void;
};

export enum PortfolioActionsMap {
  buy = 'buy',
  transaction = 'transaction',
  sell = 'sell',
  dividend = 'dividend',
  coupon = 'coupon',
  expense = 'expense',
}

const PortfolioTableToolbar = ({ onChooseTransaction }: Props) => {
  return (
    <>
      <SelectInput
        sx={{
          minWidth: '200px',
        }}
        onChange={onChooseTransaction}
        label={'Добавить'}
        options={[
          { id: PortfolioActionsMap.buy, name: 'BUY' },
          { id: PortfolioActionsMap.sell, name: 'SELL' },
          { id: PortfolioActionsMap.transaction, name: 'Депозит/Кэшаут' },
          { id: PortfolioActionsMap.dividend, name: 'Дивиденд' },
          { id: PortfolioActionsMap.coupon, name: 'Купон' },
          { id: PortfolioActionsMap.expense, name: 'Траты' },
        ]}
        value={''}
      />
    </>
  );
};

export default PortfolioTableToolbar;
