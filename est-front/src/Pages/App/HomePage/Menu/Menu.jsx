import { Menu, message } from 'antd';
import React, { useEffect, useState } from 'react';

import { AppstoreOutlined, LinkOutlined, SettingOutlined, UserOutlined } from '@ant-design/icons';
import { UserApi } from 'est_proxy_api';

const items = (data) => [
    {
        key: "head_group",
        type: "group",
        label: "e-Sketch",
        children: [
            {
                key: "account",
                icon: <UserOutlined />,
                label: data.username,
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

const apiInstance = new UserApi();

const AppMenu = () => {
    const [userData, setUserData] = useState({
        email: "",
    });

    const [messageApi, context] = message.useMessage(); 

    useEffect(() => {
        apiInstance.getSelf().then((data) => {
            console.log(data)
            setUserData(data);
          }, (error) => {
            console.log(error);
            messageApi.open({
                type: "error",
                content: error,
            })
          });
    }, [messageApi]);

    return (
        <>
            {context}
            <Menu
                className='h100p'
                style={{
                    width: 256,
                    minHeight: 600,
                }}
                items={items(userData)}
            />
        </>
    );
};

export default AppMenu;