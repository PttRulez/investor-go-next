import { SelectInput } from '@pttrulez/mui-based-ui';
import { SelectChangeEvent } from '@mui/material';

type Props = {
  onChooseTransaction: (e: SelectChangeEvent) => void;
};

export enum PortfolioActionsMap {
  buy = 'buy',
  transaction = 'transaction',
  sell = 'sell',
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
        ]}
        value={''}
      />
    </>
  );
};

export default PortfolioTableToolbar;
