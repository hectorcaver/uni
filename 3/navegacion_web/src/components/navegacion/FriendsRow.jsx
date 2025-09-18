/**
 * @file FriendsRow.jsx
 * @description Componente que representa una fila de amigos en la tabla de amigos.
 * 
 * Este componente muestra la imagen de perfil, el nombre de usuario y un botón de opciones.
 * Al hacer clic en el botón de opciones, se muestra un modal con varias opciones para el amigo.
 * 
 * @param {string} img - URL de la imagen de perfil del amigo.
 * @param {string} username - Nombre de usuario del amigo.
 * @param {string} mail - Correo electrónico del amigo.
 * @param {function} onDelete - Callback para notificar al componente padre cuando se elimina un amigo.
 * 
 * @returns {JSX.Element} Componente de fila de amigos.
 */

import React, { useState, useRef } from 'react';
import '/src/styles/FriendsRow.css';
import usePost from '../../customHooks/usePost';

const assetsUrl = '/src/assets/';
const avataresUrl = '/src/assets/avatares/';

const {postData}  = usePost('https://guinyoteonline-hkio.onrender.com/amigos/eliminarAmigo/');

const FriendsRow = ({ img, username, mail, usrMail, onDelete }) => {
    const [showModal, setShowModal] = useState(false);
    const [isDeleted, setIsDeleted] = useState(false);
    const [modalPosition, setModalPosition] = useState({ left: 0, top: 0 });
    const buttonRef = useRef(null);

    const handleOnClick = () => {
        if (buttonRef.current) {
            const rect = buttonRef.current.getBoundingClientRect();
            setModalPosition({
                left: rect.left - 150, // Ajusta la posición a la izquierda del botón
                top: rect.top
            });
        }
        setShowModal(!showModal);
    };

    const onClickProfile = () => {
        console.log("Pulsado perfil");
        setShowModal(false);
    }

    const onClickDeleteFriend = async () => {
        const { error } = await postData({ idEliminador: usrMail, idEliminado: mail }, '');
        if (error) {
            alert("Error al eliminar amigo " + username);
            return;
        }
        alert("Amigo " + username + " eliminado correctamente");
        setIsDeleted(true);
        onDelete(mail); // Notify parent to remove the row
    };

    const onClickInviteGroup = () => {
        console.log("Pulsado invitar a grupo");
        setShowModal(false);
    }

    if(isDeleted) {
        return null; // Si el amigo ha sido eliminado, no se muestra nada
    }

    return (
        <>
            <tr key={mail} className='friend-row'>
                <td><img src={avataresUrl + img} alt="avatar" /></td>
                <td><p>{username}</p></td>
                <td>
                    <button onClick={handleOnClick} ref={buttonRef} className='friend-row-options'>
                        <img src={assetsUrl + "options.png"} alt="options.png" />
                    </button>
                </td>
            </tr>

            {showModal && (
                <div className="modal" style={{ left: modalPosition.left, top: modalPosition.top }}>
                    <p><b>Opciones de {username}</b></p>
                    <button onClick={onClickProfile} className="btn-profile">Ver perfil</button>
                    <button onClick={onClickInviteGroup} className="btn-invite-friend">Invitar a grupo</button>
                    <button onClick={onClickDeleteFriend} className="btn-delete-friend">Eliminar amigo</button>
                </div>
            )}
        </>
    );
};

export default FriendsRow;


