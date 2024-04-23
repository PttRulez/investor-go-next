'use client';
import { Box, Dialog, IconButton, Typography } from '@mui/material';
import {
  MRT_ColumnDef,
  MaterialReactTable,
  useMaterialReactTable,
} from 'material-react-table';
import { SyntheticEvent, useMemo, useState } from 'react';
import RemoveRedEyeIcon from '@mui/icons-material/RemoveRedEye';
import { OpinionType } from '@/types/enums';
import { IOpinionResponse } from '@/types/apis/go-api';

type Props = {
  opinions: IOpinionResponse[];
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

const OpinionsTable = ({ opinions }: Props) => {
  const [opinionText, setOpinionText] = useState<null | string>(null);

  const columns = useMemo<Array<MRT_ColumnDef<IOpinionResponse>>>(
    () => [
      {
        header: 'Дата',
        accessorKey: 'date',
      },
      {
        header: 'Автор',
        accessorKey: 'expert.name',
      },
      {
        header: 'Прогноз?',
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
  });

  return (
    <>
      <MaterialReactTable table={table} />
      <Dialog open={!!opinionText} onClose={() => setOpinionText(null)}>
        <Box>{opinionText}</Box>
      </Dialog>
    </>
  );
};

// export default OpinionsTable;
