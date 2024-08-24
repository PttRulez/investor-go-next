import NextAuthProvider from './next-auth';
import ReactQueryProvider from './react-query';
import { ThemeProvider } from '@mui/material/styles';
import theme from '../../theme';

const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <NextAuthProvider>
      <ReactQueryProvider>
        <ThemeProvider theme={theme}>{children}</ThemeProvider>
      </ReactQueryProvider>
    </NextAuthProvider>
  );
};

export default Providers;
