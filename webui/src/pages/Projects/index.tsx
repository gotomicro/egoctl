import {Button, Card, Divider, Form, message, Modal, Tag} from 'antd';
import {PageHeaderWrapper} from '@ant-design/pro-layout';
import React, {Fragment, useRef, useState} from 'react';
import ListForm from "./components/ListForm"
import Editor from "./components/Editor"
import Render from "./components/Render"
import {PlusOutlined} from '@ant-design/icons';
import SearchTable, {SearchTableInstance} from '@/components/SearchTable';
import api from "@/services/api";
import moment from "moment";

const handleCreate = async (values) => {
  const hide = message.loading('正在添加');
  try {
    const resp = await api.ProjectCreate(values)
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
  const hide = message.loading('正在更新');
  try {
    const resp = await api.ProjectUpdate(values)
    if (resp.code !== 0) {
      hide();
      message.error('更新失败，错误信息：' + resp.msg);
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

const handleDSL = async (values) => {
  const hide = message.loading('正在更新DSL');
  try {
    const resp = await api.ProjectUpdateDSL(values)
    if (resp.code !== 0) {
      hide();
      message.error('更新DSL失败，错误信息：' + resp.msg);
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
  const [editorModalVisible, handleEditorModalVisible] = useState<boolean>(false);
  const [renderModalVisible, handleRenderModalVisible] = useState<boolean>(false);
  const [initialValues, setInitialValues] = useState({});
  const [form] = Form.useForm();
  const actionRef = useRef<SearchTableInstance>();
  const columns = [
    {
      title: "项目名",
      dataIndex: "name",
      key: "name",
    }, {
      title: "路径",
      dataIndex: "path",
      key: "path",
    }, {
      title: "模板名称",
      dataIndex: "templateName",
      key: "templateName",
    }, {
      title: "语言",
      dataIndex: "language",
      key: "language",
      render(val): JSX.Element {
        return (<Tag color={"green"} key={val}>
          {val}
        </Tag>)
      }
    }, {
      title: "创建时间",
      dataIndex: "ctime",
      key: "ctime",
      render(val) {
        return moment(val, "X").format('YYYY-MM-DD HH:mm:ss')
      },
    }, {
      title: "更新时间",
      dataIndex: "utime",
      key: "utime",
      render(val) {
        return moment(val, "X").format('YYYY-MM-DD HH:mm:ss')
      },
    }, {
      title: '操作',
      dataIndex: 'operating',
      key: 'operating',
      valueType: "option",
      render: (value, record) => (
        <Fragment>
          <a
            onClick={() => {
              setInitialValues(record);
              handleEditorModalVisible(true);
            }}
          >
            DSL描述
          </a>
          <Divider type="vertical"/>
          <a
            onClick={() => {
              api.ProjectGen(record).then((res) => {
                if (res.code !== 0) {
                  message.error(res.msg);
                  return false;
                }
                message.success("生成代码成功，请查看目录：" + record.path)
                return true;
              });
            }}
          >
            生成代码
          </a>
          <Divider type="vertical"/>
          <a
            onClick={() => {
              setInitialValues(record);
              handleRenderModalVisible(true);
            }}
          >
            获取数据
          </a>
          <Divider type="vertical"/>
          <a
            onClick={() => {
              setInitialValues(record);
              handleUpdateModalVisible(true);
            }}
          >
            编辑
          </a>
          <Divider type="vertical"/>
          <a
            onClick={() => {
              Modal.confirm({
                title: '确认删除？',
                okText: '确认',
                cancelText: '取消',
                onOk: () => {
                  api.ProjectDelete(record).then((res) => {
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
                <Form style={{
                  marginTop: '10px'
                }}>
                  <Button type="primary" onClick={() => {
                    setInitialValues({
                      proType: "default",
                      apiPrefix: "/api",
                    });
                    handleCreateModalVisible(true)
                  }}>
                    <PlusOutlined/> 新建
                  </Button>
                </Form>
              </div>
            );
          }}
          request={(params) => api.ProjectList({...params})}
        />
      </Card>
      <ListForm
        formTitle={"创建"}
        onSubmit={async (value) => {
          const success = handleCreate(value);
          if (success) {
            handleCreateModalVisible(false);
            setInitialValues({});
            actionRef.current?.refresh();
          }
        }}
        initialValues={initialValues}
        onCancel={() => {
          handleCreateModalVisible(false)
          setInitialValues({});
        }}
        modalVisible={createModalVisible}
      />
      <ListForm
        formTitle={"编辑"}
        onSubmit={async (value) => {
          const success = await handleUpdate(value);
          if (success) {
            handleUpdateModalVisible(false);
            setInitialValues({});
            actionRef.current?.refresh();
          }
        }}
        onCancel={() => {
          setInitialValues({})
          handleUpdateModalVisible(false)
        }}
        modalVisible={updateModalVisible}
        initialValues={initialValues}
      />
      <Editor
        formTitle={"编辑DSL"}
        onSubmit={async (value) => {
          const success = await handleDSL(value);
          if (success) {
            setInitialValues({});
            handleEditorModalVisible(false);
            actionRef.current?.refresh();
          }
        }}
        onCancel={() => {
          setInitialValues({})
          handleEditorModalVisible(false)
        }}
        modalVisible={editorModalVisible}
        initialValues={initialValues}
      />
      <Render
        formTitle={"展示渲染数据"}
        onCancel={() => {
          handleRenderModalVisible(false)
        }}
        modalVisible={renderModalVisible}
        initialValues={initialValues}
      />
    </PageHeaderWrapper>
  );
}
export default TableList;
