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
  const router = useRouter();

  const columns = useMemo<Array<MRT_ColumnDef<IPositionResponse>>>(
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
      },
      {
        header: 'Стоимость',
        accessorKey: 'currentCost',
        size: 5,
        accessorFn: position => {
          return position.currentCost?.toLocaleString('ru-RU');
        },
      },
      // {
      //   header: 'Коммент',
      //   accessorKey: 'comment',
      //   Cell: ({ row }) => {
      //     return (
      //       <LightTooltip title={row.original.comment}>
      //         <Typography
      //           sx={{
      //             overflow: 'hidden',
      //             textOverflow: 'ellipsis',
      //             display: '-webkit-box',
      //             WebkitLineClamp: '2',
      //             WebkitBoxOrient: 'vertical',
      //           }}
      //         >
      //           {row.original.comment}
      //         </Typography>
      //       </LightTooltip>
      //     );
      //   },
      //   size: 5,
      // },
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

  const muiTableBodyRowProps = ({
    row,
  }: {
    row: MRT_Row<IPositionResponse>;
  }) => {
    let sx: SxProps = {
      backgroundColor: 'inherit',
    };

    if (row.original.shortName == undefined) {
      sx.backgroundColor = 'rgba(0, 0, 0, 0.05)';
    } else {
      sx.cursor = 'pointer';
    }

    return {
      sx,
      hover: false,
      onClick: () => {
        if (row.original.securityType) {
          let secType;
          switch (row.original.securityType) {
            case SecurityType.BOND:
              secType = 'bonds';
              break;
            case SecurityType.SHARE:
              secType = 'shares';
              break;
            default:
              secType = 'shares';
          }

          router.push(`/${secType}/moex/${row.original.ticker}`);
        }
      },
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
    renderDetailPanel: ({ row }) =>
      row.original.comment || row.original.opinions.length > 0
        ? row.original.comment
        : null,
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
