'use client';
import PositionForm from '@/app/position/components/PositionForm';
import { IPortfolioResponse, IPositionResponse } from '@/types/apis/go-api';
import { getDefaultMRTOptions } from '@/utils/mrt-default-options';
import EditIcon from '@mui/icons-material/Edit';
import {
  Box,
  Dialog,
  IconButton,
  SelectChangeEvent,
  SxProps,
  Typography,
} from '@mui/material';
import {
  MRT_Row,
  MaterialReactTable,
  useMaterialReactTable,
  type MRT_ColumnDef,
} from 'material-react-table';
import { useRouter } from 'next/navigation';
import { SyntheticEvent, useMemo, useState } from 'react';
import PortfolioTableFooter from './PortfolioTableFooter';
import PortfolioTableToolbar from './PortfolioTableToolbar';
import { SecurityType } from '@/types/enums';
import NextLink from 'next/link';
import Link from '@mui/material/Link';
import PositionDetails from './PositionDetails';

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

  const columns = useMemo<Array<MRT_ColumnDef<IPositionResponse>>>(
    () => [
      {
        header: 'Название',
        accessorKey: 'shortName',
        Cell: ({ row }) => {
          const p = row.original;
          let href = '';
          if (p.securityType) {
            let secType;
            switch (p.securityType) {
              case SecurityType.BOND:
                secType = 'bonds';
                break;
              case SecurityType.SHARE:
                secType = 'shares';
                break;
              default:
                secType = 'shares';
            }
            href = `/${secType}/moex/${p.ticker}`;
          }

          return (
            <Typography variant="body1">
              {p.shortName && (
                <>
                  <Link
                    component={NextLink}
                    sx={{
                      color: 'info.main',
                      textDecoration: 'none',
                    }}
                    href={href}
                  >
                    {p.shortName && `${p.shortName}`}
                  </Link>
                  {' - '}
                </>
              )}
              <Box
                component="span"
                sx={{
                  fontSize: '0.8rem',
                  color: 'gray',
                  paddingLeft: '5px',
                }}
              >
                {p.ticker}
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
        header: 'Средняя цена',
        accessorKey: 'averagePrice',
        size: 5,
      },
      {
        header: 'Текущая цена',
        accessorKey: 'currentPrice',
        size: 5,
        accessorFn: position => {
          if (!position.currentPrice) return '';
          if (position.targetPrice) {
            return `${position.currentPrice} (${position.targetPrice})`;
          }
          return position.currentPrice;
        },
      },
      {
        header: 'Стоимость',
        accessorKey: 'currentCost',
        size: 5,
        accessorFn: position => {
          return position.currentCost?.toLocaleString('ru-RU');
        },
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

  const portfolioData = useMemo<IPositionResponse[]>(() => {
    if (!portfolio) return [];
    return [
      {
        currentCost: portfolio?.bondPositions?.reduce((acc, p) => {
          return acc + p.currentCost;
        }, 0),
        subRows: portfolio?.bondPositions,
        ticker: 'Облигации',
      },
      {
        currentCost: portfolio?.sharePositions?.reduce((acc, p) => {
          return acc + p.currentCost;
        }, 0),
        subRows: portfolio?.sharePositions,
        ticker: 'Акции',
      },
      {
        currentCost: portfolio?.cash,
        ticker: 'Рубли',
      },
    ] as IPositionResponse[];
  }, [portfolio]);

  // Свойства одной строчки таблицы
  const muiTableBodyRowProps = ({
    row,
  }: {
    row: MRT_Row<IPositionResponse>;
  }) => {
    const sx: SxProps = {
      backgroundColor: 'inherit',
    };

    if (row.original.shortName == undefined) {
      sx.backgroundColor = 'rgba(0, 0, 0, 0.05)';
    }

    return {
      sx,
      hover: false,
    };
  };

  const table = useMaterialReactTable<IPositionResponse>({
    columns,
    data: portfolioData,
    enableColumnActions: false,
    enableColumnDragging: false,
    enableExpanding: true,
    enableTopToolbar: false,
    enableToolbarInternalActions: false,
    initialState: {
      expanded: {
        0: true,
        1: true,
      },
    },
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

    muiTableBodyRowProps: muiTableBodyRowProps,
    muiTablePaperProps: { sx: { marginBottom: '100px' } },
    renderBottomToolbar: () => (
      <PortfolioTableFooter
        cashoutsSum={portfolio.cashoutsSum}
        currentValue={portfolio.totalCost}
        depositsSum={portfolio.depositsSum}
        profitability={portfolio.profitability}
      />
    ),
    renderDetailPanel: ({ row }) => {
      return row.original.comment || row.original.opinions?.length > 0 ? (
        <PositionDetails position={row.original} />
      ) : null;
    },
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
