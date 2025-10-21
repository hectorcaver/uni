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

