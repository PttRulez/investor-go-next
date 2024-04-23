import NextAuthProvider from './next-auth';
import ReactQueryProvider from './react-query';

const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <NextAuthProvider>
      <ReactQueryProvider>{children}</ReactQueryProvider>
    </NextAuthProvider>
  );
};

export default Providers;
