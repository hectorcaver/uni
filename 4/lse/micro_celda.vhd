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

entity micro_celda is
    Port (  
        info_in  : in  STD_LOGIC_VECTOR (1 downto 0);
        celda    : in  STD_LOGIC_VECTOR (1 downto 0);
        info_out : out STD_LOGIC_VECTOR (1 downto 0);
        SALIDA   : out std_logic
    );
end micro_celda;

architecture Behavioral of micro_celda is
begin

    SALIDA <= '1' when (info_in = "10" and celda = "00") else '0';

    info_out <= "00" when (celda = "00") or (info_in = "00" and celda = "10") else
                "01" when (celda = "01") else
                "10" when (celda = "10" and (info_in = "01" or info_in = "10"));

end Behavioral;
