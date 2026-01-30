library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

-- Este dise�o es para un tama�o fijo de 8x8 pero ser�a trivial hacerlo para un tama�o variable, 
-- bastar�a con definir dos par�metros gen�ricos que indicasen el tama�o de cada dimensi�n y usarlos en los generates
-- Notaci�n: los dos bits menos significativos del tablero son la casilla (0,0), los siguientes la (0,1)...
-- Igualmente el bit menos significativo de Mov_posibles es la salida de la casilla (0,0)...

entity macro_celda is
    Port (  up_in : in STD_LOGIC_VECTOR (1 downto 0);
           down_in : in STD_LOGIC_VECTOR (1 downto 0);
           left_in : in STD_LOGIC_VECTOR (1 downto 0);
           right_in : in STD_LOGIC_VECTOR (1 downto 0);
           up_left_in : in STD_LOGIC_VECTOR (1 downto 0);
           up_right_in : in STD_LOGIC_VECTOR (1 downto 0);
           down_left_in : in STD_LOGIC_VECTOR (1 downto 0);
           down_right_in : in STD_LOGIC_VECTOR (1 downto 0);
           up_out : out STD_LOGIC_VECTOR (1 downto 0);
           down_out : out STD_LOGIC_VECTOR (1 downto 0);
           left_out : out STD_LOGIC_VECTOR (1 downto 0);
           right_out : out STD_LOGIC_VECTOR (1 downto 0);
           up_left_out : out STD_LOGIC_VECTOR (1 downto 0);
           up_right_out : out STD_LOGIC_VECTOR (1 downto 0);
           down_left_out : out STD_LOGIC_VECTOR (1 downto 0);
           down_right_out : out STD_LOGIC_VECTOR (1 downto 0);
           input : in STD_LOGIC_VECTOR (1 downto 0);
           output : out std_logic);
end macro_celda;

architecture Behavioral of macro_celda is
-- Os pongo las entradas de la macrocelda. Recibe informaci�n de las ocho entradas y 
-- a partir de cada entrada y la celda 8input) genera la salida contraria y los movimientos posibles
component micro_celda 
Port (  info_in : in STD_LOGIC_VECTOR (1 downto 0);
        celda : in STD_LOGIC_VECTOR (1 downto 0);
        info_out : out STD_LOGIC_VECTOR (1 downto 0);
        SALIDA : out std_logic);
end component;

-- signal matrix es un tipo que nos permite definir las matrices de se�ales que van a conectar nuestra red.
-- Su tama�o es uno mayor que el n�mero de celdas porque tambi�n incluimos las entradas de las fronteras
-- Esas entradas ser�n siempre "00" 
type signal_wire is STD_LOGIC_Vector(1 downto 0);
signal wire:  signal_wire;

begin
    
    microc_up_l: micro_celda port map (
        up_left_in => info_in,
        info_out => down_right_out,
    );

    microc_up_l: micro_celda port map (
        up_left_in => info_in,
        info_out => down_right_out,
    );

    microc_up_l: micro_celda port map (
        up_left_in => info_in,
        info_out => down_right_out,
    );

    microc_up_l: micro_celda port map (
        up_left_in => info_in,
        info_out => down_right_out,
    );

    microc_up_l: micro_celda port map (
        up_left_in => info_in,
        info_out => down_right_out,
    );

    microc_up_l: micro_celda port map (
        up_left_in => info_in,
        info_out => down_right_out,
    );

    microc_up_l: micro_celda port map (
        up_left_in => info_in,
        info_out => down_right_out,
    );

    microc_up_l: micro_celda port map (
        up_left_in => info_in,
        info_out => down_right_out,
    );
               
end Behavioral;
