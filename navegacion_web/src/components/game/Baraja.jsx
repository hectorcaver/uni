import React, { useState } from "react";
import { useUser } from '../../context/UserContext';

const Baraja = ({ controller }) => {
    
    const { cartas } = useUser();
    const spriteSrc = `/assets/stacks/${cartas}.png`;

    return (
        <div className="baraja">
            <img src={spriteSrc} alt={`Baraja`} />
        </div>
    );
}

export default Baraja;
