import React from 'react';

import { Flex, Typography } from 'antd';

import SignUpForm from './SignUpForm';
import FormPage from '../../FormPage/FormPage';


const SignUpPage = () => {
    return (
        <FormPage title={"Регистрация"} form={<SignUpForm />}>
            <Flex justify='center' style={{ width: "100%", justifyContent: "center", height: "fit-content" }}>
                <Typography.Text>
                    Уже есть аккаунт? <Typography.Link href='signin' underline>Вход</Typography.Link>
                </Typography.Text>
            </Flex>
        </FormPage>
    );
};

export default SignUpPage;