import React, { useState } from 'react';

import { Avatar, Button, Card, Divider, Flex, Modal, Radio, Select, Typography } from 'antd';
import { EllipsisOutlined, LockOutlined, UploadOutlined } from '@ant-design/icons';
import { BoardApi, UserApi } from 'est_proxy_api';
import useMessage from 'antd/es/message/useMessage';

import classes from './HeadMenu.module.css';

import PropTypes from 'prop-types';

const boardApiInstance = new BoardApi();
const userApiInstance = new UserApi();

const defaultAccessSettingsData = {
  userId: null,
  access: 'read',
};

const accesses = {
  'owner': 'Владелец',
  'admin': 'Админ',
  'read': 'Чтение',
  'write': 'Редактирование',
};

const userAccessOptions = [
  {
    value: 'read',
    label: 'Чтение',
  },
  {
    value: 'write',
    label: 'Запись',
  },
  {
    value: 'admin',
    label: 'Полный доступ',
  },
  {
    value: 'delete',
    label: 'Закрыть доступ',
  },
];

const UserAvatar = ({ user, access, boardId, refreshData }) => {
  const [messageApi, contextHolder] = useMessage();
  const [active] = useState(true);

  const onChange = (newAccess) => {
    if (newAccess === 'delete') {
      boardApiInstance.unshare(boardId, {
        'unshareRequest': {
          userId: user.id,
        },
      }).then(() => {
        refreshData();
      }).catch((error) => {
        messageApi.open({
          type: 'error',
          content: error.response ? error.response.text : 'Попробуйте позже',
        });
      });

      return;
    }

    boardApiInstance.changeAccess(boardId, {
      'shareBoardDto': {
        userId: user.id,
        access: newAccess,
      },
    }).then(() => {
      messageApi.open({
        type: 'success',
        content: 'Вы успешно изменили права!',
      });
    }).catch((error) => {
      messageApi.open({
        type: 'error',
        content: error.response ? error.response.text : 'Попробуйте позже',
      });
    });

    refreshData();
  };

  if (!active) {
    return (
      <div></div>
    );
  }

  return (
    <Flex align="center" gap="large">
      {contextHolder}
      <Avatar src={'https://api.dicebear.com/7.x/miniavs/svg?seed=' + user.id} />
      <Flex className="w100p" align="center" justify="space-between">
        <Typography.Text>{user.username}</Typography.Text>
        {
          access === 'owner' ?
            <Typography.Text type="secondary">{accesses[access]}</Typography.Text> :
            <Select options={userAccessOptions} defaultValue={access} style={{ width: 140.47 }}
                    allowClear={false} onChange={onChange} />
        }
      </Flex>
    </Flex>
  );
};

UserAvatar.propTypes = {
  user: PropTypes.shape({
    id: PropTypes.string.isRequired,
    username: PropTypes.string.isRequired,
  }).isRequired,
  access: PropTypes.oneOf(['read', 'write', 'admin', 'owner']).isRequired,
  boardId: PropTypes.string,
  refreshData: PropTypes.func,
};

const HeadMenu = ({ data, updateData, refreshData }) => {
  const [accessSettingsOpened, setAccessSettingsOpened] = useState(false);
  const [accessSettingsData, setAccessSettingsData] = useState(defaultAccessSettingsData);
  const [userList, setUserList] = useState([]);
  const [searchValue, setSearchValue] = useState();
  const [messageApi, contextHolder] = useMessage();

  const onNameChange = (name) => {
    const newData = data;
    newData.name = name;

    updateData(newData);
  };

  const handleOk = () => {
    setAccessSettingsData(defaultAccessSettingsData);
    setSearchValue();
  };

  const handleCancel = () => {
    setAccessSettingsOpened(false);
    setAccessSettingsData(defaultAccessSettingsData);
    setSearchValue();
  };

  const searchUsers = (text) => {
    if (text === '') return;

    userApiInstance.search(text).then((data) => {
      setUserList(data);
    }).catch((error) => {
      messageApi.open({
        type: 'error',
        content: error.response ? error.response.text : 'Что-то пошло не так ;(',
      });
    });
  };

  const share = () => {
    if (!accessSettingsData.userId) {
      messageApi.open({
        type: 'error',
        content: 'Пожалуйста, выберите пользователя',
      });

      return;
    }

    boardApiInstance.share(data.id, {
      'shareBoardDto': accessSettingsData,
    }).then(() => {
      refreshData();

      messageApi.open({
        type: 'success',
        content: 'Вы успешно поделились доской!',
      });

      handleOk();
    }).catch((error) => {
      messageApi.open({
        type: 'error',
        content: error.response ? error.response.text : 'Что-то пошло не так ;(',
      });
    });

    setAccessSettingsData(defaultAccessSettingsData);
  };

  return (
    <Flex className={`w100p ${classes.top}`} justify="space-between" style={{ zIndex: 15 }}>
      {contextHolder}
      <Modal
        title={
          <Flex>
            <LockOutlined style={{ margin: '0 8px 0 0' }} key={1} />
            <div key={2}>Настройки Доступа</div>
          </Flex>
        }
        open={accessSettingsOpened}
        onOk={share}
        onCancel={handleCancel}
        footer={[
          <Button type="default" onClick={handleCancel} key="cancel">Отмена</Button>,
          <Button type="primary" onClick={share} key="share">Поделиться</Button>,
        ]}
      >
        <Flex vertical gap="middle">
          <Divider style={{ margin: '8px 0 0 0' }} />
          <Typography.Title level={5} style={{ margin: 0 }}>Пользователи, имеющие доступ</Typography.Title>
          <Flex vertical gap="small">
            <div key={"owner-" + data.ownerInfo.id}>
              <UserAvatar user={data.ownerInfo} access={'owner'} />
            </div>
            {
              data.sharedWith.map(
                (userData) => (
                  <div key={userData.userInfo.id}>
                    <UserAvatar user={userData.userInfo} access={userData.access} boardId={data.id}
                              refreshData={refreshData} />
                  </div>
                  
                ),
              )
            }
          </Flex>
          <Divider style={{ margin: '8px 0 0 0' }} />
          <Typography.Title level={5} style={{ margin: 0 }}>Добавить пользователя</Typography.Title>
          <Select
            className="w100p"
            showSearch
            value={searchValue}
            defaultActiveFirstOption={false}
            suffixIcon={null}
            filterOption={false}
            onSearch={searchUsers}
            onChange={(val) => setSearchValue(val)}
            onSelect={(id) => {
              accessSettingsData.userId = id;
              setAccessSettingsData({...accessSettingsData});
            }}
            notFoundContent={null}
            placeholder="Выберите пользователя, которому хотите дать доступ"
            options={userList.map((d) => ({
              value: d.id,
              label: d.username,
            })).filter((d) => (d.value !== data.ownerInfo.id) &&
              !data.sharedWith.map((dd) => dd.userInfo.id).includes(d.value),
            )}//&& ))}
          />
          <Typography.Title level={5} style={{ margin: 0 }}>Уровень доступа</Typography.Title>
          <Radio.Group
            block
            options={[
              {
                label: 'Чтение',
                value: 'read',
              },
              {
                label: 'Редактирование',
                value: 'write',
              },
              {
                label: 'Полный доступ',
                value: 'admin',
              },
            ]}
            defaultValue={accessSettingsData.access}
            optionType="button"
            buttonStyle="solid"
            onChange={(e) => {
              accessSettingsData.access = e.target.value;
              setAccessSettingsData({...accessSettingsData});
            }}
          />
          <div />
        </Flex>

      </Modal>
      <Card size="small" className="shadow">
        <Flex gap="small" align="center">
          <Typography.Title level={4} style={{ margin: 0 }}
                            editable={{ 'onChange': (e) => onNameChange(e) }}>{data.name}</Typography.Title>
          <Divider type="vertical" style={{ height: 30 }} />
          <Button
            icon={<EllipsisOutlined />}
            key="ellipsis"
          ></Button>
          <Button icon={<UploadOutlined />}></Button>
        </Flex>
      </Card>
      <Card size="small" className="shadow" style={{width: 300}}>
        <Flex gap="small" align="center">
          <Typography.Link href="/app/home/my">
            <Avatar src={'https://api.dicebear.com/7.x/miniavs/svg?seed=' + data.ownerInfo.id}
                    shape="circle" />
          </Typography.Link>
          <Divider type="vertical" style={{ height: 30 }} />
          <Button
            type="primary"
            icon={<LockOutlined />}
            onClick={() => setAccessSettingsOpened(true)}
            style={{width: "-webkit-fill-available"}}
          >
            Настройки Доступа
          </Button>
        </Flex>
      </Card>
    </Flex>
  );
};

HeadMenu.propTypes = {
  data: PropTypes.shape({
    id: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    ownerInfo: PropTypes.shape({
      id: PropTypes.string.isRequired,
      username: PropTypes.string.isRequired,
    }).isRequired,
    sharedWith: PropTypes.arrayOf(
      PropTypes.shape({
        userInfo: PropTypes.shape({
          id: PropTypes.string.isRequired,
          username: PropTypes.string.isRequired,
        }).isRequired,
        access: PropTypes.oneOf(['read', 'write', 'admin', 'owner']).isRequired,
      })
    ).isRequired,
  }).isRequired,
  updateData: PropTypes.func.isRequired,
  refreshData: PropTypes.func.isRequired,
};

export default HeadMenu;