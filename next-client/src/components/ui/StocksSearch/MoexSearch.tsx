'use client';
import { useEffect, useState } from 'react';
import TextField, { TextFieldProps } from '@mui/material/TextField';
import Autocomplete, {
  AutocompleteInputChangeReason,
  AutocompleteProps,
} from '@mui/material/Autocomplete';
import useDebounce from '@/hooks/useDebounce';
import { useQuery } from '@tanstack/react-query';
import { moexService } from '@/axios/moex/moex.service';
import { Box, Typography } from '@mui/material';
import { moexStockTypesRU } from '@/constants/moex';
import { MoexSearchAutocompleteOption } from './types';
import { AxiosError } from 'axios';
import { IMoexISSSearchResults } from '@/types/apis/moex-iss-api';

type MoexSearchProps = Omit<
  AutocompleteProps<MoexSearchAutocompleteOption, false, false, false>,
  'renderInput' | 'options'
> &
  Pick<TextFieldProps, 'helperText' | 'error'>;

const MoexSearch = ({
  error,
  helperText,
  ...autocompleteProps
}: MoexSearchProps) => {
  const [submittedTicker, setSubmittedTicker] = useState<string>('');
  const debouncedValue = useDebounce<string>(submittedTicker, 500);

  const { data: stocksOptions, refetch } = useQuery<
    IMoexISSSearchResults,
    AxiosError,
    MoexSearchAutocompleteOption[],
    [string, string]
  >({
    queryKey: ['search', debouncedValue],
    queryFn: ({ queryKey }) => moexService.search(queryKey[1]),
    enabled: false,
    select: data => {
      return data.securities.data.map(
        (sec, index) =>
          ({
            board: sec[14],
            group: sec[13],
            jsxKey: `${sec[0]} ${sec[1]} ${index}`,
            name: sec[4],
            shortName: sec[2],
            ticker: sec[1],
            type: sec[12],
          }) satisfies MoexSearchAutocompleteOption,
      );
    },
  });

  useEffect(() => {
    if (debouncedValue) {
      refetch();
    }
  }, [debouncedValue, refetch]);

  const inputHandler = (
    e: React.SyntheticEvent,
    value: string,
    reason: AutocompleteInputChangeReason,
  ) => {
    if (reason === 'input' && value) {
      setSubmittedTicker(value);
    }
  };

  const typesUnique = new Set();

  if (stocksOptions) {
    for (const sec of stocksOptions) {
      typesUnique.add(sec.group);
    }
  }

  return (
    <Autocomplete
      {...autocompleteProps}
      disablePortal
      id="combo-box-demo"
      options={stocksOptions ?? []}
      onInputChange={inputHandler}
      getOptionLabel={option => option.shortName ?? option.name ?? ''}
      renderOption={(props, option) => {
        return (
          <Box component="li" {...props} key={option.jsxKey}>
            {option.shortName}{' '}
            <Typography
              variant="body1"
              sx={{ fontStyle: 'italic', color: 'grey.500' }}
            >
              - {moexStockTypesRU[option.type]}
            </Typography>
          </Box>
        );
      }}
      sx={{ minWidth: 500 }}
      // @ts-ignore
      renderInput={params => (
        <TextField
          {...params}
          label={'Название бумаги'}
          error={error}
          helperText={helperText}
        />
      )}
    />
  );
};

export default MoexSearch;
