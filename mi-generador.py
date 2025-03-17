import sys

ERROR_MESSAGE = "No se han ingresado los argumentos correctos: <archivo de salida> <cantidad clientes>"
FOOTER = "networks:\n  testing_net:\n    ipam:\n      driver: default\n      config:\n        - subnet: 172.25.125.0/24\n"

def get_header(amount_clients):
    return f"name: tp0\nservices:\n  server:\n    container_name: server\n    image: server:latest\n    entrypoint: python3 /main.py\n    environment:\n      - PYTHONUNBUFFERED=1\n      - AMOUNT_CLIENTS={amount_clients}\n    networks:\n      - testing_net\n    volumes:\n      - ./server/config.ini:/config.ini:ro\n"   

def do_generate_compose_file(args):
    amount_clients = int(args["clients"])
    file = open(args["file"], "w+")
    file.write(get_header(amount_clients))
    for client_id in range(1, int(args["clients"])+1):
        file.write(generate_client_code(client_id))
    file.write(FOOTER)
    file.close()

def generate_client_code(client_id):
    return f'  client{client_id}:\n    container_name: client{client_id}\n    image: client:latest\n    entrypoint: /client\n{generate_client_environment(client_id)}    networks:\n      - testing_net\n    volumes:\n      - ./client/config.yaml:/config.yaml:ro\n      - ./.data/agency-{client_id}.csv:/agency-{client_id}.csv:ro\n    depends_on:\n      - server\n'

def generate_client_environment(client_id):
    return f'    environment:\n      - CLI_ID={client_id}\n      - CLI_LOG_LEVEL=DEBUG\n      - CLI_DATA=agency-{client_id}.csv\n'

def generate_compose_file():
    args = parse_args()
    if args:
        do_generate_compose_file(args)

def parse_args():
    if len(sys.argv) >= 3:
        return {
            "file": sys.argv[1],
            "clients": sys.argv[2]
        }
    else:
        print(ERROR_MESSAGE)
        return
    
generate_compose_file()