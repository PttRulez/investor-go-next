'use client';

import { useState } from 'react';
import investorService from '@/axios/investor/investor.service';
import DealsTable from '@/app/portfolios/components/DealsTable';
import PortfolioTable from '../components/PortfolioTable/PortfolioTable';
import { Dialog, SelectChangeEvent, Tab } from '@mui/material';
import CreateDealForm from '../components/DealForm/CreateDealForm';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { TabContext, TabList, TabPanel } from '@mui/lab';
import { PortfolioActionsMap } from '../components/PortfolioTable/PortfolioTableToolbar';
import TransactionForm from '../components/TransactionForm/TransactionForm';
import { useRouter } from 'next/navigation';
import { DealType } from '@/types/enums';

export default function Portfolio({ params }: { params: { id: string } }) {
  const [openModal, setOpenModal] = useState<PortfolioActionsMap | false>(
    false,
  );
  const [activeTab, setActiveTab] = useState<string>('portfolio');
  const router = useRouter();

  const client = useQueryClient();
  const { data: portfolio } = useQuery({
    queryKey: ['portfolio', parseInt(params.id)],
    queryFn: async () => {
      const portfolio = await investorService.portfolio.getPortfolio(params.id);
      if (!portfolio) {
        // router.push(`/portfolios`);
        console.log('portfolio', portfolio);
      } else {
        return {
          ...portfolio,
        };
      }
    },
    onSuccess: () => {
      client.invalidateQueries(['prices']);
    },
    // initialData: initPortfolio,
  });

  const chooseTransactionHandler = (e: SelectChangeEvent) => {
    switch (e.target.value) {
      case PortfolioActionsMap.buy:
        setOpenModal(PortfolioActionsMap.buy);
        break;
      case PortfolioActionsMap.transaction:
        setOpenModal(PortfolioActionsMap.transaction);
        break;
      case PortfolioActionsMap.sell:
        setOpenModal(PortfolioActionsMap.sell);
        break;
    }
  };

  return (
    <>
      <TabContext value={activeTab}>
        <TabList onChange={(_, v) => setActiveTab(v)}>
          <Tab label="Портфолио" value="portfolio"></Tab>
          <Tab label="Сделки" value="deals"></Tab>
          <Tab label="Транзакции" value="transactions"></Tab>
        </TabList>
        <TabPanel value="portfolio">
          {portfolio && (
            <PortfolioTable
              onChooseTransaction={chooseTransactionHandler}
              portfolio={portfolio}
            />
          )}
        </TabPanel>
        {/* <TabPanel value="deals">
          {portfolio && <DealsTable deals={portfolio.deals ?? []} />}
        </TabPanel>
        <TabPanel value="transactions">
          {portfolio && (
            <TransactionsTable
              portfolioId={portfolio.id}
              transactions={portfolio.transactions ?? []}
            />
          )}
        </TabPanel> */}
      </TabContext>

      {/* <AddNewButton onClick={() => setCreateDeal(true)} /> */}

      <Dialog open={!!openModal} onClose={() => setOpenModal(false)}>
        {openModal &&
          {
            [PortfolioActionsMap.buy]: (
              <CreateDealForm
                afterSuccessfulSubmit={() => setOpenModal(false)}
                dealType={DealType.BUY}
                portfolioId={Number(params.id)}
              />
            ),
            [PortfolioActionsMap.transaction]: (
              <TransactionForm
                afterSuccessfulSubmit={() => setOpenModal(false)}
                portfolioId={Number(params.id)}
              />
            ),
            [PortfolioActionsMap.sell]: (
              <CreateDealForm
                afterSuccessfulSubmit={() => setOpenModal(false)}
                dealType={DealType.SELL}
                portfolioId={Number(params.id)}
              />
            ),
          }[openModal]}
      </Dialog>
    </>
  );
}
