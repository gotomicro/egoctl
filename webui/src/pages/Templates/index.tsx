import {Button, Card, Divider, Form, message, Modal} from 'antd';
import {PageHeaderWrapper} from '@ant-design/pro-layout';
import React, {Fragment, useRef, useState} from 'react';
import ListForm from "./components/ListForm"
import {PlusOutlined} from '@ant-design/icons';
import SearchTable, {SearchTableInstance} from '@/components/SearchTable';
import {history} from "umi";
import api from "@/services/api";

const handleCreate = async (values) => {
  const hide = message.loading('正在添加');
  try {
    const resp = await api.TemplateCreate(values)
    if (resp.code !== 0) {
      hide();
      message.error('添加失败，错误信息：' + resp.msg);
      return true
    }
    hide();
    message.success('添加成功');
    return true;
  } catch (error) {
    hide();
    message.error('添加失败请重试！' + error);
    return false;
  }
};

const handleUpdate = async (values) => {
  const hide = message.loading('正在排序');
  try {
    const resp = await api.TemplateUpdate(values)
    if (resp.code !== 0) {
      hide();
      message.error('排序失败，错误信息：' + resp.msg);
      return true
    }
    hide();
    message.success('更新成功');
    return true;
  } catch (error) {
    hide();
    message.error('更新失败请重试！' + error);
    return false;
  }
};

const TableList: React.FC<{}> = () => {
  const [createModalVisible, handleCreateModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [updateValues, setUpdateValues] = useState({});
  const [form] = Form.useForm();
  const actionRef = useRef<SearchTableInstance>();
  const columns = [
    {
      title: "模板名",
      dataIndex: "name",
      key: "name",
    }, {
      title: "项目地址",
      dataIndex: "gitRemotePath",
      key: "gitRemotePath",
    }, {
      title: "路径",
      dataIndex: "path",
      key: "path",
    },{
      title: "状态",
      dataIndex: "statusText",
      key: "statusText",
    },
    {
      title: '操作',
      dataIndex: 'operating',
      key: 'operating',
      valueType: "option",
      render: (value, record) => (
        <Fragment>
          <a
            onClick={() => {
              if (record.resourceType === 3) {
                history.push(`/resource/document/write?columnId=${record.id}`)
              } else {
                handleUpdateModalVisible(true);
                setUpdateValues(record);
              }
            }}
          >
            编辑
          </a>
          <Divider type="vertical"/>
          <a
            onClick={() => {
              Modal.confirm({
                title: '第一次同步模板时间较长，请耐心等待？',
                okText: '确认',
                cancelText: '取消',
                onOk: () => {
                  api.TemplateSync(record).then((res) => {
                    if (res.code !== 0) {
                      message.error(res.msg);
                      return false;
                    }
                    actionRef.current?.refresh();
                    return true;
                  });
                },
              });
            }}
          >
            同步模板
          </a>
          <Divider type="vertical"/>
          <a
            onClick={() => {
              Modal.confirm({
                title: '确认删除？',
                okText: '确认',
                cancelText: '取消',
                onOk: () => {
                  api.TemplateDelete(record).then((res) => {
                    if (res.code !== 0) {
                      message.error(res.msg);
                      return false;
                    }
                    actionRef.current?.refresh();
                    return true;
                  });
                },
              });
            }}
          >
            删除
          </a>
        </Fragment>
      ),
    },
  ];
  return (
    <PageHeaderWrapper>
      <Card>
        <SearchTable
          ref={actionRef}
          form={form}
          columns={columns}
          rowKey="id"
          pagination={false}
          formContent={search => {
            return (
              <div>
                <Form style={
                  {
                    marginTop: '10px'
                  }
                }>
                  <Button type="primary" onClick={() => handleCreateModalVisible(true)}>
                    <PlusOutlined/> 新建
                  </Button>
                </Form>
              </div>
            );
          }}
          request={(params) => api.TemplateList({...params})}
        />
      </Card>
      <ListForm
        formTitle={"创建"}
        onSubmit={async (value) => {
          const success = handleCreate(value);
          if (success) {
            handleCreateModalVisible(false);
            actionRef.current?.refresh();
          }
        }}
        onCancel={() => handleCreateModalVisible(false)}
        modalVisible={createModalVisible}
      />
      <ListForm
        formTitle={"编辑"}
        onSubmit={async (value) => {
          const success = await handleUpdate(value);
          if (success) {
            handleUpdateModalVisible(false);
            setUpdateValues({});
            actionRef.current?.refresh();
          }
        }}
        onCancel={() => {
          setUpdateValues({})
          handleUpdateModalVisible(false)
        }}
        modalVisible={updateModalVisible}
        initialValues={updateValues}
      />
    </PageHeaderWrapper>
  );
}
export default TableList;
