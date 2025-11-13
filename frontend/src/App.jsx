import React, { useState } from "react";
import { Card, Col, Form, Input, Row, Button, Space } from "antd";
import { Typography } from "antd";
import { DeleteOutlined, CaretUpOutlined } from "@ant-design/icons";
import ModalEditarCard from "./componentes/ModalEditarCard";
import FormNovoCard from "./componentes/FormNovoCard";
import "./style.css";

const { Title } = Typography;

const KanbanBoard = () => {
  const [formAdicionarColuna] = Form.useForm();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [columnToEdit, setColumnToEdit] = useState(null);
  const [cardToEdit, setCardToEdit] = useState(null);

  const [columns, setColumns] = useState([
    {
      id: 1,
      title: "A Fazer",
      cards: [
        { id: 1, title: "Card 1", description: "Descrição do Card 1" },
        { id: 2, title: "Card 2", description: "Descrição do Card 2" },
      ],
    },
    { id: 2, title: "Em progresso", cards: [{ id: 3, title: "Card 3", description: "Descrição do Card 3" }] },
    { id: 3, title: "Concluídas", cards: [{ id: 4, title: "Card 4", description: "Descrição do Card 4" }] },
  ]);

  const handleDragStart = (event, cardId) => {
    event.dataTransfer.setData("text/plain", cardId.toString());
  };

  const handleDragOver = (event) => {
    event.preventDefault();
  };

  function handleDrop(event, column) {
    event.preventDefault();

    const cardId = Number(event.dataTransfer.getData("text/plain"));

    const isInDestinationColumn = columns.some(
      (col) =>
        col.title === column && col.cards.some((card) => card.id === cardId)
    );

    if (isInDestinationColumn) {
      console.log("O card já está na coluna de destino.");
      return;
    }

    const updatedColumns = columns.map((col) => {
      if (col.title === column) {
        const cardToMove = columns
          .find((col) => col.cards.some((card) => card.id === cardId))
          .cards.find((card) => card.id === cardId);
        const updatedCards = [...col.cards, cardToMove];
        return { ...col, cards: updatedCards };
      } else {
        const updatedCards = col.cards.filter((card) => card.id !== cardId);
        return { ...col, cards: updatedCards };
      }
    });

    setColumns(updatedColumns);
  }

  const handleFormSubmit = (values, column) => {
    const newCardId = Date.now();
    const newCard = { id: newCardId, title: values.cardTitle, description: values.cardDescription || "" };
    const updatedColumns = columns.map((col) => {
      if (col.title === column) {
        const updatedCards = [...col.cards, newCard];
        return { ...col, cards: updatedCards };
      }
      return col;
    });
    setColumns(updatedColumns);
  };

  const handleAddColumn = (values) => {
    const columnExists = columns.some(
      (coluna) => coluna.title === values.columnTitle
    );

    if (columnExists) {
      alert("Essa coluna já existe");
      return;
    }

    const newColumn = {
      id: Date.now(),
      title: values.columnTitle,
      cards: [],
    };
    setColumns([...columns, newColumn]);
    formAdicionarColuna.resetFields();
  };

  const handleCardDelete = (columnId, cardId) => {
    const updatedColumns = columns.map((col) => {
      if (col.id === columnId) {
        const updatedCards = col.cards.filter((card) => card.id !== cardId);
        return { ...col, cards: updatedCards };
      }
      return col;
    });
    setColumns(updatedColumns);
  };

  const handleCardEdit = (card, column) => {
    setCardToEdit(card);
    setColumnToEdit(column);
    setIsModalOpen(true);
  };

  const handleOk = (cardId, newValues) => {
    const updatedColumns = columns.map((col) => {
      const updatedCards = col.cards.map((card) =>
        card.id === cardId
          ? { ...card, title: newValues.title, description: newValues.description }
          : card
      );
      return { ...col, cards: updatedCards };
    });
    setColumns(updatedColumns);
    setIsModalOpen(false);
  };

  const handleCancel = () => {
    setIsModalOpen(false);
  };

  const deleteColumn = (columnId) => {
    const updatedColumns = columns.filter((column) => column.id !== columnId);
    setColumns(updatedColumns);
  };

  function moverCardTopo(columnId, cardId, novaPosicao) {
    const columnIndex = columns.findIndex((column) => column.id === columnId);

    if (columnIndex === -1) {
      console.log("Coluna não encontrada.");
      return;
    }

    const cardIndex = columns[columnIndex].cards.findIndex(
      (card) => card.id === cardId
    );

    if (cardIndex === -1) {
      console.log("Card não encontrado.");
      return;
    }

    const card = columns[columnIndex].cards.splice(cardIndex, 1)[0];
    columns[columnIndex].cards.splice(novaPosicao, 0, card);

    const updatedColumns = [...columns];

    setColumns(updatedColumns);
  }

  return (
    <div className="kanban-container">
      <Title>Quadro Kanban</Title>
      <Row gutter={16} className="kanban-board">
        {columns.map((column) => (
          <Col key={column.id} span={6}>
            <div className="column-header-wrapper">
              {column.id !== 1 && column.id !== 2 && column.id !== 3 && (
                <DeleteOutlined
                  onClick={() => deleteColumn(column.id)}
                  className="delete-column-icon"
                />
              )}
            </div>

            <Card title={column.title} className="column">
              <div
                className="card-container"
                onDragOver={handleDragOver}
                onDrop={(event) => handleDrop(event, column.title)}
              >
                <Space direction="vertical" size="middle" className="card-space">
                  {column.cards.map((card) => (
                    <div
                      key={card.id}
                      className="card-outer-wrapper"
                      draggable
                      onDragStart={(event) => handleDragStart(event, card.id)}
                    >
                      <div className="card-wrapper">
                        <Card
                          className="card"
                          onClick={() => handleCardEdit(card, column)}
                        >
                          <span className="card-title">{card.title}</span>
                          <p className="card-description">{card.description}</p>
                        </Card>

                        <div className="card-actions">
                          <CaretUpOutlined
                            className="card-action-icon"
                            onClick={() => moverCardTopo(column.id, card.id, 0)}
                          />
                          <DeleteOutlined
                            className="card-action-icon"
                            onClick={() => handleCardDelete(column.id, card.id)}
                          />
                        </div>
                      </div>
                    </div>
                  ))}
                </Space>

                <FormNovoCard
                  handleFormSubmit={handleFormSubmit}
                  columnTitle={column.title}
                />
              </div>
            </Card>
          </Col>
        ))}

        <Form
          onFinish={handleAddColumn}
          initialValues={{ addColumn: "" }}
          form={formAdicionarColuna}
          className="add-column-form"
        >
          <Form.Item
            name="columnTitle"
            rules={[
              {
                required: true,
                min: 3,
                message: "O campo deve ter pelo menos 3 caracteres.",
              },
            ]}
          >
            <Input
              placeholder="Digite o título da coluna"
              name="addColumn"
              maxLength={35}
            />
          </Form.Item>
          <Button type="link" htmlType="submit">
            Adicionar Coluna
          </Button>
        </Form>
      </Row>

      <ModalEditarCard
        isModalOpen={isModalOpen}
        handleOk={handleOk}
        handleCancel={handleCancel}
        card={cardToEdit}
        column={columnToEdit}
      />
    </div>
  );
};

export default KanbanBoard;