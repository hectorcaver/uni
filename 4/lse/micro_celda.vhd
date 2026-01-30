----------------------------------------------------------------------------------
-- Company: 
-- Engineer: 
-- 
-- Create Date: 17.02.2016 19:24:40
-- Design Name: 
-- Module Name: matriz_celdas - Behavioral
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

-- Este dise�o es para un tama�o fijo de 8x8 pero ser�a trivial hacerlo para un tama�o variable, 
-- bastar�a con definir dos par�metros gen�ricos que indicasen el tama�o de cada dimensi�n y usarlos en los generates
-- Notaci�n: los dos bits menos significativos del tablero son la casilla (0,0), los siguientes la (0,1)...
-- Igualmente el bit menos significativo de Mov_posibles es la salida de la casilla (0,0)...

entity micro_celda is
    Port (  info_in : in STD_LOGIC_VECTOR (1 downto 0);
        celda : in STD_LOGIC_VECTOR (1 downto 0);
        info_out : out STD_LOGIC_VECTOR (1 downto 0);
        SALIDA : out std_logic);
end micro_celda;

architecture Behavioral of micro_celda is
begin

    

end Behavioral;
