import { IMenuItem } from './MenuItem';

const menu = (loggedIn: boolean): IMenuItem[] => {
  let menuItems: IMenuItem[];

  if (loggedIn) {
    menuItems = [
      {
        title: 'Скринер',
        iconName: 'ShowChart',
        link: '/',
        active: false,
      },
      {
        title: 'Портфолио',
        iconName: 'BusinessCenter',
        link: '/portfolios',
        active: false,
      },
      {
        title: 'Эксперты',
        iconName: 'SentimentVerySatisfied',
        link: '/experts',
        active: false,
      },
      {
        title: 'Выйти',
        iconName: 'Logout',
        link: '/api/auth/signout',
        active: false,
      },
    ];
  } else {
    menuItems = [
      {
        title: 'Логин',
        iconName: 'BusinessCenterIcon',
        link: '/login',
        active: false,
      },
      {
        title: 'Регистрация',
        iconName: 'BusinessCenterIcon',
        link: '/register',
        active: false,
      },
    ];
  }

  return menuItems;
};

export default menu;
