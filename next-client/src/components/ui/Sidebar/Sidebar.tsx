'use client';
import { FC, useMemo } from 'react';
import { AppBar, List } from '@mui/material';
import MenuItem from './MenuItem';
import menu from './menu';
import { redirect, usePathname } from 'next/navigation';
import { useSession } from 'next-auth/react';

const Sidebar: FC = () => {
  const { data: session } = useSession();
  const pathname = usePathname();
  const loggedIn = !!session;

  const menuItems = useMemo(() => {
    return menu(loggedIn).map(i => {
      i.active = pathname === i.link;
      return i;
    });
  }, [loggedIn, pathname]);

  return (
    <AppBar
      component="aside"
      sx={{
        bgcolor: 'grey.A100',
        height: '100%',
        width: 240,
        position: 'fixed',
        left: 0,
        top: 0,
        bottom: 0,
        overflowY: 'auto',
        '&::-webkit-scrollbar': {
          width: '0.4em',
        },
        '&::-webkit-scrollbar-track': {
          boxShadow: 'inset 0 0 6px rgba(0,0,0,0.00)',
          webkitBoxShadow: 'inset 0 0 6px rgba(0,0,0,0.00)',
        },
        '&::-webkit-scrollbar-thumb': {
          // backgroundColor: 'rgba(0,0,0,.1)',
          backgroundColor: 'secondary.light',
          outline: '1px solid slategrey',
        },
      }}
    >
      <List component="nav" sx={{ padding: 0 }}>
        {menuItems.map(item => (
          <MenuItem item={item} key={`${item.link} ${item.title}`} />
        ))}
      </List>
    </AppBar>
  );
};

export default Sidebar;
