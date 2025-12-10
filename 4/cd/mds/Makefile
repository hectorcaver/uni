HEAD = cd .. && ./genPDF.sh -u

T22: T22_RISC-V.md
	${HEAD} T22_RISC-V

T23: T23_C_A+B.md
	${HEAD} T23_C_A+B

T24: T24_Potencia_Y_Paralelismo.md
	${HEAD} T24_Potencia_Y_Paralelismo

T41: T41_Certificación_Energy_Star.md
	${HEAD} T41_Certificación_Energy_Star

T43: T43_IBM_Telum2.md
	${HEAD} T43_IBM_Telum2

add:
	@echo "Añadiendo entrada para $(NAME)..."
	@echo "" >> Makefile
	@echo "$(NAME): $(NAME).md" >> Makefile
	@printf "\t\$${HEAD} $(NAME)\n" >> Makefile
	@echo "Entrada añadida con éxito."

T51: T51_Eficiencia2008.md
	${HEAD} T51_Eficiencia2008

T91: T91_Estudio_Mercado.md
	${HEAD} T91_Estudio_Mercado

T61: T61_Consumos_AWS_Aragón.md
	${HEAD} T61_Consumos_AWS_Aragón

T92: T92_RAIDS.md
	${HEAD} T92_RAIDS

T93_Articulo: T93_Articulo.md
	${HEAD} T93_Articulo
