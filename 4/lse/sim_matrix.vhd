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
    wait for 1 ns;
   -- siguiendo este esquema pod�is poner las casillas que quer�is:
   -- primero se pone el tablero vacio 
    board <= (others => '0');
   -- despu�s colocas las casillas una a una
   -- en este ejemplo pongo una casilla negra en la posici�n (4,3)
    board(16*4 + 3*2 +1  downto 16*4 + 3*2) <= "10";
    -- en este una casilla blanca en la posici�n (4,2)
    board(16*4 + 2*2 +1  downto 16*4 + 2*2) <= "01";
    --finalmente pones un wait para que el resultado se vea
    wait for 1 ns;
    -- Poned unas cuantas casillas con sentido y comprobad que la salida es correcta
    wait;
end process;

end Behavioral;
