import { MoexSecurityGroup } from '@/types/apis/go-api';
import { MoexBoard, MoexSecurityType } from '@/types/enums';
import {
  AutocompleteChangeDetails,
  AutocompleteChangeReason,
} from '@mui/material';

export type MoexSearchHandler = (
  event: React.SyntheticEvent,
  value: MoexSearchAutocompleteOption | null, //Record<string, any> | null | Array<Record<string, any> | null>,
  reason: AutocompleteChangeReason,
  details?: AutocompleteChangeDetails<Record<string, any>>,
) => void;

// onChange?: (
//     event: React.SyntheticEvent,
//     value: AutocompleteValue<Value, Multiple, DisableClearable, FreeSolo>,
//     reason: AutocompleteChangeReason,
//     details?: AutocompleteChangeDetails<Value>,
//   ) => void;

export type MoexSearchAutocompleteOption = {
  board: MoexBoard;
  group: MoexSecurityGroup;
  jsxKey: string;
  name: string;
  shortName: string;
  ticker: string;
  type: MoexSecurityType;
};
