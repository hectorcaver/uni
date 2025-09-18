import React, { useEffect } from "react";
import Carta from "./Carta";
import '/src/styles/Game.css';

const IA_Player = ({ controller, numIA, cartaJugada, handleCartaClick, handleCambiarSiete }) => {
  const spriteSrc = `/assets/Mano.png`;
  const esMiTurno = controller.state.esMiTurno;

  useEffect(() => {
    if (esMiTurno) {
      let index = controller.turnoLogic();
      handleCartaClick(index);
    }
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

export default IA_Player;
