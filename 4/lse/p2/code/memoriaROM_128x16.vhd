----------------------------------------------------------------------------------
-- Company: 
-- Engineer: 
-- 
-- Create Date:    14:12:11 04/04/2014 
-- Design Name: 
-- Module Name:    memoriaRAM - Behavioral 
-- Project Name: 
-- Target Devices: 
-- Tool versions: 
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
use ieee.std_logic_unsigned.all;
use ieee.numeric_std.all;
-- Uncomment the following library declaration if using
-- arithmetic functions with Signed or Unsigned values
--use IEEE.NUMERIC_STD.ALL;

-- Uncomment the following library declaration if instantiating
-- any Xilinx primitives in this code.
--library UNISIM;
--use UNISIM.VComponents.all;
-- Memoria RAM de 128 palabras de 16 bits
entity memoriaROM_128x16 is port (
		  	ADDR : in std_logic_vector (6 downto 0); --Dir 
        	Dout : out std_logic_vector (15 downto 0));
end memoriaROM_128x16;

-- Esta RAM está inicializada con 4 árboles tipo GDBT, cada uno de ellos procesa dos features y pone 1  en el bit correspondiente de la salida si su valor es mayor que x"0F"
-- Ejemplo: 
-- Entrada= X"A0B1C30E00F10200"; Salida final = "11100100" 

architecture Behavioral of memoriaROM_128x16 is
type RomType is array(0 to 127) of std_logic_vector(15 downto 0);
signal ROM : RomType := (  			"0"&"001"&x"4"&x"0F", --@00 first tree
									"0"&"000"&x"2"&x"0F",
									"1"&"000"&x"5"&x"03",
									"1"&"000"&x"4"&x"02",
									"0"&"000"&x"2"&x"0F",
									"1"&"000"&x"2"&x"01",
									"1"&"000"&x"1"&x"00",
									"0"&"011"&x"4"&x"0F", --@07 second tree
									"0"&"010"&x"2"&x"0F",
									"1"&"000"&x"5"&x"0C",
									"1"&"000"&x"4"&x"08",
									"0"&"010"&x"2"&x"0F",
									"1"&"000"&x"2"&x"04",
									"1"&"000"&x"1"&x"00",
									"0"&"101"&x"4"&x"0F", --@0E third tree
									"0"&"100"&x"2"&x"0F",
									"1"&"000"&x"5"&x"30",
									"1"&"000"&x"4"&x"20",
									"0"&"100"&x"2"&x"0F",
									"1"&"000"&x"2"&x"10",
									"1"&"000"&x"1"&x"00",
									"0"&"111"&x"4"&x"0F", --@15fourth tree
									"0"&"110"&x"2"&x"0F",
									"1"&"000"&x"F"&x"C0", --final node: "F" means that it is the last node of the last tree
									"1"&"000"&x"F"&x"80", --final node: "F" means that it is the last node of the last tree
									"0"&"110"&x"2"&x"0F",
									"1"&"000"&x"F"&x"40", --final node: "F" means that it is the last node of the last tree
									"1"&"000"&x"F"&x"00", --final node: "F" means that it is the last node of the last tree
									X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000",
									X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000", X"0000");
begin
 
 Dout <= ROM(conv_integer(addr));

end Behavioral;

