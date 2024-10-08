'use client';

import investorService from '@/axios/investor/investor.service';
import { DealType, SecurityType } from '@/types/enums';
import { TabContext, TabList, TabPanel } from '@mui/lab';
import { Dialog, SelectChangeEvent, Tab } from '@mui/material';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect, useMemo, useState } from 'react';
import CreateCouponForm from '../components/CouponForm/CreateCouponForm';
import CreateDealForm from '../components/DealForm/CreateDealForm';
import CreateDividendForm from '../components/DividendForm/CreateDividendForm';
import PortfolioTable from '../components/PortfolioTable/PortfolioTable';
import { PortfolioActionsMap } from '../components/PortfolioTable/PortfolioTableToolbar';
import TransactionForm from '../components/TransactionForm/TransactionForm';
import CreateExpenseForm from '../components/ExpenseForm/CreateExpenseForm';

export default function Portfolio({ params }: { params: { id: string } }) {
  const [openModal, setOpenModal] = useState<PortfolioActionsMap | false>(
    false,
  );
  const [activeTab, setActiveTab] = useState<string>('portfolio');

  const client = useQueryClient();
  const { data: portfolio } = useQuery({
    queryKey: ['portfolio', parseInt(params.id)],
    queryFn: async () => {
      const portfolio = await investorService.portfolio.getPortfolio(params.id);
      if (!portfolio) {
        // router.push(`/portfolios`);
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

  const shareTickers = useMemo<SelectOption[]>(() => {
    if (!portfolio) {
      return [];
    }
    return portfolio.deals
      .reduce<string[]>((acc, deal) => {
        if (
          !acc.includes(deal.ticker) &&
          deal.securityType === SecurityType.SHARE
        ) {
          acc.push(deal.ticker);
        }
        return acc;
      }, [])
      .map<SelectOption>(ticker => {
        return {
          id: ticker,
          name: ticker,
        };
      })
      .sort((a, b) => (a.id > b.id ? 1 : -1));
  }, [portfolio]);

  const bondTickers = useMemo<SelectOption[]>(() => {
    if (!portfolio) {
      return [];
    }
    return portfolio.deals
      .reduce<string[]>((acc, deal) => {
        if (
          !acc.includes(deal.ticker) &&
          deal.securityType === SecurityType.BOND
        ) {
          acc.push(deal.ticker);
        }
        return acc;
      }, [])
      .map<SelectOption>(ticker => {
        return {
          id: ticker,
          name: ticker,
        };
      })
      .sort((a, b) => (a.id > b.id ? 1 : -1));
  }, [portfolio]);

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
      case PortfolioActionsMap.dividend:
        setOpenModal(PortfolioActionsMap.dividend);
        break;
      case PortfolioActionsMap.coupon:
        setOpenModal(PortfolioActionsMap.coupon);
        break;
      case PortfolioActionsMap.expense:
        setOpenModal(PortfolioActionsMap.expense);
        break;
    }
  };

  useEffect(() => {
    console.log('openModal', openModal);
  }, [openModal]);

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
            [PortfolioActionsMap.dividend]: (
              <CreateDividendForm
                afterSuccessfulSubmit={() => setOpenModal(false)}
                portfolioId={Number(params.id)}
                tickerList={shareTickers}
              />
            ),
            [PortfolioActionsMap.coupon]: (
              <CreateCouponForm
                afterSuccessfulSubmit={() => setOpenModal(false)}
                portfolioId={Number(params.id)}
                tickerList={bondTickers}
              />
            ),
            [PortfolioActionsMap.expense]: (
              <CreateExpenseForm
                afterSuccessfulSubmit={() => setOpenModal(false)}
                portfolioId={Number(params.id)}
              />
            ),
          }[openModal]}
      </Dialog>
    </>
  );
}
