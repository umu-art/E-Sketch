import { Menu, message } from 'antd';
import React, { useEffect, useState } from 'react';

import { AppstoreOutlined, LinkOutlined, SettingOutlined, UserOutlined } from '@ant-design/icons';
import { UserApi } from 'est_proxy_api';
import { Link } from 'react-router-dom';


const items = (data) => [
    {
        key: "head_group",
        type: "group",
        label: "e-Sketch",
        children: [
            {
                key: "account",
                icon: <UserOutlined/>,
                label: data.username,
            },
            {
                key: "settings",
                icon: <SettingOutlined/>,
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
                icon: <Link to="my"><AppstoreOutlined/></Link>,
                label: "Мои доски",
            },
            {
                key: "shared",
                icon: <Link to="shared"><LinkOutlined /></Link>,
                label: "Поделились со мной",
            },
        ]
    },
]

const apiInstance = new UserApi();

const AppMenu = () => {
    const [userData, setUserData] = useState(null);

    const [messageApi, context] = message.useMessage();

    useEffect(() => {
        apiInstance.getSelf().then((data) => {
            setUserData(data);
        }).catch((error) => {
            console.log(error);
            messageApi.open({
                type: "error",
                content: error,
            })
        });
    }, [messageApi]);

    if (!userData) {
        return (
            <></>
        );
    }

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