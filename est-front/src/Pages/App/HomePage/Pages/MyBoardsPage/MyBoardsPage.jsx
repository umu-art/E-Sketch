import React, { useEffect, useState } from 'react';

import { Button, Divider, Flex, message, Modal, Typography } from 'antd';
import { LoadingOutlined, PlusOutlined } from '@ant-design/icons';

import BoardCard from '../../BoardCard/BoardCard';

import { BoardApi } from 'est_proxy_api';

import { Config } from '../../../../../config.js';
import { useNavigate } from 'react-router-dom';
import CreateBoardForm from '../../CreateBoardForm/CreateBoardForm';


const MyBoardsPage = () => {
    const [boards, setBoards] = useState(null);
    const [recentBoards, setRecentBoards] = useState(null);

    const [createBoardModalOpen, setCreateBoardModalOpen] = useState(false);

    const [messageApi, contextHolder] = message.useMessage();
    const navigate = useNavigate();

    const apiInstance = new BoardApi();
    apiInstance.apiClient.basePath = Config.back_url

    useEffect(() => {
        apiInstance.list().then((data) => {
            setBoards(data.mine);
            setRecentBoards(data.recent);
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
            <Modal
                centered
                title="Создать доску"
                open={createBoardModalOpen}
                onCancel={() => setCreateBoardModalOpen(false)}
                footer={null}
            >
                <CreateBoardForm/>
            </Modal>

            {
                recentBoards === null || recentBoards === undefined || recentBoards.length === 0 || (boards === null || boards.length <= 4) ? null :
                <>
                <Flex className="w100p" align='left' justify='space-between'>
                    <Typography.Title style={{ margin: 0 }}>
                        Недавние
                    </Typography.Title>
                </Flex>

                <Divider></Divider>

                <Flex className='w100p' align='left' vertical>
                    <Flex wrap gap="large" style={{ width: 'auto' }}>
                        {
                            recentBoards === null ?
                                <LoadingOutlined/>
                                :
                                recentBoards.map(
                                    (board, i) =>
                                        <BoardCard data={board} key={i}/>
                                )
                        }
                    </Flex>
                </Flex>
                <Divider></Divider>
                </>
            }
            <Flex className="w100p" align='center' justify='space-between'>
                <Typography.Title style={{ margin: 0 }}>
                    Мои доски
                </Typography.Title>
                <Button icon={<PlusOutlined/>} type='primary' onClick={() => setCreateBoardModalOpen(true)}>Новая доска</Button>
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
                                    <BoardCard data={board} key={i}/>
                            )
                    }
                </Flex>
            </Flex>
        </Flex>
    );
};

export default MyBoardsPage;