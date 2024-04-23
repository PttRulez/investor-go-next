import Grid from '@mui/material/Unstable_Grid2';

type Props = {
  depositsSum: number;
  cashoutsSum: number;
  currentValue: number;
  profitability: string;
};

const PortfolioTableFooter = ({ depositsSum, cashoutsSum, currentValue, profitability }: Props) => {
  return (
    <Grid justifyContent={'space-between'} container sx={{ width: '100%', '& .MuiGrid2-root': { padding: '20px' } }}>
      <Grid>Сумма депозитов: {depositsSum.toLocaleString('RU-ru')}</Grid>
      <Grid>Сумма кэшаутов: {cashoutsSum.toLocaleString('RU-ru')}</Grid>
      <Grid>Текущая стоимость: {currentValue.toLocaleString('RU-ru')}</Grid>
      <Grid>Общая доходность: {profitability}</Grid>
    </Grid>
  );
};

export default PortfolioTableFooter;
