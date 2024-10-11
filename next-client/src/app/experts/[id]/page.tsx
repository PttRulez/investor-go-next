import OpinionsTable from '@/app/opinions/components/OpinionsTable';

const Expert = ({ params }: { params: { id: string } }) => {
  return (
    <>
      <h1>Expert Page</h1>
      {/* <OpinionsTable filters={{ expertId: +params.id }} /> */}
    </>
  );
};

export default Expert;
