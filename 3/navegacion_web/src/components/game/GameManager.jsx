/**
 * GameManager class
 * - Manages the game state and logic for a IA game.
 * 
 */

import TurnManager from "./TurnManager";
import PlayerBase from "./PlayerBase";
import Online_PlayerBase from "./Online_PlayerBase";
import IA_PlayerBase from "./IA_PlayerBase";
import BarajaClass from "./BarajaBase";

class GameManager {
    constructor(_numPlayers, esOnline) {
        this.state = {
            turnManager: null,
            players: Array(_numPlayers).fill(null),
            orden: Array(_numPlayers).fill(null),
            puntosJugadores: Array(_numPlayers).fill(0),
            numPlayers: _numPlayers,
            esOnline: esOnline,
            baraja: null,
            cartasJugadas: Array(_numPlayers).fill(null),
            triunfo: null,
            ganador: null,
            segundaBaraja: false,
            arrastre: false,
            finRonda: false,
            finJuego: false,
            /**
             * ! indexGanador: null,
             * ! ganador: null,
             * ! mostrarCantar: null,
             * ! cartasMoviendo: null,
             */
        };

        this.Evaluar = this.Evaluar.bind(this);
        this.TurnChange = this.TurnChange.bind(this);
    }

    Init(arrayDeCartas, primero, index) {
        this.state.arrastre = false;
        this.state.segundaBaraja = false;
        this.state.finRonda = false;
        this.state.finJuego = false;

        // * -------------------------- INICIAR BARAJA ----------------------------
        if (this.state.esOnline && arrayDeCartas) {
            this.state.baraja = new BarajaClass(arrayDeCartas) ;
        } else {
            this.state.baraja = new BarajaClass() ;
            //this.state.baraja.barajar();
        }
        // * ----------------------------------------------------------------------

        
        this.InitJugadores(index);

        this.state.triunfo = this.state.baraja.darCarta();
        this.state.baraja.anyadirAlFinal(this.state.triunfo);
            
        this.state.turnManager = new TurnManager(this.state.numPlayers, this.Evaluar, this.TurnChange, this.state.players);
        this.state.turnManager.reset();


        // * -------------------------- INICIAR ORDEN ----------------------------
        if(primero != null) {
            for (let i = 0; i < this.state.numPlayers; i++) {
                this.state.cartasJugadas[i] = i;
            }
        } else {
            for (let i = 0; i < this.state.numPlayers; i++) {
                let j = primero;
                this.state.orden[i] = j;
                primero = (primero + 1) % this.state.numPlayers;
            }
        }

        // * ----------------------------------------------------------------------

        this.state.turnManager.tick(); 

    }

    InitJugadores(index) {

        if (this.state.esOnline && index != null) {
            // * Si es online y recibimos el index del server inicializamos los jugadores
            // * El jugador 0 es el local y el resto son online
            this.state.players[0] = new PlayerBase(this);
            for (let j = 0; j < 6; j++) {
                this.state.players[0].anyadirCarta(this.state.baraja.darCarta());
            }

            // * Estos usuarios se inicializan dependiendo del número de jugadores
            // * Cada jugador se inicializa a partir del index recibido
            for (let i = 1; j < _numPlayers; i++) {

                index = (index + i) % this.state.numPlayers;
                this.state.players[i] = new Online_PlayerBase(this.state.numPlayers, this, index);
                for (let j = 0; j < 6; j++) {
                    this.state.players[i].anyadirCarta(this.state.baraja.darCarta());
                }

            }

        }
        else {
            this.state.players[0] = new PlayerBase(this);
            for (let j = 0; j < 6; j++) {
                this.state.players[0].anyadirCarta(this.state.baraja.darCarta());
            }
            if (this.state.numPlayers === 2) {
                this.state.players[1] = new IA_PlayerBase(this.state.numPlayers, this, 1);
                for (let j = 0; j < 6; j++) {
                    this.state.players[1].anyadirCarta(this.state.baraja.darCarta());
                }
            } else if (this.state.numPlayers === 4) {
                this.state.players[1] = new IA_PlayerBase(this.state.numPlayers, this, 1);
                for (let j = 0; j < 6; j++) {
                    this.state.players[1].anyadirCarta(this.state.baraja.darCarta());
                }
                this.state.players[2] = new IA_PlayerBase(this.state.numPlayers, this, 2);
                for (let j = 0; j < 6; j++) {
                    this.state.players[2].anyadirCarta(this.state.baraja.darCarta());
                }
                this.state.players[3] = new IA_PlayerBase(this.state.numPlayers, this, 3);
                for (let j = 0; j < 6; j++) {
                    this.state.players[3].anyadirCarta(this.state.baraja.darCarta());
                }
            }
        }
    }

    TurnChange() {
        this.state.players[this.state.orden[this.state.turnManager.state.playerTurn]].state.esMiTurno = true;
        
        console.log("Turno de: " + this.state.orden[this.state.turnManager.state.playerTurn]);
    }

    Evaluar() {
        this.evaluarLogic();
        console.log(this.state.orden);
        this.state.cartasJugadas = Array(this.state.numPlayers).fill(null);
        this.state.turnManager.tick();
    }

    evaluarLogic() {
        let indexGanador = this.state.orden[0];
        let maxPuntos = this.state.cartasJugadas[this.state.orden[0]].puntos;

        let sumaPuntos = this.state.cartasJugadas[this.state.orden[0]].puntos;
        let boolTriunfo = (this.state.cartasJugadas[this.state.orden[0]].palo === this.state.triunfo.palo);
        let paloJugado = this.state.cartasJugadas[this.state.orden[0]].palo;

        for (let i = 1; i < this.state.numPlayers; i++) {
            let aux = this.state.cartasJugadas[this.state.orden[i]].puntos;
            sumaPuntos += aux;
            if (boolTriunfo) {
                if (this.state.cartasJugadas[this.state.orden[i]].palo === this.state.triunfo.palo) {
                    if (aux > maxPuntos) {
                        maxPuntos = aux;
                        indexGanador = this.state.orden[i];
                    } else if (aux === maxPuntos && this.state.cartasJugadas[this.state.orden[i]].numero > this.state.cartasJugadas[indexGanador].numero) {
                        indexGanador = this.state.orden[i];
                    }
                }
            } else {
                if (this.state.cartasJugadas[this.state.orden[i]].palo === this.state.triunfo.palo) {
                    boolTriunfo = true;
                    maxPuntos = aux;
                    indexGanador = this.state.orden[i];
                } else if (this.state.cartasJugadas[this.state.orden[i]].palo === paloJugado) {
                    if (aux > maxPuntos) {
                        maxPuntos = aux;
                        indexGanador = this.state.orden[i];
                    } else if (aux === maxPuntos && this.state.cartasJugadas[this.state.orden[i]].numero > this.state.cartasJugadas[indexGanador].numero) {
                        indexGanador = this.state.orden[i];
                    }
                }
            }
        }

        // FINALIZACION DE TURNO
        this.state.players[indexGanador].state.ganador = true;
        this.state.players[(indexGanador + 1) % this.state.numPlayers].state.ganador = false;

        if (this.state.numPlayers === 4) {
            this.state.players[(indexGanador + 3) % 4].state.ganador = false;
            this.state.players[(indexGanador + 2) % 4].state.ganador = true;
        }

        this.state.players[indexGanador].state.puntos += sumaPuntos;

        for (let i = 0; i < this.state.numPlayers; i++) {
            this.state.orden[i] = (i + indexGanador) % this.state.numPlayers;
        }
        for (let i = 0; i < this.state.numPlayers; i++) {
            this.state.players[this.state.orden[i]].anyadirCarta(this.state.baraja.darCarta());
        }

        // COMPROBAR SI HA ACABADO LA RONDA
        this.state.finRonda = true;
        for (let i of this.state.players[0].state.mano) {
            if (i != null) {
                this.state.finRonda = false;
                break;
            }
        }

        if (this.state.finRonda) {
            this.state.players[this.state.orden[0]].state.puntos += 10;
            this.terminarRonda();
            // Si no es online, barajamos y reportimos
            if (!this.state.esOnline) {
                this.barajarYRepartir();
            } else {
                // * Si es online, 
                // * - se espera a recibir la baraja del server
                // * - se reparte a los jugadores
            }
            return;
        } // 10 ultimas

        if (!this.state.segundaBaraja) return;
        if (this.state.numPlayers === 4) {
            if (this.state.players[0].state.puntos + this.state.players[2].state.puntos > 100) {
                this.state.ganador = 1;
                this.state.finJuego = true;
            }
            if (this.state.players[1].state.puntos + this.state.players[3].state.puntos > 100) {
                this.state.ganador = 2;
                this.state.finJuego = true;
            }
        } else {
            if (this.state.players[0].state.puntos > 100) {
                this.state.ganador = 1;
                this.state.finJuego = true;
            }
            if (this.state.players[1].state.puntos > 100) {
                this.state.ganador = 2;
                this.state.finJuego = true;
            }
        }
    }

    terminarRonda() {
        console.log("==============================");
        console.log("Pasamos a segunda baraja");
        console.log("==============================");
        this.state.segundaBaraja = false;
        if (this.state.numPlayers == 4) {
            if (this.state.players[0].state.puntos + this.state.players[2].state.puntos > 100) this.state.ganador = 1;
            else if (this.state.players[1].state.puntos + this.state.players[3].state.puntos > 100) this.state.ganador = 2;
            else this.state.segundaBaraja = true;
        }
        else {
            if (this.state.players[0].state.puntos > 100) this.state.ganador = 1;
            else if (this.state.players[1].state.puntos > 100) this.state.ganador = 2;
            else {this.state.segundaBaraja = true; console.log("Segunda baraja");}
        }
    }

    barajarYRepartir() {
        console.log("==============================");
        console.log("Reparto segunda baraja");
        console.log("==============================");
        this.state.arrastre = false;
        this.state.baraja.recogerCartas();

        this.state.baraja.barajar();
        for (let i = 0; i < this.state.numPlayers; i++)
        {
            for (let j = 0; j < 6; j++)
            {
                this.state.players[i].anyadirCarta(this.state.baraja.darCarta());
            }
            this.state.players[i].reset();
        }
        this.state.triunfo = this.state.baraja.darCarta();
        this.state.baraja.anyadirAlFinal(this.state.triunfo);

        this.state.turnManager.reset();


        this.state.orden = Array(this.state.numPlayers).fill(null);
        for (let i = 0; i < this.state.numPlayers; i++)
        {
            this.state.orden[i] = i;
        }
    }

    setTriunfo(newTriunfo) {
        this.state.triunfo = newTriunfo;
    }
}

export default GameManager;