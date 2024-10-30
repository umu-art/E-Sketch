import React from 'react';

import { Card, Flex, Typography, Divider } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

import 'antd/dist/reset.css';

import classes from './FormPage.module.scss';


const FormPage = ({title, form, children}) => {
    return (
        <div className={classes.container}>
            <Card>
                <Flex>
                    <Flex style={{ width: 300 }} align='center' justify="center" vertical>
                        <LoadingOutlined style={{fontSize: "70px"}}/>
                        <Typography.Title level={2}>{title}</Typography.Title>
                    </Flex>
                    <Divider type='vertical' style={{height: "auto", margin: "0 24px"}}/>
                    <Flex style={{ width: 300 }} vertical>
                        {form}
                        {children}
                    </Flex>
                </Flex>
            </Card>
        </div>
    );
};

export default FormPage;