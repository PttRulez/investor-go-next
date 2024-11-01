import OpinionsTable from '@/app/opinions/components/OpinionsTable';
import { IPositionResponse } from '@/types/apis/go-api';
import { Box, Typography } from '@mui/material';

type Props = {
  position: IPositionResponse;
};

const PositionDetails = ({ position: p }: Props) => {
  return (
    <>
      {p.comment && (
        <Box
          sx={{
            marginBottom: '50px',
            padding: '20px',
            border: '1px primary.dark',
            boxShadow: '0px 0px 5px 0px rgba(0,0,0,0.4)',
          }}
        >
          {p.comment}
        </Box>
      )}
      {p.opinions.length > 0 && <OpinionsTable opinions={p.opinions} ticker={p.ticker} />}
    </>
  );
};

export default PositionDetails;
