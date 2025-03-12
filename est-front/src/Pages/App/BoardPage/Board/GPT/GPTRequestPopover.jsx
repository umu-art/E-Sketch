import React, { useState } from 'react';
import { Modal, Input, Button } from 'antd';
import { useDispatch, useSelector } from 'react-redux';
import { addMessage, hideGPTPopover, setGPTStatus } from './../../../../../redux/toolSettings/actions';
import { GptApi } from 'est_proxy_api';
import { Config } from '../../../../../config';
import PropTypes from 'prop-types';

const GPT_STATUS = {
    SUCCESS: 'success',
    ERROR: 'error',
};

const apiInstance = new GptApi();
apiInstance.apiClient.basePath = Config.back_url;

const useGPTRequest = (boardId, popover) => {
    const dispatch = useDispatch();
    const [message, setMessage] = useState('');

    const handleSubmit = async () => {
        if (!message.trim()) return;

        dispatch(hideGPTPopover());

        const opts = {
            gPTRequestDto: {
                boardId,
                prompt: message,
                leftUp: popover.request.leftUp,
                rightDown: popover.request.rightDown,
            },
        };

        try {
            const data = await apiInstance.request(opts);

            dispatch(setGPTStatus(GPT_STATUS.SUCCESS));
            dispatch(
                addMessage({
                id: Date.now(),
                title: 'GPT',
                content: data.text,
                })
            );
        } catch (error) {
            dispatch(setGPTStatus(GPT_STATUS.ERROR));
            dispatch(
                addMessage({
                id: Date.now(),
                title: 'GPT',
                content: 'Произошла ошибка. Попробуйте позже.',
                })
            );
        } finally {
            setTimeout(() => dispatch(setGPTStatus(null)), 3000);
            setMessage('');
        }
    };

    const handleCancel = () => {
        dispatch(hideGPTPopover());
        dispatch(setGPTStatus(null));
        setMessage('');
    };

    return { message, setMessage, handleSubmit, handleCancel };
};

const GPTRequestPopover = ({ boardId }) => {
    const popover = useSelector((state) => state.popover);
    const { message, setMessage, handleSubmit, handleCancel } = useGPTRequest(boardId, popover);

    return (
        <Modal
        title="Введите запрос"
        open={popover.visible}
        onOk={handleSubmit}
        onCancel={handleCancel}
        centered
        footer={[
            <Button key="cancel" onClick={handleCancel}>
            Отмена
            </Button>,
            <Button key="submit" type="primary" onClick={handleSubmit}>
            Отправить
            </Button>,
        ]}
        >
        <Input.TextArea
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Введите текст..."
            autoSize
        />
        </Modal>
    );
};

GPTRequestPopover.propTypes = {
    boardId: PropTypes.string.isRequired,
};

export default GPTRequestPopover;