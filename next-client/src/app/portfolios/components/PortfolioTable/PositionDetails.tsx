import OpinionsTable from '@/app/opinions/components/OpinionsTable';
import { IPositionResponse } from '@/types/apis/go-api';
import { Typography } from '@mui/material';

type Props = {
  position: IPositionResponse;
};

const PositionDetails = ({ position: p }: Props) => {
  return (
    <>
      <Typography variant="body1">{p.comment}</Typography>
      <OpinionsTable opinions={p.opinions} ticker={p.ticker} />
    </>
  );
};

export default PositionDetails;
