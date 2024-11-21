import React from 'react';

import { Button, Form, Input, message } from 'antd';
import 'antd/dist/reset.css';

import { BoardApi } from 'est_proxy_api';

import { Config } from '../../../../config.js';


const CreateBoardForm = () => {
    const [messageApi, contextHolder] = message.useMessage();

    const apiInstance = new BoardApi()
    apiInstance.apiClient.basePath = Config.back_url
    apiInstance.apiClient.defaultHeaders = {
        ...apiInstance.apiClient.defaultHeaders,
    };

    const onFinish = async (values) => {
        const opts = {
            createRequest: {
                name: values.name,
                description: values.description,
                linkSharedMode: 'none_by_link'
            },
        };

        try {
            await apiInstance.create(opts);

            messageApi.open({
                type: 'success',
                content: 'Доска создана!'
            })

            window.location.reload();
        } catch (error) {
            messageApi.open({
                type: 'error',
                content: error.response === undefined ? null : error.response.text,
            });
        }
    };

    return (
        <Form
            name="create board"
            onFinish={onFinish}
            initialValues={{
                remember: true,
            }}
            layout="vertical"
            requiredMark={false}
        >
            {contextHolder}
            <Form.Item
                name="name"
                label="Название"
                rules={[
                    {
                        required: true,
                        message: 'Пожалуйста, введите название доски!',
                    },
                ]}
            >
                <Input placeholder='Введите название доски...'/>
            </Form.Item>
            <Form.Item
                name="description"
                label="Описание"
            >
                <Input.TextArea placeholder='Введите описание доски...' autoSize/>
            </Form.Item>
            <Button type="primary" htmlType="submit" style={{width: '100%'}}>
                Создать
            </Button>
        </Form>
    );
};

export default CreateBoardForm;