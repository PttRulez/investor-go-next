interface SxProp {
  [key: string]: string | number;
}

interface SelectList {
  [key: number | string]: string;
}

interface SelectOption {
  id: number | string;
  name: string | number;

  [key: string]: string | number;
}
