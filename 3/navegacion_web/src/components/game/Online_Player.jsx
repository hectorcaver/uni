import React, { useEffect } from "react";
import Carta from "./Carta";
import '/src/styles/Game.css';

const Online_Player = ({ controller, numPlayer, cartaJugada, handleCartaClick, handleCambiarSiete, handleCantar }) => {
  const spriteSrc = `/src/assets/Mano.png`;
  const esMiTurno = controller.state.esMiTurno;

  const esperar = (ms) => new Promise((resolve) => setTimeout(resolve, ms)); // Espera

  useEffect(() => {
    const ejecutarTurnoOnline = async () => {
      if (esMiTurno) {
        // RESET
        controller.reset();
        await esperar(1000);

        if(controller.state.ganador){
          // CANTAR
          controller.intentarCantar();
          if (controller.state.cantadoEsteTurno) {
            console.log("CANTANDO IA");
            handleCantar(controller.state.paloCantadoEsteTurno);
            await esperar(1000);
          }

          // CAMBIO SIETE
          controller.intentarCambiarSiete();
          if (controller.state.sieteCambiado) {
            console.log("SIETE CAMBIADO IA");
            handleCambiarSiete();
            await esperar(1000);
          }
        }

        // JUGAR CARTA
        let index = controller.turnoLogic();
        handleCartaClick(index);
      }
    };

    ejecutarTurnoIA(); // Llamar a la función asíncrona
  }, [esMiTurno, controller.state.gameManager.state.turnManager.state.playerTurn]);

  return (
    <>
      <div className={"cartaJugadaIA_" + numIA}>
        {cartaJugada && (
          <Carta
            id={cartaJugada.id}
            key={cartaJugada.id}
            palo={cartaJugada.palo}
            numero={cartaJugada.numero}
            callbackClick={() => { }}
            enMano={false}
          />
        )}
      </div>
      <div className={"manoIA_" + numIA}>
        <img src={spriteSrc} alt='Mano' />
        <div className={"cartasIA_" + numIA}>
        {controller.state.mano.map((carta, index) => (
            carta && (
              <div key={index} className={"carta " +  index}>
                {esMiTurno ? (
                  <Carta
                    id={carta.palo + "_" + carta.numero}
                    key={carta.palo + "_" + carta.numero}
                    palo={carta.palo}
                    numero={carta.numero}
                    callbackClick={() => {}}
                    enMano={true}
                  />
                ) : (
                  <Carta
                    id={carta.palo + "_" + carta.numero}
                    key={carta.palo + "_" + carta.numero}
                    palo={carta.palo}
                    numero={carta.numero}
                    callbackClick={() => {}}
                    enMano={false}
                  />
                )}
              </div>
            )
          ))}
        </div>
      </div>
    </>
  );
};

export default Online_Player;
