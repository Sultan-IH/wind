import asyncio
import websockets
import time
import argparse

async def transmit(target, id):
    
    async with websockets.connect(
            f'ws://{target}/transmit/{id}') as websocket:

        for i in range(100):
            time.sleep(1)
            
            await websocket.send(f"{i}")
            print(f"> {i}")

if __name__ == "__main__":
    # ARGS
    parser = argparse.ArgumentParser()
    parser.add_argument("-i", "--id", nargs='?', type=str,default="1", help="id of the data channel")
    parser.add_argument("-t", "--target", nargs='?', type=str,default="localhost:9009", help="target url")
    args = parser.parse_args()

    asyncio.get_event_loop().run_until_complete(transmit(args.target, args.id))