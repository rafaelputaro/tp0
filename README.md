# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Además, definir un script de bash `generar-compose.sh` que permita crear una definición de DockerCompose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```
#### Resolución:

Se creo un script llamado "generar-compose-sh" el cual llama a un script en python llamado "mi-generador.py" que genera efectivamente el docker-compose con la cantidad de clientes solicitada.

#### Ejecución:
```
. generar-compose.sh docker-compose-dev.yaml <nro clientes>
make docker-compose-up
docker ps
make docker-compose-logs
make docker-compose-down
```

***

### Ejercicio N°2:

Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (config.ini y config.yaml, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: docker volumes).

#### Resolución:

Modifique "mi-generador.py" con el objetivo de que cree volumes para cada tipo de servicio (cliente y servidor) sobre los archivos de configuración.

#### Ejecución:
```
. generar-compose.sh docker-compose-dev.yaml <nro clientes>
make docker-compose-up
docker ps
make docker-compose-logs
make docker-compose-down
```

***

### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `

#### Resolución:

Creo el archivo "docker-compose-validar-dev.yaml" el cual levanta un cliente netcat en la misma red que server, asumiendo que el contenedor de "server" ya esta corriendo. Este contenedor, "netcat-cli", corre un script llamado "do-validate-echo-server" cuyo log es parseado por el script "validar-echo-server.sh" el cual imprime el resultado de el mensaje al servidor.

#### Ejecución:

```
. generar-compose.sh docker-compose-dev.yaml 0
make docker-compose-up
. validar-echo-server.sh
make docker-compose-logs
make docker-compose-down
```

***

### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

#### Resolución:

En el cliente se crea un canal en el "main" al cual se notifican las señales las cuales son tomadas del canal en el loop principal de cliente en "client.go" para finalizar dicho loop correctamente.
En el servidor se crea una función que introduce un callback para el manejo de señales llamado "__handle_a_signal" el cual efectúa el log del evento cierra los socket's de cada uno de los clientes para luego finalizar el proceso servidor.

#### Ejecución:

```
. generar-compose.sh docker-compose-dev.yaml 1
make docker-compose-up
make docker-compose-logs
make docker-compose-down
```

***

## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).

#### Resolución:

* Para enviar la apuesta se utiliza el siguiente formato de mensaje:
```
<cant bytes><id agencia + apuesta (con datos apostador) como string utf8 separados por comas>
```
* Para confirmar la apuesta se utiliza el siguiente formato de mensaje:
```
<número de la apuesta como string utf8>
```
* Se tienen módulos que modelan la apuesta, aplican el parseo y realizan la comunicación tanto en el cliente (bet, parser y protocol) como en el servidor (utils/Bet y protocol/Protocol).

* Se modifico el script de generación del docker-compose para generar las variables entorno de la apuesta de cada cliente.

#### Ejecución:

```
. generar-compose.sh docker-compose-dev.yaml 5
make docker-compose-up
make docker-compose-logs
make docker-compose-down
```
NOTA: El script generar-compose permite exceder los 5 clientes generando clientes genéricos.
***

### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

#### Resolución:

* Para enviar la apuesta se utiliza el siguiente formato de mensaje ejemplificado a continuación con dos chunks desde el cliente:
```
<cant bytes><id agencia>
chunk 1: <cant bytes><apuesta 1 como string utf8>;......<apuesta m como string utf8>
chunk 2: <cant bytes><apuesta m+1 como string utf8>;......<apuesta n como string utf8>
<cant bytes><EOF,numero total de apuestas enviadas>
```
* Para confirmar las apuestas se utiliza el siguiente formato de mensaje desde el servidor:
```
<cantidad de apuestas como string>
```
* Por otro lado se modifica mi-generador.py para que pase el path del archivo para cada cliente el cuál se toma de la descompresión .data/dataset.zip del cuá se generan los volumenes correspondientes de cada archivo correspondiente al contenedor de cada cliente.

* Además se coloca en el main del cliente el código ncesario para levantar del config.yaml el número máximo de apuestas de cada lote.

#### Ejecución:

```
. generar-compose.sh docker-compose-dev.yaml 5
make docker-compose-up
make docker-compose-logs
make docker-compose-down
```

***

### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.

#### Resolución:

Del punto anterior ya arrastro lo siguiente:

* Para enviar la apuesta se utiliza el siguiente formato de mensaje por ejemplo con dos chunks desde el cliente:
```
<cant bytes><id agencia>
chunk 1: <cant bytes><apuesta 1 como string utf8>;......<apuesta m como string utf8>
chunk 2: <cant bytes><apuesta m+1 como string utf8>;......<apuesta n como string utf8>
<cant bytes><EOF,numero total de apuestas enviadas>
```
* Para confirmar las apuestas se utiliza el siguiente formato de mensaje desde el servidor:
```
<cantidad de apuestas como string>
```
* El cliente para consultar sobre los ganadores envía el siguiente mensaje:
```
<cant bytes><WINNERS>
```
* El servidor responde de la siguiente manera (suponiendo que los dni's de todos los ganadores entran en un sólo mensaje):
```
<dni winner 1>;<dni winner 2>;.....<dni winner n>
```
* En caso de no poder responder porque aún faltan clientes por terminar envía el siguiente mensaje:
```
<KEEP WAITING WINNERS>
```
* En el servidor hay una instancia de una clase de Lottery que almacena las agencias que esperan por
el resultado del sorteo, de esta manera al desconectar y conectar los clientes se sabe como continuar
ante la espera de mensajes desde el cliente.

* Modifique el generador de docker-compose para que le pase al servidor la cantidad de clientes por medio de una variable de entorno.

***

#### Ejecución:

```
. generar-compose.sh docker-compose-dev.yaml 5
make docker-compose-up
make docker-compose-logs
make docker-compose-down
```
***

#### Test:

```
REPO_PATH=/home/<usuario>/Workspace/tp0 pytest -s
```
***