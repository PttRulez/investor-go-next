import {
  AdvancedTable,
  type AdvancedTableColumn,
} from '@pttrulez/mui-based-ui';
import dayjs from '@/dayjs.config';
import { IconButton } from '@mui/material';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import investorService from '@/axios/investor/investor.service';
import CloseIcon from '@mui/icons-material/Close';
import ArrowCircleUpIcon from '@mui/icons-material/ArrowCircleUp';
import ArrowCircleDownIcon from '@mui/icons-material/ArrowCircleDown';
import ShoppingCartCheckoutIcon from '@mui/icons-material/ShoppingCartCheckout';
import AddShoppingCartIcon from '@mui/icons-material/AddShoppingCart';
import { ITransactionResponse } from '@/types/apis/go-api';
import { TransactionType } from '@/types/enums';

type Props = {
  portfolioId: number;
  transactions: ITransactionResponse[];
};

const TransactionsTable = ({ portfolioId, transactions }: Props) => {
  const client = useQueryClient();

  const deleteTransaction = useMutation(
    (transaction: ITransactionResponse): Promise<ITransactionResponse> => {
      return investorService.transaction.deleteTransaction(transaction.id);
    },
    {
      onSuccess: _ => {
        client.invalidateQueries({
          queryKey: ['portfolio', portfolioId],
        });
      },
    },
  );

  const columns: AdvancedTableColumn<ITransactionResponse>[] = [
    {
      label: 'Тип',
      name: 'type',
      render: (value: TransactionType) =>
        value === TransactionType.CASHOUT ? (
          <ShoppingCartCheckoutIcon sx={{ color: 'error.main' }} />
        ) : (
          <AddShoppingCartIcon sx={{ color: 'green' }} />
        ),
    },
    {
      label: 'Сумма',
      name: 'amount',
      render: (value: number) => value.toLocaleString('RU-ru'),
    },
    {
      label: 'Дата',
      name: 'date',
      format: value => dayjs(value).format('DD MMMM YYYY'),
    },
    {
      label: '',
      name: 'actions',
      render: (_, transaction) => {
        return (
          <IconButton
            onClick={(e: React.SyntheticEvent) => {
              e.stopPropagation();
              deleteTransaction.mutate(transaction);
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
      rows={transactions}
      columns={columns}
      sx={{
        padding: '20px 100px',
      }}
    />
  );
};

export default TransactionsTable;
