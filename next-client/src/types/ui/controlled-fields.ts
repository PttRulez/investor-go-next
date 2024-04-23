import { Control, FieldPath, FieldValues } from 'react-hook-form';

type ValidationRules<ValueType = any, FormValuesType = Record<string, any>> = {
  required?: boolean;
  min?: number;
  max?: number;
  minLength?: number;
  maxLength?: number;
  pattern?: RegExp;
  validate?:
    | ((value: ValueType, formValues: FormValuesType) => boolean)
    | Record<string, (v?: ValueType, formValues?: FormValuesType) => boolean>;
};

export type ControlledField<T extends FieldValues = FieldValues> = {
  control: Control<T>;
  name: FieldPath<T>;
  rules?: ValidationRules;
};
