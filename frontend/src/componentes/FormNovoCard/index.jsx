import React from "react";
import { Form, Input, Button } from "antd";

const FormNovoCard = ({ handleFormSubmit, columnTitle }) => {
  const [form] = Form.useForm();

  const onFinish = (values) => {
    handleFormSubmit(values, columnTitle);
    form.resetFields();
  };

  return (
    <Form
      className="new-card-form"
      onFinish={onFinish}
      initialValues={{ cardTitle: "" }}
      form={form}
    >
      <Form.Item
        name="cardTitle"
        rules={[
          {
            required: true,
            min: 3,
            message: "O campo deve ter pelo menos 3 caracteres.",
          },
        ]}
      >
        <Input
          placeholder="Digite o título do card"
          maxLength={65}
        />
      </Form.Item>
      <Form.Item
        name="cardDescription" 
        rules={[{ required: false }]}
      >
        <Input.TextArea
          placeholder="Digite a descrição (opcional)"
          rows={3}
        />
      </Form.Item>
      <Button type="primary" htmlType="submit">
        Adicionar Card
      </Button>
    </Form>
  );
};

export default FormNovoCard;