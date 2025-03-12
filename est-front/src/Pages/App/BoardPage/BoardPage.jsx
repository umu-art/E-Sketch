import { Flex } from 'antd';
import React, { useEffect, useState } from 'react';
import Board from './Board/Board';
import { useNavigate, useParams } from 'react-router-dom';
import ToolPanel from './ToolPanel/ToolPanel';
import { BoardApi } from 'est_proxy_api';
import LoadingPage from '../../LoadingPage/LoadingPage';
import HeadMenu from './HeadMenu/HeadMenu';

import ErrorPage from '../../ErrorPages/ErrorPage';
import ScaleElement from './Scale/Scale';
import Messages from './Messages/Messages';
import GPTRequestPopover from './Board/GPT/GPTRequestPopover';

const apiInstance = new BoardApi();

const BoardPage = () => {
  const { boardId } = useParams();
  const [data, setData] = useState(null);

  const [err, setErr] = useState(null);

  const navigate = useNavigate();

  const updateData = (newData) => {
    apiInstance.update(boardId, { 'createRequest': newData }).then((respData) => {
      setData(respData);
    }).catch((error) => {
      console.log(error);
    });
  };

  const refreshBoardData = () => {
    apiInstance.getByUuid(boardId).then((data) => {
      setData(data);
    }).catch((error) => {
      if (error.status === 401) {
        navigate(`/auth/signin?to=${window.location.pathname}${window.location.search}`);
      } else {
        setErr(error);
      }
    });
  };


  useEffect(() => {
    if (!data && !err) {
      refreshBoardData();
    }
  });

  if (err) {
    return (
      <ErrorPage status={err.status === 400 ? 404 : err.status} />
    );
  }

  if (!data) {
    return (
      <LoadingPage />
    );
  }

  return (
      <div style={{ overflow: 'hidden', height: '100vh', width: '100vw', position: 'absolute' }}>
        <GPTRequestPopover boardId={boardId}/>
        <Board className="h100vh w100vw" style={{ position: 'absolute' }} boardId={boardId} />
        <ScaleElement style={{ position: 'absolute', right: '0', bottom: '0', zIndex: '6' }}></ScaleElement>
        { /* Menu wrap */}
        <Flex className="h100vh w100vw" style={{ padding: '20px 20px', position: 'absolute'}} vertical
              align="center" justify="space-between">
          { /* Top */}
          <Flex className='w100p' gap="small" vertical>
            <HeadMenu data={data} updateData={updateData} refreshData={refreshBoardData} />
            <Flex className='w100p' justify="end">
              <Messages />
            </Flex>
          </Flex>
          { /* Bottom */}
          <Flex className="w100p" justify="center">
            <ToolPanel/>
          </Flex>
        </Flex>
      </div>
  );
};

export default BoardPage;