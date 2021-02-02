import {Form, Input, Modal, notification, Select, TreeSelect} from 'antd';
import React, { useEffect, useState } from "react";

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
      </Form>
    </Modal>
  );
};
export default ListForm;

