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

const cluster = require('cluster');
const numCPUs = require('os').cpus().length;
const VERSAO = "0.1.0";

if (cluster.isMaster || cluster.isPrimary) {
    console.log(`INFO: ZAPTIDGEN NODEJS - V ${VERSAO} - Master ${process.pid} está rodando`);

    // Fork workers.
    for (let i = 0; i < numCPUs; i++) {
        cluster.fork();
    }

    cluster.on('exit', (worker, code, signal) => {
        console.log(`INFO: Worker ${worker.process.pid} morreu. Código: ${code}, Sinal: ${signal}`);
    });
} else {
    // Workers importam o servidor
    require('./server.js');
    console.log(`INFO: Worker ${process.pid} iniciado`);
}