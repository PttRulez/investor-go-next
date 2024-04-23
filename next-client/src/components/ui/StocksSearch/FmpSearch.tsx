// 'use client';
// import { ChangeEvent, FC, FormEventHandler, useEffect, useState } from 'react';
// import TextField from '@mui/material/TextField';
// import Autocomplete from '@mui/material/Autocomplete';
// import { searchStock } from '@/types/fmpTypes/fmp.type';
// import { fmpService } from '@/axios/fmp/fmp.service';
// import { QueryFunctionContext, useQuery } from '@tanstack/react-query';
// import useDebounce from '@/hooks/useDebounce';
// import { useRouter } from 'next/navigation';

// interface ISearchField {
//   searchTerm: string;
//   handleSubmit?: FormEventHandler<HTMLFormElement> | undefined;
//   handleChange?: (event: ChangeEvent<HTMLInputElement>) => void;
// }

// const FmpSearch: FC = () => {
//   const router = useRouter();
//   const [submittedTicker, setSubmittedTicker] = useState<string>('');
//   const debouncedValue = useDebounce<string>(submittedTicker, 500);

//   const { data: stocksOptions, refetch } = useQuery({
//     queryKey: ['search', debouncedValue],
//     queryFn: ({ queryKey }) => fmpService.search(queryKey[1] as string),
//     enabled: false,
//   });

//   useEffect(() => {
//     if (debouncedValue) {
//       refetch();
//     }
//   }, [debouncedValue, refetch]);

//   const inputHandler = (e, value, reason) => {
//     if (reason === 'input' && value) {
//       setSubmittedTicker(value);
//     }
//   };

//   const changeHandler = (e, value, reason) => {
//     router.push(`/ticker/${value.symbol}`);
//   };

//   return (
//     <Autocomplete
//       disablePortal
//       id="combo-box-demo"
//       options={stocksOptions ?? []}
//       onInputChange={inputHandler}
//       onChange={changeHandler}
//       getOptionLabel={option =>
//         option.name ? `${option.symbol} - ${option.name}` : ''
//       }
//       sx={{ width: 300 }}
//       // @ts-ignore
//       renderInput={params => <TextField {...params} label="Ticker" />}
//     />
//   );
// };

// export default FmpSearch;
