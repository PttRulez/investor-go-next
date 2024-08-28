import { useEffect, useRef } from 'react';
import { CandlestickData, ColorType, createChart } from 'lightweight-charts';
import { Box, BoxProps } from '@mui/material';
import { SxProps } from '@mui/material';

interface CandleChartProps extends BoxProps {
  data: CandlestickData[];
}

const boxPadding = 20;
const backgroundColor = 'white';

const CandlestickChart = ({ sx, data }: CandleChartProps) => {
  const chartContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!chartContainerRef.current || !data) return;

    const handleResize = () => {
      chart.applyOptions({ width: chartContainerRef?.current?.clientWidth });
    };

    const chartOptions = {
      autoSize: true,
      layout: {
        textColor: 'black',
        background: {
          type: ColorType.Solid,
          color: backgroundColor,
        },
      },
      // width: 'inherit', // chartContainerRef.current.clientWidth - boxPadding * 2,
      height: 300,
    };

    const chart = createChart(chartContainerRef.current, chartOptions);

    const candlestickSeries = chart.addCandlestickSeries({
      upColor: '#26a69a',
      downColor: '#ef5350',
      borderVisible: false,
      wickUpColor: '#26a69a',
      wickDownColor: '#ef5350',
      autoscaleInfoProvider: (original: any) => {
        const res = original();
        if (res !== null) {
          res.priceRange.minValue -= 10;
          res.priceRange.maxValue += 10;
        }
        return res;
      },
    });

    candlestickSeries.setData(data);

    chart.timeScale().fitContent();

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      chart.remove();
    };
  }, [data]);

  return (
    <Box
      ref={chartContainerRef}
      sx={{
        padding: `${boxPadding}px`,
        backgroundColor: backgroundColor,
        ...sx,
      }}
    />
  );
};

export default CandlestickChart;
