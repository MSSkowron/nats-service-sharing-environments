import sys
import asyncio
import nats
import struct
import os

light_on = False


async def main(light_name, light_room):
    global light_on

    nc = await nats.connect("nats://admin:admin@nats3:4222")
    js = nc.jetstream()

    async def light_handler(msg):
        global light_on
        light_on = struct.unpack('?', msg.data)[0]
        print(f"Received message: {light_on}")

    await js.subscribe(f"lights.{light_room}.{light_name}.change", cb=light_handler)
    await js.subscribe(f"lights.{light_room}.change", cb=light_handler)
    await js.subscribe("lights.change", cb=light_handler)

    async def status_handler(msg):
        global light_on
        await js.publish("publishers.device", f"lights.{light_room}.{light_name}.{light_on}".encode())

    await js.subscribe("publishers.publisher", cb=status_handler)

    await js.publish("publishers.device", f"lights.{light_room}.{light_name}.{light_on}".encode())

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        print("Subscriber stopped.")

    await nc.close()


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: python main.py <light_name> <light_room>")
        sys.exit(1)

    name = os.environ.get('LIGHT_NAME')
    room = os.environ.get('LIGHT_ROOM')
    print(name,room)
    asyncio.run(main(name, room))
