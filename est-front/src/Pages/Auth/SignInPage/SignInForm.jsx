import React from 'react';
import { useNavigate } from 'react-router-dom';

import { Button, Checkbox, Flex, Form, Input, message, Typography } from 'antd';
import 'antd/dist/reset.css';

import { UserApi } from 'est_proxy_api';

import { Config } from '../../../config';


const SignInForm = ({ redirectTo }) => {
    const navigate = useNavigate();
    const [messageApi, contextHolder] = message.useMessage();

    const apiInstance = new UserApi()
    apiInstance.apiClient.basePath = Config.back_url
    apiInstance.apiClient.defaultHeaders = {
        ...apiInstance.apiClient.defaultHeaders,
    };

    const onFinish = async (values) => {
        const opts = {
            authDto: {
                email: values.email,
                passwordHash: values.password,
            },
        };

        try {
            await apiInstance.login(opts);

            messageApi.open({
                type: 'success',
                content: 'Авторизация прошла успешно!'
            })

            navigate(redirectTo);
        } catch (error) {
            messageApi.open({
                type: 'error',
                content: error.response ? error.response.text : "Произошла ошибка, попробуйте позже",
            });
        }
    };

    return (
        <Form
            name="login"
            onFinish={onFinish}
            initialValues={{
                remember: true,
            }}
            layout="vertical"
            requiredMark={false}
        >
            {contextHolder}
            <Form.Item
                name="email"
                label="Почта"
                rules={[
                    {
                        type: 'email',
                        message: 'Некорректная почта!',
                    },
                    {
                        required: true,
                        message: 'Пожалуйста, введите вашу почту!',
                    },
                ]}
            >
                <Input placeholder='Введите почту...'/>
            </Form.Item>
            <Form.Item
                name="password"
                label="Пароль"
                rules={[
                    {
                        required: true,
                        message: 'Пожалуйста, введите ваш пароль!',
                    },
                ]}
            >
                <Input.Password placeholder='Введите пароль...'/>
            </Form.Item>
            <Form.Item
                name="rememberme"
            >
                <Flex align='center' justify='space-between'>
                    <Checkbox>Запомнить меня?</Checkbox>

                    <Typography.Link>
                        Забыли пароль?
                    </Typography.Link>
                </Flex>
            </Form.Item>
            <Form.Item>
                <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
                    Войти
                </Button>
            </Form.Item>
        </Form>
    );
};

export default SignInForm;