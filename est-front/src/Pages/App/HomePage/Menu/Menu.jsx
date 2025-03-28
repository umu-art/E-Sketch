import { Button, Flex, Menu, message } from 'antd';
import React, { useEffect, useState } from 'react';

import { AppstoreOutlined, LinkOutlined, LogoutOutlined, SettingOutlined, UserOutlined } from '@ant-design/icons';
import { UserApi } from 'est_proxy_api';
import { Link, useNavigate } from 'react-router-dom';


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

const AppMenu = () => {
    const [userData, setUserData] = useState(null);

    const [messageApi, context] = message.useMessage();

    const navigate = useNavigate();

    const apiInstance = new UserApi();
    
    useEffect(() => {
        async function fetchUserData() {
            try {
                const data = await apiInstance.getSelf();
                setUserData(data);
            } catch (error) {
                console.log(error);
                messageApi.open({
                    type: "error",
                    content: error,
                });
            }
        }

        fetchUserData();
    });

    if (!userData) {
        return (
            <>Loading...</>
        );
    }

    return (
        <div style={{display: "flex", flexDirection: "column", alignContent: "space-between"}}>
            {context}
            <Menu
                className='h100p'
                style={{
                    width: 256,
                    minHeight: 600,
                }}
                items={items(userData)}
            />
            <Flex style={{width: "100%"}} align='center' alignContent="center">
                <Button icon={<LogoutOutlined />} onClick={
                    () => {
                        apiInstance.logout().then(() => {
                            navigate('/');
                        }).catch((error) => {
                            console.log(error);
                            messageApi.open({
                                type: "error",
                                content: error,
                            })
                        });
                    }
                }>Выйти</Button>
            </Flex>
        </div>
    );
};

export default AppMenu;