
import asyncio
import websockets
import argparse


async def receive_client(target, id):
    async with websockets.connect(
            f'ws://{target}/receive/{id}') as websocket:

        for _ in range(100):
            msg = await websocket.recv()
            print(f"< {msg}")
if __name__ == "__main__":
    # ARGS
    parser = argparse.ArgumentParser()
    parser.add_argument("-i" ,"--id", nargs='?',type=str, default="1", help="id of the data channel")
    parser.add_argument("-t" ,"--target",  nargs='?',type=str,default="localhost:9009", help="target url")
    args = parser.parse_args()

    # RUN MAIN
    asyncio.get_event_loop().run_until_complete(receive_client(args.target, args.id))