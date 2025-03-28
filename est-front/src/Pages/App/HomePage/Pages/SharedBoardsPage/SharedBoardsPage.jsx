import React, { useEffect, useState } from 'react';

import { Divider, Flex, message, Typography } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

import BoardCard from '../../BoardCard/BoardCard';

import { BoardApi } from 'est_proxy_api';

import { Config } from '../../../../../config.js';
import { useNavigate } from 'react-router-dom';


const SharedBoardsPage = () => {
    const [boards, setBoards] = useState(null);

    const [messageApi, contextHolder] = message.useMessage();
    const navigate = useNavigate();

    const apiInstance = new BoardApi();
    apiInstance.apiClient.basePath = Config.back_url

    useEffect(() => {
        apiInstance.list().then((data) => {
            setBoards(data.shared);
            console.log(data.shared);
        }).catch((error) => {
            navigate("/auth/signin");
            console.log(error);

            if (error.statusCode === 401) {
                messageApi.open({
                    type: 'error',
                    content: error.rawResponse,
                })

                navigate("/auth/signin");
            }
        })
    })

    return (
        <Flex vertical style={{ padding: "20px 50px", width: '-webkit-fill-available' }}>
            {contextHolder}
            <Flex className="w100p" align='center' justify='space-between'>
                <Typography.Title style={{ margin: 0 }}>
                    Поделились со мной
                </Typography.Title>
            </Flex>

            <Divider></Divider>
            <Flex className='w100p' align='left' vertical>
                <Flex wrap gap="large" style={{ width: 'auto' }}>
                    {
                        boards === null ?
                            <LoadingOutlined/>
                            :
                            boards.map(
                                (board, i) =>
                                    <BoardCard data={board} key={i} editable={false}/>
                            )
                    }
                </Flex>
            </Flex>
        </Flex>
    );
};

export default SharedBoardsPage;