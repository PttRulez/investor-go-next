'use client';
import { useEffect, useState } from 'react';
import { moexService } from '@/axios/moex/moex.service';
import investorService from '@/axios/investor/investor.service';
import { useQuery } from '@tanstack/react-query';
import CandlestickChart from '@/components/ui/Charts/Candlestick/CandlestickChart';
import { Typography } from '@mui/material';
import { useParams } from 'next/navigation';
import { dependOn } from '@/utils/react-query';
import { CandlestickData } from 'lightweight-charts';

const MoexBondTicker = (): JSX.Element => {
  const { ticker } = useParams<{ ticker: string }>();

  const [chartData, setChartData] = useState<CandlestickData[] | null>(null);
  const [stockName, setStockName] = useState('');

  const { data: bondData } = useQuery({
    queryKey: ['info', ticker],
    queryFn: () => investorService.moexBond.getByTicker(ticker),
  });

  const { data: historyData } = useQuery({
    queryKey: ['history', ticker],
    enabled: !!bondData,
    queryFn: dependOn(bondData, bondData =>
      moexService.getStockHistoryByTicker({
        market: bondData.market,
        board: bondData.board,
        ticker,
      }),
    ),
  });

  useEffect(() => {
    if (historyData) {
      setChartData(
        historyData.data.history.data.map(arr => ({
          open: arr[1],
          high: arr[2],
          low: arr[3],
          close: arr[4],
          time: arr[5],
        })),
      );

      setStockName(historyData.data.history.data[0][0]);
    }
  }, [historyData]);

  return (
    <>
      <Typography
        variant={'h1'}
        sx={{ backgroundColor: 'red', margin: 'auto' }}
      >
        {stockName}
      </Typography>
      {chartData && <CandlestickChart data={chartData} />}
    </>
  );
};

export default MoexBondTicker;
