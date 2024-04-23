import OpinionsTable from '@/app/opinions/components/OpinionsTable';

const Expert = ({ params }: { params: { id: string } }) => {
  return (
    <>
      <OpinionsTable filters={{ expertId: +params.id }} />
    </>
  );
};

export default Expert;
