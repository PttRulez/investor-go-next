import { AdvancedTable, AdvancedTableColumn } from '@pttrulez/mui-based-ui';
import dayjs from '@/dayjs.config';
import { IconButton } from '@mui/material';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import investorService from '@/axios/investor/investor.service';
import CloseIcon from '@mui/icons-material/Close';
import ArrowCircleUpIcon from '@mui/icons-material/ArrowCircleUp';
import ArrowCircleDownIcon from '@mui/icons-material/ArrowCircleDown';
import { IDealResponse } from '@/types/apis/go-api';
import { DealType, SecurityType } from '@/types/enums';

type Props = {
  deals: IDealResponse[];
  portfolioId: number;
};

const DealsTable = ({ deals, portfolioId }: Props) => {
  const client = useQueryClient();
  const deleteDeal = useMutation(
    (dealId: number) => investorService.deal.deleteDeal(dealId),
    {
      onSuccess: _ => {
        client.invalidateQueries({ queryKey: ['portfolio', portfolioId] });
      },
    },
  );

  const columns: AdvancedTableColumn<IDealResponse>[] = [
    {
      label: 'Название',
      name: 'shortName',
      render: (value, row) => value,
    },
    {
      label: 'Тип',
      name: 'type',
      render: (value: DealType) =>
        value === DealType.BUY ? (
          <ArrowCircleUpIcon sx={{ color: 'success.main' }} />
        ) : (
          <ArrowCircleDownIcon sx={{ color: 'error.main' }} />
        ),
    },
    {
      label: 'Кол-во',
      name: 'amount',
    },
    {
      label: 'Цена',
      name: 'price',
      format: price => price.toLocaleString('RU-ru'),
    },
    {
      label: 'Сумма',
      name: 'total',
      format: (_, deal) => (deal.price * deal.amount).toLocaleString('RU-ru'),
    },
    {
      label: 'НКД',
      name: 'nkd',
    },
    {
      label: 'Дата',
      name: 'date',
      format: value => dayjs(value).format('DD MMMM YYYY'),
    },
    {
      label: '',
      name: 'actions',
      render: (_, deal) => {
        return (
          <IconButton
            onClick={(e: React.SyntheticEvent) => {
              e.stopPropagation();
              deleteDeal.mutate(deal.id);
            }}
          >
            <CloseIcon sx={{ color: 'error.main' }} />
          </IconButton>
        );
      },
    },
  ];

  return (
    <AdvancedTable
      rows={deals}
      columns={columns}
      sx={{
        height: '90vh',
        padding: '20px',
      }}
    />
  );
};

export default DealsTable;
