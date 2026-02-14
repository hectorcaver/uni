----------------------------------------------------------------------------------
-- Company: Universidad de Zaragoza
-- Engineer: Pablo Plumed (874167) y Hector Lacueva (869637) (ft. Javier Resano)
-- 
-- Create Date:    22/02/2022 
-- Design Name: 
-- Module Name:    gbdt 
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
use IEEE.std_logic_arith.all;
use IEEE.std_logic_unsigned.all;
-- Uncomment the following library declaration if using
-- arithmetic functions with Signed or Unsigned values
-- use IEEE.NUMERIC_STD.ALL;

-- Uncomment the following library declaration if instantiating
-- any Xilinx primitives in this code.
--library UNISIM;
--use UNISIM.VComponents.all;

entity gbdt is
    Generic(DATA_LENGTH: natural);-- number of bits of the inputs
	port(	-- Generic control signals
         	Clk	:   in std_logic;
        	Reset: in std_logic;
    	    Start: in std_logic;
    	    Done: out std_logic;
			--features vector (Entradas)
    	    Din:  in std_logic_vector(8 * DATA_LENGTH - 1 downto 0); 
        	-- Output (Salidas)
    	    Dout: out std_logic_vector(DATA_LENGTH - 1 downto 0)); 
end gbdt;

architecture Behavioral of gbdt is

component memoriaROM_128x16 is port ( 
		  	ADDR : in std_logic_vector (6 downto 0); --Dir 
        	Dout : out std_logic_vector (15 downto 0));
end component;

component reg is
        generic(BITS: positive);
        port(-- Control signals
             Clk:   in std_logic;
             Reset: in std_logic;
             Load:  in std_logic;         
             -- Input
             Din: in std_logic_vector(BITS - 1 downto 0);             
             -- Output
             Dout: out std_logic_vector(BITS - 1 downto 0));
    end component;
    
component mux8_1 is
    Generic(DATA_LENGTH: natural);-- number of bits of the inputs
	port(Ctrl: in std_logic_vector(2 downto 0);
         Din:  in std_logic_vector(8 * DATA_LENGTH - 1 downto 0);
         Dout: out std_logic_vector(DATA_LENGTH - 1 downto 0));
    end component;
    
-- Señales orientativas. Podéis usarlas o crear las vuestras
-- Si alguna no la usais la podéis quitar
signal node : STD_LOGIC_VECTOR (15 downto 0);
signal feature_addr: STD_LOGIC_VECTOR (2 downto 0);
signal addr_distance: STD_LOGIC_VECTOR (3 downto 0);
signal comparisom_value, Feature_selected, leaf_value, Int_Dout, addition: STD_LOGIC_VECTOR (7 downto 0);
signal curr_addr, next_addr, left_addr, right_addr: STD_LOGIC_VECTOR (6 downto 0);
signal next_node, reg_reset, node_type, internal_done, rom_finished, load_addr, Trees_finished, load_output, int_reset, cmp_result: std_logic;
type state_type is (Processing, Finished);
signal state, next_state : state_type; 
begin

    -- Añadido    internal_done <= addr_distance = "1111" 
    -- La señal de start resetea los registros
    int_reset <= start;
    -- Fin añadido

    -- Reset for the registers
    reg_reset <= reset or int_reset;

    -- Addr register
    node_counter: reg
        generic map(BITS => 7)
        port map(Clk   => Clk,
                 Reset => reg_reset,
                 Load  => load_addr,
                 Din   => next_addr,
                 Dout  => curr_addr);

	GBDT_ROM: memoriaROM_128x16  port map (addr=> curr_addr, dout=> node);
	-- common fields
	addr_distance <= node(11 downto 8);
	node_type <= node(15); --'0': Non-leaf node; '1': leaf node
	-- Non-leaf node fields
	comparisom_value <= node(7 downto 0);	
	feature_addr <= node(14 downto 12);
	-- Leaf node fields
	leaf_value <= node(7 downto 0);
-- Mux to select the corresponding feature
    features_mux: mux8_1
        generic map(DATA_LENGTH => 8)
        Port map(Ctrl => feature_addr,
                 Din  => Din,
                 Dout => Feature_selected);


-- CMP
    -- Cmp result = 1 if entry <= Value
    --cmp_result <= '1' when (Feature_selected <= comparisom_value) else '0';
    cmp_result <= '1' when (IEEE.std_logic_unsigned."<="(Feature_selected, comparisom_value)) else '0';
    -- Addr left node (curr_addr + 1)
    left_addr <= (curr_addr + 1);
    -- Addr right node (curr_addr + addr_distance)
    right_addr <= (curr_addr + addr_distance);
    -- Logica para saber si un árbol ya ha acabado
    Trees_Finished <= node_type;
    -- Addr logic: In these trees the left son is for the greater condition, and the right son for lesser equal.
    next_addr <= left_addr when (cmp_result = '0' and Trees_Finished = '0') else right_addr;
    -- Accumulating the trees results
    addition <= (leaf_value + Int_Dout);
    -- Load output (if node_type = 1 (leaf), update reg with addition)
    load_output <= (node_type);
    -- Si addr_distance es "1111" --> Done = 1
    --internal_done <= '1' when (addr_distance = "1111") else '0';
    rom_finished <= '1' when (IEEE.std_logic_unsigned."="(addr_distance, "1111")) else '0';

-- Output register
output_reg: reg
        generic map(BITS => 8)
        port map(Clk   => Clk,
                 Reset => reg_reset,
                 Load  => load_output,
                 Din   => addition,
                 Dout  => Int_Dout);
	Dout <= Int_Dout; -- en vhdl no puedes leer las salidas. Si necesitas leerlas, declaras una señal interna que puedes leer, y la usas para generar la salida
	Done <= internal_done;

----------------------------
---------
--UC
-------------------------------------
-- state register
 SYNC_PROC: process (clk)
   begin
      if (clk'event and clk = '1') then
         if (reset = '1') then
            state <= Processing;
         else
            state <= next_state;
         end if;        
      end if;
   end process;

   --State-Machine
   OUTPUT_DECODE: process (state, Trees_finished, start)-- recordad que en la lista de sensibilidad deben estar todas las señales que provoquen que el proceso pueda cambiar sus salidas. Si metéis nuevas entradas, incluidlas. 
   begin
 -- valores por defecto, si no se asigna otro valor en un estado valdrán lo que se asigna aquí
 -- si necesitáis más señales añadirlas
		next_state <= state; 
		load_addr <= '0'; 
		internal_done <= '0';
		int_reset <= '0';
 	-- Estado Inicio          
 	-- Esta máquina de estados y sus salidas la tenéis que diseñar vosotros. Estas líneas son sólo un ejemplo.
        if (state = Processing) then
        	load_addr <= '1';
            if (rom_finished = '1') then
                next_state <= Finished;
            end if;
        else 
        	internal_done <= '1';

            if (start = '1') then
                int_reset <= '1';
                next_state <= Processing;
            end if;
     	end if;
  	end process;	     
end Behavioral;

