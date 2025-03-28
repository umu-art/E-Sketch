import React, { useEffect, useMemo, useState } from 'react';

import { Button, Flex, Form, Input, message, Typography } from 'antd';

import FormPage from '../../FormPage/FormPage';
import { useLocation, useNavigate } from 'react-router-dom';

import { UserApi } from 'est_proxy_api';
import { Config } from '../../../config';


function useQuery() {
    const { search } = useLocation();

    return useMemo(() => new URLSearchParams(search), [search]);
}

const EmailConfirmPage = () => {
    let query = useQuery();
    const tokenParam = query.get("token");

    const [token, setToken] = useState(tokenParam);

    const [messageApi, contextHolder] = message.useMessage();
    const navigate = useNavigate();

    const apiInstance = new UserApi();
    apiInstance.apiClient.basePath = Config.back_url;
    apiInstance.apiClient.defaultHeaders = {
        ...apiInstance.apiClient.defaultHeaders,
    };

    const onFinish = async () => {
        try {
            await apiInstance.confirm({
                confirmationDto: {
                    token: token,
                }
            })

            messageApi.open({
                type: 'success',
                content: 'Регистрация прошла успешно!'
            })

            navigate("/app");
        } catch(error) {
            messageApi.open({
                type: 'error',
                content: error.response ? error.response.text : 'Что-то пошло не так ;(',
            });
        }
    }

    useEffect(() => {
        if (token) {
            onFinish();
        }
    })

    return (
        <FormPage title={"Подтверждение почты"} form={
            <Flex vertical gap="large">
                <Typography.Text>
                    Письмо с подтверждением отправлено на вашу электронную почту.
                </Typography.Text>
                <Form
                    name="confirm"
                    layout="vertical"
                    requiredMark={false}
                    onFinish={onFinish}
                >
                    <Form.Item
                        name="token"
                        label="Код"
                    >
                        <Input placeholder='Введите код подтверждения...' onChange={(e) => setToken(e.target.value)}/>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType='submit' style={{ width: '100%' }}>
                            Подтвердить
                        </Button>
                    </Form.Item>
                </Form>
            </Flex>
        }>
            {contextHolder}
            <Flex justify='center' style={{ width: "100%", justifyContent: "center", height: "fit-content" }}>
                <Typography.Text>
                    Если вы не получили письма с подтверждением, проверьте папку спам.
                </Typography.Text>
            </Flex>
        </FormPage>
    );
};

export default EmailConfirmPage;