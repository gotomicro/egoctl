import { Form, Input, Modal, InputNumber } from "antd";
import React, { useEffect, useState } from "react";
import MonacoEditor from "react-monaco-editor";

interface ListFormProps {
  modalVisible: boolean;
  formTitle: string;
  initialValues: {};
  onSubmit: () => void;
  onCancel: () => void;
}

const formLayout = {
  labelCol: { span: 4 },
  wrapperCol: { span: 20 },
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

  const modalFooter = {
    okText: "保存",
    onOk: handleSubmit,
    onCancel,
  };

  return (
    <Modal
      destroyOnClose
      title={formTitle}
      visible={modalVisible}
      width={"1200px"}
      {...modalFooter}
    >
      <Form {...formLayout} form={form} onFinish={onSubmit} scrollToFirstError>
        <Form.Item name="path" label="path" hidden>
          <Input />
        </Form.Item>
        <Form.Item
          label="描述DSL" name="dsl"
        >
          <MonacoEditor
            height={"500px"}
            width={"80%"}
            language={'go'}
            options={{
              theme: "vs-dark"
              // wordWrap: 'on',
              // automaticLayout: true,
              // tabSize: 2
            }}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};
export default ListForm;
