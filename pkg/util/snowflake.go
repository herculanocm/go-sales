package util

import (
	"hash/fnv"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
)

var node *snowflake.Node

// init inicializa o nó Snowflake.
// A melhor prática é definir a variável de ambiente SNOWFLAKE_NODE_ID.
// Como fallback, ele gera um ID a partir do hostname do contêiner/máquina
// para reduzir o risco de colisões em ambientes não configurados explicitamente.
func init() {
	var nodeID int64

	nodeIDStr := os.Getenv("SNOWFLAKE_NODE_ID")
	if nodeIDStr != "" {
		// 1. Opção preferencial: Usar a variável de ambiente se definida.
		id, err := strconv.ParseInt(nodeIDStr, 10, 64)
		if err != nil {
			log.Fatal().Msgf("SNOWFLAKE_NODE_ID inválido: %v", err)
		}
		nodeID = id
	} else {
		// 2. Fallback: Gerar um ID a partir do hostname.
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatal().Msgf("Não foi possível obter o hostname para gerar o ID do nó Snowflake: %v", err)
		}

		// Usa um hash FNV para converter o nome do host em um número.
		h := fnv.New64a()
		h.Write([]byte(hostname))
		// O ID do nó Snowflake é de 10 bits (0-1023). Usamos o módulo para garantir que o ID esteja nesse intervalo.
		const maxNodeID = 1023
		nodeID = int64(h.Sum64() & maxNodeID)
		log.Warn().Msgf("Aviso: SNOWFLAKE_NODE_ID não definido. Usando ID de nó %d derivado do hostname '%s'. Defina a variável de ambiente em produção.", nodeID, hostname)
	}

	const maxNodeID = 1023
	if nodeID < 0 || nodeID > maxNodeID {
		log.Fatal().Msgf("ID do nó Snowflake (%d) está fora do intervalo permitido [0, %d]", nodeID, maxNodeID)
	}

	n, err := snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatal().Msgf("Falha ao criar o nó Snowflake: %v", err)
	}
	node = n
}

// NewSnowflake gera um novo ID Snowflake como int64.
func NewSnowflake() int64 {
	return node.Generate().Int64()
}

// NewSnowflakePtr gera um ponteiro para um novo ID Snowflake.
func NewSnowflakePtr() *int64 {
	id := NewSnowflake()
	return &id
}

func ParseSnowflake(s string) (int64, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
