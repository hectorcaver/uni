% Interconexión de redes IPv4\
  Diseño y Administración de redes
% Héctor Lacueva Sacristán
% 9/2025

\newpage

# Internet 

Internet es un conjunto mundial de redes interconectadas con protocolos
comunes (TCP/IP) y un direccionamiento universal (IP).

Cada red se incorpora voluntariamente a Internet, y se gestiona de
manera autónoma. Sin embargo, existen cierta organización:

- ISOC (Internet Society): asociación internacional para la promoción de la tecnología y servicios de internet. (IAB, IETF, IRTF).
- IANA (Internet assigned Numbers Authority), ahora ICANN (Internet Corporation of Assigned direccionamientoes and Numbers).

# Direccionamiento IPv4

## Dirección IP

Son identificadores virtuales universales. Identifican una conexión de un nodo. **Un nodo tendrá tantas direcciones IP como interfaces de conexión a red**. Son interpretados por el software e independiente del direccionamiento hardware.

- **CLASSFUL**: estructura de clases que define el tamaño de red y número de hosts posibles. Los algoritmos de encaminamiento determinan la dirección de red en función de la Classful_network. Para más info mirar en [link](https://en.wikipedia.org/wiki/Classful_network).
- **CLASSLESS**: en este caso, las **máscaras** definen el tamaño de red y número de hosts posibles. Los algoritmos de encaminamiento usan la máscara para identificar la dirección de red.  

Internet funciona con CIDR, Classless Inter Domain Routing. La notación de red es la siguiente, **X.X.X.X/P**, donde **X** son números entre 0-255 y **P** es el **prefijo**, esto representa el número de 1's de la máscara de red, de esta forma, direcciones consecutivas pueden agregarse con un prefijo común (para reducir el número de entradas de las tablas de encaminamiento). Por ejemplo, el acceso a la red 200.25.16.0/20, puede implicar acceder a las subredes 200.25.16.0/21, 200.25.24.0/22, ...

### Direcciones reservadas

|Red o rango|Uso|
|:-:|:-:|
| 0.0.0.0          | Sin especificar (arranque)               |
| 000...000.hostid | Uso en arranque                          |
| netid.000...000  | @ de red (netid)                         |
| netid.111...111  | Difusión (todos los nodos de netid)      |
| 255.255.255.255  | Difusión limitada (arranque, red física) |
| 127.X.Y.Z        | *loopback* (uso en pruebas)              |
| 127.0.0.0        | Reservado (fin clase A)                  |
| 128.0.0.0        | Reservado (inicio clase B)               |
| 191.255.000      | Reservado (fin clase B)                  |
| 192.0.0.0        | Reservado (inicio clase C)               |
| 224.0.0.0        | Reservado (inicio clase D)               |
| 240.0.0.0 - 255.255.255.254 | Reservado (clase E)           |
| 10.0.0.0         | Privado (clase A)                        |
| 172.16.0.0 - 172.31.0.0 | Privado (clase B)                 |
| 192.168.0.0 - 192.168.255.0 | Privado (clase C)             |


## Asignación de direcciones 

Inicialmente NIC (Network Information Center), en los 90s RIR (Regional Internet Registry) y, por último, delegación en LIRs (Local Internet Registry):

- ARIN (American Registry for Internet Numbers): EEUU y Canadá.
- APNIC (Asia Pacific Network Information Centre): Asia Oriental, Pacífico
- RIPE (Réseaux IP Européenes): Europa, Oriente Medio, Asia Centre
- LACNIC (Latin American and Caribbean Network Information Centre): América (excepto EEUU y Canadá) y el Caribe.
- AFRINIC (African Network Information Center): África, Oceano Indico.

Direcciones agotadas. Siguiente paso IPv6.


## Encaminamiento

Sirve para definir el camino o ruta a seguir por los datagramas, a través de una o más redes, para que estos alcancen su destino.

- **Encaminamiento directo**: llegar a "su propia red"
  - Asociación @IP - @PHY.
  - Resolución de dirección (ej. ARP).
- **Encaminamiento indirecto**: llegar a "otra red": siguiente salto (router).
  - Se traduce en encaminamiento directo hacia dicho siguiente salto.
  - Tabla de encaminamiento:
  
  |Red destino|Máscara|Métrica|Siguiente salto|Interfaz de salida|
  |:-:|:-:|:-:|:-:|:-:|

### Host 

- Tiene una tabla con **su dirección de red (directo) y el router por defecto**.
- Si recibe un datagrama que no es para él, lo descarta.
- Si desea enviar un datagrama:
  - **A su propia red $\rightarrow$ ARP**.
  - **A otra red destino** $\rightarrow$ siguiente salto: router "por defecto".

### Router

- Tiene una tabla que contiene **todas las redes a las que está conectado, redes externas + router por defecto**.
- Si reciben un datagrama que no es para él, **intenta reenviarlo consultando su tabla de rutas**.
- Si tienen que reenviar:
  - En una **red propia**: ARP.
  - En una **red externa** se consulta la tabla de enrutamiento y se compara @IP/máscara:
    - Existe una **única entrada** a la red destino: **enviar**.
    - Existen **varias entradas** a la red destino: enviar a la **métrica menor**.
    - **No existe** la entrada explícita: **enviar al router por defecto**.

> Si hay coincidencia con varias, la de máscara más larga (más bits coincidentes) $\rightarrow$ long prefix match.

### Construcción de la tabla de rutas

- Dependiendo del tamaño o complejidad de Internet:
  - **Encaminamiento estático**:
    - Rutas fijas establecidas durante el arranque (boot)
    - Útil en casos muy simples, cuando los cambios de encaminamiento son lentos y poco frecuentes.
  - **Encaminamiento dinámico**:
    - Inicialización en arranque y actualización por protocolo.
    - Protocolos de encaminamiento (intercambio de información entre los router)
    - Necesario en grandes redes, con cambios frecuentes y rápidos.

- En definitiva, dos fuentes de información:
  - **Inicialización** (ej. de disco) $\rightarrow$ Host normalmente "congelan" la tabla tras inicializar.
  - **Actualización** (ej. a partir de protocolos) $\rightarrow$ Los router aprenden información nueva y actualizan las tablas.

### Entradas "especiales" de la tabla de rutas

- **Ruta basada en host** (Host-Specific)
  - Se corresponde con un valor completo de 32-bit: @IP de host, no de red.
  - Se puede utilizar para enviar tráfico a un host específico a través de un camino concreto.
- **Ruta basada en net** (Host-Specific)
  - Se corresponde con un valor de red y su máscara correspondiente.
  - Se puede utilzar para enviar tráfico a una net específica a través de un camino concreto.
- **Ruta por defecto** (default) 
  - Únicamente se permite una entrada de este tipo.
  - Se corresponde con "cualquier" dirección de destino.
  - Únicamente se utiliza si no hay otra correspondencia en la tabla.

## Funcionalidad del protocolo IPv4

### PDU (Protocol Data Unit): Datagrama IP

- **NO ORIENTADO A CONEXIÓN: QoS Best effort**
  - No hay conexión de extremo a extremo (sin establecimiento, sin información de estado...)
  - Cada paquete tratado de forma individual (encaminamiento según dirección destino)
  - No hay garantías de entrega.
- Tamaño máximo / mínimo = 65535 / 28 (20 cabecera completa + 8 de fragmento mínimo).
- Tamaño mínimo MTU recomendable = 576 bytes.
- Tamaño variable de cabecera: opciones.

Para más información del protocolo o del datagrama buscar información en los apuntes (diapositiva 23-25) o buscar en internet.

El comando `ping` tiene varias opciones:

- **Record route** (`ping -r`): anota en la cabecera IP la ruta seguida por el datagrama.
- **Timestamp** (`ping -s`): anota la ruta y la marca de tiempo de cada salto.
- **Strict source routing** (`ping -k`): la cabecera contine la ruta paso a paso que debe seguir el datagrama.
- **Loose source routing** (`ping -j`): la cabecera lleva una lista de router por los que debe pasar el datagrama, pero puede pasar por otros.

### Fragmentación y reensamblado



