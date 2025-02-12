import React from 'react';
import { useNavigate } from 'react-router-dom';

import { Button, Form, Input, message } from 'antd';
import 'antd/dist/reset.css';

import { UserApi } from 'est_proxy_api';

import { Config } from '../../../config';


const SignUpForm = () => {
    const navigate = useNavigate();
    const [messageApi, contextHolder] = message.useMessage();

    const apiInstance = new UserApi();

    apiInstance.apiClient.basePath = Config.back_url
    apiInstance.apiClient.defaultHeaders = {
        ...apiInstance.apiClient.defaultHeaders,
    };

    const onFinish = async (values) => {
        const opts = {
            registerDto: {
                email: values.email,
                passwordHash: values.password,
                username: values.username,
            },
        };

        try {
            await apiInstance.register(opts);

            messageApi.open({
                type: 'success',
                content: 'Регистрация прошла успешно!',
            });

            navigate("/auth/confirm");
        } catch (error) {
            messageApi.open({
                type: 'error',
                content: error.response.text,
            });
        }
    };

    return (
        <Form
            name="register"
            onFinish={onFinish}
            initialValues={{
                remember: true,
            }}
            layout="vertical"
        >
            {contextHolder}
            <Form.Item
                name="username"
                label="Имя пользователя"
                rules={[
                    {
                        required: true,
                        message: 'Пожалуйста, введите ваше имя пользователя!',
                    },
                ]}
            >
                <Input placeholder='Введите имя...'/>
            </Form.Item>
            <Form.Item
                name="email"
                label="Email"
                rules={[
                    {
                        type: 'email',
                        message: 'Некорректный email!',
                    },
                    {
                        required: true,
                        message: 'Пожалуйста, введите ваш email!',
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
            {/* <Form.Item>
                <Checkbox>
                    Я прочитал и согласен с <Typography.Link href='terms'>Правилами и Условиями</Typography.Link>
                </Checkbox>
            </Form.Item> */}
            <Form.Item>
                <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
                    Зарегистрироваться
                </Button>
            </Form.Item>
        </Form>
    );
};

export default SignUpForm;