HEAD = cd .. && ./genPDF.sh -u

T1: T1_Interconexión_de_Redes_IPv4.md
	${HEAD} T1_Interconexión_de_Redes_IPv4

P11: P11_Diseño_Gestion_Escenarios_IPv4.md
	${HEAD} P11_Diseño_Gestion_Escenarios_IPv4


add:
	@echo "Añadiendo entrada para $(NAME)..."
	@echo "" >> Makefile
	@echo "$(NAME): $(NAME).md" >> Makefile
	@printf "\t\$${HEAD} $(NAME)\n" >> Makefile
	@echo "Entrada añadida con éxito."


P12: P12_Diseño_Gestion_Escenarios_IPv4.md
	${HEAD} P12_Diseño_Gestion_Escenarios_IPv4
