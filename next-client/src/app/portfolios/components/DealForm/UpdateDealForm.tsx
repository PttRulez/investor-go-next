// 'use client';

// import { Box, Button, Stack } from '@mui/material';
// import FormText from '@/components/ui/Forms/FormText';
// import { useForm } from 'react-hook-form';
// import { FC } from 'react';
// import { useMutation, useQueryClient } from '@tanstack/react-query';
// import MoexSearch from '@/components/ui/StocksSearch/MoexSearch';
// import FormDatePicker from '@/components/ui/Forms/FormDatePicker';
// import investorService from '@/axios/investor';
// import { moexStockTypeToGeneralType } from '@/constants/moex';
// import { UpdateDealDto, IDealResponse  } from 'contracts';

// interface DealFormProps {
//   deal: IDealResponse;
//   afterSuccessfulSubmit: () => void;
// }

// const UpdateDealForm: FC<DealFormProps> = ({ deal, afterSuccessfulSubmit }) => {
//   const { control, handleSubmit, setValue, watch } = useForm<IDealResponse>({
//     defaultValues: deal,
//   });
//   const watchAll = watch();
//   const client = useQueryClient();

//   const changeHandler = (e, value: Deal, _) => {
//     // if (value) {
//     //   setValue('ticker', value.id);
//     //   setValue('name', value.label);
//     //   setValue('exchangeName', 'MOEX');
//     //   setValue('stockType', moexStockTypeToGeneralType[value.type]);
//     //   setValue('board', value.primary_boardid);
//     // }
//   };

//   const updateDeal = useMutation((formData: UpdateDealDto) => investorService.updateDeal(formData), {
//     onSuccess: () => {
//       afterSuccessfulSubmit();
//       client.invalidateQueries(['portfolio', +deal.portfolioId]);
//     },
//   });

//   const onSubmit = (data: IDealResponse) => {
//     const dataToSend: UpdateDealDto = {} as UpdateDealDto;
//     updateDeal.mutate(dataToSend);
//   };

//   return (
//     <Box onSubmit={handleSubmit(onSubmit)} component="form" sx={{ padding: '10px', minHeight: '700px' }}>
//       <Stack gap={'20px'}>
//         <MoexSearch
//           //@ts-ignore
//           label={'Бумага'}
//           control={control}
//           name={'ticker'}
//           onChange={changeHandler}
//         />
//         <FormText
//           type="number"
//           //@ts-ignore
//           control={control}
//           label={'Кол-во бумаг'}
//           name={'amount'}
//           value={watchAll.amount}
//         />
//         <FormText
//           type="number"
//           //@ts-ignore
//           control={control}
//           label={'Цена покупки'}
//           name={'price'}
//           value={watchAll.price}
//         />
//         <FormDatePicker
//           //@ts-ignore
//           control={control}
//           name={'date'}
//           handleClear={() => setValue('date', '')}
//           onChange={newValue => setValue('date', newValue)}
//           label={'Дата покупки'}
//           value={watchAll.date}
//         />
//         <Button variant="outlined" color="primary" type="submit" sx={{ color: 'grey.700' }}>
//           Сохранить
//         </Button>
//       </Stack>
//     </Box>
//   );
// };

// export default UpdateDealForm;
