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
    Port (  
        info_in  : in  STD_LOGIC_VECTOR (1 downto 0);
        celda    : in  STD_LOGIC_VECTOR (1 downto 0);
        info_out : out STD_LOGIC_VECTOR (1 downto 0);
        SALIDA   : out std_logic
    );
    end component;

    -- Señales para recoger la SALIDA de cada micro_celda
    signal s : std_logic_vector(7 downto 0);

begin

    -- 1. Arriba Izquierda -> Abajo Derecha
    mc_ul: micro_celda port map (
        info_in  => up_left_in,
        celda    => input,
        info_out => down_right_out,
        SALIDA   => s(0)
    );

    -- 2. Arriba -> Abajo
    mc_u: micro_celda port map (
        info_in  => up_in,
        celda    => input,
        info_out => down_out,
        SALIDA   => s(1)
    );

    -- 3. Arriba Derecha -> Abajo Izquierda
    mc_ur: micro_celda port map (
        info_in  => up_right_in,
        celda    => input,
        info_out => down_left_out,
        SALIDA   => s(2)
    );

    -- 4. Izquierda -> Derecha
    mc_l: micro_celda port map (
        info_in  => left_in,
        celda    => input,
        info_out => right_out,
        SALIDA   => s(3)
    );

    -- 5. Derecha -> Izquierda
    mc_r: micro_celda port map (
        info_in  => right_in,
        celda    => input,
        info_out => left_out,
        SALIDA   => s(4)
    );

    -- 6. Abajo Izquierda -> Arriba Derecha
    mc_dl: micro_celda port map (
        info_in  => down_left_in,
        celda    => input,
        info_out => up_right_out,
        SALIDA   => s(5)
    );

    -- 7. Abajo -> Arriba
    mc_d: micro_celda port map (
        info_in  => down_in,
        celda    => input,
        info_out => up_out,
        SALIDA   => s(6)
    );

    -- 8. Abajo Derecha -> Arriba Izquierda
    mc_dr: micro_celda port map (
        info_in  => down_right_in,
        celda    => input,
        info_out => up_left_out,
        SALIDA   => s(7)
    );

    -- La salida final es '1' si CUALQUIERA de las direcciones genera un movimiento válido
    output <= s(0) or s(1) or s(2) or s(3) or s(4) or s(5) or s(6) or s(7);
               
end Behavioral;
