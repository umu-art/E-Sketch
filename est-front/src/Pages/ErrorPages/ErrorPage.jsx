import { Button, Flex, Result } from 'antd';
import React from 'react';
import { Link } from 'react-router-dom';

const errorPageSubTitles = {
    403: "У вас нет доступа к этой странице.",
    404: "Страница не найдена.",
    500: "Что-то пошло не так.",
};

const ErrorPage = ({status}) => {
    const validStatus = [403, 404, 500].includes(status) ? status : 500;

    return (
        <Flex className='w100vw h100vh' align='center' justify='center'>
            <Result
                status={`${validStatus}`}
                title={`${validStatus}`}
                subTitle={errorPageSubTitles[validStatus]}
                extra={<Link to="/app"><Button type="primary">Вернуться домой</Button></Link>}
            />
        </Flex>
    );
};

export default ErrorPage;