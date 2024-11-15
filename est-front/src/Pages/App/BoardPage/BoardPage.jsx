import { Avatar, Button, Card, Divider, Flex, Typography } from 'antd';
import React, { useEffect, useState } from 'react';
import Board from './Board/Board';
import { useNavigate, useParams } from 'react-router-dom';
import ToolPanel from './ToolPanel/ToolPanel';
import { BoardApi } from 'est_proxy_api';
import { EllipsisOutlined, UploadOutlined, UsergroupAddOutlined } from '@ant-design/icons';
import LoadingPage from '../../LoadingPage/LoadingPage';

const apiInstance = new BoardApi();

const BoardPage = () => {
    const { boardId } = useParams();
    const [data, setData] = useState(null)

    const navigate = useNavigate();

    const updateData = (newData) => {
        apiInstance.update(boardId, {'createRequest': newData}).then((respData) => {
            setData(respData);
        }).catch((error) => {
            console.log(error);
        });
    };

    const onNameChange = (name) => {
        const newData = data; 
        newData.name = name;

        updateData(newData);
    }

    useEffect(() => {
        apiInstance.getByUuid(boardId).then((data) => {
            setData(data);
            console.log(data);
        }).catch((error) => {
            navigate("/app/home/my");

            console.log(error);
        })
    }, [])

    if (!data) {
        return (
            <LoadingPage />
        );
    }

    return (
        <>  
            <Board className="h100vh w100vw" style={{position: 'absolute'}} boardId={boardId}/>
            { /* Menu wrap */}
            <Flex className="h100vh w100vw" style={{padding: "20px 20px", position: 'absolute'}} vertical align='center' justify='space-between'>
                { /* Top */}
                <Flex className='w100p' justify='space-between'>
                    <Card size='small' className='shadow'>
                    <Flex gap="small" align='center'>
                        <Typography.Title level={4} style={{margin: 0}} editable={{"onChange": (e) => onNameChange(e)}}>{data.name}</Typography.Title>
                        <Divider type='vertical' style={{height: 30}}/>
                        <Button 
                            icon={<EllipsisOutlined />}
                            key='ellipsis'
                        ></Button>
                        <Button icon={<UploadOutlined />}></Button>
                    </Flex>
                    </Card>
                    <Card size='small' className='shadow'>
                    <Flex gap="small" align='center'>
                        <Typography.Link href='/app/home/my'>
                            <Avatar src={"https://api.dicebear.com/7.x/miniavs/svg?seed=" + data.ownerInfo.id} shape='circle'/>
                        </Typography.Link>
                        <Divider type='vertical' style={{height: 30}}/>
                        <Button
                            type='primary'  
                            icon={<UsergroupAddOutlined />}
                        >
                            Поделиться
                        </Button>
                    </Flex>
                    </Card>
                </Flex>
                { /* Bottom */}
                <Flex className='w100p' justify='center'>
                    <ToolPanel onToolChange={(tool) => console.log(tool)} />
                </Flex>
            </Flex>
        </>
        
    );
};

export default BoardPage;