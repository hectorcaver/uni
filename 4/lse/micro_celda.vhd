----------------------------------------------------------------------------------
-- Company: 
-- Engineer:  Pablo Plumed (874167) y Hector Lacueva (869637)
-- 
-- Create Date: 1.02.2026 19:50
-- Design Name: 
-- Module Name: micro_celda - Behavioral
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
use IEEE.STD_LOGIC_1164.all;

entity micro_celda is
  port (
    info_in  : in std_logic_vector (1 downto 0);
    celda    : in std_logic_vector (1 downto 0);
    info_out : out std_logic_vector (1 downto 0);
    SALIDA   : out std_logic
  );
end micro_celda;

architecture Behavioral of micro_celda is
begin

-- Si la celda está vacía ("00") y se han encontrado fichas blancas después de una neegra ("10"),
-- entonces se puede poner una ficha negra a esa posición, por lo que SALIDA se pone a '1'.
  SALIDA <= '1' when (info_in = "10" and celda = "00") else
    '0';

-- Devuelve "estado nada" si la celda está vacía o si le llega "estado nada" a un celda blanca ("10").
-- Devuelve "estado hay_negro" si la celda es negra ("01").
-- Devuelve "estado negro_blanco" si la celda es blanca y le llega "estado hay_negro" o "estado negro_blanco".
  info_out <= "00" when (celda = "00") or (info_in = "00" and celda = "10") else
              "01" when (celda = "01")                                      else
              "10" when (celda = "10" and (info_in = "01" or info_in = "10"));

end Behavioral;
