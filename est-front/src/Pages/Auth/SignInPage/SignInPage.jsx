import React from 'react';

import { Flex, Typography } from 'antd';

import SignInForm from './SignInForm';
import FormPage from '../../FormPage/FormPage';


const SignInPage = () => {
    return (
        <FormPage title={"Авторизация"} form={<SignInForm />}>
            <Flex justify='center' style={{width: "100%", justifyContent: "center"}}>
                <Typography.Text>
                    Ещё нет аккаунта? <Typography.Link href='signup' underline>Регистрация</Typography.Link>
                </Typography.Text>
            </Flex>
        </FormPage>
    );
};

export default SignInPage;