// Copyright 2023 Cleuton Sampaio.
//
// Licenciado sob a Licença Apache, Versão 2.0 (a "Licença");
// você não pode usar este arquivo exceto em conformidade com a Licença.
// Você pode obter uma cópia da Licença em
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// A menos que exigido por lei aplicável ou acordado por escrito, software
// distribuído sob a Licença é distribuído em uma base "COMO ESTÁ",
// SEM GARANTIAS OU CONDIÇÕES DE QUALQUER TIPO, expressas ou implícitas.
// Veja a Licença para o idioma específico que rege as permissões e
// limitações sob a Licença.

const express = require('express');
const { Sonyflake } = require('sonyflake');

const VERSAO = "0.1.0";
const PORTA = 8888;

// Inicializar instância do Sonyflake
const flake = new Sonyflake();

if (!flake) {
    console.error("ERRO FATAL: Não foi possível criar a instância do Sonyflake.");
    process.exit(1);
}

const app = express();

// Middleware para analisar corpos JSON
app.use(express.json());

// Definir o endpoint /nextid
app.get('/nextid', (req, res) => {
    try {
        const id = flake.nextId();
        res.setHeader('Content-Type', 'application/json');
        res.status(200).json({ erro: false, id: id.toString() }); // Converter para string para compatibilidade
    } catch (err) {
        console.error("ERRO: flake.nextId() falhou com", err);
        res.setHeader('Content-Type', 'application/json');
        res.status(500).json({ erro: true });
    }
});

// Capturar todas as outras rotas não suportadas
app.use((req, res) => {
    res.status(404).send('Caminho não suportado');
});

// Iniciar o servidor
const servidor = app.listen(PORTA, '0.0.0.0', () => {
    console.log(`INFO: ZAPTIDGEN NODEJS - V ${VERSAO} - worker ${process.pid} escutando em http://0.0.0.0:${PORTA}`);
});

// Desligamento gracioso
process.on('SIGINT', () => {
    console.log(`INFO: ZAPTIDGEN NODEJS worker ${process.pid} desligando gracioso...`);
    servidor.close(() => {
        console.log(`INFO: ZAPTIDGEN NODEJS worker ${process.pid} fechado.`);
    });
});

module.exports = app; // Exportar o app para uso no cluster