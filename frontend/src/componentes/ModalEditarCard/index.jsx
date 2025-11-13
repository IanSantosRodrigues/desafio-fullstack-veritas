import { useState, useEffect } from 'react';
import { Modal, Input, Typography, Tag } from 'antd';
import { TagFilled } from '@ant-design/icons';

const { Text } = Typography;
const { TextArea } = Input;

const ModalEditarCard = ({ isModalOpen, handleOk, handleCancel, card, column }) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');

  useEffect(() => {
    if (card) {
      setTitle(card.title);
      setDescription(card.description || '');
    }
  }, [card]);

  const handleTitleChange = (event) => {
    setTitle(event.target.value);
  };

  const handleDescriptionChange = (event) => {
    setDescription(event.target.value);
  };

  if (!card) {
    return null;
  }

  const onOkClick = () => {
    handleOk(card.id, { title: title, description: description });
  };

  return (
    <Modal
      title="Editar Card"
      visible={isModalOpen}
      onOk={onOkClick}
      onCancel={handleCancel}
    >
      {/* Campo Título */}
      <Text>Escreva um novo título para o card:</Text>
      <Input
        value={title}
        onChange={handleTitleChange}
        maxLength={35}
        style={{ marginBottom: '16px' }}
      />

      {/* Campo Descrição */}
      <Text>Descrição:</Text>
      <TextArea
        value={description}
        onChange={handleDescriptionChange}
        rows={4}
        placeholder="Adicione uma descrição (opcional)"
      />

      <Tag color='processing' icon={<TagFilled />} style={{ marginTop: '16px' }}>
        Coluna: {column.title}
      </Tag>
    </Modal>
  );
};

export default ModalEditarCard;