import React, { useMemo } from 'react';

import { Flex, Typography } from 'antd';

import SignInForm from './SignInForm';
import FormPage from '../../FormPage/FormPage';
import { useLocation } from 'react-router-dom';

function useQuery() {
    const { search } = useLocation();

    return useMemo(() => new URLSearchParams(search), [search]);
}

const SignInPage = () => {
    let query = useQuery();
    const to = query.get("to");

    return (
        <FormPage title={"Авторизация"} form={<SignInForm redirectTo={to || '/app'}/>}>
            <Flex justify='center' style={{ width: "100%", justifyContent: "center" }}>
                <Typography.Text>
                    Ещё нет аккаунта? <Typography.Link href='signup' underline>Регистрация</Typography.Link>
                </Typography.Text>
            </Flex>
        </FormPage>
    );
};

export default SignInPage;