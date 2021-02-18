import {Form, Input, Modal, InputNumber, message} from "antd";
import React, { useEffect, useState } from "react";
import MonacoEditor from "react-monaco-editor";
import api from "@/services/api";

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
    if (initialValues && initialValues.path != undefined) {
      form.resetFields();
      api.ProjectRender({path:initialValues.path}).then((res) => {
        if (res.code !== 0) {
          message.error(res.msg);
          return false;
        }
        form.setFieldsValue({
          ...initialValues,
          render: JSON.stringify(res.data, null, 2),
        });
        return true;
      });

      form.setFieldsValue({
        ...initialValues,
      });
    }
  }, [initialValues]);

  const modalFooter = {
    okText: "OK",
    onOk: onCancel,
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
      <Form {...formLayout} form={form} scrollToFirstError>
        <Form.Item name="path" label="path" hidden>
          <Input />
        </Form.Item>
        <Form.Item
          label="渲染数据" name="render"
        >
          <MonacoEditor
            height={"800px"}
            width={"80%"}
            language={'json'}
            options={{
              theme: "vs-dark",
              wordWrap: 'on',
              automaticLayout: true,
              tabSize: 2
            }}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};
export default ListForm;
