'use client';
import { Roboto } from 'next/font/google';
import { createTheme } from '@mui/material/styles';

const roboto = Roboto({
  weight: ['300', '400', '500', '700'],
  subsets: ['latin'],
  display: 'swap',
});

const theme = createTheme({
  components: {
    MuiTextField: {
      styleOverrides: {
        root: ({ theme }) => ({
          '.MuiInputLabel-formControl.Mui-focused': {
            color: theme.palette.primary.contrastText,
          },
          '& .MuiInput-underline:after': {
            borderBottomColor: theme.palette.primary.contrastText,
          },
        }),
      },
    },
    MuiButton: {
      styleOverrides: {
        root: ({ theme }) => ({
          '&.MuiButton-outlined.MuiButton-colorPrimary': {
            color: theme.palette.primary.contrastText,
            borderColor: theme.palette.primary.contrastText,
            ':hover': {
              color: theme.palette.primary.contrastText,
              borderColor: theme.palette.primary.contrastText,
              backgroundColor: theme.palette.primary.main,
            },
          },
        }),
      },
    },
  },
  palette: {
    primary: {
      light: '#eeeeee',
      main: '#eeeeee',
      dark: '#bdbdbd', //
      contrastText: '#424242',
    },
    secondary: {
      light: '#ff7961',
      main: '#f44336',
      dark: '#ba000d',
      contrastText: '#000',
    },
    success: {
      light: '#5cb85c',
      main: '#53a653',
      dark: '#4a934a',
      contrastText: '#376e37',
    },
    error: {
      light: '#ff7961',
      main: '#f44336',
      dark: '#ba000d',
      contrastText: '#000',
    },
    myAwesomeColor: {
      light: 'purple',
    },
  },
  typography: {
    fontFamily: roboto.style.fontFamily,
  },
});

declare module '@mui/material/styles' {
  interface Palette {
    myAwesomeColor: {
      light?: string;
      main?: string;
      dark?: string;
      contrastText?: string;
    };
  }
  interface PaletteOptions {
    myAwesomeColor: {
      light?: string;
      main?: string;
      dark?: string;
      contrastText?: string;
    };
  }
}

export default theme;
