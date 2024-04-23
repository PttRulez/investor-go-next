'use client';
import { SyntheticEvent, useMemo, useState } from 'react';
import {
  MaterialReactTable,
  type MRT_ColumnDef,
  useMaterialReactTable,
} from 'material-react-table';
import PortfolioTableFooter from './PortfolioTableFooter';
import { getDefaultMRTOptions } from '@/utils/mrt-default-options';
import PortfolioTableToolbar from './PortfolioTableToolbar';
import {
  Dialog,
  IconButton,
  SelectChangeEvent,
  Tooltip,
  Typography,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import PositionForm from '@/app/position/components/PositionForm';
import { LightTooltip } from '@pttrulez/mui-based-ui';
import { IPortfolioResponse, IPositionResponse } from '@/types/apis/go-api';
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

  const columns = useMemo<Array<MRT_ColumnDef<IPositionResponse>>>(
    () => [
      {
        header: 'Тип',
        Header: <></>,
        accessorKey: 'security.securityType',
        accessorFn: p =>
          p.security.securityType === SecurityType.SHARE
            ? 'Акция'
            : 'Облигация',
        size: 5,
        enableSorting: false,
        enableColumnActions: false,
      },
      {
        header: 'Тикер',
        accessorKey: 'security.ticker',
        size: 5,
      },
      {
        header: 'Название',
        accessorKey: 'security.name',
        size: 50,
      },
      {
        header: 'Кол-во',
        accessorKey: 'amount',
        size: 5,
      },
      {
        header: 'Цена',
        accessorKey: 'currentPrice',
        size: 5,
      },
      {
        header: 'Стоимость',
        accessorFn: position =>
          position.total.toLocaleString('ru-RU', { maximumFractionDigits: 0 }),
        size: 5,
      },
      {
        header: 'Комент',
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
        accessorKey: 'total',
        Cell: ({ row }) => {
          return (
            <IconButton
              onClick={(e: SyntheticEvent) => {
                e.stopPropagation();
                setPositionToEdit(row.original);
              }}
            >
              <EditIcon />
            </IconButton>
          );
        },
      },
    ],
    [],
  );

  const table = useMaterialReactTable<IPositionResponse>({
    // ...defaultMRTOptions,
    columns,
    data: portfolio?.positions?.allPositions ?? [],
    enableColumnActions: false,
    enableColumnDragging: false,
    enableGrouping: true,
    // enableToolbarInternalActions: false,
    renderToolbarInternalActions: props => (
      <PortfolioTableToolbar
        onChooseTransaction={e => {
          onChooseTransaction(e);
        }}
      />
    ),
    icons: { SortIcon: <></> },
    muiTableHeadRowProps: {
      sx: {
        '.Mui-TableHeadCell-Content-Wrapper': { whiteSpace: 'pre-wrap' },
        '.MuiTableSortLabel-icon': {
          display: 'none',
        },
      },
    },
    muiTableHeadCellProps: {
      align: 'center',
    },
    muiTableBodyCellProps: {
      align: 'center',
    },
    initialState: { expanded: true, grouping: ['security.securityType'] },
    // renderColumnActionsMenuItems: () => [<></>],
    muiTablePaperProps: { sx: { marginBottom: '100px' } },
    // enableRowPinning: false,
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
