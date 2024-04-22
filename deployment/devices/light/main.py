import sys
import asyncio
import nats
import struct

light_on = False


async def main(light_name, light_room):
    global light_on

    nc = await nats.connect("nats://admin:admin@nats3:4222")
    js = nc.jetstream()

    # to remove
    try:
        await js.delete_stream(name="device-stream")
    except:
        pass
    finally:
        await js.add_stream(name="device-stream",
                            subjects=[f"lights/{light_room}/{light_name}/change", f"lights/{light_room}/change",
                                      "lights/change", "publisher/add", "publisher/detect"])

    async def light_handler(msg):
        global light_on
        light_on = struct.unpack('?', msg.data)[0]
        print(f"Received message: {light_on}")

    async def status_handler(msg):
        global light_on
        await js.publish("publisher/add", f"lights/{light_room}/{light_name}/{light_on}".encode())

    await js.subscribe(f"lights/{light_room}/{light_name}/change", cb=light_handler)
    await js.subscribe(f"lights/{light_room}/change", cb=light_handler)
    await js.subscribe("lights/change", cb=light_handler)
    await js.subscribe("publisher/detect", cb=status_handler)

    await js.publish("publisher/add", f"lights/{light_room}/{light_name}/{light_on}".encode())

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        print("Subscriber stopped.")

    await nc.close()


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: python light.py <light_name> <light_room>")
        sys.exit(1)

    name = sys.argv[1]
    room = sys.argv[2]

    asyncio.run(main(name, room))
