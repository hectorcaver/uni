----------------------------------------------------------------------------------
-- Company: 
-- Engineer: 
-- 
-- Create Date: 18.02.2016 14:33:12
-- Design Name: 
-- Module Name: sim_matrix - Behavioral
-- Project Name: 
-- Target Devices: 
-- Tool Versions: 
-- Description: 
-- 
-- Dependencies: 
-- 
-- Revision:
-- Revision 0.01 - File Created
-- Additional Comments:
-- 
----------------------------------------------------------------------------------


library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

-- Uncomment the following library declaration if using
-- arithmetic functions with Signed or Unsigned values
--use IEEE.NUMERIC_STD.ALL;

-- Uncomment the following library declaration if instantiating
-- any Xilinx leaf cells in this code.
--library UNISIM;
--use UNISIM.VComponents.all;

entity sim_matrix is
    Port ( salida : out STD_LOGIC_VECTOR (63 downto 0));
end sim_matrix;

architecture Behavioral of sim_matrix is
component matriz_celdas
    Port (  tablero : in STD_LOGIC_VECTOR (127 downto 0);
            Mov_posibles: out STD_LOGIC_VECTOR (63 downto 0));
end component;
signal board: STD_LOGIC_VECTOR (127 downto 0);
begin
matriz: matriz_celdas port map (board, salida);
process
begin
    --tablero vacio
    board <= (others => '0');
    wait for 10 ns;
   -- siguiendo este esquema pod�is poner las casillas que quer�is:
   -- primero se pone el tablero vacio 
    board <= (others => '0');
   -- despu�s colocas las casillas una a una
   -- en este ejemplo pongo una casilla negra en la posici�n (4,3)

    -- Blancas (01) en cruz y diagonales desde (2,2) 
    board(16*1 + 2*2 + 1 downto 16*1 + 2*2) <= "01"; -- Arriba (1,2)
    board(16*3 + 2*2 + 1 downto 16*3 + 2*2) <= "01"; -- Abajo (3,2)
    board(16*2 + 1*2 + 1 downto 16*2 + 1*2) <= "01"; -- Izquierda (2,1)
    board(16*2 + 3*2 + 1 downto 16*2 + 3*2) <= "01"; -- Derecha (2,3)
    board(16*1 + 1*2 + 1 downto 16*1 + 1*2) <= "01"; -- Diag UL (1,1)
    board(16*3 + 3*2 + 1 downto 16*3 + 3*2) <= "01"; -- Diag DR (3,3)

    -- Negras (10) al final de cada línea para "cerrar" el flanqueo 
    board(16*0 + 2*2 + 1 downto 16*0 + 2*2) <= "10"; -- Cierra arriba (0,2)
    board(16*4 + 2*2 + 1 downto 16*4 + 2*2) <= "10"; -- Cierra abajo (4,2)
    board(16*2 + 0*2 + 1 downto 16*2 + 0*2) <= "10"; -- Cierra izquierda (2,0)
    board(16*2 + 4*2 + 1 downto 16*2 + 4*2) <= "10"; -- Cierra derecha (2,4)
    board(16*0 + 0*2 + 1 downto 16*0 + 0*2) <= "10"; -- Cierra UL (0,0)
    board(16*4 + 4*2 + 1 downto 16*4 + 4*2) <= "10"; -- Cierra DR (4,4)

    -- salida(2*8 + 2) DEBERÍA SER '1' porque se puede jugar en (2,2)

    --finalmente pones un wait para que el resultado se vea
    wait for 10 ns;
    -- Poned unas cuantas casillas con sentido y comprobad que la salida es correcta
    wait;
end process;

end Behavioral;
