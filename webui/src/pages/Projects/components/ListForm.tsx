import {Form, Input, Modal, notification, Select, TreeSelect} from 'antd';
import React, { useEffect, useState } from "react";
import api from "@/services/api";

interface ListFormProps {
  modalVisible: boolean;
  formTitle: string;
  initialValues: {};
  onSubmit: () => void;
  onCancel: () => void;
}

const formLayout = {
  labelCol: { span: 7 },
  wrapperCol: { span: 13 },
};

const ListForm: React.FC<ListFormProps> = (props) => {
  const { modalVisible, onCancel, onSubmit, initialValues, formTitle } = props;
  const [form] = Form.useForm();
  const [selectData, setSelectData] = useState([]); // 设置select

  useEffect(() => {
    if (initialValues) {
      form.resetFields();
      form.setFieldsValue({
        ...initialValues,
      });
    }
  }, [initialValues]);

  const handleSubmit = () => {
    if (!form) return;
    form.submit();
  };

  // 选择select
  useEffect(() => {
    api.TemplateSelect().then((r)=>{
      if (r.code !== 0) {
        notification.error({
          message: "加载失败",
        });
        return;
      }
      setSelectData(r.data);
    })
  }, []);

  const laguageSelect = [{
    title: "Go",
    value: "Go",
  },{
    title: "React",
    value: "React",
  },{
    title: "Vue",
    value: "Vue",
  },{
    title: "其他",
    value: "其他",
  }]

  const modalFooter = { okText: '保存', onOk: handleSubmit, onCancel }

  return (
    <Modal
      destroyOnClose
      title={formTitle}
      visible={modalVisible}
      {...modalFooter}
    >
      <Form
        {...formLayout}
        form={form}
        onFinish={onSubmit}
        scrollToFirstError
      >
        <Form.Item
          name="name"
          label="标题"
        >
          <Input />
        </Form.Item>
        <Form.Item
          name="path"
          label="路径"
        >
          <Input />
        </Form.Item>
        <Form.Item
          name="gitRemotePath"
          label="模板"
        >
          <Select
            style={{ width: '100%' }}
            placeholder="模板"
            optionFilterProp={"name"}
          >
            {  (selectData || []).map((item,index)=>{
              return (<Select.Option key={index} name={item.title} value={item.value}>{item.title}</Select.Option>)
            })}
          </Select>
        </Form.Item>
        <Form.Item
          name="language"
          label="语言"
        >
          <Select
            style={{ width: '100%' }}
            placeholder="语言"
            optionFilterProp={"name"}
          >
            {  laguageSelect.map((item,index)=>{
              return (<Select.Option key={index} name={item.title} value={item.value}>{item.title}</Select.Option>)
            })}
          </Select>
        </Form.Item>
        <Form.Item
          name="proType"
          label="模板类型"
        >
          <Input />
        </Form.Item>
        <Form.Item
          name="apiPrefix"
          label="API前缀"
        >
          <Input />
        </Form.Item>
      </Form>
    </Modal>
  );
};
export default ListForm;

