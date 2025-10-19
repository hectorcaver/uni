% Public Key Infraestructure (PKI) \
  Seguridad Informática
% 16/10/2025

\newpage

# Ataque Man-In-The_Middle (MITM)

El problema fundamental se da porque cuando recibes una clave pública de X, no sabes si realmente es de X o de otra persona.

La manera de solucionar este problema es tener un tercero de confianza que verifique la identidad en base a la clave pública y se asegure que el certificado no se pueda modificar ni fabricar (mediante firma digital).

Se debe obtener un certificado de un tercero:

- Para ello se debe acudir a dicho tercero.
- El tercero verifica la identidad de X y la asocia a su clave pública.
- X envía el certificado a Y.
- Y verifica el certificado usando la clave pública del tercero de confianza.


