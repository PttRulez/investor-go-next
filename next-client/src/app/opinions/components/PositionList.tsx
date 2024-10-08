import { AttachInfo } from '@/axios/investor/domains';
import investorService from '@/axios/investor/investor.service';
import { Box, IconButton, Paper, Typography } from '@mui/material';
import Grid from '@mui/material/Unstable_Grid2/Grid2';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import LockOpenIcon from '@mui/icons-material/LockOpen';
import LockIcon from '@mui/icons-material/Lock';
import { SyntheticEvent } from 'react';
import { AxiosError } from 'axios';

type Props = {
  opinionId: number;
  ticker: string;
};

export default function PositionList({ opinionId, ticker }: Props) {
  const queryClient = useQueryClient();
  const { data: positions } = useQuery({
    queryKey: ['positions'],
    queryFn: () => investorService.position.getAllUserPosition(),
  });

  const attachOpinion = useMutation(
    (i: AttachInfo) => investorService.opinion.attachToPosition(i),
    {
      onSuccess(data: any, variables: any, context: any) {
        queryClient.invalidateQueries({ queryKey: ['positions'] });
      },
    },
  );

  return (
    <>
      {positions
        ? positions
            .filter(p => p.ticker == ticker)
            .map(p => {
              return (
                <Paper
                  sx={{
                    marginBottom: '10px',
                  }}
                  key={p.id}
                >
                  <Grid
                    container
                    columns={12}
                    justifyContent="space-between"
                    alignItems={'center'}
                    direction={'row'}
                  >
                    <Grid xs={9} spacing={2}>
                      <Typography sx={{ width: '100%' }}>
                        {`${p.portfolioName} (${p.amount})`}
                      </Typography>
                    </Grid>
                    <Grid xs={3} spacing={2}>
                      {p.opinionIds.includes(opinionId) ? (
                        <IconButton
                          onClick={(e: SyntheticEvent) => {
                            e.stopPropagation();
                            attachOpinion.mutate({
                              opinionId,
                              positionId: p.id,
                            });
                          }}
                        >
                          <LockIcon sx={{ color: 'green' }} />
                        </IconButton>
                      ) : (
                        <IconButton
                          onClick={(e: SyntheticEvent) => {
                            e.stopPropagation();
                            attachOpinion.mutate({
                              opinionId,
                              positionId: p.id,
                            });
                          }}
                        >
                          <LockOpenIcon sx={{ color: 'primary.dark' }} />
                        </IconButton>
                      )}
                    </Grid>
                  </Grid>
                </Paper>
              );
            })
        : 'Нет позиций'}
    </>
  );
}
