'use client';
import { useEffect, useState } from 'react';
import { moexService } from '@/axios/moex/moex.service';
import investorService from '@/axios/investor/investor.service';
import { useQuery } from '@tanstack/react-query';
import CandlestickChart from '@/components/ui/Charts/Candlestick/CandlestickChart';
import { Button, Dialog, Typography } from '@mui/material';
import { useParams } from 'next/navigation';
import { CandlestickData } from 'lightweight-charts';
import { dependOn } from '@/utils/react-query';
import OpinionForm from '@/app/opinions/components/OpinionForm/OpinionForm';
import OpinionsTable from '@/app/opinions/components/OpinionsTable';
import { Exchange, SecurityType } from '@/types/enums';

const MoexSharePage = (): JSX.Element => {
  const { ticker } = useParams<{ ticker: string }>();
  const [chartData, setChartData] = useState<CandlestickData[] | null>(null);
  const [stockName, setStockName] = useState('');
  const [opinionModalOpen, setOpinionModalOpen] = useState(false);

  const { data: shareData } = useQuery({
    queryKey: ['info', ticker],
    queryFn: () => investorService.moexShare.getByTicker(ticker),
  });

  const { data: historyData } = useQuery({
    queryKey: ['history', ticker],
    enabled: !!shareData,
    queryFn: dependOn(shareData, shareData =>
      moexService.getStockHistoryByTicker({
        market: shareData.market,
        board: shareData.board,
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
        variant={'h2'}
        sx={{
          color: 'grey.700',
          textAlign: 'center',
          marginBottom: '50px',
        }}
      >
        {stockName}
      </Typography>
      <Button
        variant="outlined"
        sx={{
          color: 'grey.700',
          borderColor: 'grey.700',
          marginBottom: '50px',
        }}
        onClick={() => setOpinionModalOpen(true)}
      >
        + Мнение
      </Button>
      {chartData && (
        <CandlestickChart sx={{ marginBottom: '50px' }} data={chartData} />
      )}
      {shareData && (
        <OpinionsTable
          filters={{
            exchange: Exchange.MOEX,
            securityType: shareData.securityType,
            securityId: shareData.id,
          }}
        />
      )}
      {shareData && (
        <Dialog
          fullWidth
          maxWidth="md"
          open={opinionModalOpen}
          onClose={() => setOpinionModalOpen(false)}
        >
          <OpinionForm
            afterSuccessfulSubmit={() => setOpinionModalOpen(false)}
            defaultName={shareData.name}
            securityId={shareData?.id}
            securityType={SecurityType.SHARE}
          />
        </Dialog>
      )}
    </>
  );
};

export default MoexSharePage;
