import React from 'react';

import {Card, Divider, Flex, Typography} from 'antd';
import 'antd/dist/reset.css';

import classes from './FormPage.module.scss';


const FormPage = ({title, form, children}) => {
    return (
        <div className={classes.container}>
            <img className={classes.inory} src="/inory.png" alt="inory"/>
            <Card>
                <Flex>
                    <Flex style={{width: 300}} align='center' justify="center" vertical>
                        <img className={classes.logo} src="/logo.png" alt="logo"
                            style={{width: '120px', height: 'auto'}}/>
                        <Typography.Title level={2} style={{textAlign: "center"}}>{title}</Typography.Title>
                    </Flex>
                    <Divider type='vertical' style={{height: "auto", margin: "0 24px"}}/>
                    <Flex style={{width: 300}} vertical>
                        {form}
                        {children}
                    </Flex>
                </Flex>
            </Card>
        </div>
    );
};

export default FormPage;