# 🍽️ Food Shop API

Uma API robusta e escalável desenvolvida para simular funcionalidades semelhantes ao iFood, proporcionando uma solução completa para gestão de produtos, processamento de pedidos e interações com clientes em ambientes de food service.

## 🚀 Funcionalidades

- **🛡️ Gestão de Permissões (RBAC)**: Todos os endpoints são protegidos por Controle de Acesso Baseado em Funções (RBAC), garantindo que apenas Gerentes e Clientes tenham acesso às ações permitidas de acordo com seus papéis.

- **📦 Gestão de Produtos**: 
  - Gerentes podem criar, atualizar e excluir produtos do sistema.
  - Possibilidade de listar produtos com suporte a paginação.
  - **🌟 Produtos mais populares**: Permite recuperar os produtos mais populares de um restaurante com base nas vendas.

- **📋 Processamento de Pedidos**:
  - Clientes podem realizar pedidos facilmente.
  - Gerentes podem visualizar e atualizar o status de pedidos, que incluem os estados:
    - ⏳ Pendente
    - ❌ Cancelado
    - 🔄 Processando
    - 🚚 Em Entrega
    - ✅ Entregue
  - Suporte a filtros e paginação para gerenciar grandes volumes de pedidos.

- **⭐ Avaliações de Clientes**:
  - Clientes podem enviar avaliações sobre produtos ou serviços.
  - Gerentes podem visualizar e responder avaliações com filtros e paginação.

- **🔐 Gerenciamento de Sessões**:
  - Suporte a múltiplas sessões simultâneas.
  - Gerenciamento de sessões utilizando Redis para garantir eficiência e escalabilidade.
  - **Magic Links**: Utiliza links mágicos para autenticação de usuários sem necessidade de senha.
  - **Cookies HTTP-Only**: Implementa cookies seguros para gerenciamento de autenticação, protegendo contra ataques XSS (Cross-Site Scripting).

- **📧 Processos Assíncronos**:
  - Uso do RabbitMQ para envio assíncrono de e-mails e outras tarefas.
  - Sistema flexível para abstrair processos futuros que requerem assincronismo.

- **⚡ Cache Inteligente**:
  - Utiliza Redis para armazenar respostas de requisições recorrentes, reduzindo o tempo de resposta e diminuindo a carga no banco de dados.

## 🛠️ Tecnologias Utilizadas

- **🔵 Linguagem de Programação**: Go (1.23.3)
- **📂 Banco de Dados**: MySQL, gerenciado com GORM
- **📨 Message Broker**: RabbitMQ
- **⚡ Cache**: Redis
- **🔧 Infraestrutura**: Docker para ambientes isolados

## 📥 Instalação e Configuração

1. **Clone o repositório**:
   ```bash
   $ git clone https://github.com/G-Villarinho/food-shop-api.git
2. **Execute docker**:
     ```bash
   $ make docker up
3. **Execute as migrations e crie as secrets keys**:
     ```bash
   $ make migration && make generate-keys
3. **Execute o programa**:
     ```bash
   $ make start