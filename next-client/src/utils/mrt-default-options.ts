import { type MRT_RowData, type MRT_TableOptions } from 'material-react-table';

//define re-useable default table options for all tables in your app
export const getDefaultMRTOptions = <TData extends MRT_RowData>(): Partial<
  MRT_TableOptions<TData>
> => ({
  //list all of your default table options here
  enableGlobalFilter: false,
  enableRowPinning: false,
  initialState: { showColumnFilters: true },
  manualFiltering: false,
  manualPagination: true,
  manualSorting: true,
  muiTableHeadCellProps: {
    sx: { fontSize: '1.1rem' },
  },
  paginationDisplayMode: 'pages',
  enablePagination: false,
  enableColumnActions: false,
  enableSorting: false,
  enableGrouping: true,
  enableColumnOrdering: false,
  //etc...
  defaultColumn: {
    enableColumnOrdering: false,
    //you can even list default column options here
  },
});
