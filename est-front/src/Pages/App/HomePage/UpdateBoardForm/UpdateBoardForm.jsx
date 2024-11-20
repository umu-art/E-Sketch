import React from 'react';

import { Button, Form, Input, message } from 'antd';
import 'antd/dist/reset.css';

import { BoardApi } from 'est_proxy_api';

import { Config } from '../../../../config.js';


const UpdateBoardForm = ({ data, onDataChange, closeModal }) => {
    const [messageApi, contextHolder] = message.useMessage();

    const apiInstance = new BoardApi()
    apiInstance.apiClient.basePath = Config.back_url
    apiInstance.apiClient.defaultHeaders = {
        ...apiInstance.apiClient.defaultHeaders,
    };

    const onFinish = async (values) => {
        const opts = {
            "createRequest": {
                ...data,
                ...values,
            }
        };

        try {
            const newData = await apiInstance.update(data.id, opts);

            onDataChange(newData);

            messageApi.open({
                type: 'success',
                content: 'Изменения сохранены!'
            })

            closeModal();
        } catch (error) {
            messageApi.open({
                type: 'error',
                content: error.response,
            });
        }
    };

    return (
        <Form
            name="update board"
            onFinish={onFinish}
            initialValues={{
                name: data.name,
                description: data.description,
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
            <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
                Сохранить
            </Button>
        </Form>
    );
};

export default UpdateBoardForm;