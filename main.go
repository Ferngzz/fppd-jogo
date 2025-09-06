// main.go - Loop principal do jogo
package main

import "os"

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

	coinRespawnChannel := make(chan bool, 1)
	spawnCoin(&jogo)

	go func(jogo *Jogo, ch <-chan bool) {
		for canSpawn := range ch {
			if canSpawn {
				spawnCoin(jogo)
			}
		}
	}(&jogo, coinRespawnChannel)

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}

		// Verificando a coleta da coin
		if jogo.Mapa[jogo.PosY][jogo.PosX] == Moeda {
			jogo.StatusMsg = "Moeda coletada!"
			jogo.Mapa[jogo.PosY][jogo.PosX] = Vazio

			// envia sinal pro channel
			select {
			case coinRespawnChannel <- true:
			default:

			}
		}

		interfaceDesenharJogo(&jogo)
	}
}
