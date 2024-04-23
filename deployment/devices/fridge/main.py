import sys
import asyncio
import nats
import struct

is_off = True


async def main(name, room):
    global is_off

    nc = await nats.connect("nats://admin:admin@nats3:4222")
    js = nc.jetstream()

    print("Fridge connected")

    async def fridge_handler(msg):
          global is_off
          is_off = struct.unpack('?', msg.data)[0]
          print(f"Received message: {is_off}")

    async def status_handler(msg):
        global is_off
        await js.publish("publisher.add", f"fridge.{name}.{not is_off}".encode())


    await js.subscribe(f"fridge.{room}.{name}.change", cb=fridge_handler)
    await js.subscribe(f"fridge.{room}.change", cb=fridge_handler)
    await js.subscribe("fridge.change", cb=fridge_handler)

    await js.subscribe("publishers.publisher", cb=status_handler)

    await js.publish("publishers.device", f"fridge.{room}.{name}.{is_off}".encode())

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        print("Subscriber stopped.")

    await nc.close()


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: python fridge.py <fridge_name> <fridge_room>")
        sys.exit(1)
    name = sys.argv[1]
    room = sys.argv[1]
    asyncio.run(main(name,room))
