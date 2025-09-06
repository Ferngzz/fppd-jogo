// main.go - Loop principal do jogo
package main

import (
	"os"
	"time"
)

var Score = 0

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	coinRespawnChannel := make(chan bool)
	monsterSpawnChannel := make(chan bool)
	powerSpawnChannel := make(chan bool, 1)

	// Go routine para spawn contínuo de moedas,
	// com timeout indicado no parâmetro da função
	go func() {
		for range coinRespawnChannel {
			spawnCoin(&jogo, 15*time.Second, coinRespawnChannel, monsterSpawnChannel)
		}
	}()

	// Go routine para spawnar um monstro
	// toda vez que uma moeda expirar
	go func() {
		for range monsterSpawnChannel {
			spawnMonster(&jogo)
		}
	}()
	// Go routine para spawnar um poder periodicamente
	go func() {
		<-time.After(1 * time.Minute)
		powerSpawnChannel <- true

		for range powerSpawnChannel {
			spawnPower(&jogo, powerSpawnChannel)
			<-time.After(1 * time.Minute)
		}
	}()

	// Spawna a primeira moeda
	coinRespawnChannel <- true

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo, coinRespawnChannel, powerSpawnChannel); !continuar {
			break
		}

		if jogo.Mapa[jogo.PosY][jogo.PosX].simbolo == Moeda.simbolo {
			jogo.Mapa[jogo.PosY][jogo.PosX] = Vazio
			jogo.StatusMsg = "Moeda coletada!"
			coinRespawnChannel <- true // spawn da próxima moeda
		}

		interfaceDesenharJogo(&jogo)
	}
}
