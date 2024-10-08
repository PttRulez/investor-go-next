import { SubmitHandler, useForm } from 'react-hook-form';
import dayjs, { Dayjs } from 'dayjs';
import {
  DefaultFormBox,
  FormDatePicker,
  FormSelect,
  FormText,
} from '@pttrulez/mui-based-ui';

import { Button } from '@mui/material';
import { useMutation, useQuery } from '@tanstack/react-query';
import investorService from '@/axios/investor/investor.service';
import { useMemo } from 'react';
import Grid from '@mui/material/Unstable_Grid2/Grid2';
import MoexSearch from '@/components/ui/StocksSearch/MoexSearch';
import {
  MoexSearchAutocompleteOption,
  MoexSearchHandler,
} from '@/components/ui/StocksSearch/types';
import { getSecurityTypeFromMoexSecType } from '@/utils/helpers';
import { useQueryClient } from '@tanstack/react-query';
import { Exchange, OpinionType, SecurityType } from '@/types/enums';
import { CreateOpinionData } from '@/validation';
import { SecurityResponse } from '@/types/apis/go-api/security';
import { IMoexSecurtiyResponse } from '@/types/apis/go-api';

interface ExpertFormProps {
  afterSuccessfulSubmit: () => void;
  defaultData?: IMoexSecurtiyResponse;
}

const OpinionForm = ({
  afterSuccessfulSubmit,
  defaultData,
}: ExpertFormProps) => {
  let defaultValues: Partial<CreateOpinionData> = {
    date: dayjs().format('YYYY-MM-DD'),
    exchange: Exchange.MOEX,
  };

  if (defaultData) {
    defaultValues.securityId = defaultData.id;
    defaultValues.securityType = defaultData.securityType;
    defaultValues.ticker = defaultData.ticker;
  }

  const { control, formState, handleSubmit, resetField, setValue, watch } =
    useForm<CreateOpinionData>({
      defaultValues,
    });

  const queryClient = useQueryClient();

  const watchAll = watch();

  const { data: expertsListResponse } = useQuery({
    queryKey: ['expertList'],
    queryFn: () => investorService.expert.getExpertsList(),
  });

  const expertsList = useMemo(() => {
    if (!expertsListResponse) return [];

    return expertsListResponse.map(e => ({
      id: e.id,
      name: e.name,
    }));
  }, [expertsListResponse]);

  const createDeal = useMutation(
    (formData: CreateOpinionData) =>
      investorService.opinion.createOpinion(formData),
    {
      onSuccess: opinion => {
        queryClient.invalidateQueries({
          queryKey: [
            'opinions',
            {
              exchange: Exchange.MOEX,
              securityId: watchAll.securityId,
              securityType: watchAll.securityType,
            },
          ],
        });
        afterSuccessfulSubmit();
      },
    },
  );

  const onSubmit: SubmitHandler<CreateOpinionData> = data => {
    createDeal.mutate(data);
  };

  const defaultAutocompleteValue: MoexSearchAutocompleteOption = (
    defaultData ? { name: defaultData.shortName } : {}
  ) as MoexSearchAutocompleteOption;

  const onMoexChange: MoexSearchHandler = async (e, value, reason) => {
    if (!value) {
      resetField('securityId');
      resetField('securityType');
      return;
    }
    const secType = getSecurityTypeFromMoexSecType(value.type);

    let security: Omit<SecurityResponse, 'exchange'>;
    if (secType === SecurityType.SHARE) {
      security = await investorService.moexShare.getByTicker(value.ticker);
    } else {
      security = await investorService.moexBond.getByTicker(value.ticker);
    }

    setValue('securityId', security.id);
    setValue('securityType', security.securityType);
    setValue('ticker', value.ticker);
  };

  return (
    <DefaultFormBox onSubmit={handleSubmit(onSubmit)}>
      <MoexSearch
        onChange={onMoexChange}
        defaultValue={defaultAutocompleteValue}
      />
      <Grid container spacing={3} justifyContent="space-between">
        <Grid xs={6}>
          <FormDatePicker
            control={control}
            handleClear={() => resetField('date')}
            onChange={(newValue: Dayjs | null) => {
              if (newValue) {
                setValue('date', newValue?.format('YYYY-MM-DD'));
              } else {
                resetField('date');
              }
            }}
            label={'Дата'}
            name={'date'}
            value={watchAll.date}
          />
        </Grid>
        <Grid xs={6}>
          <FormSelect
            control={control}
            name={'expertId'}
            label="Эксперт"
            options={expertsList}
            value={watchAll.expertId}
            variant="outlined"
          />
        </Grid>
      </Grid>

      <FormText
        control={control}
        error={!!formState.errors.sourceLink}
        handleClear={() => setValue('sourceLink', '')}
        helperText={formState.errors.sourceLink?.message}
        label={'Ссылка на источник'}
        name={'sourceLink'}
        value={watchAll.sourceLink}
        multiline
      />

      <Grid container spacing={3} justifyContent="space-between">
        <Grid xs={6}>
          <FormSelect
            control={control}
            name={'type'}
            label="Прогноз"
            options={{
              [OpinionType.FLAT]: 'Флэт',
              [OpinionType.GENERAL]: 'Без прогноза',
              [OpinionType.GROWTH]: 'Рост',
              [OpinionType.REDUCTION]: 'Снижение',
            }}
            value={watchAll.type}
          />
        </Grid>
        <Grid xs={6}>
          <FormText
            control={control}
            decimal
            error={!!formState.errors.targetPrice}
            handleClear={() => resetField('targetPrice')}
            helperText={formState.errors.targetPrice?.message}
            label={'Целевая цена'}
            name={'targetPrice'}
            onChange={(e: any) => {
              if (e.target.value != '') {
                setValue('targetPrice', parseFloat(e.target.value));
              }
            }}
            type="number"
            value={watchAll.targetPrice}
          />
        </Grid>
      </Grid>

      <FormText
        control={control}
        error={!!formState.errors.text}
        handleClear={() => setValue('text', '')}
        helperText={formState.errors.text?.message}
        label={'Текст мнения'}
        multiline
        name={'text'}
        value={watchAll.text}
        variant="outlined"
      />
      <Button
        variant="outlined"
        color="primary"
        type="submit"
        sx={{ color: 'grey.700' }}
      >
        Сохранить
      </Button>
    </DefaultFormBox>
  );
};

export default OpinionForm;
