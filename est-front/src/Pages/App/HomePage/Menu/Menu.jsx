import { Menu } from 'antd';
import React from 'react';

import { AppstoreOutlined, LinkOutlined, SettingOutlined, UserOutlined, ClockCircleOutlined } from '@ant-design/icons';

const items = [
    {
        key: "head_group",
        type: "group",
        label: "e-Sketch",
        children: [
            {
                key: "account",
                icon: <UserOutlined />,
                label: "Павел В.",
            },
            {
                key: "settings",
                icon: <SettingOutlined />,
                label: "Настройки",
            }
        ]
    },
    {
        key: "main_group",
        type: "group",
        label: "Доски",
        children: [
            {
                key: "all",
                icon: <AppstoreOutlined />,
                label: "Мои доски",
            },
            // {
            //     key: "recent",
            //     icon: <ClockCircleOutlined />,
            //     label: "Недавние",
            // },
            {
                key: "shared",
                icon: <LinkOutlined />,
                label: "Поделились со мной",
            },
        ]
    },
]

const AppMenu = () => {
    return (
        <Menu
            className='h100p'
            style={{
                width: 256,
                minHeight: 600,
            }}
            items={items}
        />
    );
};

export default AppMenu;