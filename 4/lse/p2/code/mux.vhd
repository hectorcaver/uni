-------------------------------------------------------------------------------
-- Generic multiplexer
-------------------------------------------------------------------------------

library IEEE;
use IEEE.std_logic_1164.all;
use IEEE.numeric_std.all;

entity mux8_1 is
    Generic(DATA_LENGTH: natural);-- number of bits of the inputs
	port(Ctrl: in std_logic_vector(2 downto 0);
         Din:  in std_logic_vector(8 * DATA_LENGTH - 1 downto 0);
         Dout: out std_logic_vector(DATA_LENGTH - 1 downto 0));
end mux8_1;

architecture Behavioral of mux8_1 is
begin
    Dout <= Din((to_integer(unsigned(Ctrl)) + 1) * DATA_LENGTH - 1
                downto to_integer(unsigned(Ctrl)) * DATA_LENGTH);
end Behavioral;

