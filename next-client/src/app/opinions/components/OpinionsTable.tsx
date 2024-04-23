'use client';
import { Box, Dialog, IconButton, Typography } from '@mui/material';
import {
  MRT_ColumnDef,
  MaterialReactTable,
  useMaterialReactTable,
} from 'material-react-table';
import { SyntheticEvent, useMemo, useState } from 'react';
import RemoveRedEyeIcon from '@mui/icons-material/RemoveRedEye';
import { useQuery } from '@tanstack/react-query';
import investorService from '@/axios/investor/investor.service';
import dayjs from '@/dayjs.config';
import { OpinionType } from '@/types/enums';
import { OpinionFilters } from '@/validation';
import { IOpinionResponse } from '@/types/apis/go-api';

const opinionTypeText = (type: OpinionType) => {
  switch (type) {
    case OpinionType.FLAT:
      return { text: 'Флэт', color: 'grey.700' };
    case OpinionType.GROWTH:
      return { text: 'Рост', color: 'green' };
    case OpinionType.REDUCTION:
      return { text: 'Падение', color: 'red' };
    case OpinionType.GENERAL:
      return { text: 'Информация', color: 'blue' };
  }
};

const OpinionsTable = ({ filters }: { filters: OpinionFilters }) => {
  const [opinionText, setOpinionText] = useState<null | string>(null);

  const { data: opinions } = useQuery({
    queryKey: [
      'opinions',
      {
        exchange: filters.exchange,
        securityId: filters.securityId,
        securityType: filters.securityType,
      },
    ],
    queryFn: () => investorService.opinion.getOpinionsList(filters),
  });

  const columns = useMemo<Array<MRT_ColumnDef<IOpinionResponse>>>(
    () => [
      {
        header: 'Дата',
        accessorKey: 'date',
        accessorFn: row => dayjs(row.date).format('DD MMMM YYYY'),
        enableSorting: false,
        enableColumnActions: false,
      },
      {
        header: 'Автор',
        accessorKey: 'expert.name',
      },
      {
        header: 'Прогноз',
        id: 'type',
        accessorFn: row => {
          const { text, color } = opinionTypeText(row.type);

          return (
            <Typography variant="body1" sx={{ color }}>
              {text}
            </Typography>
          );
        },
      },
      {
        header: 'Целевая',
        accessorKey: 'targetPrice',
      },
      {
        header: '',
        accessorKey: 'text',
        Cell: ({ cell }) => {
          return (
            <IconButton
              onClick={(e: SyntheticEvent) => {
                e.stopPropagation();
                setOpinionText(cell.getValue<string>());
              }}
            >
              <RemoveRedEyeIcon />
            </IconButton>
          );
        },
      },
    ],
    [],
  );
  const table = useMaterialReactTable<IOpinionResponse>({
    columns,
    data: opinions ?? [],
    enableBottomToolbar: false,
    enableRowActions: false,
    enableSorting: false,
    enableColumnActions: false,
    enableToolbarInternalActions: false,
    renderTopToolbar: false,
    // renderToolbarInternalActions: _ => <div>toolbar</div>,
    // renderBottomToolbar: _ => <div>bottom_toolbar</div>,
  });

  return (
    <>
      <MaterialReactTable table={table} />
      <Dialog open={!!opinionText} onClose={() => setOpinionText(null)}>
        <Box sx={{ padding: '10px' }}>{opinionText}</Box>
      </Dialog>
    </>
  );
};

export default OpinionsTable;
