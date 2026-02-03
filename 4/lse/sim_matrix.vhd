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

    -- Posición inicial de Reversi, debería haber cuatro aperturas posibles
     -- Blancas en posiciones (3, 3) y (4, 4)
    board(16*3 + 3*2 + 1 downto 16*3 + 3*2) <= "10";
    board(16*4 + 4*2 + 1 downto 16*4 + 4*2) <= "10";

    -- Negras en posiciones (3, 4) y (4, 3)
    board(16*3 + 4*2 + 1 downto 16*3 + 4*2) <= "01";
    board(16*4 + 3*2 + 1 downto 16*4 + 3*2) <= "01";

    -- Blanca en (2, 2) para probar diagonales
    board(16*2 + 2*2 + 1 downto 16*2 + 2*2) <= "01";

    --finalmente pones un wait para que el resultado se vea
    wait for 10 ns;
    -- Poned unas cuantas casillas con sentido y comprobad que la salida es correcta
    -- Salida esperada:
    -- Negras en posiciones (2, 3), (3, 2), (4, 5), (5, 4) y (5, 5)
    -- Bits de salida a '1' -> 19, 26, 37, 44, 45 
    wait;
end process;

end Behavioral;
