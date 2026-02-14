-- TestBench Template 

  LIBRARY ieee;
  USE ieee.std_logic_1164.ALL;
  USE ieee.numeric_std.ALL;

  ENTITY testbench IS
  END testbench;

  ARCHITECTURE behavior OF testbench IS 

  -- Component Declaration
          COMPONENT gbdt is
    		Generic(DATA_LENGTH: natural);-- number of bits of the inputs
			port(	-- Generic control signals
	         	Clk	:   in std_logic;
	        	Reset: in std_logic;
	    	    Start: in std_logic;
	    	    Done: out std_logic;
				--features vector
	    	    Din:  in std_logic_vector(8 * DATA_LENGTH - 1 downto 0); 
	        	-- Output
	    	    Dout: out std_logic_vector(DATA_LENGTH - 1 downto 0)); 
			END COMPONENT;

          SIGNAL clk, reset, start, Done :  std_logic;
          SIGNAL Din :  std_logic_vector(8*8 -1 downto 0);
          SIGNAL Dout :  std_logic_vector(7 downto 0);
          
  -- Clock period definitions
   constant CLK_period : time := 10 ns;
  BEGIN

  -- Component Instantiation
   uut: gbdt 
   		Generic MAP (DATA_LENGTH => 8)
   		PORT MAP(clk => clk, reset => reset, start => start, Done => Done, Din => Din, Dout => Dout);

-- Clock process definitions
   CLK_process :process
   begin
		CLK <= '0';
		wait for CLK_period/2;
		CLK <= '1';
		wait for CLK_period/2;
   end process;

 stim_proc: process
   begin		
    reset <= '1';
    Din <= X"A0B1C30E00F10200"; --Salida final = "11100100" 
    wait for CLK_period*2;
	reset <= '0';
	start <= '1';
	wait for CLK_period;
	start <= '0';
	wait until (Done = '1');
	wait for CLK_period*2;
	Din <= X"0A0B0C131E00F1F2"; --Salida final = "00011011" 
	start <= '1';
	wait for CLK_period;
	start <= '0';
	wait;
   end process;

  END;
