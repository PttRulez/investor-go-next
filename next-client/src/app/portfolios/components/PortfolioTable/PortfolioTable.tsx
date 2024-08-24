'use client';
import PositionForm from '@/app/position/components/PositionForm';
import {
  IMoexBondPositionResponse,
  IMoexSharePositionResponse,
  IPortfolioResponse,
  IPositionResponse,
} from '@/types/apis/go-api';
import { SecurityType } from '@/types/enums';
import { getDefaultMRTOptions } from '@/utils/mrt-default-options';
import EditIcon from '@mui/icons-material/Edit';
import {
  Box,
  Dialog,
  IconButton,
  SelectChangeEvent,
  Typography,
} from '@mui/material';
import { LightTooltip } from '@pttrulez/mui-based-ui';
import {
  MRT_ExpandedState,
  MaterialReactTable,
  useMaterialReactTable,
  type MRT_ColumnDef,
} from 'material-react-table';
import { SyntheticEvent, useMemo, useState } from 'react';
import PortfolioTableFooter from './PortfolioTableFooter';
import PortfolioTableToolbar from './PortfolioTableToolbar';

const defaultMRTOptions = getDefaultMRTOptions<IPositionResponse>();

const PortfolioTable = ({
  onChooseTransaction,
  portfolio,
}: {
  onChooseTransaction: (e: SelectChangeEvent) => void;
  portfolio: IPortfolioResponse;
}) => {
  const [positionToEdit, setPositionToEdit] =
    useState<IPositionResponse | null>(null);
  const [expanded, setExpanded] = useState<MRT_ExpandedState>({});

  const columns = useMemo<
    Array<MRT_ColumnDef<IMoexSharePositionResponse | IMoexBondPositionResponse>>
  >(
    () => [
      {
        header: 'Название',
        accessorKey: 'shortName',
        Cell: ({ row }) => {
          const p = row.original;
          return (
            <Typography variant="body1">
              {p.shortName}
              <Box
                sx={{
                  fontSize: '0.8rem',
                  color: 'gray',
                }}
              >
                {'ticker' in p ? p.ticker : p.isin}
              </Box>
            </Typography>
          );
        },
        size: 50,
      },
      {
        header: 'Кол-во',
        accessorKey: 'amount',
        size: 5,
      },
      {
        header: 'Цена покупки',
        accessorKey: 'averagePrice',
        size: 5,
      },
      {
        header: 'Текущая цена',
        accessorKey: 'currentPrice',
        size: 5,
      },
      {
        header: 'Стоимость',
        accessorKey: 'currentCost',
        size: 5,
      },
      {
        header: 'Коммент',
        accessorKey: 'comment',
        Cell: ({ row }) => {
          return (
            <LightTooltip title={row.original.comment}>
              <Typography
                sx={{
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  display: '-webkit-box',
                  WebkitLineClamp: '2',
                  WebkitBoxOrient: 'vertical',
                }}
              >
                {row.original.comment}
              </Typography>
            </LightTooltip>
          );
        },
        size: 5,
      },
      {
        header: '',
        Header: () => (
          <PortfolioTableToolbar
            onChooseTransaction={e => {
              onChooseTransaction(e);
            }}
          />
        ),
        accessorKey: 'total',
        Cell: ({ row }) => {
          return row.original.shortName ? (
            <IconButton
              onClick={(e: SyntheticEvent) => {
                e.stopPropagation();
                setPositionToEdit(row.original);
              }}
            >
              <EditIcon />
            </IconButton>
          ) : (
            <></>
          );
        },
      },
    ],
    [onChooseTransaction],
  );
  type PositionRow = IMoexBondPositionResponse | IMoexSharePositionResponse;

  const portfolioData = useMemo<PositionRow[]>(() => {
    if (!portfolio) return [];
    return [
      {
        amount: '' as unknown as number,
        averagePrice: '' as unknown as number,
        comment: '',
        id: '' as unknown as number,
        currentPrice: '' as unknown as number,
        currentCost: portfolio?.bondPositions.reduce((acc, p) => {
          console.log('acc', acc, p);
          return acc + p.currentCost;
        }, 0),
        targetPrice: '' as unknown as number,
        shortName: '',
        ticker: 'Облигации',
        subRows: portfolio?.bondPositions,
      },
      {
        amount: '' as unknown as number,
        averagePrice: '' as unknown as number,
        comment: '',
        id: '' as unknown as number,
        currentPrice: '' as unknown as number,
        currentCost: portfolio?.sharePositions.reduce((acc, p) => {
          return acc + p.currentCost;
        }, 0),
        targetPrice: '' as unknown as number,
        shortName: '',
        ticker: 'Акции',
        subRows: portfolio?.sharePositions,
      },
      {
        amount: '' as unknown as number,
        averagePrice: '' as unknown as number,
        comment: '',
        id: '' as unknown as number,
        currentPrice: '' as unknown as number,
        currentCost: portfolio?.cash,
        targetPrice: '' as unknown as number,
        shortName: '',
        ticker: 'Рубли',
      },
    ];
  }, [portfolio]);

  const table = useMaterialReactTable<PositionRow>({
    // ...defaultMRTOptions,
    columns,
    data: portfolioData,
    enableColumnActions: false,
    enableColumnDragging: false,
    enableExpanding: true,
    enableTopToolbar: false,
    enableToolbarInternalActions: false,
    // icons: { SortIcon: <></> },
    muiExpandAllButtonProps: {
      sx: {
        display: 'none',
      },
    },
    renderRowActions: () => <></>,
    muiTableHeadRowProps: {
      sx: {
        '.Mui-TableHeadCell-Content-Wrapper': { whiteSpace: 'pre-wrap' },
        '.MuiTableSortLabel-icon': {
          display: 'none',
        },
        '.MuiTableCell-head': {
          textAlign: 'center',
          verticalAlign: 'middle',
        },
        '.Mui-TableHeadCell-Content': {
          justifyContent: 'center',
          '.MuiBadge-root': {
            display: 'none',
          },
        },
      },
    },
    muiTableBodyCellProps: {
      sx: {
        textAlign: 'center',
      },
    },
    muiTableBodyRowProps: ({ row }) => {
      if (!row.original.shortName) {
        return {
          sx: {
            backgroundColor: 'rgba(0, 0, 0, 0.05)',
          },
        };
      }
      return {
        sx: {
          backgroundColor: 'inherit',
        },
      };
    },
    muiTablePaperProps: { sx: { marginBottom: '100px' } },
    renderBottomToolbar: () => (
      <PortfolioTableFooter
        cashoutsSum={portfolio.cashoutsSum}
        currentValue={portfolio.totalCost}
        depositsSum={portfolio.depositsSum}
        profitability={portfolio.profitability}
      />
    ),
  });

  return (
    <>
      <MaterialReactTable table={table} />
      <Dialog open={!!positionToEdit} onClose={() => setPositionToEdit(null)}>
        {positionToEdit && (
          <PositionForm
            afterSuccessfulSubmit={() => setPositionToEdit(null)}
            id={positionToEdit.id}
            position={{
              comment: positionToEdit.comment,
              targetPrice: positionToEdit.targetPrice,
            }}
            portfolioId={portfolio.id}
          />
        )}
      </Dialog>
    </>
  );
};

export default PortfolioTable;
