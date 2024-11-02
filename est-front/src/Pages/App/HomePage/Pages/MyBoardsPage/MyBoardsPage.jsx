import React, { useEffect, useState } from 'react';

import { Divider, Flex, message, Typography } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

import BoardCard from '../../BoardCard/BoardCard';

import { BoardApi } from 'est_proxy_api';

import { Config } from '../../../../../config.js';
import { useNavigate } from 'react-router-dom';


const apiInstance = new BoardApi();
apiInstance.apiClient.basePath = Config.back_url
// apiInstance.apiClient.defaultHeaders = {
//     ...apiInstance.apiClient.defaultHeaders,
// };


const MyBoardsPage = () => {
    const [boards, setBoards] = useState(null);

    const [messageApi, contextHolder] = message.useMessage();
    const navigate = useNavigate();

    useEffect(() => {
        apiInstance.list().then((error, data, response) => {
            if (error) {
                console.log("!!!")
                console.log(error);
                console.log("!!!")
                messageApi.open({
                    type: 'error',
                    content: error.response.text,
                })

                return;
            }

            setBoards(data.mine);
        }).catch((error) => {
            console.log(error);
            if (error.error.statusCode === 401) {
                messageApi.open({
                    type: 'error',
                    content: error.rawResponse,
                })
                navigate("/auth/signin");
            }
        })
    }, [])

    return (
        <Flex vertical style={{padding: "20px 50px", width: '-webkit-fill-available'}}>
            {contextHolder}
            <Typography.Title style={{margin: 0}}>Мои доски</Typography.Title>
            <Divider></Divider>
            <Flex style={{padding: "", width: '-webkit-fill-available'}} wrap gap="large">
                {
                    boards === null ?
                    <LoadingOutlined /> 
                    :
                    <Typography.Paragraph>
                        {boards}
                    </Typography.Paragraph>
                    // boards.map(
                    //     (board, i) =>
                    //         <BoardCard id={board.id} key={board.id}/>
                    // )
                }
            </Flex>
        </Flex>
    );
};

export default MyBoardsPage;