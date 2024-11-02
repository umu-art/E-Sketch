import React, { useEffect, useState } from 'react';

import { Avatar, Card, message } from 'antd';
import { EditOutlined, EllipsisOutlined, SettingOutlined, LoadingOutlined, Loa } from '@ant-design/icons';

import classes from './BoardCard.module.scss';

import { BoardApi } from 'est_proxy_api';
import Meta from 'antd/es/card/Meta';


const apiInstance = new BoardApi();

const BoardCard = (
    id,
) => {
    const [data, setData] = useState(
        {
            id: "",
            name: "",
            description: "",
            ownerInfo: {
              id: "",
              username: "",
              avatar: ""
            },
            sharedWith: [
              {
                userInfo: {
                  id: "",
                  username: "",
                  avatar: ""
                },
                access: ""
              }
            ],
            linkSharedMode: "",
            preview: ""
        }
    );

    const [messageApi, contextHolder] = message.useMessage();

    // useEffect(() => {
    //     apiInstance.getByUuid(id, (error, data, response) => {
    //         if (error) {
    //             messageApi.open({
    //                 type: "error",
    //                 content: error.response.text,
    //             })

    //             return;
    //         }

    //         setData(data);
    //     })
    // }, [id]);

    if (data.id !== id) {
        return (
            <Card 
                className={classes.card}
                loading
                style={{
                    height: 'fit-content',
                }}
                cover={
                    <div className='flex w100p align-items-center justify-items-center justify-content-center' style={{height: "225px"}}>
                        <LoadingOutlined />
                    </div>
                }
                actions={[
                    <LoadingOutlined key="1" />,
                    <LoadingOutlined key="2" />,
                    <LoadingOutlined key="3" />,
                ]}
            >
                {contextHolder}
                <Meta
                    avatar={<LoadingOutlined />}
                    title={""}
                    description={""}
                    loading
                />
            </Card>
        );
    }

    return (
        <Card 
            className={classes.card}
            style={{
                height: 'fit-content',
            }}
            cover={
                <img alt="Preview" src={data.preview} style={{height: "225px", display: "block", objectFit: "cover"}}/>
            }
            actions={[
                <SettingOutlined key="setting" />,
                <EditOutlined key="edit" />,
                <EllipsisOutlined key="ellipsis" />,
            ]}
        >
            {contextHolder}
            <Meta
                avatar={<Avatar src={"https://api.dicebear.com/7.x/miniavs/svg?seed=" + data.ownerInfo.id} />}
                title={data.name}
                description={data.description}
            />
        </Card>
    );
};

export default BoardCard;