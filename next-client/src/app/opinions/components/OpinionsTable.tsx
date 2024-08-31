'use client';
import { Box, Dialog, IconButton, Typography } from '@mui/material';
import {
  MRT_ColumnDef,
  MaterialReactTable,
  useMaterialReactTable,
} from 'material-react-table';
import { SyntheticEvent, useMemo, useState } from 'react';
import RemoveRedEyeIcon from '@mui/icons-material/RemoveRedEye';
import BackHandIcon from '@mui/icons-material/BackHand';
import { useQuery } from '@tanstack/react-query';
import investorService from '@/axios/investor/investor.service';
import dayjs from '@/dayjs.config';
import { OpinionType } from '@/types/enums';
import { OpinionFilters } from '@/validation';
import { IOpinionResponse } from '@/types/apis/go-api';
import Grid from '@mui/material/Unstable_Grid2/Grid2';
import PositionList from './PositionList';

type Props = {
  opinions: IOpinionResponse[];
  // filters: OpinionFilters;
  ticker: string;
};

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

const OpinionsTable = ({ opinions, ticker }: Props) => {
  const [opinionText, setOpinionText] = useState<null | string>(null);
  const [opinionIdToAttach, setOpinionIdToAttach] = useState<null | number>(
    null,
  );

  // const { data: opinions } = useQuery({
  //   queryKey: [
  //     'opinions',
  //     {
  //       exchange: filters.exchange,
  //       securityId: filters.securityId,
  //       securityType: filters.securityType,
  //     },
  //   ],
  //   queryFn: () => investorService.opinion.getOpinionsList(filters),
  // });

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
        accessorFn(opinion) {
          return opinion.targetPrice ?? '-';
        },
      },
      {
        header: '',
        accessorKey: 'text',
        Cell: ({ cell, row }) => {
          return (
            <Grid container spacing={3} justifyContent="space-between">
              <Grid xs={6}>
                <IconButton
                  onClick={(e: SyntheticEvent) => {
                    e.stopPropagation();
                    setOpinionText(cell.getValue<string>());
                  }}
                >
                  <RemoveRedEyeIcon />
                </IconButton>
              </Grid>
              <Grid xs={6}>
                <IconButton
                  onClick={(e: SyntheticEvent) => {
                    e.stopPropagation();
                    setOpinionIdToAttach(row.original.id);
                  }}
                >
                  <BackHandIcon />
                </IconButton>
              </Grid>
            </Grid>
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
      <Dialog
        open={!!opinionIdToAttach}
        onClose={() => setOpinionIdToAttach(null)}
        sx={{
          '& .MuiPaper-root': {
            padding: '20px',
          },
        }}
      >
        {opinionIdToAttach && (
          <PositionList opinionId={opinionIdToAttach} ticker={ticker} />
        )}
      </Dialog>
    </>
  );
};

export default OpinionsTable;
