import { AdvancedTable, AdvancedTableColumn } from '@pttrulez/mui-based-ui';
import dayjs from '@/dayjs.config';
import { IconButton } from '@mui/material';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import investorService from '@/axios/investor/investor.service';
import CloseIcon from '@mui/icons-material/Close';
import ArrowCircleUpIcon from '@mui/icons-material/ArrowCircleUp';
import ArrowCircleDownIcon from '@mui/icons-material/ArrowCircleDown';
import { IDealResponse } from '@/types/apis/go-api';
import { DealType } from '@/types/enums';

type Props = {
  deals: IDealResponse[];
};

const DealsTable = ({ deals }: Props) => {
  const client = useQueryClient();
  const deleteDeal = useMutation(
    (dealId: number) => investorService.deal.deleteDeal(dealId),
    {
      onSuccess: deal => {
        client.invalidateQueries({ queryKey: ['portfolio', deal.portfolioId] });
      },
    },
  );

  const columns: AdvancedTableColumn<IDealResponse>[] = [
    {
      label: 'Тикер',
      name: 'ticker',
    },
    {
      label: 'Тип',
      name: 'type',
      render: (value: DealType) =>
        value === DealType.BUY ? (
          <ArrowCircleUpIcon sx={{ color: 'green' }} />
        ) : (
          <ArrowCircleDownIcon sx={{ color: 'red' }} />
        ),
    },
    {
      label: 'Кол-во',
      name: 'amount',
    },
    {
      label: 'Цена',
      name: 'price',
    },
    {
      label: 'Сумма',
      name: 'total',
      format: (_, deal) => deal.price * deal.amount,
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

  return <AdvancedTable rows={deals} columns={columns} />;
};

export default DealsTable;
